// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var commentCount = plugins_aladino.PluginBuiltIns().Functions["commentCount"].Code

func TestCommentCount(t *testing.T) {
	wantCommentCount := aladino.BuildIntValue(6)

	mockedEnv, err := aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{}
	gotCommentCount, err := commentCount(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantCommentCount, gotCommentCount)
}
