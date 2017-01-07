package tracker

import (
	"errors"
	"log"
	"net"
	"net/url"
	"strings"
)

type Peer struct {
	IP   net.IP
	Port int
}

type ScrapeRequest struct {
	InfoHashes []Hash
}
type ScrapeTorrentInfo struct {
	Seeders   int32 `bencode:"downloaded"`
	Completed int32 `bencode:"complete"`
	Leechers  int32 `bencode:"incomplete"`
}
type ScrapeResponse struct {
	ScrapeDatas map[string]ScrapeTorrentInfo `bencode:"files"`
}

var (
	ErrBadScheme = errors.New("unknown scheme")
)

func Scrape(tracker string, infoHashes []string) (r ScrapeResponse, err error) {
	_url, err := url.Parse(strings.Replace(tracker, "announce", "scrape", -1))
	if err != nil {
		return
	}
	var req ScrapeRequest
	hashesLen := len(infoHashes)
	for i := 0; i < hashesLen; i++ {
		var hash Hash
		err := hash.FromHexString(infoHashes[i])
		if err != nil {
			log.Println("Error: ", err.Error())
		}
		req.InfoHashes = append(req.InfoHashes, hash)
	}
	switch _url.Scheme {
	case "http", "https":
		return scrapeHTTP(&req, _url)
		//break
	case "udp":
		return scrapeUDP(&req, _url)
	default:
		err = ErrBadScheme
		return
	}
	return
}
