// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/v49/github"
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

type FirstCommitAndReviewDateQuery struct {
	Repository struct {
		PullRequest struct {
			Commits struct {
				Nodes []struct {
					Commit struct {
						AuthoredDate time.Time
					}
				}
			} `graphql:"commits(first: 1)"`
			Reviews struct {
				Nodes []struct {
					CreatedAt time.Time
				}
			} `graphql:"reviews(first: 1)"`
		} `graphql:"pullRequest(number: $pullRequestNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type GetObjectQuery struct {
	Repository struct {
		Object struct {
			Blog struct {
				IsBinary bool
			} `graphql:"... on Blob"`
		} `graphql:"object(expression: $expression)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type GetLastCommitSHAQuery struct {
	Repository struct {
		PullRequest struct {
			Commits struct {
				Nodes []struct {
					Commit struct {
						OID string
					}
				}
			} `graphql:"commits(last: 1)"`
		} `graphql:"pullRequest(number: $pullRequestNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type PullRequestLinkedProjectsQuery struct {
	Repository struct {
		PullRequest struct {
			ProjectItems *struct {
				Nodes    []GQLProjectV2Item
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"projectItems(first: 10, after: $projectItemsCursor)"`
		} `graphql:"pullRequest(number: $issueNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type GetApprovalsCountQuery struct {
	Repository struct {
		PullRequest struct {
			Reviews struct {
				TotalCount int `graphql:"totalCount"`
			} `graphql:"reviews(first: 1, states: [APPROVED])"`
		} `graphql:"pullRequest(number: $pullRequestNumber)"`
	} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
}

type GetRefIDQuery struct {
	Repository struct {
		Ref struct {
			ID string
		} `graphql:"ref(qualifiedName: $ref)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

type LastFiftyOpenedPullRequestsQuery struct {
	Repository struct {
		PullRequests struct {
			Nodes []struct {
				ReviewRequests struct {
					Nodes []struct {
						RequestedReviewer struct {
							TypeName string `graphql:"__typename"`
							AsUser   struct {
								Login githubv4.String
							} `graphql:"... on User"`
						}
					}
				} `graphql:"reviewRequests(first: 50)"`
				Reviews struct {
					Nodes []struct {
						Author struct {
							Login githubv4.String
						}
					}
				} `graphql:"reviews(first: 50)"`
			}
		} `graphql:"pullRequests(states: OPEN, last: 50)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
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

func (c *GithubClient) Review(ctx context.Context, owner string, repo string, number int, review *github.PullRequestReviewRequest) (*github.PullRequestReview, *github.Response, error) {
	return c.clientREST.PullRequests.CreateReview(ctx, owner, repo, number, review)
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
		if pushedDate := event.PullRequestCommit.Commit.PushedDate; pushedDate != nil {
			pushDate = pushedDate
		} else {
			pushDate = event.PullRequestCommit.Commit.CommittedDate
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

func (c *GithubClient) DeleteReference(ctx context.Context, owner, repo, ref string) error {
	_, err := c.clientREST.Git.DeleteRef(ctx, owner, repo, ref)
	return err
}

func (c *GithubClient) GetFirstCommitAndReviewDate(ctx context.Context, owner, repo string, number int) (*time.Time, *time.Time, error) {
	var firstCommitAndReviewDateQuery FirstCommitAndReviewDateQuery
	varGQLFirstCommitDate := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(number),
	}

	err := c.GetClientGraphQL().Query(ctx, &firstCommitAndReviewDateQuery, varGQLFirstCommitDate)
	if err != nil {
		return nil, nil, err
	}

	var firstCommit, firstReview *time.Time

	if len(firstCommitAndReviewDateQuery.Repository.PullRequest.Commits.Nodes) == 1 {
		firstCommit = &firstCommitAndReviewDateQuery.Repository.PullRequest.Commits.Nodes[0].Commit.AuthoredDate
	}

	if len(firstCommitAndReviewDateQuery.Repository.PullRequest.Reviews.Nodes) == 1 {
		firstReview = &firstCommitAndReviewDateQuery.Repository.PullRequest.Reviews.Nodes[0].CreatedAt
	}

	return firstCommit, firstReview, nil
}

func (c *GithubClient) GetCheckRunsForRef(ctx context.Context, owner string, repo string, number int, ref string, opts *github.ListCheckRunsOptions) ([]*github.CheckRun, error) {
	checks, err := PaginatedRequest(
		func() interface{} {
			return []*github.CheckRun{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentChecks := i.([]*github.CheckRun)
			checks, resp, err := c.clientREST.Checks.ListCheckRunsForRef(ctx, owner, repo, ref, &github.ListCheckRunsOptions{
				CheckName: opts.CheckName,
				Status:    opts.Status,
				AppID:     opts.AppID,
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: maxPerPage,
				},
			})
			if err != nil {
				return nil, nil, err
			}
			currentChecks = append(currentChecks, checks.CheckRuns...)
			return currentChecks, resp, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return checks.([]*github.CheckRun), nil
}

func (c *GithubClient) IsFileBinary(ctx context.Context, owner, repo, branch, file string) (bool, error) {
	var getObjectQuery GetObjectQuery
	varGQLGetObject := map[string]interface{}{
		"repositoryOwner": githubv4.String(owner),
		"repositoryName":  githubv4.String(repo),
		"expression":      githubv4.String(fmt.Sprintf("%s:%s", branch, file)),
	}

	err := c.GetClientGraphQL().Query(ctx, &getObjectQuery, varGQLGetObject)
	if err != nil {
		return false, err
	}

	return getObjectQuery.Repository.Object.Blog.IsBinary, nil
}

func (c *GithubClient) GetLinkedProjectsForPullRequest(ctx context.Context, owner, repo string, number int, retryCount int) ([]GQLProjectV2Item, error) {
	projectItems := []GQLProjectV2Item{}
	hasNextPage := true
	currentRequestRetry := 1

	varGQLGetProjectFieldsQuery := map[string]interface{}{
		"repositoryOwner":    githubv4.String(owner),
		"repositoryName":     githubv4.String(repo),
		"issueNumber":        githubv4.Int(number),
		"projectItemsCursor": githubv4.String(""),
	}

	var getLinkedProjects PullRequestLinkedProjectsQuery

	for hasNextPage {
		if err := c.clientGQL.Query(ctx, &getLinkedProjects, varGQLGetProjectFieldsQuery); err != nil {
			currentRequestRetry++
			if currentRequestRetry <= retryCount {
				continue
			}
			return nil, err
		}

		items := getLinkedProjects.Repository.PullRequest.ProjectItems
		if items == nil {
			return nil, ErrProjectItemsNotFound
		}

		projectItems = append(projectItems, items.Nodes...)

		hasNextPage = items.PageInfo.HasNextPage

		varGQLGetProjectFieldsQuery["projectItemsCursor"] = githubv4.String(items.PageInfo.EndCursor)
	}

	return projectItems, nil
}

func (c *GithubClient) GetLastCommitSHA(ctx context.Context, owner, repo string, number int) (string, error) {
	var getLastCommitQuery GetLastCommitSHAQuery
	varGQLLastCommitSHAData := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(number),
	}

	err := c.GetClientGraphQL().Query(ctx, &getLastCommitQuery, varGQLLastCommitSHAData)
	if err != nil {
		return "", err
	}

	if len(getLastCommitQuery.Repository.PullRequest.Commits.Nodes) < 1 {
		return "", nil
	}

	return getLastCommitQuery.Repository.PullRequest.Commits.Nodes[0].Commit.OID, nil
}

func (c *GithubClient) GetApprovalsCount(ctx context.Context, owner, repo string, number int) (int, error) {
	var getApprovalsCountQuery GetApprovalsCountQuery
	varGQLApprovalCountData := map[string]interface{}{
		"repositoryOwner":   githubv4.String(owner),
		"repositoryName":    githubv4.String(repo),
		"pullRequestNumber": githubv4.Int(number),
	}

	err := c.GetClientGraphQL().Query(ctx, &getApprovalsCountQuery, varGQLApprovalCountData)
	if err != nil {
		return 0, err
	}

	return getApprovalsCountQuery.Repository.PullRequest.Reviews.TotalCount, nil
}

func (c *GithubClient) RefExists(ctx context.Context, owner, repo, ref string) (bool, error) {
	var getRefIDQuery GetRefIDQuery
	varGetRefIDQueryData := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
		"ref":   githubv4.String(ref),
	}

	err := c.GetClientGraphQL().Query(ctx, &getRefIDQuery, varGetRefIDQueryData)
	if err != nil {
		return false, err
	}

	return getRefIDQuery.Repository.Ref.ID != "", nil
}

func (c *GithubClient) GetOpenPullRequestsAsReviewer(ctx context.Context, owner string, repo string, usernames []string) (map[string]int, error) {
	if usernames == nil {
		usernames = make([]string, 0)
	}

	if len(usernames) == 0 {
		repoCollaborators, err := c.GetRepoCollaborators(ctx, owner, repo)
		if err != nil {
			return nil, err
		}

		for _, collaborator := range repoCollaborators {
			usernames = append(usernames, collaborator.GetLogin())
		}
	}

	totalOpenPullRequestsByUser := make(map[string]int)
	for _, username := range usernames {
		totalOpenPullRequestsByUser[username] = 0
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	}

	var lastFiftyOpenedPullRequests LastFiftyOpenedPullRequestsQuery

	err := c.GetClientGraphQL().Query(ctx, &lastFiftyOpenedPullRequests, variables)
	if err != nil {
		return nil, err
	}

	for _, pullRequest := range lastFiftyOpenedPullRequests.Repository.PullRequests.Nodes {
		isRequestedReviewerByUser := make(map[string]bool)
		for _, reviewRequest := range pullRequest.ReviewRequests.Nodes {
			requestedReviewer := string(reviewRequest.RequestedReviewer.AsUser.Login)
			if contains(usernames, requestedReviewer) {
				isRequestedReviewerByUser[requestedReviewer] = true
				totalOpenPullRequestsByUser[requestedReviewer]++
			}
		}

		hasReviewedByUser := make(map[string]bool)
		for _, review := range pullRequest.Reviews.Nodes {
			reviewer := string(review.Author.Login)
			if contains(usernames, reviewer) {
				isRequestedReviewer, ok := isRequestedReviewerByUser[reviewer]
				if ok && isRequestedReviewer {
					continue
				}

				userHasReviewed, ok := hasReviewedByUser[reviewer]
				if ok && userHasReviewed {
					continue
				}

				hasReviewedByUser[reviewer] = true
				totalOpenPullRequestsByUser[reviewer]++
			}
		}
	}

	return totalOpenPullRequestsByUser, nil
}

func contains(slice []string, s string) bool {
	for _, element := range slice {
		if element == s {
			return true
		}
	}
	return false
}
