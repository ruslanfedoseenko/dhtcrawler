package Config

import (

	"github.com/jbrukh/bayesian"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
)

type Classifier struct {
	refreshMutex sync.Mutex
	classifier   *bayesian.Classifier
	Categories   []Models.GeneralCategory
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

func SetupBayesianClassification(Db *gorm.DB) (c *Classifier) {

	var classifier *bayesian.Classifier
	var categories []Models.GeneralCategory
	if _, err := os.Stat("/etc/dhtcrawler/cls.json"); os.IsNotExist(err) {
		classifier, categories = newClassifierFromDb(Db)
	} else {
		classifier, err = bayesian.NewClassifierFromFile("/etc/dhtcrawler/cls.json")
		if err != nil {
			log.Println("Failed to deserialize classifier", err)
			classifier, categories = newClassifierFromDb(Db)
		} else {
			categories, _ = readCategoriesFromDb(Db)
		}
	}

	c = &Classifier{
		Categories: categories,
		classifier: classifier,
	}
	return c
}

func readCategoriesFromDb(Db *gorm.DB) (categories []Models.GeneralCategory, categoryNames []bayesian.Class) {
	Db.Model(&Models.GeneralCategory{}).Where("parent_id IS NULL").Find(&categories)
	if len(categories) > 0 {
		for i, category := range categories {
			Db.Model(categories[i]).Related(&categories[i].TrainingData, "TrainingData")
			categories[i].Children = make([]Models.GeneralCategory, 10)
			Db.Model(categories[i]).Association("Children").Find(&categories[i].Children)
			for j, subCategory := range categories[i].Children {
				Db.Model(categories[i].Children[j]).Related(&categories[i].Children[j].TrainingData, "TrainingData")
				categories[i].Children[j].Name = category.Name + "\\" + subCategory.Name
				categoryNames = append(categoryNames, bayesian.Class(categories[i].Children[j].Name))
			}
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

			for _, subCategory := range category.Children {
				var trainingDataStr []string
				for _, trainingData := range subCategory.TrainingData {
					trainingDataStr = append(trainingDataStr, trainingData.Data)
				}
				log.Println("Learning category", subCategory.Name, "Data", trainingDataStr)
				classifier.Learn(trainingDataStr, bayesian.Class(subCategory.Name))
			}
		}
	}
	return
}
func tokenizer(c rune) bool {
	return strings.ContainsRune(";., -_|[]()/\\", c)
}
func (c *Classifier) Refresh(Db *gorm.DB) {
	c.refreshMutex.Lock()
	defer c.refreshMutex.Unlock()
	var categories []Models.GeneralCategory
	Db.Model(&Models.GeneralCategory{}).Where("parent_id IS NULL").Find(&categories)
	log.Println("Classifier General Categories:", categories)
	if len(categories) > 0 {
		var categoryNames []bayesian.Class
		for i, category := range categories {
			Db.Model(categories[i]).Related(&categories[i].TrainingData, "TrainingData")
			categories[i].Children = make([]Models.GeneralCategory, 10)
			Db.Model(categories[i]).Association("Children").Find(&categories[i].Children)
			for j, subCategory := range categories[i].Children {
				Db.Model(categories[i].Children[j]).Related(&categories[i].Children[j].TrainingData, "TrainingData")
				categories[i].Children[j].Name = category.Name + "\\" + subCategory.Name
				categoryNames = append(categoryNames, bayesian.Class(categories[i].Children[j].Name))
			}
			categoryNames = append(categoryNames, bayesian.Class(category.Name))
		}

		log.Println("Classifier General Categories:", categoryNames)

		classifier := bayesian.NewClassifier(categoryNames...)

		for _, category := range categories {

			var trainingDataStr []string

			for _, trainingData := range category.TrainingData {
				trainingDataStr = append(trainingDataStr, trainingData.Data)
			}
			log.Println("Learning category", category.Name, "Data", trainingDataStr)
			classifier.Learn(trainingDataStr, bayesian.Class(category.Name))

			for _, subCategory := range category.Children {
				var trainingDataStr []string
				for _, trainingData := range subCategory.TrainingData {
					trainingDataStr = append(trainingDataStr, trainingData.Data)
				}
				log.Println("Learning category", subCategory.Name, "Data", trainingDataStr)
				classifier.Learn(trainingDataStr, bayesian.Class(subCategory.Name))
			}
		}
		c.Categories = categories
		c.classifier = classifier
	}
}

func (c *Classifier) Classify(torrent Models.Torrent) Models.GeneralCategory {
	c.refreshMutex.Lock()
	defer c.refreshMutex.Unlock()
	var torrentTokens []string
	for _, file := range torrent.Files {
		ext := filepath.Ext(file.Path)
		if Utils.IndexOf(extStopList, ext) < 0 {
			torrentTokens = append(torrentTokens, strings.FieldsFunc(strings.ToLower(file.Path), tokenizer)...)
		}
	}
	torrentTokens = append(torrentTokens, strings.FieldsFunc(strings.ToLower(torrent.Name), tokenizer)...)
	//scores, classIndex, isStrict := c.classifier.LogScores(torrentTokens)
	_, classIndex, _ := c.classifier.LogScores(torrentTokens)

	var chosenClass string = string(c.classifier.Classes[classIndex])

	//log.Println(torrentTokens, "classifed as", chosenClass, "scores =", scores, "isStrict", isStrict)

	for _, category := range c.Categories {
		if category.Name == chosenClass {
			return category
		}
		for _, subCategory := range category.Children {
			if subCategory.Name == chosenClass {
				return subCategory
			}
		}
	}

	return Models.GeneralCategory{}

}
