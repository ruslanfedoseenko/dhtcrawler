package Rules

import (
	"github.com/ruslanfedoseenko/dhtcrawler/TorrentClassification"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/jbrukh/bayesian"
	"os"
	"log"
	"github.com/jinzhu/gorm"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"path/filepath"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
	"strings"
)

type BasicTypeClassification struct {
	classifier   *bayesian.Classifier
	Categories   []Models.GeneralCategory
	TorrentClassification.Rule
}

func NewBasicTypeClassificationRule(app *Config.App) (b BasicTypeClassification) {
	TorrentClassification.RuleID = TorrentClassification.RuleID+1;
	b = BasicTypeClassification{
		TorrentClassification.Rule{
			ID: TorrentClassification.RuleID,
			Name: "BasicTypeClassification",

		},
	}
	b.Children = make(map[int]TorrentClassification.Rule)
	var classifier *bayesian.Classifier
	var categories []Models.GeneralCategory
	if _, err := os.Stat("/etc/dhtcrawler/cls.json"); os.IsNotExist(err) {
		classifier, categories = newClassifierFromDb(app.Db)
	} else {
		classifier, err = bayesian.NewClassifierFromFile("/etc/dhtcrawler/cls.json")
		if err != nil {
			log.Println("Failed to deserialize classifier", err)
			classifier, categories = newClassifierFromDb(app.Db)
		} else {
			categories, _ = readCategoriesFromDb(app.Db)
		}
	}

	b.Categories = categories
	b.classifier = classifier
	return
}


func readCategoriesFromDb(Db *gorm.DB) (categories []Models.GeneralCategory, categoryNames []bayesian.Class) {
	Db.Model(&Models.GeneralCategory{}).Where("parent_id IS NULL").Find(&categories)
	if len(categories) > 0 {
		for i, category := range categories {
			Db.Model(categories[i]).Related(&categories[i].TrainingData, "TrainingData")
			categories[i].Children = make([]Models.GeneralCategory, 10)
			categoryNames = append(categoryNames, bayesian.Class(category.Name))
		}
	}
	return
}

func newClassifierFromDb(Db *gorm.DB) (classifier *bayesian.Classifier, categories []Models.GeneralCategory) {

	categories, categoryNames := readCategoriesFromDb(Db)

	log.Println("Classifier General Categories:", categoryNames)

	if len(categories) > 0 {
		classifier = bayesian.NewClassifier(categoryNames...)

		for _, category := range categories {

			var trainingDataStr []string

			for _, trainingData := range category.TrainingData {
				trainingDataStr = append(trainingDataStr, trainingData.Data)
			}
			log.Println("Learning category", category.Name, "Data", trainingDataStr)
			classifier.Learn(trainingDataStr, bayesian.Class(category.Name))

		}
	}
	return
}
func tokenizer(c rune) bool {
	return strings.ContainsRune(";., -_|[]()/\\", c)
}

var extStopList []string = []string{
	".url", ".txt", ".ico", ".srt", ".gif", ".log",
	".nfo", ".ass", ".lnk", ".rtf", ".bc!",
	".bmp", ".m3u", ".mht", ".cue", ".sfv", ".diz",
	".azw3", ".odt", ".chm", ".md5", ".idx", ".sub",
	".ini", ".html", ".ssa", ".lit", ".xml", ".clpi",
	".bup", ".ifo", ".htm", ".info", ".css", ".php",
	".js", ".jar", ".json", ".sha", ".docx", ".csv",
	".scr", ".inf", ".hdr", ".prq", ".isn", ".inx", ".tpl",
	".aco", ".opa", ".dpc", ".qdl2", ".acf", ".cdx",
	".iwd", ".ff", ".tmp", ".asi", ".flt", ".cfg",
	".tdl", ".tta", ".ape", ".btn", ".sig", ".sql", ".db",
	".zdct", ".bak", ".fxp", ".nxp", ".nsk", ".256",
	".mpls", ".clpi", ".bdmv", ".cdd", ".dbf",
	".vmx", ".vmsd", ".vmxf", ".nvram",
}
func (b *BasicTypeClassification) ApplyRule(torrent *Models.Torrent) int {


	var torrentTokens []string
	for _, file := range torrent.Files {
		ext := filepath.Ext(file.Path)
		if Utils.IndexOf(extStopList, ext) < 0 {
			torrentTokens = append(torrentTokens, strings.FieldsFunc(strings.ToLower(file.Path), tokenizer)...)
		}
	}
	torrentTokens = append(torrentTokens, strings.FieldsFunc(strings.ToLower(torrent.Name), tokenizer)...)
	//scores, classIndex, isStrict := c.classifier.LogScores(torrentTokens)
	scores, classIndex, isStrict := b.classifier.LogScores(torrentTokens)

	var chosenClass string = string(b.classifier.Classes[classIndex])

	log.Println(torrentTokens, "classifed as", chosenClass, "scores =", scores, "isStrict", isStrict)



	return Models.GeneralCategory{}
}
func (b *BasicTypeClassification)IsFinalRule() bool{
	return false
}

