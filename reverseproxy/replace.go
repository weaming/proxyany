package reverseproxy

import (
	"bytes"
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
	for k := range map[string][]string(*head) {
		head.Set(k, p.ReplaceStr(head.Get(k)))
	}
}
