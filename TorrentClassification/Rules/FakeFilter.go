package Rules

import (
	"github.com/ruslanfedoseenko/dhtcrawler/TorrentClassification"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

type FakeFilterRule struct{
	FakeCateGoryId int
	isLastFake bool
	TorrentClassification.Rule
}

func NewFakeFilter(app Config.App) (f FakeFilterRule){
	TorrentClassification.RuleID++;
	f = FakeFilterRule{
		TorrentClassification.Rule{
			Name: "FakeFilter",
			ID: TorrentClassification.RuleID,
		},
	}
	return
}

func  (f FakeFilterRule)ApplyRule(t *Models.Torrent) int{
	return -1;
}
func (f FakeFilterRule)IsFinalRule() bool{
	return false;
}