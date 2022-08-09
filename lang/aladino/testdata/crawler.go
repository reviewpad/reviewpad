// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type SafeCounter struct {
	v   map[string]bool
	mux sync.Mutex
}

var c SafeCounter = SafeCounter{v: make(map[string]bool)}

func (s *SafeCounter) checkvisited(url string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, ok := s.v[url]
	if ok == false {
		s.v[url] = true
		return false
	}
	return true

}

// First version as presented at:
// https://gist.github.com/harryhare/6a4979aa7f8b90db6cbc74400d0beb49#file-exercise-web-crawler-go
func Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}
	if c.checkvisited(url) {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		go Crawl(u, depth-1, fetcher)
	}
	return
}
