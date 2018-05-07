package reverseproxy

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type DomainMapping struct {
	from string
	to   string
}

func (p *DomainMapping) Reverse() *DomainMapping {
	return &DomainMapping{
		from: p.to,
		to:   p.from,
	}
}

func (p *DomainMapping) ReplaceStr(content string) string {
	return strings.Replace(content, p.from, p.to, -1)
}

func (p *DomainMapping) ReplaceBytes(content []byte) []byte {
	return bytes.Replace(content, []byte(p.from), []byte(p.to), -1)
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

type BodyReplace struct {
	reqIn     *http.Request       // client request
	resIn     *http.Response      // server response
	writerOut http.ResponseWriter // proxy response
}

func (p *BodyReplace) HandleCompression() (readerIn io.Reader, writerOut io.Writer, err error) {
	readerIn = p.resIn.Body
	writerOut = p.writerOut

	reqAcceptEncoding := p.reqIn.Header.Get("Accept-Encoding")
	// We are ignoring any q-value here, so this is wrong for the case q=0
	clientAcceptsGzip := strings.Contains(reqAcceptEncoding, "gzip")
	log.Println("client request accept encoding:", reqAcceptEncoding, clientAcceptsGzip)

	p.writerOut.Header().Del("Content-Encoding")

	if p.resIn.Header.Get("Content-Encoding") == "gzip" {
		var e error
		readerIn, e = gzip.NewReader(readerIn)
		if e != nil {
			// Work around the closed-body-on-redirect bug in the runtime
			// https://github.com/golang/go/issues/10069
			readerIn = p.resIn.Body
			return
		}

		if clientAcceptsGzip {
			writerOut = gzip.NewWriter(writerOut)
			p.writerOut.Header().Set("Content-Encoding", "gzip")
		}
	}

	return
}
