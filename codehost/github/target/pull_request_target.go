// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package target

import (
	"context"
	"time"

	"github.com/google/go-github/v49/github"
	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/reviewpad/v4/codehost"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/shurcooL/githubv4"
)

type Patch map[string]*codehost.File

type PullRequestTarget struct {
	*CommonTarget

	ctx          context.Context
	PullRequest  *pbe.ExternalCodeReview
	githubClient *gh.GithubClient
	Patch        Patch
}

// ensure PullRequestTarget conforms to Target interface
var _ codehost.Target = (*PullRequestTarget)(nil)

func getPullRequestPatch(ctx context.Context, pullRequest *pbe.ExternalCodeReview, githubClient *gh.GithubClient) (Patch, error) {
	owner := pullRequest.GetBase().GetRepo().GetOwner()
	repo := pullRequest.GetBase().GetRepo().GetName()
	number := pullRequest.GetNumber()

	files, err := githubClient.GetPullRequestFiles(ctx, owner, repo, int(number))
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

func NewPullRequestTarget(ctx context.Context, targetEntity *handler.TargetEntity, githubClient *gh.GithubClient, pr *pbe.ExternalCodeReview) (*PullRequestTarget, error) {
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
	return t.PullRequest.GetNodeId()
}

func (t *PullRequestTarget) Close(comment string, _ string) error {
	ctx := t.ctx
	pr := t.PullRequest

	if pr.GetStatus() == pbe.CodeReviewStatus_CLOSED {
		return nil
	}

	var closePullRequestMutation struct {
		ClosePullRequest struct {
			ClientMutationID string
		} `graphql:"closePullRequest(input: $input)"`
	}

	input := githubv4.ClosePullRequestInput{
		PullRequestID: githubv4.ID(pr.GetNodeId()),
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
		Login: pr.GetAuthor().GetLogin(),
	}, nil
}

func (t *PullRequestTarget) GetLabels() []*pbe.Label {
	return t.PullRequest.GetLabels()
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

func (t *PullRequestTarget) GetRequestedReviewers() []*pbe.ExternalUser {
	return t.PullRequest.RequestedReviewers.Users
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

func (t *PullRequestTarget) GetAssignees() []*pbe.ExternalUser {
	return t.PullRequest.GetAssignees()
}

func (t *PullRequestTarget) GetBase() string {
	return t.PullRequest.GetBaseBranch()
}

func (t *PullRequestTarget) GetCommentCount() int64 {
	return t.PullRequest.GetCommentsCount()
}

func (t *PullRequestTarget) GetCommitCount() int64 {
	return t.PullRequest.GetCommitsCount()
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

func (t *PullRequestTarget) GetCreatedAt() string {
	return time.Unix(0, t.PullRequest.GetCreatedAt()).String()
}

func (t *PullRequestTarget) GetUpdatedAt() string {
	return time.Unix(0, t.PullRequest.GetUpdatedAt()).String()
}

func (t *PullRequestTarget) GetDescription() string {
	return t.PullRequest.GetDescription()
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

func (t *PullRequestTarget) GetHead() string {
	return t.PullRequest.GetHeadBranch()
}

func (t *PullRequestTarget) IsDraft() bool {
	return t.PullRequest.GetIsDraft()
}

func (t *PullRequestTarget) GetState() pbe.CodeReviewStatus {
	return t.PullRequest.GetStatus()
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
	head := t.PullRequest.GetHeadBranch()

	_, err := t.githubClient.TriggerWorkflowByFileName(ctx, owner, repo, head, workflowFileName)

	return err
}

func (t *PullRequestTarget) JSON() (string, error) {
	return t.PullRequest.Raw, nil
}

func (t *PullRequestTarget) GetLatestReviewFromReviewer(author string) (*codehost.Review, error) {
	clientGQL := t.githubClient.GetClientGraphQL()
	targetEntity := t.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	var reviewsQuery struct {
		Repository struct {
			PullRequest struct {
				Reviews struct {
					Nodes []struct {
						Author struct {
							Login githubv4.String
						}
						Body        githubv4.String
						State       githubv4.String
						SubmittedAt *time.Time
					}
				} `graphql:"reviews(last: 1, author: $author)"`
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}
	varGQLReviews := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(number),
		"author":            githubv4.String(author),
	}

	err := clientGQL.Query(context.Background(), &reviewsQuery, varGQLReviews)
	if err != nil {
		return nil, err
	}

	reviews := reviewsQuery.Repository.PullRequest.Reviews.Nodes
	if len(reviews) == 0 {
		return nil, nil
	}

	latestReview := reviews[0]
	return &codehost.Review{
		User: &codehost.User{
			Login: string(latestReview.Author.Login),
		},
		Body:        string(latestReview.Body),
		State:       string(latestReview.State),
		SubmittedAt: latestReview.SubmittedAt,
	}, nil
}

func (target *PullRequestTarget) GetProjectV2ItemID(projectID string) (string, error) {
	ctx := target.ctx
	targetEntity := target.targetEntity
	owner := targetEntity.Owner
	repo := targetEntity.Repo

	return target.githubClient.GetPullRequestProjectV2ItemID(ctx, owner, repo, projectID, targetEntity.Number)
}

func (target *PullRequestTarget) GetLatestApprovedReviews() ([]string, error) {
	clientGQL := target.githubClient.GetClientGraphQL()
	targetEntity := target.targetEntity
	ctx := target.ctx
	owner := targetEntity.Owner
	repo := targetEntity.Repo
	number := targetEntity.Number

	var latestReviewsQuery struct {
		Repository struct {
			PullRequest struct {
				Reviews struct {
					Nodes []struct {
						Author struct {
							Login githubv4.String
						}
						State githubv4.String
					}
				} `graphql:"latestReviews(last: 100)"`
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	varGQLReviews := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(number),
	}

	err := clientGQL.Query(ctx, &latestReviewsQuery, varGQLReviews)
	if err != nil {
		return nil, err
	}

	approvedBy := make([]string, 0)
	for _, review := range latestReviewsQuery.Repository.PullRequest.Reviews.Nodes {
		if review.State == "APPROVED" {
			approvedBy = append(approvedBy, string(review.Author.Login))
		}
	}

	return approvedBy, nil
}
