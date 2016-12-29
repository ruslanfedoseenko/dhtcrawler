package Services

import (
	"encoding/hex"
	"fmt"
	"github.com/shiyanhui/dht"
	"strings"
	"github.com/lib/pq"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
	"github.com/op/go-logging"
)

type DhtCrawlingService struct {
	dht 	    *dht.DHT
	wire       *dht.Wire
	config     *dht.Config
	shouldStop bool
}
var dhtLog = logging.MustGetLogger("DhtCrawler")
var App *Config.App
var videExtractor *VideoInfoExtractor
var dhtCrawlingSvc DhtCrawlingService

func SetupDhtCrawling(app *Config.App) {
	App = app
	videExtractor = NewVideoInfoExtractor()
	dhtCrawlingSvc = DhtCrawlingService{}
	App.AddService(dhtCrawlingSvc)

}

func (svc DhtCrawlingService) Start() {
	for i := 0; i < App.Config.DhtConfig.Workers; i++ {
		w := dht.NewWire(65536, 6144, 3072)
		var config *dht.Config = dht.NewCrawlConfig()
		config.MaxNodes = 70000;
		config.Address = fmt.Sprintf(":%d", App.Config.DhtConfig.StartPort)
		dhtLog.Info("Starting Dht Crawling On ", App.Config.DhtConfig.StartPort)
		App.Config.DhtConfig.StartPort++
		config.OnAnnouncePeer = func(infoHash, ip string, port int) {
			var count int
			var infoHashStr = hex.EncodeToString([]byte(infoHash));
			App.Db.Model(&Models.Torrent{}).Where(Models.Torrent{Infohash: infoHashStr}).Count(&count)

			if count == 0 {
				w.Request([]byte(infoHash), ip, port)
			}
		}
		svc.dht = dht.New(config)
		go func() {
			for resp := range w.Response() {

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
				var name string;

				if name, ok = info["name"].(string); !ok {
					dhtLog.Error("Info section has name but it is not a string", info)
					continue
				}
				bt := Models.Torrent{
					Infohash: hex.EncodeToString(resp.InfoHash),
					Name: name,
					Leechers: -1,
					Seeds:    -1,
					LastScrape: pq.NullTime{
						Valid: false,
					},
				}
				var v interface{}
				if v, ok = info["files"]; ok {

					var files []interface{}
					if files, ok = v.([]interface{}); !ok{
						continue
					}

					bt.Files = make([]Models.File, len(files))

					for i, item := range files {
						f := item.(map[string]interface{})
						if f == nil || f["path"] == nil {
							continue
						}
						pathString := Utils.SliceToPathString(f["path"].([]interface{}))
						if strings.HasPrefix(pathString, "____padding") {
							continue
						}
						bt.Files[i] = Models.File{
							Path: pathString,
							Size: f["length"].(int),
						}

					}

				} else if _, ok := info["length"]; ok {
					bt.Files = make([]Models.File, 1)
					bt.Files[0] = Models.File{
						Path: bt.Name,
						Size: info["length"].(int),
					}

				}

				group := App.Classifier.Classify(bt)

				bt.GroupId = group.Id

				var count int

				App.Db.Model(&Models.Torrent{}).Where(Models.Torrent{Infohash: bt.Infohash}).Count(&count)

				if count > 0 {
					continue
				} else {
					//dhtLog.Println("Found new torrent", bt.Infohash, "group", bt.Group.Name, "group_id", bt.GroupId)
					/*for _, title := range bt.Titles {
						dhtLog.Println("Inserting title", title)
						App.Db.Debug().Exec("INSERT INTO `titles`(`id`, `title`, `poster`, `title_type`, `description`) " +
							"VALUES (?,?,?,?,?) ON DUPLICATE KEY " +
							"UPDATE `title`=VALUES(`title`),`poster`=VALUES(`poster`),`title_type`=VALUES(`title_type`),`description`=VALUES(`description`)", title.Id, title.Title,title.PosterUrl, title.TitleType, title.Description)
					}*/
					err = App.Db.Create(&bt).Error
					if err != nil {
						dhtLog.Error("Torrent Add to DB err:", err.Error())
					}
				}

			}
		}()
		go w.Run()
		go svc.dht .Run()

	}
}
