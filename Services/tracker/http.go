package tracker

import (

	"fmt"
	"net/http"
	"net/url"
	"github.com/anacrolix/missinggo/httptoo"
	"log"

 	"github.com/zeebo/bencode"
)

type HttpScrapeResponse struct {
	Files map[string]ScrapeTorrentInfo `bencode:"files"`
}

func setAnnounceParams(uri *url.URL, sr *ScrapeRequest) {

	q := uri.Query()
	for _, infoHash := range sr.InfoHashes {
		q.Add("info_hash", string(infoHash[:]))
	}
	uri.RawQuery = q.Encode()
	log.Println("Final tracker Url", uri.RequestURI())

}

func scrapeHTTP(sr *ScrapeRequest, _url *url.URL) (ret ScrapeResponse, err error) {
	uri := httptoo.CopyURL(_url)
	setAnnounceParams(uri, sr)
	log.Println("Creating Http Req with URI:", uri.String())
	req, err := http.NewRequest("GET", uri.String(), nil)
	log.Println("Request:", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("response from tracker: %s:", resp.Status)
		return
	}
	var httpResponse HttpScrapeResponse;

	decoder := bencode.NewDecoder(resp.Body)
	err = decoder.Decode(&httpResponse)

	ret.ScrapeDatas = make(map[string]ScrapeTorrentInfo)
	//log.Println("Bytes:", string(bytesArr), "decoded as:")
	for key, value := range  httpResponse.Files {
		log.Println("InfoHash:", key, "(S C L)", value)
		var infoHash Hash;
		infoHash.FromString(key)
		ret.ScrapeDatas[infoHash.HexString()] = value
	}
	if err != nil {
		err = fmt.Errorf("error decoding  %s", err)
		return
	}
	return
}
