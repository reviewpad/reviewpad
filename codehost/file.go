// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package codehost

import (
	"fmt"
	"regexp"

	pbc "github.com/reviewpad/api/go/codehost"
)

type File struct {
	Repr *pbc.File
	Diff []*diffBlock
}

func (f *File) AppendToDiff(
	isContext bool,
	oldStart, oldEnd, newStart, newEnd int,
	oldLine, newLine string,
) {
	f.Diff = append(f.Diff, &diffBlock{
		IsContext: isContext,
		Old: &diffSpan{
			int32(oldStart),
			int32(oldEnd),
		},
		New: &diffSpan{
			int32(newStart),
			int32(newEnd),
		},
		oldLine: oldLine,
		newLine: newLine,
	})
}

func NewFile(file *pbc.File) (*File, error) {
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
		if !block.IsContext {
			if r.Match([]byte(block.newLine)) {
				return true, nil
			}
		}
	}
	return false, nil
}
