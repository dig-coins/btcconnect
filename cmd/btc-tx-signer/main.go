package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dig-coins/btcconnect/internal/btctx"
	"github.com/dig-coins/btcconnect/internal/config"
	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/dig-coins/btcconnect/internal/txsigner"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/l"
)

func main() {
	var configFile, commandS string

	flag.StringVar(&configFile, "config", "btc-tx-signer-1.yaml", "config file path")
	flag.StringVar(&commandS, "command", "", "hex command string")
	flag.Parse()

	logger := l.NewConsoleLoggerWrapper()
	logger.GetLogger().SetLevel(l.LevelDebug)

	if configFile != "" {
		logger.Info("config file is " + configFile)
	}

	cfg := config.GetBTCTxSignerConfig(configFile)

	signer, err := txsigner.NewTxSigner(config.GetHDWalletCoinType(cfg.CoinType),
		config.GetBTCNetParams(cfg.CoinType), cfg.SeedFileName, cfg.SeedSecKey,
		cfg.MultiSignAddressInfos)
	if err != nil {
		logger.Fatal(err)
	}

	if commandS != "" {
		cr := processCommandS(signer, config.GetBTCNetParams(cfg.CoinType), commandS, logger)
		if cr.ErrMessage != "" {
			logger.Error("exception: " + cr.ErrMessage)

			return
		}

		var result string

		result, err = share.MarshalCommandResult(cr)
		if err != nil {
			logger.Error("marshal command result failed: " + err.Error())

			return
		}

		logger.Info("result: " + result + "\n")

		return
	}

	runCommandServer(signer, config.GetBTCNetParams(cfg.CoinType), cfg.Listens, logger)
}

func dbgCommand(command *share.Command, netParams *chaincfg.Params, logger l.Wrapper) (ok bool) {
	fnDbgUnsignedTx := func(unsignedTx *btctx.UnsignedTx) {
		for idx, input := range unsignedTx.Inputs {
			logger.Infof("  INPUT[%d]  %s %s - %d amount is %d\n", idx, input.Address,
				input.TxID, input.VOut, input.Amount)
		}

		for idx, output := range unsignedTx.Outputs {
			logger.Infof("  OUTPUT[%d]  %s amount is %d\n", idx, output.Address, output.Amount)
		}
	}

	switch command.CommandType {
	case share.CommandTypeGenTx:
		unsignedTx, err := btctx.UnmarshalUnsignedTx(command.Input)
		if err != nil {
			logger.Error("failed on UnmarshalUnsignedTx")

			return
		}

		logger.Info("dbg incoming command: unsignedTx")
		fnDbgUnsignedTx(unsignedTx)
	case share.CommandTypeUpdateTx:
		msTx, err := txsigner.UnmarshalMiddleSignMidTx(command.Input)
		if err != nil {
			logger.Error("failed on UnmarshalMiddleSignMidTx")

			return
		}

		if !msTx.Check(netParams) {
			logger.Error("invalid msTx")

			return
		}

		fnDbgUnsignedTx(msTx.GetUnsignedTx())
	default:
		logger.Error(fmt.Sprintf("invald command type: %d", command.CommandType))

		return
	}

	ok = true

	return
}

func processCommandS(signer *txsigner.TxSigner, netParams *chaincfg.Params, commandS string, logger l.Wrapper) (cr share.CommandResult) {
	command, err := share.UnmarshalCommand(commandS)
	if err != nil {
		cr.ErrMessage = err.Error()

		return
	}

	return processCommand(signer, netParams, &command, logger)
}

func processCommand(signer *txsigner.TxSigner, netParams *chaincfg.Params, command *share.Command, logger l.Wrapper) (cr share.CommandResult) {
	if !command.Valid() {
		cr.ErrMessage = "bad command"

		return
	}

	dbgCommand(command, netParams, logger)

	var err error

	switch command.CommandType {
	case share.CommandTypeGenTx:
		cr.Tx, cr.AllSignedFlag, err = signer.SignTx(command.Input)
	case share.CommandTypeUpdateTx:
		cr.Tx, cr.AllSignedFlag, err = signer.UpdateMiddleSignedTxHex(command.Input)
	default:
		err = commerr.ErrUnimplemented
	}

	if err != nil {
		cr.ErrMessage = err.Error()
	}

	return
}

func runCommandServer(signer *txsigner.TxSigner, netParams *chaincfg.Params, listenAddresses string, logger l.Wrapper) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Recovery())
	r.Use(requestid.New())

	r.Any("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	})

	r.POST("/command/encrypt", func(c *gin.Context) {
		var command share.Command

		err := c.BindJSON(&command)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		if !command.Valid() {
			c.String(http.StatusBadRequest, err.Error())

			return
		}

		s, err := share.MarshalCommand(command)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.String(http.StatusOK, s)
	})

	r.POST("/command-result/decrypt", func(c *gin.Context) {
		defer c.Request.Body.Close()

		d, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		cr, err := share.UnmarshalCommandResult(string(d))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.JSON(http.StatusOK, cr)
	})

	r.POST("/sign/json", func(c *gin.Context) {
		var command share.Command

		err := c.BindJSON(&command)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		if !command.Valid() {
			c.String(http.StatusBadRequest, err.Error())

			return
		}

		cr := processCommand(signer, netParams, &command, logger)
		c.JSON(http.StatusOK, cr)
	})

	r.POST("/sign", func(c *gin.Context) {
		defer c.Request.Body.Close()

		d, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		s, err := share.MarshalCommandResult(processCommandS(signer, netParams, string(d), logger))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.String(http.StatusOK, s)
	})

	fnListen := func(listen string) {
		srv := &http.Server{
			Addr:        listen,
			ReadTimeout: time.Second,
			Handler:     r,
		}

		logger.Info("listen on ", listen)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}

	listens := strings.Split(listenAddresses, " ")

	for idx := 0; idx < len(listens)-1; idx++ {
		go fnListen(listens[idx])
	}

	fnListen(listens[len(listens)-1])
}
