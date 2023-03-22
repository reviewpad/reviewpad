// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_semantic

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v49/github"
	"github.com/migueleliasweb/go-github-mock/src/mock"
	"github.com/reviewpad/api/go/entities"
	pbe "github.com/reviewpad/api/go/entities"
	"github.com/reviewpad/api/go/services"
	"github.com/reviewpad/api/go/services_mocks"
	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v4/plugins/aladino/services"
	"github.com/reviewpad/reviewpad/v4/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetBlocks(t *testing.T) {
	fileName := "crawlerWithRemovedBlock.go"
	patchData := `
@@ -2,9 +2,11 @@ package main
 
 import "fmt"
 
-// First version as presented at:
-// https://gist.github.com/harryhare/6a4979aa7f8b90db6cbc74400d0beb49#file-exercise-web-crawler-go
 func Crawl(url string, depth int, fetcher Fetcher) {
+	defer  c.wg.Done()
+
 	if depth <= 0 {
 		return
 	}
@@ -18,6 +20,7 @@ func Crawl(url string, depth int, fetcher Fetcher) {
 	}
 	fmt.Printf("found: %s %q\n", url, body)
 	for _, u := range urls {
+		c.wg.Add(1)
 		go Crawl(u, depth-1, fetcher)
 	}
 	return
`
	ghFile := &pbe.CommitFile{
		Patch:    patchData,
		Filename: fileName,
	}

	patchFile, err := codehost.NewFile(ghFile)
	assert.Nil(t, err)

	expectedRes := []*entities.ResolveBlockDiff{
		{
			Start: 2,
			End:   4,
		},
		{
			Start: 5,
			End:   6,
		},
		{
			Start: 5,
			End:   5,
		},
		{
			Start: 6,
			End:   7,
		},
		{
			Start: 8,
			End:   12,
		},
	}

	actualRes := getBlocks(patchFile)

	assert.Equal(t, &expectedRes, &actualRes)
}

func TestGetSymbolsFromPatch_WhenDownloadContentsFails(t *testing.T) {
	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedHeadSHA := "abc123"
	mockedPatchFilePath := "test"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
	// Since the patch is simply passed around, it can be an empty string
	mockedPatch := ""
	mockedBlobId := "1234"
	codeReview := &pbe.CodeReview{
		Number:   6,
		IsMerged: true,
		Base: &pbe.Branch{
			Name: "master",
			Repo: &pbe.Repository{
				Name:  mockedPRRepoName,
				Owner: mockedPRRepoOwner,
				Uri:   fmt.Sprintf("https://github.com/%s/%s", mockedPRRepoOwner, mockedPRRepoName),
			},
		},
		Head: &pbe.Branch{
			Name: "new-topic",
			Sha:  mockedHeadSHA,
			Repo: &pbe.Repository{
				Name:  mockedPRRepoName,
				Owner: mockedPRRepoOwner,
				Uri:   fmt.Sprintf("https://github.com/%s/%s", mockedPRRepoName, mockedPRRepoOwner),
			},
		},
		RequestedReviewers: &pbe.RequestedReviewers{
			Users: []*pbe.ExternalUser{},
			Teams: []*pbe.Team{},
		},
		Files: []*pbe.CommitFile{
			{
				Sha:      mockedBlobId,
				Filename: mockedPatchFileRelativeName,
				Patch:    mockedPatch,
			},
		},
	}

	failMessage := "DownloadContents"

	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mock.WriteError(
						w,
						http.StatusInternalServerError,
						failMessage,
					)
				}),
			),
		},
		nil,
		codeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	gotVal, gotErr := GetSymbolsFromPatch(mockedEnv)

	assert.Nil(t, gotVal)
	assert.Equal(t, gotErr.(*github.ErrorResponse).Message, failMessage)
}

