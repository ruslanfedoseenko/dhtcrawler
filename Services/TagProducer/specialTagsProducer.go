package TagProducer

import "github.com/ruslanfedoseenko/dhtcrawler/Models"

type specialTagProducer interface {
	SatisfyTag(torrentTokens []string) bool
	GetTags(torrent *Models.Torrent) []string
}
