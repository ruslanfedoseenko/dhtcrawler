package TorrentClassification

import (
	"github.com/ruslanfedoseenko/dhtcrawler/TorrentClassification/Rules"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
)

type RuleBasedTorrentClassifier struct {
	Rules []Rule
	Categories []string
}


func NewRuleBasedTorrentClassifier(app *Config.App) (c RuleBasedTorrentClassifier)  {
	c = new(RuleBasedTorrentClassifier)
	buildRulesTree(&c, app)
	return c
}
func buildRulesTree(classifier *RuleBasedTorrentClassifier,app *Config.App) {

	var basicTypeClassificationRule Rule = Rules.NewBasicTypeClassificationRule(app)
	append(classifier.Rules, basicTypeClassificationRule)
}