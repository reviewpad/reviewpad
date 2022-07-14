// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package aladino

import (
	"fmt"
	"regexp"

	"github.com/google/go-github/v42/github"
)

type File struct {
	Repr *github.CommitFile
	Diff []*DiffBlock
}

func NewFile(file *github.CommitFile) (*File, error) {
	diffBlocks, err := parseFilePatch(file.GetPatch())
	if err != nil {
		return nil, fmt.Errorf("error in file patch %s: %v", file.GetFilename(), err)
	}

	return &File{
		Repr: file,
		Diff: diffBlocks,
	}, nil
}

func (f *File) Query(expr string) (bool, error) {
	r, err := regexp.Compile(expr)
	if err != nil {
		return false, fmt.Errorf("query: compile error %v", err)
	}

	for _, block := range f.Diff {
		if !block.isContext {
			if r.Match([]byte(block.newLine)) {
				return true, nil
			}
		}
	}
	return false, nil
}
