// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/tomnomnom/linkheader"
)

const maxPerPage int = 100

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

func PaginatedRequest(
	initFn func() interface{},
	reqFn func(interface{}, int) (interface{}, *github.Response, error),
) (interface{}, error) {
	page := 1
	results, resp, err := reqFn(initFn(), page)
	if err != nil {
		return nil, err
	}

	numPages := ParseNumPages(resp)
	page++
	for page <= numPages && resp.NextPage > page {
		results, _, err = reqFn(results, page)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func ParseNumPagesFromLink(link string) int {
	urlInfo := linkheader.Parse(link).FilterByRel("last")
	if len(urlInfo) < 1 {
		return 0
	}

	urlData, err := url.Parse(urlInfo[0].URL)
	if err != nil {
		return 0
	}

	numPagesStr := urlData.Query().Get("page")
	if numPagesStr == "" {
		return 0
	}

	numPages, err := strconv.ParseInt(numPagesStr, 10, 32)
	if err != nil {
		return 0
	}

	return int(numPages)
}

//ParseNumPages Given a link header string representing pagination info, returns total number of pages.
func ParseNumPages(resp *github.Response) int {
	link := resp.Header.Get("Link")
	if strings.Trim(link, " ") == "" {
		return 0
	}

	return ParseNumPagesFromLink(link)
}

func GetPullRequestComments(ctx context.Context, client *github.Client, owner string, repo string, number int, opts *github.IssueListCommentsOptions) ([]*github.IssueComment, error) {
	fs, err := PaginatedRequest(
		func() interface{} {
			return []*github.IssueComment{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			fls := i.([]*github.IssueComment)
			fs, resp, err := client.Issues.ListComments(ctx, owner, repo, number, &github.IssueListCommentsOptions{
				Sort:      opts.Sort,
				Direction: opts.Direction,
				Since:     opts.Since,
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: maxPerPage,
				},
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

	return fs.([]*github.IssueComment), nil
}

func GetPullRequestFiles(ctx context.Context, client *github.Client, owner string, repo string, number int) ([]*github.CommitFile, error) {
	fs, err := PaginatedRequest(
		func() interface{} {
			return []*github.CommitFile{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			fls := i.([]*github.CommitFile)
			fs, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, number, &github.ListOptions{
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

func GetPullRequestReviewers(ctx context.Context, client *github.Client, owner string, repo string, number int, opts *github.ListOptions) (*github.Reviewers, error) {
	reviewers, err := PaginatedRequest(
		func() interface{} {
			return &github.Reviewers{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentReviewers := i.(*github.Reviewers)
			reviewers, resp, err := client.PullRequests.ListReviewers(ctx, owner, repo, number, &github.ListOptions{
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

func GetRepoCollaborators(ctx context.Context, client *github.Client, owner string, repo string) ([]*github.User, error) {
	collaborators, err := PaginatedRequest(
		func() interface{} {
			return []*github.User{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentCollaborators := i.([]*github.User)
			collaborators, resp, err := client.Repositories.ListCollaborators(ctx, owner, repo, &github.ListCollaboratorsOptions{
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

func GetIssuesAvailableAssignees(ctx context.Context, client *github.Client, owner string, repo string) ([]*github.User, error) {
	assignees, err := PaginatedRequest(
		func() interface{} {
			return []*github.User{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentAssignees := i.([]*github.User)
			assignees, resp, err := client.Issues.ListAssignees(ctx, owner, repo, &github.ListOptions{
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

func GetPullRequestCommits(ctx context.Context, client *github.Client, owner string, repo string, number int) ([]*github.RepositoryCommit, error) {
	commits, err := PaginatedRequest(
		func() interface{} {
			return []*github.RepositoryCommit{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentCommits := i.([]*github.RepositoryCommit)
			commits, resp, err := client.PullRequests.ListCommits(ctx, owner, repo, number, &github.ListOptions{
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

// GetPullRequestReviews fetch pull request reviews
func GetPullRequestReviews(ctx context.Context, client *github.Client, owner string, repo string, number int, opts *github.ListOptions) ([]*github.PullRequestReview, error) {
	reviews, err := PaginatedRequest(
		func() interface{} {
			return []*github.PullRequestReview{}
		},
		func(i interface{}, page int) (interface{}, *github.Response, error) {
			currentReviews := i.([]*github.PullRequestReview)
			reviews, resp, err := client.PullRequests.ListReviews(ctx, owner, repo, number, &github.ListOptions{
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
