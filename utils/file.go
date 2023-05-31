// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import (
	"bytes"
	"context"
	"strings"

	pbc "github.com/reviewpad/api/go/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github"
	gh "github.com/reviewpad/reviewpad/v4/codehost/github"
	"github.com/sirupsen/logrus"
)

const pullRequestFileLimit = 50

func FileExt(fp string) string {
	strs := strings.Split(fp, ".")
	str := strings.Join(strs[1:], ".")

	if str != "" {
		str = "." + str
	}

	return str
}

// reviewpad-an: experimental
// ReviewpadFileChanges checks if a file path was changed in a pull request.
// The way this is done depends on the number of files changed in the pull request.
// If the number of files changed is greater than pullRequestFileLimit,
// then we download both files using the filePath and check their contents.
// This strategy assumes that the file path exists in the head branch.
// Otherwise, we download the pull request files and check the filePath exists in them.
func ReviewpadFileChanged(ctx context.Context, githubClient *gh.GithubClient, filePath string, pullRequest *pbc.PullRequest) (bool, error) {
	if pullRequest.ChangedFilesCount > pullRequestFileLimit {
		rawHeadFile, err := githubClient.DownloadContents(ctx, filePath, pullRequest.Head, &github.DownloadContentsOptions{
			Method: github.DownloadMethodSHA,
		})
		if err != nil {
			if strings.HasPrefix(err.Error(), "no file named") {
				return true, nil
			}
			return false, err
		}

		rawBaseFile, err := githubClient.DownloadContents(ctx, filePath, pullRequest.Base, &github.DownloadContentsOptions{
			Method: github.DownloadMethodSHA,
		})
		if err != nil {
			if strings.HasPrefix(err.Error(), "no file named") {
				return true, nil
			}
			return false, err
		}

		// FIXME: use the hashes of the files
		return !bytes.Equal(rawBaseFile, rawHeadFile), nil
	}

	branchRepoOwner := pullRequest.Base.Repo.Owner
	branchRepoName := pullRequest.Base.Repo.Name

	files, err := githubClient.GetPullRequestFiles(ctx, branchRepoOwner, branchRepoName, int(pullRequest.Number))
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if filePath == *file.Filename {
			return true, nil
		}
	}

	return false, nil
}

func DownloadReviewpadFileFromGitHub(
	ctx context.Context,
	logger *logrus.Entry,
	githubClient *gh.GithubClient,
	filePath string,
	branch *pbc.Branch,
	options *github.DownloadContentsOptions,
) (*bytes.Buffer, error) {
	reviewpadFileContent, err := githubClient.DownloadContents(ctx, filePath, branch, options)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(reviewpadFileContent), nil
}
