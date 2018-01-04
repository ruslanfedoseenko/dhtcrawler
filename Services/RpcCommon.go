package Services

import (
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/rpc"
)

var rpcLog = logging.MustGetLogger("RpcCommon")

func SetupRpc(app *Config.App) {
	if app.Config.RpcConfig.Mode == "" {
		rpcLog.Error("RpcConfig section missing")
		return
	}

	switch app.Config.RpcConfig.Mode {
	case Config.CLIENT:
		{
			rpcLog.Error("Invalid RpcConfig.Mode value")
		}
	case Config.SERVER:
		{
			Rpc.SetupRpcServer(app)
		}
	}
}
