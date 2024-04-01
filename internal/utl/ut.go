package utl

import (
	"testing"

	"github.com/sgostarter/libeasygo/ut"
	"github.com/spf13/cast"
)

type UTConfig struct {
	cfg *ut.Config
}

func Setup() *UTConfig {
	return &UTConfig{
		cfg: ut.SetupUTConfigEx("ut.yaml", []string{"../../"}),
	}
}

func (cfg *UTConfig) GetBTCRpcHost(t *testing.T) string {
	return cast.ToString(cfg.cfg.GetEnv(t, "rHost"))
}

func (cfg *UTConfig) GetBTCRpcPort(t *testing.T) int {
	return cast.ToInt(cfg.cfg.GetEnv(t, "rPort"))
}

func (cfg *UTConfig) GetBTCRpcUser(t *testing.T) string {
	return cast.ToString(cfg.cfg.GetEnv(t, "rUser"))
}

func (cfg *UTConfig) GetBTCRpcPassword(t *testing.T) string {
	return cast.ToString(cfg.cfg.GetEnv(t, "rPassword"))
}

func (cfg *UTConfig) GetBTCRpcUseTLS(t *testing.T) bool {
	return cast.ToBool(cfg.cfg.GetEnv(t, "rUseTLS"))
}

func (cfg *UTConfig) GetSSeedsFileSecKey(t *testing.T) string {
	return cast.ToString(cfg.cfg.GetEnv(t, "sSeedsKey"))
}

func (cfg *UTConfig) GetWallet(t *testing.T) string {
	return cast.ToString(cfg.cfg.GetEnv(t, "rWallet"))
}
