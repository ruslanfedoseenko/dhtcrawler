package TagProducer

import (
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
	"regexp"
)

type videoTagProducer struct {
	videoInfoExtractor *VideoInfoExtractor
	qualityInfo        *regexp.Regexp
}

func NewVideoTagProducer(app *Config.App) *videoTagProducer {
	return &videoTagProducer{
		videoInfoExtractor: NewVideoInfoExtractor(app),
		qualityInfo:        regexp.MustCompile("web-dl|webdl|season|episode|ep(\\s+)?\\d+|tc|xxx|hdrip|bdrip|dvdrip|bdrip|hdtv|1080p|1080|720p|720|480|480p|576|576p|xvid|divx|mkv|mp4|avi|brrip|ac3|mp3|x264|h264|x265|h265|aac|s\\d(\\d)?(e\\d(\\d)?)?|bluray|rip|avc"),
	}
}

func (vtp *videoTagProducer) SatisfyTag(torrentTokens []string) bool {
	if Utils.IndexOf(torrentTokens, "video") > -1 {
		return true
	}
	return false
}

func (vtp *videoTagProducer) GetTags(torrent *Models.Torrent) []string {
	TagProducerLog.Info("GetTags", vtp, vtp.videoInfoExtractor)
	titles := vtp.videoInfoExtractor.GetAssociatedVideos(VideoInfoExtractWork{Name: torrent.Name, })
	torrent.Titles = titles;
	titlesLen := len(titles)
	tags := make([]string, 0, titlesLen * 5)
	for i := 0; i < titlesLen; i++ {
		tags = append(tags, titles[i].Ganres...)
	}
	tags = append(tags, vtp.lookForQuality(torrent.Name)...)
	return tags
}
func (vtp *videoTagProducer) lookForQuality(i string) []string {
	return vtp.qualityInfo.FindAllString(i, -1)
}
