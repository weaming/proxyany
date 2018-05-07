package reverseproxy

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type DomainMapping struct {
	From   string
	To     string
	Target *url.URL
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

func (p *MapGroup) init() {
	for i, mm := range p.maps {
		url, err := url.Parse(mm.To)
		if err != nil {
			panic(err)
		}
		p.maps[i].Target = url
		p.maps[i].To = url.Host
	}
}

func (p *MapGroup) GetMapping(host string) *DomainMapping {
	for _, mm := range p.maps {
		if strings.HasSuffix(host, mm.From) {
			return &mm
		}
	}
	log.Printf("can't find mapping for %v\n", host)
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
