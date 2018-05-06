package reverseproxy

import "strings"

type HostReplaceFunc func(accessHost, targetDomain string) string

func NewHostReplaceFunc(accessDomain string) HostReplaceFunc {
	return func(accessHost, targetDomain string) string {
		if strings.HasSuffix(accessHost, accessDomain) && len(accessHost) > len(accessDomain) {
			return strings.TrimRight(accessHost, accessDomain) + targetDomain
		}
		return targetDomain
	}
}
