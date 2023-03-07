// Copyright (C) 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package handler

import "encoding/json"

// ActionEvent contains information about the workflow run and the event that triggered the run.
// For more information, visit: https://docs.github.com/en/actions/learn-github-actions/contexts#github-context
type ActionEvent struct {
	ActionName       *string          `json:"action,omitempty"`
	ActionPath       *string          `json:"action_path,omitempty"`
	ActionRef        *string          `json:"action_ref,omitempty"`
	ActionRepository *string          `json:"action_repository,omitempty"`
	ActionStatus     *string          `json:"action_status,omitempty"`
	Actor            *string          `json:"actor,omitempty"`
	ApiUrl           *string          `json:"api_url,omitempty"`
	BaseRef          *string          `json:"base_ref,omitempty"`
	HeadRef          *string          `json:"head_ref,omitempty"`
	Env              *string          `json:"env,omitempty"`
	EventPayload     *json.RawMessage `json:"event,omitempty"`
	EventName        *string          `json:"event_name,omitempty"`
	EventPath        *string          `json:"event_path,omitempty"`
	QraphqlUrl       *string          `json:"graphql_url,omitempty"`
	JobID            *string          `json:"job,omitempty"`
	Ref              *string          `json:"ref,omitempty"`
	RefName          *string          `json:"ref_name,omitempty"`
	RefProtected     *bool            `json:"ref_protected,omitempty"`
	RefType          *string          `json:"ref_type,omitempty"`
	Path             *string          `json:"path,omitempty"`
	Repository       *string          `json:"repository,omitempty"`
	RepositoryOwner  *string          `json:"repository_owner,omitempty"`
	RepositoryUrl    *string          `json:"repositoryUrl,omitempty"`
	RetentionDays    *string          `json:"retention_days,omitempty"`
	RunID            *string          `json:"run_id,omitempty"`
	RunNumber        *string          `json:"run_number,omitempty"`
	RunAttempt       *string          `json:"run_attempt,omitempty"`
	ServerUrl        *string          `json:"server_url,omitempty"`
	SHA              *string          `json:"sha,omitempty"`
	Token            *string          `json:"token,omitempty"`
	Workflow         *string          `json:"workflow,omitempty"`
	WorkflowRef      *string          `json:"workflow_ref,omitempty"`
	Workspace        *string          `json:"workspace,omitempty"`
}
