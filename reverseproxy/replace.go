package reverseproxy

import (
	"bytes"
	"compress/gzip"
	"io"
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

func NewMapGroup(maps []DomainMapping) MapGroup {
	rv := MapGroup{maps}
	rv.init()
	return rv
}

func (p *MapGroup) init() {
	for _, mm := range p.maps {
		url, err := url.Parse(mm.To)
		if err != nil {
			panic(err)
		}
		mm.Target = url
	}
}

func (p *MapGroup) GetMapping(req *http.Request) *DomainMapping {
	for _, mm := range p.maps {
		if req.Host == mm.From {
			return &mm
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
