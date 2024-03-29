// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var commentCount = plugins_aladino.PluginBuiltIns().Functions["commentCount"].Code

func TestCommentCount(t *testing.T) {
	wantCommentCount := lang.BuildIntValue(6)

	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	args := []lang.Value{}
	gotCommentCount, err := commentCount(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantCommentCount, gotCommentCount)
}
