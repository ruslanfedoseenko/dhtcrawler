package TorrentClassification

import "github.com/ruslanfedoseenko/dhtcrawler/Models"

var RuleID int = 0;


// If IsFinalRule == true than ApplyRule returns class index
type Rule struct {
	ID       int
	Name     string
	Children map[int]Rule
}

type TorrentApplicableRule interface {
	ApplyRule(t *Models.Torrent) int
	IsFinalRule() bool
}




