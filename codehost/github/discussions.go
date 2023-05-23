// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package github

import (
	"context"
	"errors"

	"github.com/shurcooL/githubv4"
)

var (
	ErrDiscussionNotFound = errors.New("discussion not found")
)

type AddDiscussionCommentMutation struct {
	AddDiscussionComment struct {
		ClientMutationID string
	} `graphql:"addDiscussionComment(input: $input)"`
}

type AddReactionMutation struct {
	AddReaction struct {
		ClientMutationID string
	} `graphql:"addReaction(input: $input)"`
}

type RemoveReactionMutation struct {
	RemoveReaction struct {
		ClientMutationID string
	} `graphql:"removeReaction(input: $input)"`
}

type DiscussionComment struct {
	ID          string
	Body        string
	AuthorLogin string
}

type GetDiscussionIDQuery struct {
	Repository struct {
		Discussion struct {
			ID string
		} `graphql:"discussion(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type GetDiscussionCommentsQuery struct {
	Repository struct {
		Discussion struct {
			ID     string
			Body   string
			Author struct {
				Login string
			}
			Comments struct {
				PageInfo PageInfo
				Nodes    []struct {
					ID     string
					Body   string
					Author struct {
						Login string
					}
				}
			} `graphql:"comments(first: 50, after: $afterCursor)"`
		} `graphql:"discussion(number: $number)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (c *GithubClient) GetDiscussionID(ctx context.Context, owner, repo string, discussionNum int) (string, error) {
	var getDiscussionIDQuery GetDiscussionIDQuery
	varGQLGetDiscussionIDQuery := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(repo),
		"number": githubv4.Int(discussionNum),
	}

	if err := c.clientGQL.Query(ctx, &getDiscussionIDQuery, varGQLGetDiscussionIDQuery); err != nil {
		return "", err
	}

	return getDiscussionIDQuery.Repository.Discussion.ID, ErrDiscussionNotFound
}

// GetDiscussionComments returns the discussion comments for a given discussion number.
// The first element of the list is the body of the discussion.
func (c *GithubClient) GetDiscussionComments(ctx context.Context, owner, repo string, discussionNum int) ([]DiscussionComment, error) {
	var getDiscussionCommentsQuery GetDiscussionCommentsQuery
	varGQLGetDiscussionCommentsQuery := map[string]interface{}{
		"owner":       githubv4.String(owner),
		"name":        githubv4.String(repo),
		"number":      githubv4.Int(discussionNum),
		"afterCursor": githubv4.String(""),
	}

	comments := make([]DiscussionComment, 0)

	i := 0
	for {
		if err := c.clientGQL.Query(ctx, &getDiscussionCommentsQuery, varGQLGetDiscussionCommentsQuery); err != nil {
			return nil, err
		}

		if i == 0 {
			comments = append(comments, DiscussionComment{
				ID:          getDiscussionCommentsQuery.Repository.Discussion.ID,
				Body:        getDiscussionCommentsQuery.Repository.Discussion.Body,
				AuthorLogin: getDiscussionCommentsQuery.Repository.Discussion.Author.Login,
			})
		}
		i++

		for _, node := range getDiscussionCommentsQuery.Repository.Discussion.Comments.Nodes {
			comments = append(comments, DiscussionComment{
				ID:          node.ID,
				Body:        node.Body,
				AuthorLogin: node.Author.Login,
			})
		}

		if !getDiscussionCommentsQuery.Repository.Discussion.Comments.PageInfo.HasNextPage {
			break
		}

		varGQLGetDiscussionCommentsQuery["afterCursor"] = githubv4.String(getDiscussionCommentsQuery.Repository.Discussion.Comments.PageInfo.EndCursor)
	}

	return comments, nil
}

func (c *GithubClient) AddCommentToDiscussion(ctx context.Context, discussionID, body string) error {
	var addDiscussionCommentMutation AddDiscussionCommentMutation
	addDiscussionCommentInput := githubv4.AddDiscussionCommentInput{
		DiscussionID: discussionID,
		Body:         githubv4.String(body),
	}

	return c.clientGQL.Mutate(ctx, &addDiscussionCommentMutation, addDiscussionCommentInput, nil)
}

func (c *GithubClient) AddReactionToDiscussionComment(ctx context.Context, subjectID string, reaction githubv4.ReactionContent) error {
	var addReactionMutation AddReactionMutation
	addReactionMutationInput := githubv4.AddReactionInput{
		SubjectID: subjectID,
		Content:   reaction,
	}

	return c.clientGQL.Mutate(ctx, &addReactionMutation, addReactionMutationInput, nil)
}

func (c *GithubClient) RemoveReactionToDiscussionComment(ctx context.Context, subjectID string, reaction githubv4.ReactionContent) error {
	var removeReactionMutation RemoveReactionMutation
	removeReactionMutationInput := githubv4.RemoveReactionInput{
		SubjectID: subjectID,
		Content:   reaction,
	}

	return c.clientGQL.Mutate(ctx, &removeReactionMutation, removeReactionMutationInput, nil)
}
