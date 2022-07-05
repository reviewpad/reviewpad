// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils_test

import (
	"errors"
	"testing"

	"github.com/google/go-github/v42/github"
	"github.com/reviewpad/reviewpad/v2/utils"
	"github.com/stretchr/testify/assert"
)

type test struct {}

func TestPaginatedRequest_WhenFirstExecutionFails(t *testing.T) {
	failMessage := "PaginatedRequestFail"
	initFn := func() interface{} {
		return test{}
	}
	reqFn := func(i interface{}, page int) (interface{}, *github.Response, error) {
		return nil, nil, errors.New(failMessage)
	}

	res, err := utils.PaginatedRequest(initFn, reqFn)

	assert.Nil(t, res)
	assert.EqualError(t, err, failMessage)
}
