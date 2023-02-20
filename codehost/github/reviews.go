// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

type ReviewsByUserTotalCountQuery struct {
	Search struct {
		PageInfo struct {
			EndCursor   githubv4.String
			HasNextPage bool
		}
		Nodes []struct {
			PullRequest struct {
				Reviews struct {
					TotalCount githubv4.Int
				} `graphql:"reviews(author: $author)"`
			} `graphql:"... on PullRequest"`
		}
	} `graphql:"search(query: $query, type: ISSUE, first: $perPage, after: $afterCursor)"`
}

func (c *GithubClient) GetReviewsCountByUserFromOpenPullRequests(ctx context.Context, userOrOrgLogin, username string) (int, error) {
	var totalCount int
	var reviewsByUserTotalCount ReviewsByUserTotalCountQuery

	reviewsByUserTotalCountData := map[string]interface{}{
		"query":       githubv4.String(fmt.Sprintf("user:%s is:pr is:open reviewed-by:%s", userOrOrgLogin, username)),
		"author":      githubv4.String(username),
		"perPage":     githubv4.Int(100),
		"afterCursor": (*githubv4.String)(nil),
	}

	for {
		err := c.GetClientGraphQL().Query(ctx, &reviewsByUserTotalCount, reviewsByUserTotalCountData)
		if err != nil {
			return 0, err
		}

		for _, node := range reviewsByUserTotalCount.Search.Nodes {
			totalCount += int(node.PullRequest.Reviews.TotalCount)
		}

		if !reviewsByUserTotalCount.Search.PageInfo.HasNextPage {
			break
		}

		reviewsByUserTotalCountData["afterCursor"] = &reviewsByUserTotalCount.Search.PageInfo.EndCursor
	}

	return totalCount, nil
}
