package TorrentClassification

import "github.com/ruslanfedoseenko/dhtcrawler/Models"

// If IsFinalRule == true than ApplyRule returns class index
type TorrentApplicableRule interface {
	ApplyRule(t *Models.Torrent) int
	IsFinalRule() bool
}

var RuleID int = 0;

type Rule struct {
	ID int
	Name string
	Children map[int]Rule
}


