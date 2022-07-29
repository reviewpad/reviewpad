// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var length = plugins_aladino.PluginBuiltIns().Functions["length"].Code

func TestLength(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil)

	slice := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("a"),
			aladino.BuildStringValue("b"),
		},
	)
	args := []aladino.Value{slice}

	wantLength := &aladino.IntValue{Val: 2}

	gotLength, err := length(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantLength, gotLength)
}
