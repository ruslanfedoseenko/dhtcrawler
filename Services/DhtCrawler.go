package Services

import (
	"encoding/hex"
	"fmt"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/rpc"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
	"github.com/ruslanfedoseenko/dht"
	"unicode/utf8"
	"strings"
	"github.com/saintfish/chardet"
	"gopkg.in/iconv.v1"
	"strconv"
)

type DhtCrawlingService struct {
	dht        []*dht.DHT
	//wire       *dht.Wire
	config     *dht.Config
	rpcClient  *Rpc.RpcClient
	shouldStop bool
}

var dhtLog = logging.MustGetLogger("DhtCrawler")
var App *Config.App
var dhtCrawlingSvc DhtCrawlingService

func SetupDhtCrawling(app *Config.App) {
	App = app
	dhtCrawlingSvc = DhtCrawlingService{
		rpcClient: Rpc.GetRpcCLientInstance(app),
	}
	App.AddService(dhtCrawlingSvc)

}
var charsetAliases = map[string]string{
	"GB-18030": "GB18030",
	"ISO-8859-8-I": "ISO-8859-8",
	"IBM420_ltr" : "IBM-420",
	"IBM420_rtl" : "IBM-420",
	"IBM424_ltr" : "IBM-424",
	"IBM424_rtl" : "IBM-424",

}
var detector = chardet.NewTextDetector()

func (svc DhtCrawlingService) Start() {
	svc.rpcClient.Start()

	svc.dht = make([]*dht.DHT, App.Config.DhtConfig.Workers)
	for i := 0; i < App.Config.DhtConfig.Workers; i++ {
		wire := dht.NewWire(65536, 6144, 6144)
		var config = dht.NewCrawlConfig()
		config.MaxNodes = 70000
		config.Address = fmt.Sprintf(":%d", App.Config.DhtConfig.StartPort+i)
		dhtLog.Info("Starting Dht Crawling On ", App.Config.DhtConfig.StartPort+i)

		config.OnAnnouncePeer = func(infoHash, ip string, port int) {
			go func(infoHash, ip string, port int) {
				var infoHashStr = hex.EncodeToString([]byte(infoHash))

				if ok, _ := svc.rpcClient.HasTorrent(infoHashStr); !ok {
					wire.Request([]byte(infoHash), ip, port)
				}
			}(infoHash, ip, port)

		}
		svc.dht[i] = dht.New(config)

		go func() {
			for resp := range wire.Response() {

				metadata, err := dht.Decode(resp.MetadataInfo)
				if err != nil {
					continue
				}
				var info map[string]interface{}
				var ok bool
				if info, ok = metadata.(map[string]interface{}); !ok {
					continue
				}

				if _, ok = info["name"]; !ok {
					continue
				}
				var name string

				if name, ok = info["name"].(string); !ok {
					dhtLog.Error("Info section has name but it is not a string", info)
					continue
				}
				r, err := detector.DetectBest([]byte(name))
				if err == nil {
					if r.Charset != "UTF-8" {
						name = convertStringToUtf8(r, name)
					}
				} else {
					dhtLog.Info("Error Detecting charset", err)
				}
				if !utf8.ValidString(name) {
					dhtLog.Error("Name string is not valid utf8 string:", name)
				}
				bt := Models.Torrent{
					Infohash: hex.EncodeToString(resp.InfoHash),
					Name:     name,
				}
				var v interface{}
				if v, ok = info["files"]; ok {

					var files []interface{}
					if files, ok = v.([]interface{}); !ok {
						continue
					}

					bt.Files = make([]Models.File, len(files))

					for i, item := range files {
						var f map[string]interface{}
						if f, ok = item.(map[string]interface{}); !ok {
							continue
						}
						if f == nil || f["path"] == nil {
							continue
						}
						pathString := Utils.SliceToPathString(f["path"].([]interface{}))
						if strings.Contains(pathString, "___padding") {
							continue
						}
						r, err := detector.DetectBest([]byte(pathString))
						if err == nil {
							if r.Charset != "UTF-8" {
								pathString = convertStringToUtf8(r, pathString)
							}
						} else {
							dhtLog.Info("Error Detecting charset", err)
						}
						if !utf8.ValidString(pathString) {
							dhtLog.Error("Path string is not valid utf8 string:", pathString)
						}
						var size int
						if size, ok = f["length"].(int); !ok {
							if sizeStr, ok := f["length"].(string); ok {
								size,err = strconv.Atoi(sizeStr)
								if err != nil {
									dhtLog.Error("Failed to convert", sizeStr, "to int")
									size = -1
								}
							} else {
								size = -1
							}
						}
						bt.Files[i] = Models.File{
							Path: pathString,
							Size: size,
						}

					}

				} else if _, ok := info["length"]; ok {
					bt.Files = make([]Models.File, 1)
					bt.Files[0] = Models.File{
						Path: bt.Name,
						Size: info["length"].(int),
					}

				}

				svc.rpcClient.AddTorrent(&bt)
				/*Avar count int

				pp.Db.Model(&Models.Torrent{}).Where(Models.Torrent{Infohash: bt.Infohash}).Count(&count)

				if count > 0 {
					continue
				} else {
					//dhtLog.Println("Found new torrent", bt.Infohash, "group", bt.Group.Name, "group_id", bt.GroupId)
					for _, title := range bt.Titles {
						dhtLog.Println("Inserting title", title)
						App.Db.Debug().Exec("INSERT INTO `titles`(`id`, `title`, `poster`, `title_type`, `description`) " +
							"VALUES (?,?,?,?,?) ON DUPLICATE KEY " +
							"UPDATE `title`=VALUES(`title`),`poster`=VALUES(`poster`),`title_type`=VALUES(`title_type`),`description`=VALUES(`description`)", title.Id, title.Title,title.PosterUrl, title.TitleType, title.Description)
					}
					err = App.Db.Create(&bt).Error
					if err != nil {
						dhtLog.Error("Torrent Add to DB err:", err.Error())
					}
				}*/

			}
		}()
		go wire.Run()
		go svc.dht[i].Run()
	}
}
func convertStringToUtf8(r *chardet.Result, name string) string {
	charset := r.Charset
	if alias, ok := charsetAliases[charset]; ok {
		charset = alias
	}
	converter, err := iconv.Open("UTF-8", charset)
	if err != nil {
		dhtLog.Error("Unable to convert string", name, " from", charset, "to UTF-8")
	} else {
		defer converter.Close()
		name = converter.ConvString(name)
	}
	return name
}
