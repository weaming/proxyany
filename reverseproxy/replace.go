package reverseproxy

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type DomainMapping struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Target *url.URL `json:"-"`
}

func (p *DomainMapping) Reverse() *DomainMapping {
	return &DomainMapping{
		From: p.To,
		To:   p.From,
	}
}

func (p *DomainMapping) ReplaceStr(content string) string {
	return strings.Replace(content, p.From, p.To, -1)
}

func (p *DomainMapping) ReplaceBytes(content []byte) []byte {
	return bytes.Replace(content, []byte(p.From), []byte(p.To), -1)
}

func (p *DomainMapping) ReplaceHeader(head *http.Header) {
	for k, vv := range *head {
		for i, v := range vv {
			v = p.ReplaceStr(v)
			if i == 0 {
				head.Set(k, v)
			} else {
				head.Add(k, v)
			}
		}
	}
}

type MapGroup struct {
	maps []DomainMapping
}

func NewMapGroup(maps []DomainMapping) *MapGroup {
	rv := &MapGroup{maps}
	rv.init()
	return rv
}

func LoadMapGroupFromJson(fp string) *MapGroup {
	raw, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	mpArr := []DomainMapping{}
	err = json.Unmarshal(raw, &mpArr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rv := &MapGroup{mpArr}
	rv.init()
	return rv
}

func (p *MapGroup) init() {
	for i, mapping := range p.maps {
		url, err := url.Parse(mapping.To)
		if err != nil {
			panic(err)
		}
		p.maps[i].Target = url
		p.maps[i].To = url.Host
	}
}

func (p *MapGroup) GetMapping(host string) *DomainMapping {
	for _, mapping := range p.maps {
		if strings.HasSuffix(host, mapping.From) {
			return &mapping
		}
	}
	return nil
}

type BodyDecompressor struct {
	requestIn  *http.Request        // client request
	responseIn *http.Response       // server response
	writerOut  *http.ResponseWriter // proxy response
}

func HandleCompression(in *http.Response, out http.ResponseWriter, clientAcceptGzip bool) (reader io.Reader, writer io.Writer, err error) {
	reader = in.Body
	writer = out

	if in.Header.Get("Content-Encoding") == "gzip" {
		rd, e := gzip.NewReader(in.Body)
		if e != nil {
			log.Println(e)
		} else {
			reader = rd
		}
	}

	if clientAcceptGzip {
		// Have bug: always write header "Content-Length: 10"
		//writer = gzip.NewWriter(writer)
	}

	return
}
