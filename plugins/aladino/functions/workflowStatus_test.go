// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v52/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/reviewpad/v4/lang"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v4/plugins/aladino"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

var workflowStatus = plugins_aladino.PluginBuiltIns().Functions["workflowStatus"].Code

func TestWorkflowStatus_WhenEventPayloadIsNotWorkflowRunEvent(t *testing.T) {
	checkName := "test-workflow"
	wantValue := lang.BuildStringValue("")

	eventPayload := &github.CheckRunEvent{}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{lang.BuildStringValue(checkName)}
	gotValue, err := workflowStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestWorkflowStatus_WhenWorkflowRunIsNil(t *testing.T) {
	checkName := "test-workflow"
	wantValue := lang.BuildStringValue("")

	eventPayload := &github.WorkflowRunEvent{
		WorkflowRun: nil,
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		nil,
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{lang.BuildStringValue(checkName)}
	gotValue, err := workflowStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestWorkflowStatus_WhenListCheckRunsForRefRequestFails(t *testing.T) {
	checkName := "test-workflow"
	failMessage := "ListCheckRunsForRefRequestFail"
	headSHA := "1234abc"

	eventPayload := &github.WorkflowRunEvent{
		WorkflowRun: &github.WorkflowRun{
			HeadSHA: &headSHA,
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{lang.BuildStringValue(checkName)}
	gotValue, err := workflowStatus(mockedEnv, args)

	assert.Nil(t, gotValue)
	assert.Equal(t, err.(*github.ErrorResponse).Message, failMessage)
}

func TestWorkflowStatus_WhenCheckRunNotFoundDueToEmptyCheckRuns(t *testing.T) {
	checkName := "test-workflow"
	headSHA := "1234abc"

	wantValue := lang.BuildStringValue("")

	eventPayload := &github.WorkflowRunEvent{
		WorkflowRun: &github.WorkflowRun{
			HeadSHA: &headSHA,
		},
	}
	emptyCheckRuns := &github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(emptyCheckRuns))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{lang.BuildStringValue(checkName)}
	gotValue, err := workflowStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestWorkflowStatus_WhenCheckRunIsMissingInNonEmptyCheckRuns(t *testing.T) {
	checkName := "test-workflow"
	headSHA := "1234abc"

	wantValue := lang.BuildStringValue("")

	eventPayload := &github.WorkflowRunEvent{
		WorkflowRun: &github.WorkflowRun{
			HeadSHA: &headSHA,
		},
	}
	dummyCheckName := "test-check"
	emptyCheckRuns := &github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{
			{
				Name: &dummyCheckName,
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(emptyCheckRuns))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{lang.BuildStringValue(checkName)}
	gotValue, err := workflowStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestWorkflowStatus_WhenEventIsCompleted(t *testing.T) {
	checkName := "test-workflow"
	checkStatus := "completed"
	checkConclusion := "success"
	headSHA := "1234abc"

	wantValue := lang.BuildStringValue(checkConclusion)

	eventPayload := &github.WorkflowRunEvent{
		WorkflowRun: &github.WorkflowRun{
			HeadSHA: &headSHA,
		},
	}
	dummyCheckName := "test-check"
	checkRuns := &github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{
			{
				Name: &dummyCheckName,
			},
			{
				Name:       &checkName,
				Status:     &checkStatus,
				Conclusion: &checkConclusion,
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(checkRuns))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{lang.BuildStringValue(checkName)}
	gotValue, err := workflowStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}

func TestWorkflowStatus_WhenEventIsNotCompleted(t *testing.T) {
	checkName := "test-workflow"
	checkStatus := "in_progress"
	headSHA := "1234abc"

	wantValue := lang.BuildStringValue(checkStatus)

	eventPayload := &github.WorkflowRunEvent{
		WorkflowRun: &github.WorkflowRun{
			HeadSHA: &headSHA,
		},
	}
	dummyCheckName := "test-check"
	checkRuns := &github.ListCheckRunsResults{
		CheckRuns: []*github.CheckRun{
			{
				Name: &dummyCheckName,
			},
			{
				Name:   &checkName,
				Status: &checkStatus,
			},
		},
	}
	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposCommitsCheckRunsByOwnerByRepoByRef,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(checkRuns))
				}),
			),
		},
		nil,
		aladino.MockBuiltIns(),
		eventPayload,
	)

	args := []lang.Value{lang.BuildStringValue(checkName)}
	gotValue, err := workflowStatus(mockedEnv, args)

	assert.Nil(t, err)
	assert.Equal(t, wantValue, gotValue)
}
