// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package codehost

import (
	"errors"
	"time"

	pbc "github.com/reviewpad/api/go/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/handler"
)

var (
	ErrNotSupported = errors.New("not supported on this entity kind")
)

type Target interface {
	AddAssignees(assignees []string) error
	AddLabels(labels []string) error
	Close(comment string, stateReason string) error
	Comment(comment string) error
	GetAssignees() []*pbc.User
	GetAvailableAssignees() ([]*User, error)
	GetAuthor() (*User, error)
	GetCommentCount() int64
	GetComments() ([]*Comment, error)
	GetCreatedAt() string
	GetUpdatedAt() string
	GetDescription() string
	GetLabels() []*pbc.Label
	GetLinkedProjects() ([]gh.GQLProjectV2Item, error)
	GetNodeID() string
	GetProjectByName(name string) (*Project, error)
	GetProjectFieldsByProjectNumber(projectNumber uint64) ([]*ProjectField, error)
	GetState() pbc.PullRequestStatus
	GetTargetEntity() *handler.TargetEntity
	GetTitle() string
	IsLinkedToProject(title string) (bool, error)
	RemoveLabel(labelName string) error
	SetProjectField(projectItems []gh.GQLProjectV2Item, projectTitle, fieldName, fieldValue string) error
	JSON() (string, error)
	GetProjectV2ItemID(projectID string) (string, error)
}

type User struct {
	Login string
}

type Team struct {
	ID   int64
	Name string
}

type Comment struct {
	Body string
}

type Reviewers struct {
	Users []User
	Teams []Team
}

type Review struct {
	ID          int64
	User        *User
	Body        string
	State       string
	SubmittedAt *time.Time
}

type Project struct {
	ID     string
	Number uint64
	Name   string
}

type ProjectField struct {
	ID      string
	Name    string
	Options []struct {
		ID   string
		Name string
	}
}

type Commit struct {
	Message      string
	ParentsCount int
}

type ReviewThread struct {
	IsResolved bool
	IsOutdated bool
}
