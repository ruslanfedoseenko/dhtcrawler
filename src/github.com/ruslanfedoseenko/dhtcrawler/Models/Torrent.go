package Models

import (
	"github.com/lib/pq"
)

type Torrent struct {
	Id         int32           `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Infohash   string          `gorm:"size:40"`
	Leechers   int32           `gorm:"column:leechers"`
	Seeds      int32           `gorm:"column:seeds"`
	LastScrape pq.NullTime     `gorm:"column:last_scrape"`
	Name       string          `gorm:"column:name;size:250"`
	Group      GeneralCategory `gorm:"ForeignKey:GroupId"`
	GroupId    int32           `gorm:"index;column:group_id" json:"-"`
	Files      []File          `gorm:"ForeignKey:TorrentId" json:"-"`
	FilesTree  []FileTreeItem  `gorm:"-" json:"FilesTree"`
	Titles     []Title         `gorm:"many2many:torrent_to_title" json:"Titles,omitempty"`
}

func (Torrent) TableName() string {
	return "torrents"
}
