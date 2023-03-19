// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils_test

import (
	"errors"
	"testing"

	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidateUrl(t *testing.T) {
	tests := map[string]struct {
		url              string
		expectedBranch   *pbe.Branch
		expectedFilePath string
		expectedErr      error
	}{
		"valid blob": {
			url: "https://github.com/reviewpad/.github/blob/main/reviewpad-models/common.yml",
			expectedBranch: &pbe.Branch{
				Repo: &pbe.Repository{
					Owner: "reviewpad",
					Name:  ".github",
				},
				Name: "main",
			},
			expectedFilePath: "reviewpad-models/common.yml",
			expectedErr:      nil,
		},
		"invalid blob": {
			url:              "https://github.com/reviewpad/.github/blo/main/reviewpad-models/common.yml",
			expectedBranch:   nil,
			expectedFilePath: "",
			expectedErr:      errors.New("fatal: url must be a link to a GitHub blob, e.g. https://github.com/reviewpad/action/blob/main/main.go"),
		},
		"url without https": {
			url:              "github.com/reviewpad/.github/blob/main/reviewpad-models/common.yml",
			expectedBranch:   nil,
			expectedFilePath: "",
			expectedErr:      errors.New("fatal: url must be a link to a GitHub blob, e.g. https://github.com/reviewpad/action/blob/main/main.go"),
		},
		"invalid github url": {
			url:              "https://gitlab.com/reviewpad/.github/blo/main/reviewpad-models/common.yml",
			expectedBranch:   nil,
			expectedFilePath: "",
			expectedErr:      errors.New("fatal: url must be a link to a GitHub blob, e.g. https://github.com/reviewpad/action/blob/main/main.go"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotBranch, gotFilepath, gotError := utils.ValidateUrl(test.url)
			assert.Equal(t, test.expectedBranch, gotBranch)
			assert.Equal(t, test.expectedFilePath, gotFilepath)
			assert.Equal(t, test.expectedErr, gotError)
		})
	}
}
