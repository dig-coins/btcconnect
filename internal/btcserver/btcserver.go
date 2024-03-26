package btcserver

import (
	"context"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/dig-coins/btcconnect/internal/share"
	"sync"

	"github.com/dig-coins/btcconnect/internal/redistorage"
	"github.com/dig-coins/btcconnect/pkg/rpclient"
	"github.com/sgostarter/i/l"
)

type Config struct {
	RPCHost     string `json:"rpc_host" yaml:"rpc_host"`
	RPCPort     int    `json:"rpc_port" yaml:"rpc_port"`
	RPCUser     string `json:"rpc_user" yaml:"rpc_user"`
	RPCPassword string `json:"rpc_password" yaml:"rpc_password"`
	RPCUseSSL   bool   `json:"rpc_use_ssl" yaml:"rpc_use_ssl"`
}

type BTCServer struct {
	logger l.Wrapper
	ctx    context.Context
	wg     sync.WaitGroup

	ctxCancel context.CancelFunc
	cfg       *Config
	stg       redistorage.Storage
	netParams *chaincfg.Params

	rpcClient            *rpclient.RpcClient
	multiSignAddressInfo map[string]*share.MultiSignAddressInfo
}

func NewBTCServer(cfg *Config, stg redistorage.Storage, netParams *chaincfg.Params,
	multiSignAddressInfo map[string]*share.MultiSignAddressInfo, logger l.Wrapper) *BTCServer {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &BTCServer{
		logger:               logger,
		ctx:                  ctx,
		ctxCancel:            cancel,
		cfg:                  cfg,
		stg:                  stg,
		netParams:            netParams,
		multiSignAddressInfo: multiSignAddressInfo,
	}

	_ = s.init()

	return s
}

func (s *BTCServer) Wait() {
	s.wg.Wait()
}

func (s *BTCServer) init() (err error) {
	s.rpcClient, err = rpclient.NewRpcClient(s.cfg.RPCHost, s.cfg.RPCPort, s.cfg.RPCUser, s.cfg.RPCPassword,
		s.cfg.RPCUseSSL, 60)
	if err != nil {
		return
	}
	/*
		s.wg.Add(1)

		go s.dataRoutine()
	*/
	return nil
}
