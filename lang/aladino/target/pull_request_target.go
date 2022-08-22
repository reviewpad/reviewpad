// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.
package target

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/host-event-handler/handler"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
)

type PullRequestTarget struct {
	*CommonTarget

	ctx          context.Context
	targetEntity *handler.TargetEntity
	pr           *github.PullRequest
	githubClient *gh.GithubClient
}

// ensure PullRequestTarget conforms to Target interface
var _ Target = PullRequestTarget{}

func NewPullRequestTarget(ctx context.Context, targetEntity *handler.TargetEntity, githubClient *gh.GithubClient, pr *github.PullRequest) *PullRequestTarget {
	return &PullRequestTarget{
		NewCommonTarget(ctx, targetEntity, githubClient),
		ctx,
		targetEntity,
		pr,
		githubClient,
	}
}

func (t PullRequestTarget) AddToProject(projectID, fieldID, optionID string) error {
	// ctx := t.env.GetCtx()
	// pr := t.env.GetPullRequest()

	// return t.env.GetGithubClient().AddToProject(ctx, projectID, *pr.NodeID, fieldID, optionID)
	return nil
}

func (t PullRequestTarget) Close() error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	pr := t.pr

	pr.State = github.String("closed")

	_, _, err := t.githubClient.EditPullRequest(ctx, owner, repo, number, pr)

	return err
}

func (t PullRequestTarget) GetAvailableAssignees() ([]User, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	assignees := make([]User, 0)

	users, err := t.githubClient.GetIssuesAvailableAssignees(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		assignees = append(assignees, User{
			Login: *user.Login,
		})
	}

	return assignees, nil
}

func (t PullRequestTarget) GetComments() ([]Comment, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	cs, err := t.githubClient.GetPullRequestComments(ctx, owner, repo, number, &github.IssueListCommentsOptions{})
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, len(cs))

	for _, comment := range cs {
		comments = append(comments, Comment{
			Body: *comment.Body,
		})
	}

	return comments, nil
}

func (t PullRequestTarget) GetLabels() ([]Label, error) {
	pr := t.pr
	labels := make([]Label, len(pr.Labels))

	for _, label := range pr.Labels {
		labels = append(labels, Label{
			ID:   *label.ID,
			Name: *label.Name,
		})
	}

	return labels, nil
}

func (t PullRequestTarget) GetProjectByName(name string) (*Project, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	project, err := t.githubClient.GetProjectV2ByName(ctx, owner, repo, name)
	if err != nil {
		return nil, err
	}

	return &Project{
		ID:     project.ID,
		Number: project.Number,
	}, nil
}

func (t PullRequestTarget) GetProjectFieldsByProjectNumber(projectNumber uint64) ([]ProjectField, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	totalRetries := 2

	ghFields, err := t.githubClient.GetProjectFieldsByProjectNumber(ctx, owner, repo, projectNumber, totalRetries)
	if err != nil {
		return nil, err
	}

	fields := make([]ProjectField, len(ghFields))
	for _, field := range ghFields {
		fields = append(fields, ProjectField{
			ID:      field.Details.ID,
			Name:    field.Details.Name,
			Options: field.Details.Options,
		})
	}

	return fields, nil
}

func (t PullRequestTarget) GetRequestedReviewers() ([]User, error) {
	pr := t.pr
	reviewers := make([]User, len(pr.RequestedReviewers))

	for _, reviewer := range pr.RequestedReviewers {
		reviewers = append(reviewers, User{
			Login: *reviewer.Login,
		})
	}

	return reviewers, nil
}