func TestGetSymbolsFromPatch_WhenSemanticServiceNotFound(t *testing.T) {
	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
	mockedHeadSHA := "abc123"

	mockedPatchFilePath := "test"
	mockedPatchFileName := "crawler.go"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
	mockedPatchLocation := fmt.Sprintf("/%v/%v/%v", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

	// Since the blob and the patch are simply passed around, they can be an empty string
	mockedBlob := ""
	mockedPatch := ""
	mockedBlobId := "1234"

	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Head: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "new-topic",
			Sha:  mockedHeadSHA,
		},
		Base: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "master",
		},
		Files: []*pbe.CommitFile{
			{
				Sha:      mockedBlobId,
				Filename: mockedPatchFileRelativeName,
				Patch:    mockedPatch,
			},
		},
	})

	failMessage := "semantic service not found"

	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
						{
							Name:        github.String(mockedPatchFileName),
							Path:        github.String(mockedPatchFilePath),
							DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%v", mockedPatchLocation)),
						},
					}))
				}),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: mockedPatchLocation,
					Method:  "GET",
				},
				mockedBlob,
			),
		},
		nil,
		mockedCodeReview,
		aladino.MockBuiltIns(),
		nil,
	)

	gotVal, gotErr := GetSymbolsFromPatch(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, gotErr, failMessage)
}

func TestGetSymbolsFromPatch_WhenGetSymbolsRequestFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedSemanticClient := services_mocks.NewMockSemanticClient(controller)

	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
	mockedHeadSHA := "abc123"

	mockedPatchFilePath := "test"
	mockedPatchFileName := "crawler.go"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
	mockedPatchLocation := fmt.Sprintf("/%v/%v/%v", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

	// Since the blob and the patch are simply passed around, they can be an empty string
	mockedBlob := ""
	mockedPatch := ""
	mockedBlobId := "1234"

	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Head: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "new-topic",
			Sha:  mockedHeadSHA,
		},
		Base: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "master",
		},
		Files: []*pbe.CommitFile{
			{
				Sha:      mockedBlobId,
				Filename: mockedPatchFileRelativeName,
				Patch:    mockedPatch,
			},
		},
	})

	failMessage := "GetSymbolsRequest"

	mockedSemanticClient.EXPECT().GetSymbols(
		gomock.Any(),
		&services.GetSymbolsRequest{
			Uri:      mockedPRUrl,
			CommitId: mockedHeadSHA,
			Filepath: mockedPatchFileRelativeName,
			Blob:     []byte(fmt.Sprintf("%#v", mockedBlob)),
			BlobId:   mockedBlobId,
			Diff:     &entities.ResolveFileDiff{Blocks: []*entities.ResolveBlockDiff{}}}).Return(
		nil, fmt.Errorf(failMessage),
	)

	mockBuiltIns := &aladino.BuiltIns{
		Services: map[string]interface{}{
			plugins_aladino_services.SEMANTIC_SERVICE_KEY: mockedSemanticClient,
		},
	}

	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
						{
							Name:        github.String(mockedPatchFileName),
							Path:        github.String(mockedPatchFilePath),
							DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%v", mockedPatchLocation)),
						},
					}))
				}),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: mockedPatchLocation,
					Method:  "GET",
				},
				mockedBlob,
			),
		},
		nil,
		mockedCodeReview,
		mockBuiltIns,
		nil,
	)
	gotVal, err := GetSymbolsFromPatch(mockedEnv)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, failMessage)
}

