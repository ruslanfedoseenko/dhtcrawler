package Services

import (
	"github.com/ryanbradynd05/go-tmdb"
	"hash/crc32"

	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
	"github.com/op/go-logging"
)
var videoExtractorLog = logging.MustGetLogger("VideoInfoExtractor")

type VideoInfoExtractWork struct {
	Id   int32
	Name string
}

type VideoInfoExtractor struct {
	tmDb              *tmdb.TMDb
	tmdbConfig        *tmdb.Configuration
		ganres			map[uint32]string
	extraTitleContent *regexp.Regexp
	titleDelimiters   *regexp.Regexp
}

func NewVideoInfoExtractor() (ve *VideoInfoExtractor) {
	tmDb := tmdb.Init("7ed1ada0530b0bbac6b697b818fc9c50")
	config, err := tmDb.GetConfiguration()
	if err != nil {
		videoExtractorLog.Error("Failed to get TMDB config", err.Error())
		return nil
	}
	extractor := VideoInfoExtractor{
		tmDb:              tmDb,
		tmdbConfig:        config,
		extraTitleContent: regexp.MustCompile("(?i)(\\d{4}|\\[([^]]+)]|\\(([^)]+)\\)|\\d(\\d)?x\\d(\\d)?(-\\d(\\d)?)?|web-dl|webdl|complete|temporada|season|episode|ep \\d+|tc|xxx|hdrip|bdrip|dvdrip|bdrip|hdtv|1080p|1080|720|720p|480|480p|576|576p|xvid|divx|mkv|mp4|avi|brrip|ac3|mp3|x264|aac|s\\d(\\d)?(e\\d(\\d)?)?|bluray|rip|avc)"),
		titleDelimiters:   regexp.MustCompile("(\\.|-|;|,|-|_|\\||\\(|\\)|\\[|]|/|\\\\)"),
	}
	var ganres *tmdb.Genre;
	ganres, err = tmDb.GetMovieGenres(map[string]string{})
	ganresCount := len(ganres.Genres)
	extractor.ganres = make(map[uint32]string)
	for i:=0; i< ganresCount; i++ {
		ganre := ganres.Genres[i]
		extractor.ganres[uint32(ganre.ID)] = ganre.Name
	}
	go extractor.processQueue()
	return &extractor
}

func (ve *VideoInfoExtractor) processQueue() {
	var counters Models.Counters
	App.Db.First(&counters)
	videoExtractorLog.Debug("Counters", counters)
	for {
		var work VideoInfoExtractWork
		var torrent Models.Torrent
		App.Db.Raw("SELECT id, Name from torrents where id > ? and group_id in (1,2) LIMIT 1", counters.LastExtractedVideoId).Scan(&torrent)
		work.Id = torrent.Id
		work.Name = torrent.Name
		videoExtractorLog.Info("Found work:", work)
		if work.Id == 0 {
			videoExtractorLog.Info("No work found. Waiting 5 sec..")
			time.Sleep(5 * time.Second)
			continue
		}

		videos := ve.GetAssociatedVideos(work)
		App.Db.Model(&Models.Counters{}).Update(&Models.Counters{
			LastExtractedVideoId: work.Id,
		})
		counters.LastExtractedVideoId = work.Id
		if len(videos) > 0 {
			for _, title := range videos {
				videoExtractorLog.Info("Inserting title", title)
				var titleCount int
				App.Db.Model(&Models.Title{}).Where("id = ?", title.Id).Count(&titleCount)
				if (titleCount == 0) {
					App.Db.Create(&title)
				}
			}
			App.Db.Model(&Models.Torrent{
				Id: work.Id,
			}).Association("Titles").Append(videos)
		}

	}

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
			Ganres:		 ve.getGanres(res.Results[i].Genres),
			TitleType:   Models.TitleType(res.Results[i].MediaType),
			Id:          crc32.ChecksumIEEE([]byte(name)),
			PosterUrl:   posterUrl,
		})
	}
	return
}
func (ve *VideoInfoExtractor) getGanres(ganre_ids []uint32) (ganreNames []string){
	len := len(ganre_ids)
	ganreNames = make([]string, len, len)
	for i := 0; i < len; i++ {
		ganreNames[i] = ve.ganres[ganre_ids[i]]
	}
	return
}

func (ve *VideoInfoExtractor) cleanupName(text string) string {

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