func (t PullRequestTarget) GetReviewers() (*Reviewers, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	ghPrRequestedReviewers, err := t.githubClient.GetPullRequestReviewers(ctx, owner, repo, number, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	users := make([]User, len(ghPrRequestedReviewers.Users))
	teams := make([]Team, len(ghPrRequestedReviewers.Teams))

	for _, ghPrRequestedReviewerUser := range ghPrRequestedReviewers.Users {
		users = append(users, User{
			Login: *ghPrRequestedReviewerUser.Login,
		})
	}

	for _, ghPrRequestedReviewerTeam := range ghPrRequestedReviewers.Teams {
		teams = append(teams, Team{
			ID:   *ghPrRequestedReviewerTeam.ID,
			Name: *ghPrRequestedReviewerTeam.Name,
		})
	}

	return &Reviewers{
		Users: users,
		Teams: teams,
	}, nil
}

func (t PullRequestTarget) GetReviews() ([]Review, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	ghPrReviews, err := t.githubClient.GetPullRequestReviews(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	reviews := make([]Review, len(ghPrReviews))

	for _, ghPrReview := range ghPrReviews {
		reviews = append(reviews, Review{
			ID:    *ghPrReview.ID,
			Body:  *ghPrReview.Body,
			State: *ghPrReview.State,
			User: &User{
				Login: *ghPrReview.User.Login,
			},
		})
	}

	return reviews, nil
}

func (t PullRequestTarget) GetUser() (*User, error) {
	pr := t.pr

	return &User{
		Login: *pr.User.Login,
	}, nil
}

func (t PullRequestTarget) Merge(mergeMethod string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.Merge(ctx, owner, repo, number, "Merged by Reviewpad", &github.PullRequestOptions{
		MergeMethod: mergeMethod,
	})

	return err
}

func (t PullRequestTarget) RemoveLabel(labelName string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, err := t.githubClient.RemoveLabelForIssue(ctx, owner, repo, number, labelName)

	return err
}

func (t PullRequestTarget) RequestReviewers(reviewers []string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.RequestReviewers(ctx, owner, repo, number, github.ReviewersRequest{
		Reviewers: reviewers,
	})

	return err
}

func (t PullRequestTarget) RequestTeamReviewers(reviewers []string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.RequestReviewers(ctx, owner, repo, number, github.ReviewersRequest{
		TeamReviewers: reviewers,
	})

	return err
}

func (t PullRequestTarget) GetAssignees() ([]User, error) {
	pr := t.pr
	assignees := make([]User, len(pr.Assignees))

	for _, assignee := range pr.Assignees {
		assignees = append(assignees, User{
			Login: *assignee.Login,
		})
	}

	return assignees, nil
}

func (t PullRequestTarget) GetBase() (string, error) {
	return t.pr.GetBase().GetRef(), nil
}

func (t PullRequestTarget) GetCommentCount() (int, error) {
	return t.pr.GetComments(), nil
}

func (t PullRequestTarget) GetCommitCount() (int, error) {
	return t.pr.GetCommits(), nil
}

func (t PullRequestTarget) GetCommits() ([]Commit, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	ghCommits, err := t.githubClient.GetPullRequestCommits(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	commits := make([]Commit, len(ghCommits))

	for _, ghCommit := range ghCommits {
		commits = append(commits, Commit{
			Message:      *ghCommit.Commit.Message,
			ParentsCount: len(ghCommit.Parents),
		})
	}

	return commits, nil
}

func (t PullRequestTarget) GetCreatedAt() (string, error) {
	return t.pr.GetCreatedAt().String(), nil
}

func (t PullRequestTarget) GetDescription() (string, error) {
	return t.pr.GetBody(), nil
}

func (t PullRequestTarget) GetLinkedIssuesCount() (int, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	return t.githubClient.GetPullRequestClosingIssuesCount(ctx, owner, repo, number)
}

func (t PullRequestTarget) GetReviewThreads() ([]ReviewThread, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	totalRetryCount := 2

	ghReviewThreads, err := t.githubClient.GetReviewThreads(ctx, owner, repo, number, totalRetryCount)
	if err != nil {
		return nil, err
	}

	reviewThreads := make([]ReviewThread, len(ghReviewThreads))
	for i, reviewThread := range ghReviewThreads {
		reviewThreads[i] = ReviewThread{
			IsResolved: bool(reviewThread.IsResolved),
			IsOutdated: bool(reviewThread.IsOutdated),
		}
	}

	return reviewThreads, nil
}

func (t PullRequestTarget) GetHead() (string, error) {
	return t.pr.GetHead().GetRef(), nil
}

func (t PullRequestTarget) IsDraft() (bool, error) {
	return t.pr.GetDraft(), nil
}
