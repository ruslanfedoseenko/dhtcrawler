package Rpc

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/valyala/gorpc"
)

var rpcSrvLog = logging.MustGetLogger("RpcServer")

type RpcServer struct {
	dispatcher       *gorpc.Dispatcher
	torrentService   *TorrentRpcService
	scrapeRpcService *ScrapeRpcService
	server           *gorpc.Server
}

const (
	TorrentRpcServerName string = "TorrentsRpc"
	ScrapeRpcServerNAme         = "ScrapeRpc"
)

func SetupRpcServer(app *Config.App) {
	var addressStr string = fmt.Sprintf("%s:%d", app.Config.RpcConfig.Host, app.Config.RpcConfig.Port)
	rpcSrvLog.Debug("Address string:", addressStr)
	rpcServer := RpcServer{
		dispatcher:       gorpc.NewDispatcher(),
		torrentService:   NewTorrentRpcService(app),
		scrapeRpcService: NewScrapeRpcService(app),
	}
	rpcServer.dispatcher.AddService(TorrentRpcServerName, rpcServer.torrentService)
	rpcServer.dispatcher.AddService(ScrapeRpcServerNAme, rpcServer.scrapeRpcService)
	rpcServer.server = gorpc.NewTCPServer(addressStr, rpcServer.dispatcher.NewHandlerFunc())
	app.AddService(rpcServer)
}

func (s RpcServer) Start() {
	rpcSrvLog.Info("Starting RpcServer %s", s.server.Addr)
	s.server.Start()
}
