// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package codehost

import (
	"errors"
	"time"

	"github.com/reviewpad/reviewpad/v3/handler"
)

var (
	ErrNotSupported = errors.New("not supported on this entity kind")
)

type Target interface {
	AddAssignees(assignees []string) error
	AddLabels(labels []string) error
	Close(comment string, stateReason string) error
	Comment(comment string) error
	GetAssignees() ([]*User, error)
	GetAvailableAssignees() ([]*User, error)
	GetAuthor() (*User, error)
	GetCommentCount() (int, error)
	GetComments() ([]*Comment, error)
	GetCreatedAt() (string, error)
	GetUpdatedAt() (string, error)
	GetDescription() (string, error)
	GetLabels() []*Label
	GetNodeID() string
	GetProjectByName(name string) (*Project, error)
	GetProjectFieldsByProjectNumber(projectNumber uint64) ([]*ProjectField, error)
	GetTargetEntity() *handler.TargetEntity
	GetTitle() string
	RemoveLabel(labelName string) error
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

type Label struct {
	ID   int64
	Name string
}

type TargetReview struct {
	ID    int64
	State string
	User  *User
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