func TestGetSymbolsFromPatch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedSemanticClient := services_mocks.NewMockSemanticClient(controller)

	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
	mockedHeadSHA := "abc123"

	mockedPatchFilePath := "test"
	mockedPatchFileName := "crawler.go"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
	mockedPatchLocation := fmt.Sprintf("/%v/%v/%v", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

	// Since the blob and the patch are simply passed around, they can be an empty string
	mockedBlob := ""
	mockedPatch := ""
	mockedBlobId := "1234"

	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Head: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "new-topic",
			Sha:  mockedHeadSHA,
		},
		Base: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "master",
		},
		Files: []*pbe.CommitFile{
			{
				Sha:      mockedBlobId,
				Filename: mockedPatchFileRelativeName,
				Patch:    mockedPatch,
			},
		},
	})

	mockedSymbols := &entities.Symbols{
		Files: map[string]*entities.File{},
		Symbols: map[string]*entities.Symbol{
			"Crawl": {
				Id:   "1",
				Name: "Crawl",
				Type: "Function",
				CodeComments: []*entities.SymbolDocumentation{
					{
						Code: "reviewpad-an: critical",
					},
				},
			},
		},
	}

	mockedSemanticClient.EXPECT().GetSymbols(
		gomock.Any(),
		&services.GetSymbolsRequest{
			Uri:      mockedPRUrl,
			CommitId: mockedHeadSHA,
			Filepath: mockedPatchFileRelativeName,
			Blob:     []byte(fmt.Sprintf("%#v", mockedBlob)),
			BlobId:   mockedBlobId,
			Diff:     &entities.ResolveFileDiff{Blocks: []*entities.ResolveBlockDiff{}}}).Return(
		&services.GetSymbolsReply{
			Symbols: mockedSymbols,
		}, nil,
	)

	mockBuiltIns := &aladino.BuiltIns{
		Services: map[string]interface{}{
			plugins_aladino_services.SEMANTIC_SERVICE_KEY: mockedSemanticClient,
		},
	}

	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
						{
							Name:        github.String(mockedPatchFileName),
							Path:        github.String(mockedPatchFilePath),
							DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%v", mockedPatchLocation)),
						},
					}))
				}),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: mockedPatchLocation,
					Method:  "GET",
				},
				mockedBlob,
			),
		},
		nil,
		mockedCodeReview,
		mockBuiltIns,
		nil,
	)
	gotVal, err := GetSymbolsFromPatch(mockedEnv)

	wantVal := map[string]*entities.Symbols{
		mockedPatchFileRelativeName: mockedSymbols,
	}

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestGetSymbolsFromFileInBranch_WhenGetSymbolsRequestFails(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedSemanticClient := services_mocks.NewMockSemanticClient(controller)

	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
	mockedHeadSHA := "abc123"

	mockedPullRequest := &github.PullRequest{
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPRRepoOwner),
				},
				URL:  github.String(mockedPRUrl),
				Name: github.String(mockedPRRepoName),
			},
			Ref: github.String("new-topic"),
			SHA: github.String(mockedHeadSHA),
		},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPRRepoOwner),
				},
				URL:  github.String(mockedPRUrl),
				Name: github.String(mockedPRRepoName),
			},
			Ref: github.String("master"),
		},
	}

	mockedPatchFilePath := "test"
	mockedPatchFileName := "crawler.go"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
	mockedPatchLocation := fmt.Sprintf("/%v/%v/%v", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

	// Since the blob and the patch are simply passed around, they can be an empty string
	mockedBlob := ""
	mockedPatch := ""
	mockedBlobId := "1234"

	mockedPullRequestFileList := []*github.CommitFile{
		{
			SHA:      github.String(mockedBlobId),
			Filename: github.String(mockedPatchFileRelativeName),
			Patch:    github.String(mockedPatch),
		},
	}

	failMessage := "GetSymbolsFromFileInBranch"

	mockedSemanticClient.EXPECT().GetSymbols(
		gomock.Any(),
		&services.GetSymbolsRequest{
			Uri:      mockedPRUrl,
			CommitId: mockedHeadSHA,
			Filepath: mockedPatchFileRelativeName,
			Blob:     []byte(fmt.Sprintf("%#v", mockedBlob)),
			BlobId:   mockedBlobId,
			Diff:     &entities.ResolveFileDiff{Blocks: []*entities.ResolveBlockDiff{}}}).Return(
		nil, fmt.Errorf(failMessage),
	)

	mockBuiltIns := &aladino.BuiltIns{
		Services: map[string]interface{}{
			plugins_aladino_services.SEMANTIC_SERVICE_KEY: mockedSemanticClient,
		},
	}

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
						{
							Name:        github.String(mockedPatchFileName),
							Path:        github.String(mockedPatchFilePath),
							DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%v", mockedPatchLocation)),
						},
					}))
				}),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: mockedPatchLocation,
					Method:  "GET",
				},
				mockedBlob,
			),
		},
		nil,
		mockBuiltIns,
		nil,
	)

	githubMockedCommitFile := &pbe.CommitFile{
		Sha:      mockedBlobId,
		Filename: mockedPatchFileRelativeName,
		Patch:    mockedPatch,
	}

	mockedCommitFile := &codehost.File{
		Repr: githubMockedCommitFile,
	}

	mockedBranch := &pbe.Branch{
		Repo: &pbe.Repository{
			Owner: mockedPRRepoOwner,
			Uri:   mockedPRUrl,
			Name:  mockedPRRepoName,
		},
		Name: "new-topic",
		Sha:  mockedHeadSHA,
	}

	gotVal, _, err := GetSymbolsFromFileInBranch(mockedEnv, mockedCommitFile, mockedBranch)

	assert.Nil(t, gotVal)
	assert.EqualError(t, err, failMessage)
}

