package Rpc

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/jinzhu/gorm"
	"errors"
	"github.com/op/go-logging"
)
var torrentRpcService = logging.MustGetLogger("TorrentRpcService")
type TorrentRpcService struct {
	db 		*gorm.DB
	classifier 	*Config.Classifier
}

func NewTorrentRpcService(app *Config.App)  *TorrentRpcService {
	return &TorrentRpcService{
		db: app.Db,
		classifier : app.Classifier,
	}
}

func (s *TorrentRpcService) HasTorrent(infoHash string) (b bool, err error){
	var count uint32 = 0
	err = s.db.Model(&Models.Torrent{}).Where(Models.Torrent{Infohash: infoHash}).Count(&count).Error

	b = (count != 0)
	torrentRpcService.Info("Handling TorrentRpcService.HasTorrent(", infoHash, ")=",b,err)
	return
}

func (s *TorrentRpcService) AddTorrent(torrent *Models.Torrent) (err error) {
	torrentRpcService.Info("Handling TorrentRpcService.AddTorrent(", torrent, ")")
	if ok, _ := s.HasTorrent(torrent.Infohash); ok {
		err = errors.New("Torrent with InfoHash" + torrent.Infohash + "already exists")
		return
	}
	group := s.classifier.Classify(*torrent)

	torrent.GroupId = group.Id
	tx:= s.db.Begin()
	tx.Exec("SET zombodb.batch_mode = true;")
	err = tx.Create(torrent).Error
	tx.Commit()

	return
}