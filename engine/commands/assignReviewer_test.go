// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package commands_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/reviewpad/reviewpad/v3/engine/commands"
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
		"when list of reviewers has invalid github username": {
			args: []string{
				"john2,jane-",
			},
			wantErr: errors.New("reviewers must be a list of comma separated valid github usernames"),
		},
		"when review policy is invalid": {
			args: []string{
				"jane,john",
				"--total-reviewers=1",
				"--review-policy=unknown",
			},
			wantErr: errors.New("invalid review policy specified: unknown"),
		},
		"when number of reviewers is not a number": {
			args: []string{
				"john,jane,john2,jane27",
				"--total-reviewers=z",
			},
			wantErr: errors.New("invalid argument \"z\" for \"-t, --total-reviewers\" flag: strconv.ParseUint: parsing \"z\": invalid syntax"),
		},
		"when missing number of reviewers and policy": {
			args: []string{
				"john",
			},
			wantAction: `$assignReviewer(["john"], 1, "reviewpad")`,
		},
		"when missing policy": {
			args: []string{
				"john-123,jane",
				"--total-reviewers=5",
			},
			wantAction: `$assignReviewer(["john-123","jane"], 5, "reviewpad")`,
		},
		"when only one reviewer is provided": {
			args: []string{
				"john-123-jane",
				"--total-reviewers=1",
				"--review-policy=reviewpad",
			},
			wantAction: `$assignReviewer(["john-123-jane"], 1, "reviewpad")`,
		},
		"when only two reviewers are provided": {
			args: []string{
				"jane,john",
				"--total-reviewers",
				"2",
				"--review-policy",
				"random",
			},
			wantAction: `$assignReviewer(["jane","john"], 2, "random")`,
		},
		"when number of provided reviewers is greater than requested reviewers": {
			args: []string{
				"jane,john",
				"--total-reviewers=1",
				"--review-policy",
				"round-robin",
			},
			wantAction: `$assignReviewer(["jane","john"], 1, "round-robin")`,
		},
		"when number of provided reviewers is less than requested reviewers": {
			args: []string{
				"jane,john,jane123",
				"--total-reviewers",
				"5",
				"--review-policy=round-robin",
			},
			wantAction: `$assignReviewer(["jane","john","jane123"], 5, "round-robin")`,
		},
		"when missing number of reviewers": {
			args: []string{
				"jane,john,jane123",
				"--review-policy=reviewpad",
			},
			wantAction: `$assignReviewer(["jane","john","jane123"], 3, "reviewpad")`,
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
