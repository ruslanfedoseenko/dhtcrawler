package Rpc

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Services/TagProducer"
)

var torrentRpcService = logging.MustGetLogger("TorrentRpcService")

type TorrentRpcService struct {
	db          *gorm.DB
	tagProducer *TagProducer.TorrentTagsProducer
}

func NewTorrentRpcService(app *Config.App) *TorrentRpcService {
	return &TorrentRpcService{
		db:          app.Db,
		tagProducer: TagProducer.NewTorrentTagsProducer(app),
	}
}

func (s *TorrentRpcService) HasTorrent(infoHash string) (b bool, err error) {
	var count uint32 = 0
	err = s.db.Model(&Models.Torrent{}).Where(Models.Torrent{Infohash: infoHash}).Count(&count).Error

	b = (count != 0)
	torrentRpcService.Info("Handling TorrentRpcService.HasTorrent(", infoHash, ")=", b, err)
	return
}

func (s *TorrentRpcService) AddTorrent(torrent *Models.Torrent) (err error) {
	torrentRpcService.Info("Handling TorrentRpcService.AddTorrent(", torrent, ")")
	if ok, _ := s.HasTorrent(torrent.Infohash); ok {
		err = errors.New("Torrent with InfoHash" + torrent.Infohash + "already exists")
		return
	}

	//s.tagProducer.FillTorrentTags(torrent)

	tx := s.db.Begin()
	tx.Exec("SET zombodb.batch_mode = true;")
	err = tx.Create(torrent).Error
	tx.Commit()

	return
}
