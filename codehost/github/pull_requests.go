// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/shurcooL/githubv4"
)

const maxPerPage int = 100

type GQLReviewThread struct {
	IsResolved githubv4.Boolean
	IsOutdated githubv4.Boolean
}

type ReviewThreadsQuery struct {
	Repository struct {
		PullRequest struct {
			ReviewThreads struct {
				Nodes    []GQLReviewThread
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"reviewThreads(first: 10, after: $reviewThreadsCursor)"`
		} `graphql:"pullRequest(number: $pullRequestNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type LastPushQuery struct {
	Repository struct {
		PullRequest struct {
			TimelineItems struct {
				Nodes []struct {
					Typename                string `graphql:"__typename"`
					HeadRefForcePushedEvent struct {
						CreatedAt *time.Time
					} `graphql:"... on HeadRefForcePushedEvent"`
					PullRequestCommit struct {
						Commit struct {
							PushedDate    *time.Time
							CommittedDate *time.Time
						}
					} `graphql:"... on PullRequestCommit"`
				}
			} `graphql:"timelineItems(last: 1, itemTypes: [HEAD_REF_FORCE_PUSHED_EVENT, PULL_REQUEST_COMMIT])"`
		} `graphql:"pullRequest(number: $pullRequestNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

func GetPullRequestHeadOwnerName(pullRequest *github.PullRequest) string {
	return pullRequest.Head.Repo.Owner.GetLogin()
}

func GetPullRequestHeadRepoName(pullRequest *github.PullRequest) string {
	return pullRequest.Head.Repo.GetName()
}

func GetPullRequestBaseOwnerName(pullRequest *github.PullRequest) string {
	return pullRequest.Base.Repo.Owner.GetLogin()
}

func GetPullRequestBaseRepoName(pullRequest *github.PullRequest) string {
	return pullRequest.Base.Repo.GetName()
}

func GetPullRequestNumber(pullRequest *github.PullRequest) int {
	return pullRequest.GetNumber()
}

func (c *GithubClient) GetPullRequestFiles(ctx context.Context, owner string, repo string, number int) ([]*github.CommitFile, error) {
	fs, err := PaginatedRequest(
		func() interface{} {
			return []*github.CommitFile{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			fls := i.([]*github.CommitFile)
			fs, resp, err := c.clientREST.PullRequests.ListFiles(ctx, owner, repo, number, &github.ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
			if err != nil {
				return nil, nil, err
			}
			fls = append(fls, fs...)
			return fls, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return fs.([]*github.CommitFile), nil
}

func (c *GithubClient) GetPullRequestReviewers(ctx context.Context, owner string, repo string, number int, opts *github.ListOptions) (*github.Reviewers, error) {
	reviewers, err := PaginatedRequest(
		func() interface{} {
			return &github.Reviewers{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentReviewers := i.(*github.Reviewers)
			reviewers, resp, err := c.clientREST.PullRequests.ListReviewers(ctx, owner, repo, number, &github.ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
			if err != nil {
				return nil, nil, err
			}
			currentReviewers.Users = append(currentReviewers.Users, reviewers.Users...)
			currentReviewers.Teams = append(currentReviewers.Teams, reviewers.Teams...)
			return currentReviewers, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return reviewers.(*github.Reviewers), nil
}

func (c *GithubClient) GetRepoCollaborators(ctx context.Context, owner string, repo string) ([]*github.User, error) {
	collaborators, err := PaginatedRequest(
		func() interface{} {
			return []*github.User{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentCollaborators := i.([]*github.User)
			collaborators, resp, err := c.clientREST.Repositories.ListCollaborators(ctx, owner, repo, &github.ListCollaboratorsOptions{
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: maxPerPage,
				},
			})
			if err != nil {
				return nil, nil, err
			}
			currentCollaborators = append(currentCollaborators, collaborators...)
			return currentCollaborators, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return collaborators.([]*github.User), nil
}

func (c *GithubClient) GetIssuesAvailableAssignees(ctx context.Context, owner string, repo string) ([]*github.User, error) {
	assignees, err := PaginatedRequest(
		func() interface{} {
			return []*github.User{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentAssignees := i.([]*github.User)
			assignees, resp, err := c.clientREST.Issues.ListAssignees(ctx, owner, repo, &github.ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
			if err != nil {
				return nil, nil, err
			}
			currentAssignees = append(currentAssignees, assignees...)
			return currentAssignees, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return assignees.([]*github.User), nil
}

func (c *GithubClient) GetPullRequestCommits(ctx context.Context, owner string, repo string, number int) ([]*github.RepositoryCommit, error) {
	commits, err := PaginatedRequest(
		func() interface{} {
			return []*github.RepositoryCommit{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentCommits := i.([]*github.RepositoryCommit)
			commits, resp, err := c.clientREST.PullRequests.ListCommits(ctx, owner, repo, number, &github.ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
			if err != nil {
				return nil, nil, err
			}
			currentCommits = append(currentCommits, commits...)
			return currentCommits, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return commits.([]*github.RepositoryCommit), nil
}

func (c *GithubClient) GetPullRequestReviews(ctx context.Context, owner string, repo string, number int) ([]*github.PullRequestReview, error) {
	reviews, err := PaginatedRequest(
		func() interface{} {
			return []*github.PullRequestReview{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentReviews := i.([]*github.PullRequestReview)
			reviews, resp, err := c.clientREST.PullRequests.ListReviews(ctx, owner, repo, number, &github.ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
			if err != nil {
				return nil, nil, err
			}
			currentReviews = append(currentReviews, reviews...)
			return currentReviews, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return reviews.([]*github.PullRequestReview), nil
}

func (c *GithubClient) GetPullRequests(ctx context.Context, owner string, repo string) ([]*github.PullRequest, error) {
	prs, err := PaginatedRequest(
		func() interface{} {
			return []*github.PullRequest{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			allPrs := i.([]*github.PullRequest)
			prs, resp, err := c.clientREST.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: maxPerPage,
				},
			})
			if err != nil {
				return nil, nil, err
			}
			allPrs = append(allPrs, prs...)
			return allPrs, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return prs.([]*github.PullRequest), nil
}

func (c *GithubClient) GetPullRequest(ctx context.Context, owner string, repo string, number int) (*github.PullRequest, *github.Response, error) {
	return c.clientREST.PullRequests.Get(ctx, owner, repo, number)
}

func (c *GithubClient) GetReviewThreads(ctx context.Context, owner string, repo string, number int, retryCount int) ([]GQLReviewThread, error) {
	var reviewThreadsQuery ReviewThreadsQuery
	reviewThreads := make([]GQLReviewThread, 0)
	hasNextPage := true

	varGQLReviewThreads := map[string]interface{}{
		"repositoryOwner":     githubv4.String(owner),
		"repositoryName":      githubv4.String(repo),
		"pullRequestNumber":   githubv4.Int(number),
		"reviewThreadsCursor": (*githubv4.String)(nil),
	}

	currentRequestRetry := 1

	for hasNextPage {
		err := c.clientGQL.Query(context.Background(), &reviewThreadsQuery, varGQLReviewThreads)
		if err != nil {
			currentRequestRetry++
			if currentRequestRetry <= retryCount {
				continue
			} else {
				return nil, err
			}
		} else {
			currentRequestRetry = 0
		}

		reviewThreads = append(reviewThreads, reviewThreadsQuery.Repository.PullRequest.ReviewThreads.Nodes...)
		hasNextPage = reviewThreadsQuery.Repository.PullRequest.ReviewThreads.PageInfo.HasNextPage
		varGQLReviewThreads["reviewThreadsCursor"] = githubv4.NewString(reviewThreadsQuery.Repository.PullRequest.ReviewThreads.PageInfo.EndCursor)
	}

	return reviewThreads, nil
}

func (c *GithubClient) GetIssueTimeline(ctx context.Context, owner string, repo string, number int) ([]*github.Timeline, error) {
	events, err := PaginatedRequest(
		func() interface{} {
			return []*github.Timeline{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentEvents := i.([]*github.Timeline)
			events, resp, err := c.clientREST.Issues.ListIssueTimeline(ctx, owner, repo, number, &github.ListOptions{
				Page:    page,
				PerPage: maxPerPage,
			})
			if err != nil {
				return nil, nil, err
			}
			currentEvents = append(currentEvents, events...)
			return currentEvents, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return events.([]*github.Timeline), nil
}

func (c *GithubClient) RequestReviewers(ctx context.Context, owner string, repo string, number int, reviewers github.ReviewersRequest) (*github.PullRequest, *github.Response, error) {
	return c.clientREST.PullRequests.RequestReviewers(ctx, owner, repo, number, reviewers)
}

func (c *GithubClient) EditPullRequest(ctx context.Context, owner string, repo string, number int, pull *github.PullRequest) (*github.PullRequest, *github.Response, error) {
	return c.clientREST.PullRequests.Edit(ctx, owner, repo, number, pull)
}

func (c *GithubClient) Merge(ctx context.Context, owner string, repo string, number int, commitMessage string, options *github.PullRequestOptions) (*github.PullRequestMergeResult, *github.Response, error) {
	return c.clientREST.PullRequests.Merge(ctx, owner, repo, number, commitMessage, options)
}

func (c *GithubClient) GetPullRequestClosingIssuesCount(ctx context.Context, owner string, repo string, number int) (int, error) {
	var pullRequestQuery struct {
		Repository struct {
			PullRequest struct {
				ClosingIssuesReferences struct {
					TotalCount githubv4.Int
				}
			} `graphql:"pullRequest(number: $pullRequestNumber)"`
		} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
	}

	varGQLPullRequestQuery := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(number),
	}

	err := c.GetClientGraphQL().Query(ctx, &pullRequestQuery, varGQLPullRequestQuery)
	if err != nil {
		return 0, err
	}

	return int(pullRequestQuery.Repository.PullRequest.ClosingIssuesReferences.TotalCount), nil
}

func (c *GithubClient) GetPullRequestLastPushDate(ctx context.Context, owner string, repo string, number int) (time.Time, error) {
	var lastPushQuery LastPushQuery
	varGQLastPush := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(number),
	}

	err := c.GetClientGraphQL().Query(ctx, &lastPushQuery, varGQLastPush)
	if err != nil {
		return time.Time{}, err
	}

	hasLastPush := len(lastPushQuery.Repository.PullRequest.TimelineItems.Nodes) > 0
	if !hasLastPush {
		return time.Time{}, errors.New("last push not found")
	}

	var pushDate *time.Time

	event := lastPushQuery.Repository.PullRequest.TimelineItems.Nodes[0]
	switch event.Typename {
	case "PullRequestCommit":
		if pushedDate := event.PullRequestCommit.Commit.PushedDate; pushDate != nil {
			pushDate = pushedDate
		} else {
			pushDate = event.PullRequestCommit.Commit.PushedDate
		}
	case "HeadRefForcePushedEvent":
		pushDate = event.HeadRefForcePushedEvent.CreatedAt
	default:
		return time.Time{}, fmt.Errorf("unknown event type %v", event.Typename)
	}

	if pushDate == nil {
		return time.Time{}, errors.New("last push not found")
	}

	return *pushDate, nil
}
