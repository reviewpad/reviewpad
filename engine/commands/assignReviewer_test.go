// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v4/engine/commands"
	"github.com/stretchr/testify/assert"
)

func TestAssignReviewer(t *testing.T) {
	testCases := map[string]struct {
		args       []string
		wantAction string
		wantErr    error
	}{
		"when invalid number of arguments": {
			args: []string{
				"john",
				"jane",
			},
			wantErr: errors.New("accepts 1 arg(s), received 2"),
		},
		"when missing required arguments": {
			args:    []string{},
			wantErr: errors.New("accepts 1 arg(s), received 0"),
		},
		"when total required reviewers is not a number": {
			args: []string{
				"a",
			},
			wantErr: errors.New("invalid argument: a, number of reviewers must be a number"),
		},
		"when total required reviewers 0": {
			args: []string{
				"0",
			},
			wantErr: errors.New("invalid argument: 0, number of reviewers must be greater than 0"),
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cmd := commands.AssignReviewerCmd()

			cmd.SetOut(out)
			cmd.SetArgs(test.args)

			err := cmd.Execute()

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantAction, out.String())
		})
	}
}
