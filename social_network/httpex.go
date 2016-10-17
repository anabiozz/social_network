package httpex

import (
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sync"
	"sync/atomic"
)

type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	hosts bool //whether any patterns contain hostnames
}

func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	var n = 0
	for k, v := range mux.m {
		if !pathMatch(k, path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = v.h
			pattern = v.pattern
		}
	}
	return
}

func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		//should not happen
		return false
	}

	n := len(pattern)
	if pattern[n-1] != '/' {
		match, _ = regexp.MatchString(pattern, path)
		return match
	}
	fullMatch, _ := regexp.MatchString(pattern, string(path[0:n]))
	return len(path) >= n && fullMatch
}
