package Models

import (
	"time"
)

type TorrentStatsPart struct {
	Id             int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	TorrentsCount  int       `gorm:"column:torrents_count"`
	FilesCount     int       `gorm:"column:files_count"`
	TotalFilesSize int64     `gorm:"column:files_size"`
	Date           time.Time `gorm:"column:date"`
}

func (TorrentStatsPart) TableName() string {
	return "torrent_stats"
}