func TestGetSymbolsFromFileInBranch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedSemanticClient := services_mocks.NewMockSemanticClient(controller)

	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
	mockedHeadSHA := "abc123"

	mockedPullRequest := &github.PullRequest{
		Head: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPRRepoOwner),
				},
				URL:  github.String(mockedPRUrl),
				Name: github.String(mockedPRRepoName),
			},
			Ref: github.String("new-topic"),
			SHA: github.String(mockedHeadSHA),
		},
		Base: &github.PullRequestBranch{
			Repo: &github.Repository{
				Owner: &github.User{
					Login: github.String(mockedPRRepoOwner),
				},
				URL:  github.String(mockedPRUrl),
				Name: github.String(mockedPRRepoName),
			},
			Ref: github.String("master"),
		},
	}

	mockedPatchFilePath := "test"
	mockedPatchFileName := "crawler.go"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
	mockedPatchLocation := fmt.Sprintf("/%v/%v/%v", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

	// Since the blob and the patch are simply passed around, they can be an empty string
	mockedBlob := ""
	mockedPatch := ""
	mockedBlobId := "1234"

	mockedPullRequestFileList := []*github.CommitFile{
		{
			SHA:      github.String(mockedBlobId),
			Filename: github.String(mockedPatchFileRelativeName),
			Patch:    github.String(mockedPatch),
		},
	}

	mockedSymbols := &entities.Symbols{
		Files: map[string]*entities.File{},
		Symbols: map[string]*entities.Symbol{
			"Crawl": {
				Id:   "1",
				Name: "Crawl",
				Type: "Function",
				CodeComments: []*entities.SymbolDocumentation{
					{
						Code: "reviewpad-an: critical",
					},
				},
			},
		},
	}

	mockedSemanticClient.EXPECT().GetSymbols(
		gomock.Any(),
		&services.GetSymbolsRequest{
			Uri:      mockedPRUrl,
			CommitId: mockedHeadSHA,
			Filepath: mockedPatchFileRelativeName,
			Blob:     []byte(fmt.Sprintf("%#v", mockedBlob)),
			BlobId:   mockedBlobId,
			Diff:     &entities.ResolveFileDiff{Blocks: []*entities.ResolveBlockDiff{}}}).Return(
		&services.GetSymbolsReply{
			Symbols: mockedSymbols,
		}, nil,
	)

	mockBuiltIns := &aladino.BuiltIns{
		Services: map[string]interface{}{
			plugins_aladino_services.SEMANTIC_SERVICE_KEY: mockedSemanticClient,
		},
	}

	mockedEnv := aladino.MockDefaultEnv(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposPullsByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequest))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetReposPullsFilesByOwnerByRepoByPullNumber,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal(mockedPullRequestFileList))
				}),
			),
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
						{
							Name:        github.String(mockedPatchFileName),
							Path:        github.String(mockedPatchFilePath),
							DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%v", mockedPatchLocation)),
						},
					}))
				}),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: mockedPatchLocation,
					Method:  "GET",
				},
				mockedBlob,
			),
		},
		nil,
		mockBuiltIns,
		nil,
	)
	githubMockedCommitFile := &pbe.CommitFile{
		Sha:      mockedBlobId,
		Filename: mockedPatchFileRelativeName,
		Patch:    mockedPatch,
	}

	mockedCommitFile := &codehost.File{
		Repr: githubMockedCommitFile,
	}

	mockedBranch := &pbe.Branch{
		Repo: &pbe.Repository{
			Owner: mockedPRRepoOwner,
			Uri:   mockedPRUrl,
			Name:  mockedPRRepoName,
		},
		Name: "new-topic",
		Sha:  mockedHeadSHA,
	}

	gotVal, _, err := GetSymbolsFromFileInBranch(mockedEnv, mockedCommitFile, mockedBranch)

	wantVal := mockedSymbols

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}

