// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/stretchr/testify/assert"
)

func TestCommentHasAnnotation(t *testing.T) {
	tests := map[string]struct {
		comment    string
		annotation string
		wantVal    bool
	}{
		"single annotation": {
			comment:    "reviewpad-an: foo",
			annotation: "foo",
			wantVal:    true,
		},
		"single annotation with comment": {
			comment:    "reviewpad-an: foo\n hello world",
			annotation: "foo",
			wantVal:    true,
		},
		"multiple annotation empty": {
			comment:    "reviewpad-an: foo     bar",
			annotation: "",
			wantVal:    false,
		},
		"multiple annotation single spaced first": {
			comment:    "reviewpad-an: foo bar",
			annotation: "foo",
			wantVal:    true,
		},
		"multiple annotation single spaced": {
			comment:    "reviewpad-an: foo bar",
			annotation: "bar",
			wantVal:    true,
		},
		"multiple annotation multi spaced": {
			comment:    "reviewpad-an: foo   bar",
			annotation: "bar",
			wantVal:    true,
		},
		"starting with empty line": {
			comment:    "\n  reviewpad-an: foo",
			annotation: "foo",
			wantVal:    true,
		},
		"starting with spaces": {
			comment:    "  reviewpad-an: foo",
			annotation: "foo",
			wantVal:    true,
		},
		"annotation not found": {
			comment:    "reviewpad-an: foo",
			annotation: "bar",
			wantVal:    false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotVal := commentHasAnnotation(test.comment, test.annotation)

			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}

func TestGetBlocks(t *testing.T) {
	fileName := "crawlerWithRemovedBlock.go"
	patchData := `
@@ -2,9 +2,11 @@ package main
 
 import "fmt"
 
-// First version as presented at:
-// https://gist.github.com/harryhare/6a4979aa7f8b90db6cbc74400d0beb49#file-exercise-web-crawler-go
 func Crawl(url string, depth int, fetcher Fetcher) {
+	defer  c.wg.Done()
+
 	if depth <= 0 {
 		return
 	}
@@ -18,6 +20,7 @@ func Crawl(url string, depth int, fetcher Fetcher) {
 	}
 	fmt.Printf("found: %s %q\n", url, body)
 	for _, u := range urls {
+		c.wg.Add(1)
 		go Crawl(u, depth-1, fetcher)
 	}
 	return
`
	ghFile := &github.CommitFile{
		Patch:    &patchData,
		Filename: &fileName,
	}

	patchFile, err := aladino.NewFile(ghFile)
	assert.Nil(t, err)

	expectedRes := []*entities.ResolveBlockDiff{
		{
			Start: 2,
			End:   4,
		},
		{
			Start: 5,
			End:   6,
		},
		{
			Start: 5,
			End:   5,
		},
		{
			Start: 6,
			End:   7,
		},
		{
			Start: 8,
			End:   12,
		},
	}

	actualRes, err := getBlocks(patchFile)

	assert.Nil(t, err)
	assert.Equal(t, &expectedRes, &actualRes)
}
