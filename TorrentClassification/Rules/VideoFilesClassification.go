package Rules

import (
	"github.com/ruslanfedoseenko/dhtcrawler/TorrentClassification"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Services"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/op/go-logging"
)

var videoFilesClassificationLog = logging.MustGetLogger("VideoFilesClassification")

type VideoFilesClassification struct{
	extractor Services.VideoInfoExtractor
	TorrentClassification.Rule
}

func NewVideoFilesClassificationRule(app *Config.App) (b VideoFilesClassification) {
	TorrentClassification.RuleID++;
	b = VideoFilesClassification{
		extractor : Services.NewVideoInfoExtractor(),
		TorrentClassification.Rule{
			ID: TorrentClassification.RuleID,
			Name: "VideoFilesClassification",
		},
	}


	return
}

func (v *VideoFilesClassification) ApplyRule(t *Models.Torrent) int {
	titles := v.extractor.GetAssociatedVideos(Services.VideoInfoExtractWork{
		Name: t.Name,
	})
	videoFilesClassificationLog.Debug(titles);
	return titles[0].Ganres[0]

}
func (v *VideoFilesClassification)IsFinalRule() bool{
	return true;
}