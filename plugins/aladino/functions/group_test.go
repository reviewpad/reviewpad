// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"testing"

	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var group = plugins_aladino.PluginBuiltIns().Functions["group"].Code

func TestGroup(t *testing.T) {
	groupName := "techLeads"
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	wantGroup := lang.BuildArrayValue([]lang.Value{lang.BuildStringValue("john"), lang.BuildStringValue("arthur")})

	mockedEnv.GetRegisterMap()[groupName] = wantGroup

	args := []lang.Value{lang.BuildStringValue(groupName)}
	gotGroup, err := group(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantGroup, gotGroup)
}

func TestGroup_WhenGroupIsNonExisting(t *testing.T) {
	groupName := "techLeads"
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	// Make sure that the group 'techLeads' doesn't exist
	delete(mockedEnv.GetRegisterMap(), groupName)

	args := []lang.Value{lang.BuildStringValue(groupName)}
	gotGroup, err := group(mockedEnv, args)

	assert.Nil(t, gotGroup)
	assert.EqualError(t, err, fmt.Sprintf("getGroup: no group with name %v in state %+q", groupName, mockedEnv.GetRegisterMap()))
}
