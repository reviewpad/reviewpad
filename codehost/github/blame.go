// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
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

type FileSliceLimits struct {
	From uint64
	To   uint64
}

type GitBlameAuthorRank struct {
	Username   string
	TotalLines int
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
	graphQLQuery, mapQueryKeyToFilePath := buildGitBlameGraphQLQuery(filePaths)
	gitBlameQueryVariables := map[string]interface{}{
		"owner":            owner,
		"name":             name,
		"objectExpression": commitSHA,
	}

	rawResponse, err := c.GetRawClientGraphQL().WithDebug(true).ExecRaw(ctx, graphQLQuery, gitBlameQueryVariables)
	if err != nil {
		return nil, fmt.Errorf("error executing blame query: %s", err.Error())
	}

	var gitBlameQuery GitBlameQuery

	if err = json.Unmarshal(rawResponse, &gitBlameQuery); err != nil {
		return nil, err
	}

	if len(gitBlameQuery.Repository.Object) == 0 {
		return nil, fmt.Errorf("no blame information found")
	}

	return mapGitBlameQueryToGitBlame(commitSHA, gitBlameQuery, mapQueryKeyToFilePath)
}

func (c *GithubClient) SliceGitBlame(blame *GitBlame, filePath string, slices []FileSliceLimits) (*GitBlame, error) {
	file, ok := blame.Files[filePath]
	if !ok {
		return nil, errors.New("file not found in blame information")
	}

	slicedBlame := &GitBlame{
		CommitSHA: blame.CommitSHA,
		Files:     map[string]GitBlameFile{},
	}

	gitBlameFile := GitBlameFile{
		FilePath:  file.FilePath,
		LineCount: file.LineCount,
		Lines:     []GitBlameFileLines{},
	}

	for _, slice := range slices {
		for _, line := range file.Lines {
			fromLine := line.FromLine
			toLine := line.ToLine

			if line.ToLine < slice.From || line.FromLine > slice.To {
				continue
			}

			if line.FromLine < slice.From {
				fromLine = slice.From
			}

			if line.ToLine > slice.To {
				toLine = slice.To
			}

			gitBlameFile.Lines = append(gitBlameFile.Lines, GitBlameFileLines{
				FromLine: fromLine,
				ToLine:   toLine,
				Author:   line.Author,
			})
		}
	}

	slicedBlame.Files[filePath] = gitBlameFile

	return slicedBlame, nil
}

func (c *GithubClient) ComputeGitBlameRank(gitBlame *GitBlame) []GitBlameAuthorRank {
	var rankedAuthors []GitBlameAuthorRank
	authorIndices := make(map[string]int)

	for _, file := range gitBlame.Files {
		for _, line := range file.Lines {
			author := line.Author
			lineCount := int(line.ToLine - line.FromLine + 1)

			authorIndex, ok := authorIndices[author]
			if !ok {
				authorIndex = len(rankedAuthors)
				authorIndices[author] = authorIndex
				rankedAuthors = append(rankedAuthors, GitBlameAuthorRank{
					Username:   author,
					TotalLines: 0,
				})
			}

			rankedAuthors[authorIndex].TotalLines += lineCount
		}
	}

	sort.Slice(rankedAuthors, func(i, j int) bool {
		return rankedAuthors[i].TotalLines > rankedAuthors[j].TotalLines
	})

	return rankedAuthors
}

func buildGitBlameGraphQLQuery(filePaths []string) (string, map[string]string) {
	graphQLQuery := `query($owner: String!, $name: String!, $objectExpression: String!) {
		repository(owner: $owner, name: $name) {
			object(expression: $objectExpression) {
				... on Commit {
					%s
				}
			}
		}
	}`

	blameQuery := strings.Builder{}

	mapQueryKeyToFilePath := map[string]string{}

	for i, filePath := range filePaths {
		key := fmt.Sprintf("blame%d", i)
		mapQueryKeyToFilePath[key] = filePath

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

	graphQLQuery = fmt.Sprintf(graphQLQuery, blameQuery.String())

	return graphQLQuery, mapQueryKeyToFilePath
}

func mapGitBlameQueryToGitBlame(commitSHA string, gitBlameQuery GitBlameQuery, pathToKeyMap map[string]string) (*GitBlame, error) {
	files := map[string]GitBlameFile{}

	for queryKey, fileGitBlame := range gitBlameQuery.Repository.Object {
		filePath := pathToKeyMap[queryKey]
		rangesLen := len(fileGitBlame.Ranges)
		var linesCount uint64

		if rangesLen > 0 {
			linesCount = fileGitBlame.Ranges[rangesLen-1].EndingLine
		}

		file := GitBlameFile{
			FilePath:  filePath,
			LineCount: linesCount,
			Lines:     []GitBlameFileLines{},
		}

		for _, ran := range fileGitBlame.Ranges {
			file.Lines = append(file.Lines, GitBlameFileLines{
				FromLine: ran.StartingLine,
				ToLine:   ran.EndingLine,
				Author:   ran.Commit.Author.User.Login,
			})
		}

		files[filePath] = file
	}

	return &GitBlame{
		CommitSHA: commitSHA,
		Files:     files,
	}, nil
}
