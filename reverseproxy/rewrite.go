package reverseproxy

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	requestIn  *http.Request       // client request
	responseIn *http.Response      // server response
	writerOut  http.ResponseWriter // proxy response
}

func (p *BodyDecompressor) HandleCompression() (readerIn io.Reader, writerOut io.Writer, err error) {
	readerIn = p.responseIn.Body
	writerOut = p.writerOut

	reqAcceptEncoding := p.requestIn.Header.Get("Accept-Encoding")
	// We are ignoring any q-value here, so this is wrong for the case q=0
	clientAcceptsGzip := strings.Contains(reqAcceptEncoding, "gzip")

	p.writerOut.Header().Del("Content-Encoding")

	if p.responseIn.Header.Get("Content-Encoding") == "gzip" {
		var e error
		readerIn, e = gzip.NewReader(readerIn)
		if e != nil {
			// Work around the closed-body-on-redirect bug in the runtime
			// https://github.com/golang/go/issues/10069
			readerIn = p.responseIn.Body
			return
		}

		if clientAcceptsGzip {
			writerOut = gzip.NewWriter(writerOut)
			p.writerOut.Header().Set("Content-Encoding", "gzip")
		}
	}

	return
}
