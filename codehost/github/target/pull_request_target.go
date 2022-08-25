// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package target

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/codehost"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/handler"
)

type Patch map[string]*codehost.File

type PullRequestTarget struct {
	*CommonTarget

	ctx          context.Context
	PullRequest  *github.PullRequest
	githubClient *gh.GithubClient
	Patch        Patch
}

// ensure PullRequestTarget conforms to Target interface
var _ codehost.Target = (*PullRequestTarget)(nil)

func getPullRequestPatch(ctx context.Context, pullRequest *github.PullRequest, githubClient *gh.GithubClient) (Patch, error) {
	owner := gh.GetPullRequestBaseOwnerName(pullRequest)
	repo := gh.GetPullRequestBaseRepoName(pullRequest)
	number := gh.GetPullRequestNumber(pullRequest)

	files, err := githubClient.GetPullRequestFiles(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	patchMap := make(map[string]*codehost.File)

	for _, file := range files {
		patchFile, err := codehost.NewFile(file)
		if err != nil {
			return nil, err
		}

		patchMap[file.GetFilename()] = patchFile
	}

	return Patch(patchMap), nil
}

func NewPullRequestTarget(ctx context.Context, targetEntity *handler.TargetEntity, githubClient *gh.GithubClient, pr *github.PullRequest) (*PullRequestTarget, error) {
	patch, err := getPullRequestPatch(ctx, pr, githubClient)
	if err != nil {
		return nil, err
	}

	return &PullRequestTarget{
		NewCommonTarget(ctx, targetEntity, githubClient),
		ctx,
		pr,
		githubClient,
		patch,
	}, nil
}

func (t *PullRequestTarget) GetNodeID() string {
	return t.PullRequest.GetNodeID()
}

func (t *PullRequestTarget) Close(comment string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	pr := t.PullRequest

	if comment != "" {
		if err := t.Comment(comment); err != nil {
			return err
		}
	}

	pr.State = github.String("closed")

	_, _, err := t.githubClient.EditPullRequest(ctx, owner, repo, number, pr)

	return err
}

func (t *PullRequestTarget) GetAuthor() (*codehost.User, error) {
	pr := t.PullRequest

	return &codehost.User{
		Login: *pr.User.Login,
	}, nil
}

func (t *PullRequestTarget) GetLabels() ([]*codehost.Label, error) {
	pr := t.PullRequest
	labels := make([]*codehost.Label, len(pr.Labels))

	for i, label := range pr.Labels {
		labels[i] = &codehost.Label{
			ID:   *label.ID,
			Name: *label.Name,
		}
	}

	return labels, nil
}

func (t *PullRequestTarget) GetProjectByName(name string) (*codehost.Project, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	project, err := t.githubClient.GetProjectV2ByName(ctx, owner, repo, name)
	if err != nil {
		return nil, err
	}

	return &codehost.Project{
		ID:     project.ID,
		Number: project.Number,
	}, nil
}

func (t *PullRequestTarget) GetRequestedReviewers() ([]*codehost.User, error) {
	pr := t.PullRequest
	reviewers := make([]*codehost.User, len(pr.RequestedReviewers))

	for i, reviewer := range pr.RequestedReviewers {
		reviewers[i] = &codehost.User{
			Login: *reviewer.Login,
		}
	}

	return reviewers, nil
}

func (t *PullRequestTarget) GetReviewers() (*codehost.Reviewers, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	ghPrRequestedReviewers, err := t.githubClient.GetPullRequestReviewers(ctx, owner, repo, number, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	users := make([]codehost.User, len(ghPrRequestedReviewers.Users))
	teams := make([]codehost.Team, len(ghPrRequestedReviewers.Teams))

	for i, ghPrRequestedReviewerUser := range ghPrRequestedReviewers.Users {
		users[i] = codehost.User{
			Login: *ghPrRequestedReviewerUser.Login,
		}
	}

	for i, ghPrRequestedReviewerTeam := range ghPrRequestedReviewers.Teams {
		teams[i] = codehost.Team{
			ID:   *ghPrRequestedReviewerTeam.ID,
			Name: *ghPrRequestedReviewerTeam.Name,
		}
	}

	return &codehost.Reviewers{
		Users: users,
		Teams: teams,
	}, nil
}

func (t *PullRequestTarget) GetReviews() ([]*codehost.Review, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	ghPrReviews, err := t.githubClient.GetPullRequestReviews(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	reviews := make([]*codehost.Review, len(ghPrReviews))

	for i, ghPrReview := range ghPrReviews {
		reviews[i] = &codehost.Review{
			ID:    *ghPrReview.ID,
			Body:  *ghPrReview.Body,
			State: *ghPrReview.State,
			User: &codehost.User{
				Login: *ghPrReview.User.Login,
			},
		}
	}

	return reviews, nil
}

func (t *PullRequestTarget) Merge(mergeMethod string) error {
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

func (t *PullRequestTarget) RemoveLabel(labelName string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, err := t.githubClient.RemoveLabelForIssue(ctx, owner, repo, number, labelName)

	return err
}

func (t *PullRequestTarget) RequestReviewers(reviewers []string) error {
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

func (t *PullRequestTarget) RequestTeamReviewers(reviewers []string) error {
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

func (t *PullRequestTarget) GetAssignees() ([]*codehost.User, error) {
	pr := t.PullRequest
	assignees := make([]*codehost.User, len(pr.Assignees))

	for i, assignee := range pr.Assignees {
		assignees[i] = &codehost.User{
			Login: *assignee.Login,
		}
	}

	return assignees, nil
}

func (t *PullRequestTarget) GetBase() (string, error) {
	return t.PullRequest.GetBase().GetRef(), nil
}

func (t *PullRequestTarget) GetCommentCount() (int, error) {
	return t.PullRequest.GetComments(), nil
}

func (t *PullRequestTarget) GetCommitCount() (int, error) {
	return t.PullRequest.GetCommits(), nil
}

func (t *PullRequestTarget) GetCommits() ([]*codehost.Commit, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	ghCommits, err := t.githubClient.GetPullRequestCommits(ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	commits := make([]*codehost.Commit, len(ghCommits))

	for i, ghCommit := range ghCommits {
		commits[i] = &codehost.Commit{
			Message:      *ghCommit.Commit.Message,
			ParentsCount: len(ghCommit.Parents),
		}
	}

	return commits, nil
}

func (t *PullRequestTarget) GetCreatedAt() (string, error) {
	return t.PullRequest.GetCreatedAt().String(), nil
}

func (t *PullRequestTarget) GetDescription() (string, error) {
	return t.PullRequest.GetBody(), nil
}

func (t *PullRequestTarget) GetLinkedIssuesCount() (int, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	return t.githubClient.GetPullRequestClosingIssuesCount(ctx, owner, repo, number)
}

func (t *PullRequestTarget) GetReviewThreads() ([]*codehost.ReviewThread, error) {
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

	reviewThreads := make([]*codehost.ReviewThread, len(ghReviewThreads))
	for i, reviewThread := range ghReviewThreads {
		reviewThreads[i] = &codehost.ReviewThread{
			IsResolved: bool(reviewThread.IsResolved),
			IsOutdated: bool(reviewThread.IsOutdated),
		}
	}

	return reviewThreads, nil
}

func (t *PullRequestTarget) GetHead() (string, error) {
	return t.PullRequest.GetHead().GetRef(), nil
}

func (t *PullRequestTarget) IsDraft() (bool, error) {
	return t.PullRequest.GetDraft(), nil
}

func (t *PullRequestTarget) GetTitle() string {
	return t.PullRequest.GetTitle()
}
