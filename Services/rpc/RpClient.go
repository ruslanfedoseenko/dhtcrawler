package Rpc

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/valyala/gorpc"
	"fmt"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/op/go-logging"
)
var rpcClsLog = logging.MustGetLogger("RpcClient")
type RpcClient struct {
	dispatcher           *gorpc.Dispatcher
	tcpClient            *gorpc.Client
	torrentServiceClient *gorpc.DispatcherClient
	scrapeServiceClient *gorpc.DispatcherClient
}

var rpcClient *RpcClient = nil

func GetRpcCLientInstance(app *Config.App) *RpcClient{
	if (rpcClient == nil) {
		rpcClient = newRpcClient(app)
	}
	return rpcClient

}

func newRpcClient(app *Config.App) *RpcClient {
	var addressStr string = fmt.Sprintf("%s:%d", app.Config.RpcConfig.Host, app.Config.RpcConfig.Port)
	gorpc.RegisterType(&Models.Torrent{})
	rpcClient := RpcClient{
		dispatcher: gorpc.NewDispatcher(),
		tcpClient: gorpc.NewTCPClient(addressStr),
	}
	rpcClient.dispatcher.AddService(TorrentRpcServerName, &TorrentRpcService{})
	rpcClient.dispatcher.AddService(ScrapeRpcServerNAme, &ScrapeRpcService{})
	rpcClient.torrentServiceClient = rpcClient.dispatcher.NewServiceClient(TorrentRpcServerName, rpcClient.tcpClient)
	rpcClient.scrapeServiceClient = rpcClient.dispatcher.NewServiceClient(ScrapeRpcServerNAme, rpcClient.tcpClient)
	return &rpcClient
}


func (c *RpcClient) Start(){
	c.tcpClient.Start();
}


func (c *RpcClient) HasTorrent(infoHash string) (b bool, err error){
	response, err := c.torrentServiceClient.Call("HasTorrent", infoHash)
	rpcClsLog.Debug("HasTorrent response:",response, "err:",err);
	b, ok := response.(bool)
	if (!ok){
		rpcClsLog.Error("Failed to convert response", response, "to bool")
	}
	return;
}

func (c *RpcClient) AddTorrent(torrent *Models.Torrent) (err error) {
	response, err := c.torrentServiceClient.CallAsync("AddTorrent", torrent)
	rpcClsLog.Debug("AddTorrent response:",response, "err:",err);
	return nil
}

func (c *RpcClient) GetNextScrapeTask() *ScrapeTask{
	response,err := c.scrapeServiceClient.Call("GetNextScrapeTask", nil)
	rpcClsLog.Debug("GetNextScrapeTask response:",response, "err:",err);
	task, ok := response.(*ScrapeTask)
	if !ok {
		rpcClsLog.Error("Failed to convert response", response, "to ScrapeTask")
	}
	return task
}

func (c *RpcClient) ReportScrapeResults(reuslt *ScrapeResult) {
	response,err := c.scrapeServiceClient.Call("ReportScrapeResults", reuslt)
	rpcClsLog.Debug("ReportScrapeResults response:",response, "err:",err);
}
func (c *RpcClient)  HasAvailableTasks() bool {
	response, err := c.scrapeServiceClient.Call("HasAvailableTasks", nil)
	rpcClsLog.Debug("HasAvailableTasks response:",response, "err:",err);
	b, ok := response.(bool)
	if (!ok){
		rpcClsLog.Error("Failed to convert response", response, "to bool")
	}
	return b;
}