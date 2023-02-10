// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type GitBlameFileLines struct {
	FromLine uint64
	ToLine   uint64
	Author   string
}

type GitBlameFile struct {
	FilePath  string
	LineCount uint64
	Lines     []GitBlameFileLines
}

type GitBlame struct {
	CommitSHA string
	Files     map[string]GitBlameFile
}

type GitBlameQuery struct {
	Repository struct {
		Object map[string]struct {
			Ranges []struct {
				StartingLine uint64
				EndingLine   uint64
				Age          uint64
				Commit       struct {
					Author struct {
						User struct {
							Login string
						}
					}
				}
			}
		}
	}
}

func (c *GithubClient) GetGitBlame(ctx context.Context, owner, name, commitSHA string, filePaths []string) (*GitBlame, error) {
	query := `query($owner: String!, $name: String!, $objectExpression: String!) {
		repository(owner: $owner, name: $name) {
			object(expression: $objectExpression) {
				... on Commit {
					%s
				}
			}
		}
	}`

	blameQuery := strings.Builder{}

	pathToKeyMap := map[string]string{}

	for i, filePath := range filePaths {
		key := fmt.Sprintf("blame%d", i)
		pathToKeyMap[key] = filePath

		blameQuery.WriteString(fmt.Sprintf(`%s: blame(path: "%s") {
			ranges {
				startingLine
				endingLine
				age
				commit {
					author {
						user {
							login
						}
					}
				}
			}
		}
		`, key, filePath))
	}

	query = fmt.Sprintf(query, blameQuery.String())

	gitBlameQueryData := map[string]interface{}{
		"owner":            owner,
		"name":             name,
		"objectExpression": commitSHA,
	}

	rawRes, err := c.GetRawClientGraphQL().ExecRaw(ctx, query, gitBlameQueryData)
	if err != nil {
		return nil, fmt.Errorf("error executing blame query: %s", err.Error())
	}

	var gitBlameQuery GitBlameQuery

	if err = json.Unmarshal(rawRes, &gitBlameQuery); err != nil {
		return nil, err
	}

	if len(gitBlameQuery.Repository.Object) == 0 {
		return nil, fmt.Errorf("error getting blame information: no blame information found")
	}

	files := map[string]GitBlameFile{}

	for key, blame := range gitBlameQuery.Repository.Object {
		path := pathToKeyMap[key]
		rangesLen := len(blame.Ranges)

		if rangesLen == 0 {
			return nil, fmt.Errorf("error getting blame information: ranges not found for %s", path)
		}

		linesCount := blame.Ranges[rangesLen-1].EndingLine

		file := GitBlameFile{
			FilePath:  path,
			LineCount: linesCount,
			Lines:     []GitBlameFileLines{},
		}

		for _, ran := range blame.Ranges {
			file.Lines = append(file.Lines, GitBlameFileLines{
				FromLine: ran.StartingLine,
				ToLine:   ran.EndingLine,
				Author:   ran.Commit.Author.User.Login,
			})
		}

		files[path] = file
	}

	return &GitBlame{
		CommitSHA: commitSHA,
		Files:     files,
	}, nil
}
