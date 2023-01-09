// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package target

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/codehost"
	gh "github.com/reviewpad/reviewpad/v3/codehost/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/shurcooL/githubv4"
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

func (t *PullRequestTarget) Close(comment string, _ string) error {
	ctx := t.ctx
	pr := t.PullRequest

	var closePullRequestMutation struct {
		ClosePullRequest struct {
			ClientMutationID string
		} `graphql:"closePullRequest(input: $input)"`
	}

	input := githubv4.ClosePullRequestInput{
		PullRequestID: githubv4.ID(pr.GetNodeID()),
	}

	if err := t.githubClient.GetClientGraphQL().Mutate(ctx, &closePullRequestMutation, input, nil); err != nil {
		return err
	}

	if comment != "" {
		if errComment := t.Comment(comment); errComment != nil {
			return errComment
		}
	}

	return nil
}

func (t *PullRequestTarget) GetAuthor() (*codehost.User, error) {
	pr := t.PullRequest

	return &codehost.User{
		Login: *pr.User.Login,
	}, nil
}

func (t *PullRequestTarget) GetLabels() []*codehost.Label {
	pr := t.PullRequest
	labels := make([]*codehost.Label, len(pr.Labels))

	for i, label := range pr.Labels {
		labels[i] = &codehost.Label{
			ID:   *label.ID,
			Name: *label.Name,
		}
	}

	return labels
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
			ID:          *ghPrReview.ID,
			Body:        *ghPrReview.Body,
			State:       *ghPrReview.State,
			SubmittedAt: ghPrReview.SubmittedAt,
			User: &codehost.User{
				Login: *ghPrReview.User.Login,
			},
		}
	}

	return reviews, nil
}

func (t *PullRequestTarget) IsLinkedToProject(title string) (bool, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	totalRetries := 2

	projects, err := t.githubClient.GetLinkedProjectsForPullRequest(ctx, owner, repo, number, totalRetries)
	if err != nil {
		return false, err
	}

	for _, project := range projects {
		if project.Project.Title == title {
			return true, nil
		}
	}

	return false, nil
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

func (t *PullRequestTarget) Review(reviewEvent, reviewBody string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	_, _, err := t.githubClient.Review(ctx, owner, repo, number, &github.PullRequestReviewRequest{
		Body:  &reviewBody,
		Event: &reviewEvent,
	})

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

func (t *PullRequestTarget) GetUpdatedAt() (string, error) {
	return t.PullRequest.GetUpdatedAt().String(), nil
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

func (t *PullRequestTarget) GetState() string {
	return t.PullRequest.GetState()
}

func (t *PullRequestTarget) GetTitle() string {
	return t.PullRequest.GetTitle()
}

func (t *PullRequestTarget) GetPullRequestLastPushDate() (time.Time, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	return t.githubClient.GetPullRequestLastPushDate(ctx, owner, repo, number)
}

func (t *PullRequestTarget) IsFileBinary(branch, file string) (bool, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	return t.githubClient.IsFileBinary(ctx, owner, repo, branch, file)
}

func (t *PullRequestTarget) GetLastCommit() (string, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	return t.githubClient.GetLastCommitSHA(ctx, owner, repo, number)
}

func (t *PullRequestTarget) GetLinkedProjects() ([]gh.GQLProjectV2Item, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number
	totalRetries := 2

	return t.githubClient.GetLinkedProjectsForPullRequest(ctx, owner, repo, number, totalRetries)
}

func (t *PullRequestTarget) GetApprovalsCount() (int, error) {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	return t.githubClient.GetApprovalsCount(ctx, owner, repo, number)
}

func (t *PullRequestTarget) TriggerWorkflowByFileName(workflowFileName string) error {
	ctx := t.ctx
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	head := t.PullRequest.GetHead().GetRef()

	_, err := t.githubClient.TriggerWorkflowByFileName(ctx, owner, repo, head, workflowFileName)

	return err
}

func (t *PullRequestTarget) JSON() (string, error) {
	j, err := json.Marshal(t.PullRequest)
	if err != nil {
		return "", err
	}

	return string(j), nil
}
