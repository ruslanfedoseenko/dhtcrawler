package Models

type File struct {
	Id        int32  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	Path      string `gorm:"column:path;size:250"`
	Size      int    `gorm:"column:size"`
	TorrentId int32  `gorm:"index;column:torrent_id" json:"-"`
}

func (File) TableName() string {
	return "files"
}
