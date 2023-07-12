// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"errors"
	"testing"

	"github.com/ohler55/ojg"
	"github.com/reviewpad/go-lib/event/event_processor"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var selectFromContext = plugins_aladino.PluginBuiltIns().Functions["selectFromContext"].Code

func TestSelectFromContext(t *testing.T) {
	mockIssueTargetEntity := &event_processor.TargetEntity{
		Kind:   event_processor.Issue,
		Owner:  "test",
		Repo:   "test",
		Number: 1,
	}
	tests := map[string]struct {
		targetEntity *event_processor.TargetEntity
		path         string
		wantErr      error
		wantRes      lang.Value
	}{
		"when error parsing expression": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.[0]",
			wantErr:      errors.New("an expression fragment can not start with a '[' at 4 in $.[0]"),
		},
		"when success getting title for pull request": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.title",
			wantRes:      lang.BuildStringValue("Amazing new feature"),
		},
		"when success getting title for issue": {
			targetEntity: mockIssueTargetEntity,
			path:         "$.title",
			wantRes:      lang.BuildStringValue("Found a bug"),
		},
		"when success getting user login for pull request": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.user.login",
			wantRes:      lang.BuildStringValue("john"),
		},
		"when success getting user login for issue": {
			targetEntity: mockIssueTargetEntity,
			path:         "$.user.login",
			wantRes:      lang.BuildStringValue("john"),
		},
		"when success getting pull request assignees": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.assignees",
			wantRes:      lang.BuildStringValue(`[{"login":"jane"}]`),
		},
		"when success getting issue assignees": {
			targetEntity: mockIssueTargetEntity,
			path:         "$.assignees",
			wantRes:      lang.BuildStringValue(`[{"login":"jane"}]`),
		},
		"when success getting pull request body": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.body",
			wantRes:      lang.BuildStringValue("Please pull these awesome changes in!"),
		},
		"when success getting pull request merged": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.merged",
			wantRes:      lang.BuildStringValue("true"),
		},
		"when success getting pull request comments": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.comments",
			wantRes:      lang.BuildStringValue("6"),
		},
		"when success getting pull request milestone title": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.milestone.title",
			wantRes:      lang.BuildStringValue("v1.0"),
		},
		"when success getting pull request first label name": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.labels[0].name",
			wantRes:      lang.BuildStringValue("enhancement"),
		},
		"when success getting pull request head repo url": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.head.repo.url",
			wantRes:      lang.BuildStringValue("https://api.github.com/repos/foobar/default-mock-repo/pulls/6"),
		},
		"when success getting issue labels": {
			targetEntity: mockIssueTargetEntity,
			path:         "$.labels",
			wantRes:      lang.BuildStringValue(`[{"id":1,"name":"bug"}]`),
		},
		"when success getting pull request reviewers": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.requested_reviewers",
			wantRes:      lang.BuildStringValue(`[{"login":"jane"}]`),
		},
		"when success getting all label names": {
			targetEntity: aladino.DefaultMockTargetEntity,
			path:         "$.labels[*].name",
			wantRes:      lang.BuildStringValue(`["enhancement","large"]`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := aladino.MockDefaultEnvWithTargetEntity(
				t,
				nil,
				nil,
				aladino.MockBuiltIns(),
				nil,
				test.targetEntity,
			)

			res, err := selectFromContext(env, []lang.Value{lang.BuildStringValue(test.path)})

			// since ojg errors contain stack traces
			// we are simplifying the error to make it easier to assert
			ojgError := &ojg.Error{}
			if errors.As(err, &ojgError) {
				err = errors.New(ojgError.Error())
			}

			assert.Equal(t, test.wantErr, err)
			assert.Equal(t, test.wantRes, res)
		})
	}
}
