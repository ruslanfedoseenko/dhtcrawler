package TagProducer

import (
	"github.com/ruslanfedoseenko/go-tmdb"
	"hash/crc32"

	"github.com/jinzhu/gorm"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Config"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"regexp"
	"strconv"
	"strings"
	"reflect"
)

var videoExtractorLog = logging.MustGetLogger("VideoInfoExtractor")

type VideoInfoExtractWork struct {
	Id   int32
	Name string
}

type VideoInfoExtractor struct {
	tmDb              *tmdb.TMDb
	tmdbConfig        *tmdb.Configuration
	ganres            map[uint32]string
	extraTitleContent *regexp.Regexp
	titleDelimiters   *regexp.Regexp
	db                *gorm.DB
}

var instance *VideoInfoExtractor = nil

func NewVideoInfoExtractor(app *Config.App) *VideoInfoExtractor {
	if instance == nil {
		videoExtractorLog.Info("Creating new VideoInfoExtractor")
		tmDb := tmdb.Init(app.Config.TmdbApi.ApiKey)
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
			db:                app.Db,
		}
		var genres *tmdb.Genre
		genres, err = tmDb.GetMovieGenres(map[string]string{})
		genresCount := len(genres.Genres)
		instance.ganres = make(map[uint32]string)
		for i := 0; i < genresCount; i++ {
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

	videoExtractorLog.Debug("SearchMulti Result:", res)
	titlesCount := minInt32(int32(res.TotalResults), 5)
	var i int32

	for i = 0; i < titlesCount; i++ {
		var name string

		var posterUrl string
		var description string
		var titleType Models.TitleType
		var genres []string
		var year int64
		var ok = false
		var tvSeriesInfo *tmdb.MultiSearchTvInfo
		var movieInfo *tmdb.MultiSearchMovieInfo
		var base = &res.Results[i]
		videoExtractorLog.Debug("SearchMulti Result:", i, reflect.TypeOf(*base), *base)
		if tvSeriesInfo, ok = (*base).(*tmdb.MultiSearchTvInfo); ok {
			name = tvSeriesInfo.OriginalName
			var yearStr string
			if len(tvSeriesInfo.FirstAirDate) > 4 {
				yearStr = tvSeriesInfo.FirstAirDate[0:4]
			}
			genres = ve.getGenres(tvSeriesInfo.GenreIDs)

			if len(tvSeriesInfo.PosterPath) > 0 {
				posterUrl = ve.tmdbConfig.Images.BaseURL + "original" + tvSeriesInfo.PosterPath
			}
			year, err = strconv.ParseInt(yearStr, 10, 32)
			if err != nil {
				year = 0
			}

			description = tvSeriesInfo.Overview
			titleType = Models.TitleType(tvSeriesInfo.MediaType)
		} else if movieInfo, ok = (*base).(*tmdb.MultiSearchMovieInfo); ok {
			videoExtractorLog.Debug("Casted to tmdb.MultiSearchMovieInfo")
			name = movieInfo.OriginalTitle
			var yearStr string
			if len(movieInfo.ReleaseDate) > 4 {
				yearStr = movieInfo.ReleaseDate[0:4]
			}

			if len(movieInfo.PosterPath) > 0 {
				posterUrl = ve.tmdbConfig.Images.BaseURL + "original" + movieInfo.PosterPath
			}
			year, err = strconv.ParseInt(yearStr, 10, 32)
			if err != nil {
				year = 0
			}

			description = movieInfo.Overview
			titleType = Models.TitleType(movieInfo.MediaType)

		}
		if ok {
			titles = append(titles, Models.Title{
				Title:       name,
				Year:        uint32(year),
				Description: description,
				Ganres:      genres,
				TitleType:   titleType,
				Id:          crc32.ChecksumIEEE([]byte(name)),
				PosterUrl:   posterUrl,
			})
		}


	}
	return
}
func (ve *VideoInfoExtractor) getGenres(ganre_ids []uint32) (genreNames []string) {
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
