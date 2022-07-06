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

var startsWith = plugins_aladino.PluginBuiltIns().Functions["startsWith"].Code

func TestStartsWithTrue(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantVal := aladino.BuildBoolValue(true)

	args := []aladino.Value{
		aladino.BuildStringValue("Title of a pull request"),
		aladino.BuildStringValue("Title of"),
	}

	gotVal, err := startsWith(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestStartsWithFalse(t *testing.T) {
	mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

	wantVal := aladino.BuildBoolValue(false)

	args := []aladino.Value{
		aladino.BuildStringValue("Pull request title"),
		aladino.BuildStringValue("title"),
	}

	gotVal, err := startsWith(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