func TestGetSymbolsFromHeadByPatch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockedSemanticClient := services_mocks.NewMockSemanticClient(controller)

	mockedPRRepoOwner := "mock-reviewpad"
	mockedPRRepoName := "test"
	mockedPRUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/6", mockedPRRepoOwner, mockedPRRepoName)
	mockedHeadSHA := "abc123"

	mockedPatchFilePath := "test"
	mockedPatchFileName := "crawler.go"
	mockedPatchFileRelativeName := fmt.Sprintf("%v/crawler.go", mockedPatchFilePath)
	mockedPatchLocation := fmt.Sprintf("/%v/%v/%v", mockedPRRepoOwner, mockedPRRepoName, mockedPatchFileName)

	// Since the blob and the patch are simply passed around, they can be an empty string
	mockedBlob := ""
	mockedPatch := ""
	mockedBlobId := "1234"

	mockedCodeReview := aladino.GetDefaultMockCodeReviewDetailsWith(&pbe.CodeReview{
		Head: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "new-topic",
			Sha:  mockedHeadSHA,
		},
		Base: &pbe.Branch{
			Repo: &pbe.Repository{
				Owner: mockedPRRepoOwner,
				Uri:   mockedPRUrl,
				Name:  mockedPRRepoName,
			},
			Name: "master",
		},
		Files: []*pbe.CommitFile{
			{
				Sha:      mockedBlobId,
				Filename: mockedPatchFileRelativeName,
				Patch:    mockedPatch,
			},
		},
	})

	mockedSymbols := &entities.Symbols{
		Files: map[string]*entities.File{},
		Symbols: map[string]*entities.Symbol{
			"Crawl": {
				Id:   "1",
				Name: "Crawl",
				Type: "Function",
				CodeComments: []*entities.SymbolDocumentation{
					{
						Code: "reviewpad-an: critical",
					},
				},
			},
		},
	}

	mockedSemanticClient.EXPECT().GetSymbols(
		gomock.Any(),
		&services.GetSymbolsRequest{
			Uri:      mockedPRUrl,
			CommitId: mockedHeadSHA,
			Filepath: mockedPatchFileRelativeName,
			Blob:     []byte(fmt.Sprintf("%#v", mockedBlob)),
			BlobId:   mockedBlobId,
			Diff:     &entities.ResolveFileDiff{Blocks: []*entities.ResolveBlockDiff{}}}).Return(
		&services.GetSymbolsReply{
			Symbols: mockedSymbols,
		}, nil,
	)

	mockBuiltIns := &aladino.BuiltIns{
		Services: map[string]interface{}{
			plugins_aladino_services.SEMANTIC_SERVICE_KEY: mockedSemanticClient,
		},
	}

	mockedEnv := aladino.MockDefaultEnvWithCodeReview(
		t,
		[]mock.MockBackendOption{
			mock.WithRequestMatchHandler(
				mock.GetReposContentsByOwnerByRepoByPath,
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					utils.MustWriteBytes(w, mock.MustMarshal([]github.RepositoryContent{
						{
							Name:        github.String(mockedPatchFileName),
							Path:        github.String(mockedPatchFilePath),
							DownloadURL: github.String(fmt.Sprintf("https://raw.githubusercontent.com/%v", mockedPatchLocation)),
						},
					}))
				}),
			),
			mock.WithRequestMatch(
				mock.EndpointPattern{
					Pattern: mockedPatchLocation,
					Method:  "GET",
				},
				mockedBlob,
			),
		},
		nil,
		mockedCodeReview,
		mockBuiltIns,
		nil,
	)

	mPullRequest := mockedEnv.GetTarget().(*target.PullRequestTarget)

	gotVal, _, err := GetSymbolsFromHeadByPatch(mockedEnv, mPullRequest.Patch)

	wantVal := mockedSymbols

	assert.Nil(t, err)
	assert.Equal(t, wantVal, gotVal)
}
