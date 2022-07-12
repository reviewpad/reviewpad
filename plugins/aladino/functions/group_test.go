// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/reviewpad/reviewpad/v2/lang/aladino"
	mocks_aladino "github.com/reviewpad/reviewpad/v2/mocks/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v2/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var group = plugins_aladino.PluginBuiltIns().Functions["group"].Code

func TestGroup(t *testing.T) {
  groupName := "techLeads"
  mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

  wantGroup := aladino.BuildArrayValue([]aladino.Value{aladino.BuildStringValue("john"), aladino.BuildStringValue("arthur")})

  mockedEnv.GetRegisterMap()[groupName] = wantGroup

  args := []aladino.Value{aladino.BuildStringValue(groupName)}
  gotGroup, err := group(mockedEnv, args)

  assert.Nil(t, err)
  assert.Equal(t, wantGroup, gotGroup)
}

func TestGroup_WhenGroupIsNonExisting(t *testing.T) {
  groupName := "techLeads"
  mockedEnv, err := mocks_aladino.MockDefaultEnv()
	if err != nil {
		log.Fatalf("mockDefaultEnv failed: %v", err)
	}

  // Make sure that the group 'techLeads' doesn't exist
  delete(mockedEnv.GetRegisterMap(), groupName)

  args := []aladino.Value{aladino.BuildStringValue(groupName)}
  gotGroup, err := group(mockedEnv, args)

  assert.Nil(t, gotGroup)
  assert.EqualError(t, err, fmt.Sprintf("getGroup: no group with name %v in state %+q", groupName, mockedEnv.GetRegisterMap()))
}
