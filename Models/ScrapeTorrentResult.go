package Models

import "github.com/lib/pq"

type ScrapeTorrentResult struct {
	Id         int64       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	TrackerUrl string      `gorm:"column:tracker_url"`
	Seeds      int32       `gorm:"column:seeds"`
	Leaches    int32       `gorm:"column:leechers"`
	TorrentId  int32       `gorm:"column:torrent_id" json:"-"`
	LastUpdate pq.NullTime `gorm:"column:last_update"`
}

func (s *ScrapeTorrentResult) TableName() string {
	return "torrent_scrape_results"
}
