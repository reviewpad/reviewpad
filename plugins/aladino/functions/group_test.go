// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var group = plugins_aladino.PluginBuiltIns().Functions["group"].Code

func TestGroup(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	groupName := "techLeads"
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	wantGroup := aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john"), aladino.BuildStringValue("arthur")})

	mockedEnv.GetRegisterMap()[groupName] = wantGroup

	args := []aladino.Value{aladino.BuildStringValue(groupName)}
	gotGroup, err := group(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantGroup, gotGroup)
}

func TestGroup_WhenGroupIsNonExisting(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	groupName := "techLeads"
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), aladino.DefaultMockEventPayload, controller)

	// Make sure that the group 'techLeads' doesn't exist
	delete(mockedEnv.GetRegisterMap(), groupName)

	args := []aladino.Value{aladino.BuildStringValue(groupName)}
	gotGroup, err := group(mockedEnv, args)

	assert.Nil(t, gotGroup)
	assert.EqualError(t, err, fmt.Sprintf("getGroup: no group with name %v in state %+q", groupName, mockedEnv.GetRegisterMap()))
}
