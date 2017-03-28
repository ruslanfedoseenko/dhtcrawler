package TagProducer

import (
	"github.com/ryanbradynd05/go-tmdb"
	"hash/crc32"

	"regexp"
	"strconv"
	"strings"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/jinzhu/gorm"
)
var videoExtractorLog = logging.MustGetLogger("VideoInfoExtractor")

type VideoInfoExtractWork struct {
	Id   int32
	Name string
}

type VideoInfoExtractor struct {
	tmDb              *tmdb.TMDb
	tmdbConfig        *tmdb.Configuration
	ganres		  map[uint32]string
	extraTitleContent *regexp.Regexp
	titleDelimiters   *regexp.Regexp
	db *gorm.DB
}
var instance *VideoInfoExtractor = nil
func NewVideoInfoExtractor(app *Config.App) (*VideoInfoExtractor) {
	if instance == nil {
		videoExtractorLog.Info("Creating new VideoInfoExtractor")
		tmDb := tmdb.Init("7ed1ada0530b0bbac6b697b818fc9c50")
		config, err := tmDb.GetConfiguration()
		if err != nil {
			videoExtractorLog.Error("Failed to get TMDB config", err.Error())
			return nil
		}
		instance = &VideoInfoExtractor{
			tmDb:              tmDb,
			tmdbConfig:        config,
			extraTitleContent: regexp.MustCompile("(?i)(\\d{4}|\\[([^]]+)]|\\(([^)]+)\\)|\\d(\\d)?x\\d(\\d)?(-\\d(\\d)?)?|web-dl|webdl|complete|temporada|season|episode|ep \\d+|tc|xxx|hdrip|dvdrip|bdrip|hdtv|1080p|1080|720|720p|480|480p|576|576p|xvid|divx|mkv|mp4|avi|brrip|ac3|mp3|x264|aac|s\\d(\\d)?(e\\d(\\d)?)?|bluray|rip|avc)"),
			titleDelimiters:   regexp.MustCompile("(\\.|-|;|,|-|_|\\||\\(|\\)|\\[|]|/|\\\\)"),
			db: app.Db,
		}
		var genres *tmdb.Genre;
		genres, err = tmDb.GetMovieGenres(map[string]string{})
		genresCount := len(genres.Genres)
		instance.ganres = make(map[uint32]string)
		for i:=0; i< genresCount; i++ {
			genre := genres.Genres[i]
			instance.ganres[uint32(genre.ID)] = genre.Name
		}

	}
	return instance
}

func minInt32(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func (ve *VideoInfoExtractor) GetAssociatedVideos(work VideoInfoExtractWork) (titles []Models.Title) {
	//videoExtractorLog.Println("Cleaning", name)
	cleanedName := ve.cleanupName(work.Name)

	videoExtractorLog.Info("Searcing for", cleanedName)
	res, err := ve.tmDb.SearchMulti(cleanedName, nil)
	if err != nil {
		videoExtractorLog.Error("SearchMulti Failed:", err.Error())
		//videoExtractorLog.Println("Re adding task")
		//ve.workQueue <-work
		return
	}

	//videoExtractorLog.Println("SearchMulti Result:", res)
	titlesCount := minInt32(int32(res.TotalResults), 5)
	var i int32
	for i = 0; i < titlesCount; i++ {
		var name string
		var yearStr string
		if res.Results[i].MediaType == "tv" {
			name = res.Results[i].OriginalName
			if len(res.Results[i].FirstAirDate) > 4 {
				yearStr = res.Results[i].FirstAirDate[0:4]
			}

		} else if res.Results[i].MediaType == "movie" {
			name = res.Results[i].OriginalTitle
			if len(res.Results[i].ReleaseDate) > 4 {
				yearStr = res.Results[i].ReleaseDate[0:4]
			}

		}

		var posterUrl string
		if len(res.Results[i].PosterPath) > 0 {
			posterUrl = ve.tmdbConfig.Images.BaseURL + "original" + res.Results[i].PosterPath
		}
		year, err := strconv.ParseInt(yearStr, 10, 32)
		if err != nil {
			year = 0
		}
		titles = append(titles, Models.Title{
			Title:       name,
			Year:        uint32(year),
			Description: res.Results[i].Overview,
			Ganres:		 ve.getGenres(res.Results[i].Genres),
			TitleType:   Models.TitleType(res.Results[i].MediaType),
			Id:          crc32.ChecksumIEEE([]byte(name)),
			PosterUrl:   posterUrl,
		})
	}
	return
}
func (ve *VideoInfoExtractor) getGenres(ganre_ids []uint32) (genreNames []string){
	ganresLen := len(ganre_ids)
	genreNames = make([]string, ganresLen, ganresLen)
	for i := 0; i < ganresLen; i++ {
		genreNames[i] = ve.ganres[ganre_ids[i]]
	}
	return
}

func (ve *VideoInfoExtractor) cleanupName(text string) string {
	videoExtractorLog.Info("data", text, ve)
	index := ve.extraTitleContent.FindStringIndex(text)
	if len(index) > 0 {
		if index[0] == 0 {
			text = text[index[1]:]
			index = ve.extraTitleContent.FindStringIndex(text)

		}
		if len(index) > 0 {
			text = text[:index[0]]
		}
	}
	cleaned := ve.extraTitleContent.ReplaceAllString(text, "")
	cleaned = ve.titleDelimiters.ReplaceAllString(cleaned, " ")
	cleaned = strings.Trim(cleaned, " ")
	return cleaned
}
