package Models

type Torrent struct {
	Id             int32                 `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Infohash       string                `gorm:"size:40"`
	Name           string                `gorm:"column:name;size:250"`
	HasTag         bool                  `gorm:"column:hastag;" json:"-"`
	ScraperResults []ScrapeTorrentResult `gorm:"ForeignKey:TorrentId" json:"TrackersInfo,omitempty"`
	Files          []File                `gorm:"ForeignKey:TorrentId" json:"Files"`
	FilesTree      []FileTreeItem        `gorm:"-" json:"FilesTree"`
	Titles         []Title               `gorm:"many2many:torrent_to_title" json:"Titles,omitempty"`
	Tags           []Tag                 `gorm:"many2many:torrent_tags" json:"Tags,omitempty"`
}

func (Torrent) TableName() string {
	return "torrents"
}
