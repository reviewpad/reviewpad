// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var appendString = plugins_aladino.PluginBuiltIns().Functions["append"].Code

func TestAppendString(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	slice := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("a"),
			aladino.BuildStringValue("b"),
		},
	)
	elems := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("c"),
		},
	)
	args := []aladino.Value{slice, elems}

	wantSlice := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("a"),
			aladino.BuildStringValue("b"),
			aladino.BuildStringValue("c"),
		},
	)

	gotSlice, err := appendString(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantSlice, gotSlice)
}
