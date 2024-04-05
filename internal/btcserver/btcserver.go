package btcserver

import (
	"context"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dig-coins/btcconnect/internal/redistorage"
	"github.com/dig-coins/btcconnect/internal/share"
	"github.com/dig-coins/btcconnect/pkg/rpclient"
	"github.com/patrickmn/go-cache"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/cuserror"
)

type Config struct {
	CoinType share.CoinType `json:"coin_type" yaml:"coin_type"`
	Listens  string         `json:"listens" yaml:"listens"`

	RPCHost     string `json:"rpc_host" yaml:"rpc_host"`
	RPCPort     int    `json:"rpc_port" yaml:"rpc_port"`
	RPCUser     string `json:"rpc_user" yaml:"rpc_user"`
	RPCPassword string `json:"rpc_password" yaml:"rpc_password"`
	RPCUseSSL   bool   `json:"rpc_use_ssl" yaml:"rpc_use_ssl"`

	MultiSignAddressInfos map[string]*share.MultiSignAddressInfo `json:"multi_sign_address_infos" yaml:"multi_sign_address_infos"`
}

func (cfg *Config) GetBTCNetParams() *chaincfg.Params {
	return share.GetBTCNetParams(cfg.CoinType)
}

type BTCServer struct {
	logger l.Wrapper
	ctx    context.Context
	wg     sync.WaitGroup

	ctxCancel context.CancelFunc
	cfg       *Config
	stg       redistorage.Storage

	rpcClient *rpclient.RpcClient

	cacheData *cache.Cache
}

func NewBTCServer(cfg *Config, stg redistorage.Storage, logger l.Wrapper) *BTCServer {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	if cfg == nil || cfg.RPCHost == "" {
		logger.Error("invalid config")

		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &BTCServer{
		logger:    logger,
		ctx:       ctx,
		ctxCancel: cancel,
		cfg:       cfg,
		stg:       stg,
		cacheData: cache.New(time.Minute, time.Minute),
	}

	err := s.init()
	if err != nil {
		return nil
	}

	return s
}

func (s *BTCServer) Wait() {
	s.wg.Wait()
}

func (s *BTCServer) init() (err error) {
	s.rpcClient, err = rpclient.NewRpcClient(s.cfg.RPCHost, s.cfg.RPCPort,
		s.cfg.RPCUser, s.cfg.RPCPassword, s.cfg.RPCUseSSL, 60)
	if err != nil {
		return
	}

	if s.cfg.Listens == "" {
		err = cuserror.NewWithErrorMsg("no listens")

		return
	}

	s.wg.Add(1)
	go s.httpServerRoutine()

	if s.stg != nil {
		s.wg.Add(1)

		go s.dataRoutine()
	}

	return
}
