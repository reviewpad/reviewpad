// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v3/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var isElementOf = plugins_aladino.PluginBuiltIns().Functions["isElementOf"].Code

func TestIsElementOf_WhenTrue(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildStringValue("elemA"),
		aladino.BuildArrayValue([]aladino.Value{
			aladino.BuildStringValue("elemA"),
			aladino.BuildStringValue("elemB"),
		}),
	}
	gotVal, err := isElementOf(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(true)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestIsElementOf_WhenFalse(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv(nil, nil)
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	args := []aladino.Value{
		aladino.BuildStringValue("elemA"),
		aladino.BuildArrayValue([]aladino.Value{
			aladino.BuildStringValue("elemB"),
			aladino.BuildStringValue("elemC"),
		}),
	}
	gotVal, err := isElementOf(mockedEnv, args)

	wantVal := aladino.BuildBoolValue(false)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
