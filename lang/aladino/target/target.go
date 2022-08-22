// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.
package target

type Target interface {
	AddAssignees(assignees []string) error
	AddLabels(labels []string) error
	AddToProject(projectID, fieldID, optionID string) error
	Close() error
	Comment(comment string) error
	GetAvailableAssignees() ([]User, error)
	GetComments() ([]Comment, error)
	GetLabels() ([]Label, error)
	GetRequestedReviewers() ([]User, error)
	GetReviewers() (*Reviewers, error)
	GetReviews() ([]Review, error)
	GetUser() (*User, error)
	Merge(mergeMethod string) error
	RemoveLabel(labelName string) error
	RequestReviewers(reviewers []string) error
	RequestTeamReviewers(reviewers []string) error
	GetProjectByName(name string) (*Project, error)
	GetProjectFieldsByProjectNumber(projectNumber uint64) ([]ProjectField, error)
	GetAssignees() ([]User, error)
	GetBase() (string, error)
	GetCommentCount() (int, error)
	GetCommitCount() (int, error)
	GetCommits() ([]Commit, error)
	GetCreatedAt() (string, error)
	GetDescription() (string, error)
	GetLinkedIssuesCount() (int, error)
	GetReviewThreads() ([]ReviewThread, error)
	GetHead() (string, error)
	IsDraft() (bool, error)
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
	ID    int64
	User  *User
	Body  string
	State string
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
