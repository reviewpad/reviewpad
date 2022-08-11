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

var sprintf = plugins_aladino.PluginBuiltIns(nil).Functions["sprintf"].Code

func TestSprintf(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	format := aladino.BuildStringValue("Hello %v!")
	slice := aladino.BuildArrayValue(
		[]aladino.Value{
			aladino.BuildStringValue("world"),
		},
	)
	args := []aladino.Value{format, slice}

	wantString := &aladino.StringValue{Val: "Hello world!"}

	gotString, err := sprintf(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantString, gotString)
}
