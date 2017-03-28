package TagProducer;

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"container/list"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"hash/crc32"
	"github.com/jinzhu/gorm"
	"strings"
	"github.com/op/go-logging"
	"gopkg.in/fatih/set.v0"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
	"time"
)

var TagProducerLog = logging.MustGetLogger("TagProducer")

type TorrentTagsProducer struct {
	tagProducers        *list.List
	specialTagProducers *list.List
	db                  *gorm.DB
}

var tagLog = logging.MustGetLogger("TagProducer")

func NewTorrentTagsProducer(app *Config.App) (p *TorrentTagsProducer) {
	p = &TorrentTagsProducer{
		tagProducers: list.New(),
		specialTagProducers: list.New(),
		db: app.Db,
	};
	p.initTokenExistenceProducers(app);
	p.initSpecialProducers(app);
	go p.processQueue();
	return;
}

func (t *TorrentTagsProducer) initSpecialProducers(app *Config.App) {
	t.specialTagProducers.PushBack(NewVideoTagProducer(app));
}
func (t *TorrentTagsProducer) initTokenExistenceProducers(app *Config.App) {
	for key, value := range app.Config.TagProducersConfig.TagProducers {
		producer := tokenExistenceTagProducer{
			tag: key,
			tokens: convertValues(value),
		}
		t.tagProducers.PushBack(producer)
	}
}
func convertValues(tokens []string) (res map[string]uint32) {
	res = make(map[string]uint32)
	tokenLen := len(tokens)
	for i := 0; i < tokenLen; i++ {
		res[tokens[i]] = 0;
	}
	return;
}
func tokenizer(c rune) bool {
	return strings.ContainsRune(";., -_|[]()/\\", c)
}

func (t *TorrentTagsProducer) processQueue() {
	for {
		var torrents []Models.Torrent;
		t.db.Table("torrents t").Preload("Files").
			Select("t.*").
			Joins("left join torrent_tags tt on tt.torrent_id = t.id").
			Where("tt.tag_id is null").
			Order("t.id asc").
			Limit(25).
			Find(&torrents)
		torrentsLen := len(torrents);
		if (torrentsLen == 0){
			TagProducerLog.Info("No torrents found waiting 5 secs....")
			<-time.After(time.Second * 5)
			continue
		}
		TagProducerLog.Info("Found", torrentsLen, "with no tags")
		for i := 0; i < torrentsLen; i++ {
			TagProducerLog.Info("Processing torrent:",torrents[i])
			t.FillTorrentTags(&torrents[i]);
			TagProducerLog.Info("Appending Tags:",torrents[i].Tags)
			t.db.Model(&torrents[i]).Association("Tags").Append(torrents[i].Tags)
		}
	}

}

func (t *TorrentTagsProducer) getTorrentTokens(torrent *Models.Torrent) []string {
	tokens := set.NewNonTS()
	for _, file := range torrent.Files {
		tokens.Add(Utils.ToInterfaceSlice(strings.FieldsFunc(strings.ToLower(file.Path), tokenizer))...)

	}
	tokens.Add(Utils.ToInterfaceSlice(strings.FieldsFunc(strings.ToLower(torrent.Name), tokenizer))...)
	return set.StringSlice(tokens);
}

func (t *TorrentTagsProducer) FillTorrentTags(torrent *Models.Torrent) {
	torrentTokens := t.getTorrentTokens(torrent);

	tagsBulkInsertArgs := make([]interface{}, 0, 20)
	tagsBulkInsertSql := "INSERT INTO tags (id, tag) VALUES "
	for e := t.tagProducers.Front(); e != nil; e = e.Next() {
		tagProducer := e.Value.(tokenExistenceTagProducer)
		if tagProducer.SatisfyTag(torrentTokens) {
			tag := tagProducer.GetTag();
			torrentTokens = append(torrentTokens, tag)
			tagModel := Models.Tag{
				Id: crc32.ChecksumIEEE(([]byte)(tag)),
				Tag: tag,
			}
			var tagCount int
			t.db.Model(&Models.Tag{}).Where(&Models.Tag{Id: tagModel.Id}).Count(&tagCount)
			if (tagCount == 0) {
				tagsBulkInsertArgs = append(tagsBulkInsertArgs, tagModel.Id, tagModel.Tag)
				tagsBulkInsertSql += " (?, ?),"
			}
			torrent.Tags = append(torrent.Tags, tagModel)

		}
	}
	for e := t.specialTagProducers.Front(); e != nil; e = e.Next() {
		tagProducer := e.Value.(specialTagProducer)

		if tagProducer.SatisfyTag(torrentTokens) {
			tags := tagProducer.GetTags(torrent);
			torrentTokens = append(torrentTokens, tags...)
			for i := 0; i < len(tags); i++ {
				tagModel := Models.Tag{
					Id: crc32.ChecksumIEEE(([]byte)(tags[i])),
					Tag: tags[i],
				}
				torrent.Tags = append(torrent.Tags, tagModel)
				var tagCount int
				t.db.Model(&Models.Tag{}).Where(&Models.Tag{Id: tagModel.Id}).Count(&tagCount)
				if (tagCount == 0) {
					tagsBulkInsertArgs = append(tagsBulkInsertArgs, tagModel.Id, tagModel.Tag)
					tagsBulkInsertSql += " (?, ?),"
				}
			}

		}

	}
	tagsBulkInsertSql = strings.TrimRight(tagsBulkInsertSql, ",")
	err := t.db.Raw(tagsBulkInsertSql, tagsBulkInsertArgs...).Error;
	if err != nil {
		tagLog.Error("Failed to pefrom bulk insert", err)
		return
	}
}

