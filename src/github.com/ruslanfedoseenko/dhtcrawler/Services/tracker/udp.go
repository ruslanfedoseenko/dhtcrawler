package tracker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/url"
	"time"
	"github.com/ruslanfedoseenko/dhtcrawler/Utils"
)

type Action int32

const (
	ActionConnect Action = iota
	ActionAnnounce
	ActionScrape
	ActionError

	connectRequestConnectionId = 0x41727101980

	// BEP 41
	optionTypeEndOfOptions = 0
	optionTypeNOP          = 1
	optionTypeURLData      = 2
)

type ConnectionRequest struct {
	ConnectionId int64
	Action       int32
	TransctionId int32
}

type ConnectionResponse struct {
	ConnectionId int64
}

type ResponseHeader struct {
	Action        Action
	TransactionId int32
}

type RequestHeader struct {
	ConnectionId  int64
	Action        Action
	TransactionId int32
} // 16 bytes

type udpScrapeResponse struct {
	torrentDatas []ScrapeTorrentInfo
}

type udpScrapeRequest struct {
	RequestHeader
	InfoHashes []Hash
}

func newTransactionId() int32 {
	return int32(rand.Uint32())
}

func timeout(contiguousTimeouts int) (d time.Duration) {
	if contiguousTimeouts > 8 {
		contiguousTimeouts = 8
	}
	d = 15 * time.Second
	for ; contiguousTimeouts > 0; contiguousTimeouts-- {
		d *= 2
	}
	return
}

type udpAnnounce struct {
	contiguousTimeouts   int
	connectionIdReceived time.Time
	connectionId         int64
	socket               net.Conn
	url                  url.URL
}

func (c *udpAnnounce) Close() error {
	if c.socket != nil {
		log.Println("Closing Udp Scraper Connetction")
		return c.socket.Close()
	}
	return nil
}

func (c *udpAnnounce) Do(req *ScrapeRequest) (res ScrapeResponse, err error) {
	err = c.connect()
	if err != nil {
		return
	}
	// Clearly this limits the request URI to 255 bytes. BEP 41 supports
	// longer but I'm not fussed.
	//options := append([]byte{optionTypeURLData, byte(len(reqURI))}, []byte(reqURI)...)
	b, err := c.request(ActionScrape, req.InfoHashes)
	if err != nil {
		return
	}
	var scrapeResponse []int32 = make([]int32, 3*len(req.InfoHashes))

	err = readBody(b, &(scrapeResponse))
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		err = fmt.Errorf("error parsing scrape response: %s", err)
		return
	}
	if len(scrapeResponse)/3 != len(req.InfoHashes) {
		err = errors.New("Reqest and respons  infohash count not match")
		return
	}

	res.ScrapeDatas = make(map[string]ScrapeTorrentInfo)
	for i, infohash := range req.InfoHashes {
		res.ScrapeDatas[infohash.HexString()] = ScrapeTorrentInfo{
			Seeders:   scrapeResponse[i*3],
			Completed: scrapeResponse[i*3+1],
			Leechers:  scrapeResponse[i*3+2],
		}
	}
	return
}

// body is the binary serializable request body. trailer is optional data
// following it, such as for BEP 41.
func (c *udpAnnounce) write(h *RequestHeader, body interface{}) (err error) {
	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, h)
	if err != nil {
		panic(err)
	}
	if body != nil {
		err = binary.Write(&buf, binary.BigEndian, body)
		if err != nil {
			panic(err)
		}
	}
	n, err := c.socket.Write(buf.Bytes())
	if err != nil {
		return
	}
	if n != buf.Len() {
		panic("write should send all or error")
	}
	return
}

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}

func write(w io.Writer, data interface{}) error {
	return binary.Write(w, binary.BigEndian, data)
}

// args is the binary serializable request body. trailer is optional data
// following it, such as for BEP 41.
func (c *udpAnnounce) request(action Action, args interface{}) (responseBody *bytes.Buffer, err error) {
	tid := newTransactionId()
	err = c.write(&RequestHeader{
		ConnectionId:  c.connectionId,
		Action:        action,
		TransactionId: tid,
	}, args)
	if err != nil {
		return
	}
	c.socket.SetReadDeadline(time.Now().Add(timeout(c.contiguousTimeouts)))
	b := make([]byte, 65536) // 2KiB
	for {
		var n int
		n, err = c.socket.Read(b)
		if opE, ok := err.(*net.OpError); ok {
			if opE.Timeout() {
				c.contiguousTimeouts++
				return
			}
		}
		if err != nil {
			return
		}
		buf := bytes.NewBuffer(b[:n])
		var h ResponseHeader
		err = binary.Read(buf, binary.BigEndian, &h)
		switch err {
		case io.ErrUnexpectedEOF:
			continue
		case nil:
		default:
			return
		}
		if h.TransactionId != tid {
			continue
		}
		c.contiguousTimeouts = 0
		if h.Action == ActionError {
			err = errors.New(buf.String())
		}
		responseBody = buf
		return
	}
}

func readBody(r io.Reader, data ...interface{}) (err error) {
	for _, datum := range data {
		err = binary.Read(r, binary.BigEndian, datum)
		if err != nil {
			break
		}
	}
	return
}

func (c *udpAnnounce) connected() bool {
	return !c.connectionIdReceived.IsZero() && time.Now().Before(c.connectionIdReceived.Add(time.Minute))
}

func (c *udpAnnounce) connect() (err error) {
	if c.connected() {
		return nil
	}
	c.connectionId = connectRequestConnectionId
	if c.socket == nil {
		hmp := Utils.SplitHostMaybePort(c.url.Host)
		if hmp.NoPort {
			hmp.NoPort = false
			hmp.Port = 80
		}
		c.socket, err = net.Dial("udp", hmp.String())
		if err != nil {
			return
		}
	}
	b, err := c.request(ActionConnect, nil)
	if err != nil {
		return
	}
	var res ConnectionResponse
	err = readBody(b, &res)
	if err != nil {
		return
	}
	c.connectionId = res.ConnectionId
	c.connectionIdReceived = time.Now()
	return
}

func scrapeUDP(ar *ScrapeRequest, _url *url.URL) (ScrapeResponse, error) {
	ua := udpAnnounce{
		url: *_url,
	}
	defer ua.Close()
	log.Println("New Udp Scraper Connetction")
	return ua.Do(ar)
}
