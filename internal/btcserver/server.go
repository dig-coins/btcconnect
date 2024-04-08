package btcserver

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sgostarter/i/l"
)

type UnsignedTxNew struct {
	Wallet             string        `json:"wallet"`
	PayAddresses       []string      `json:"pay_addresses"`
	To                 []TransOutput `json:"to"`
	ConfirmationTarget int           `json:"confirmation_target"`
	ChangeAddress      string        `json:"change_address"`
	MinTransFlag       bool          `json:"min_trans_flag"`
}

func (p *UnsignedTxNew) IsValid() string {
	if len(p.PayAddresses) == 0 {
		return "invalid pay addresses"
	}

	if len(p.To) == 0 && p.ChangeAddress == "" {
		return "no output"
	}

	if p.ConfirmationTarget <= 0 {
		p.ConfirmationTarget = 1
	}

	return ""
}

func (s *BTCServer) httpServerRoutine() {
	logger := s.logger.WithFields(l.StringField(l.RoutineKey, "httpServerRoutine"))

	logger.Info("enter")

	defer logger.Info("leave")

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Recovery())
	r.Use(requestid.New())

	r.Any("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello")
	})

	r.POST("/tx/signed/broadcast", func(c *gin.Context) {
		var m map[string]string

		err := c.BindJSON(&m)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		tx := m["tx"]

		if tx == "" {
			c.String(http.StatusBadRequest, "no tx")

			return
		}

		txID, err := s.SendRawTransaction(tx)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"tx_id": txID,
		})
	})

	r.POST("/tx/unsigned/new", func(c *gin.Context) {
		var n UnsignedTxNew

		err := c.BindJSON(&n)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		if errMsg := n.IsValid(); errMsg != "" {
			c.String(http.StatusBadRequest, errMsg)

			return
		}

		unsignedTx, err := s.genUnsignedTx(n.Wallet, n.PayAddresses, n.To, n.ConfirmationTarget, n.ChangeAddress, n.MinTransFlag)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())

			return
		}

		c.JSON(http.StatusOK, share.Command{
			CommandType: share.CommandTypeGenTx,
			Input:       unsignedTx,
		})
	})

	r.GET("/wallets", s.handlerGetWallets)
	r.POST("/balance", s.handlerGetBalance)
	r.POST("/pay/simple", s.handlerSimplePay)
	r.POST("/pay", s.handlerPay)
	r.GET("/fee", s.handlerGetNetworkFee)
	r.POST("/unspent/group_wallet_address", s.handlerGetUnspentOfGroupWalletAddress)
	r.POST("/re-unsigned-tx", s.handlerResignedTx)
	r.POST("/unsigned-tx/load", s.handlerLoadUnsignedTx)

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

	listens := strings.Split(s.cfg.Listens, " ")

	for idx := 0; idx < len(listens)-1; idx++ {
		go fnListen(listens[idx])
	}

	fnListen(listens[len(listens)-1])
}
