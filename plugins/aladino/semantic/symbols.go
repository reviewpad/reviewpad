// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_semantic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reviewpad/api/go/entities"
	pbe "github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v4/plugins/aladino/services"
)

func GetSymbolsFromPatch(e aladino.Env) (map[string]*entities.Symbols, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	res := make(map[string]*entities.Symbols)

	head := pullRequest.CodeReview.GetHead()
	base := pullRequest.CodeReview.GetBase()
	url := head.Repo.Uri
	patch := pullRequest.Patch

	lastCommit := head.Sha

	for fp, commitFile := range patch {
		blob, err := e.GetGithubClient().DownloadContents(e.GetCtx(), fp, head)
		if err != nil {
			// If fails to download the file from head, then tries to download it from base.
			// This happens when a file has been removed.
			blob, err = e.GetGithubClient().DownloadContents(e.GetCtx(), fp, base)

			if err != nil {
				return nil, err
			}
		}

		// Convert the bytes array to string, remove the escape and convert again to bytes
		str, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(blob)), `\\u`, `\u`, -1))
		blob = []byte(str)

		blocks := getBlocks(commitFile)

		service, ok := e.GetBuiltIns().Services[plugins_aladino_services.SEMANTIC_SERVICE_KEY]
		if !ok {
			return nil, fmt.Errorf("semantic service not found")
		}

		semanticClient := service.(api.SemanticClient)
		req := &api.GetSymbolsRequest{
			Uri:      url,
			CommitId: lastCommit,
			Filepath: fp,
			Blob:     blob,
			BlobId:   commitFile.Repr.Sha,
			Diff:     &entities.ResolveFileDiff{Blocks: blocks},
		}
		reply, err := semanticClient.GetSymbols(e.GetCtx(), req)
		if err != nil {
			return nil, err
		}

		res[fp] = reply.Symbols
	}

	return res, nil
}

func GetSymbolsFromHeadByPatch(e aladino.Env, patch target.Patch) (*entities.Symbols, map[string]string, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	head := pullRequest.CodeReview.GetHead()

	res := &entities.Symbols{
		Files:   make(map[string]*entities.File),
		Symbols: make(map[string]*entities.Symbol),
	}

	files := make(map[string]string)

	for fp, commitFile := range patch {
		if commitFile.Repr.Status == pbe.CommitFile_REMOVED {
			// in this case, the file is not in the head branch
			continue
		}

		symbols, contents, err := GetSymbolsFromFileInBranch(e, commitFile, head)
		if err != nil {
			return nil, nil, err
		}

		files[fp] = contents
		joinSymbols(res, symbols)
	}

	return res, files, nil
}

func GetSymbolsFromBaseByPatch(e aladino.Env, patch target.Patch) (*entities.Symbols, map[string]string, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	base := pullRequest.CodeReview.GetBase()

	res := &entities.Symbols{
		Files:   make(map[string]*entities.File),
		Symbols: make(map[string]*entities.Symbol),
	}

	files := make(map[string]string)

	for fp, commitFile := range patch {
		if commitFile.Repr.GetStatus() == pbe.CommitFile_ADDED {
			// in this case, the file is not in the base branch
			continue
		}

		symbols, contents, err := GetSymbolsFromFileInBranch(e, commitFile, base)
		if err != nil {
			return nil, nil, err
		}

		files[fp] = contents
		joinSymbols(res, symbols)
	}

	return res, files, nil
}

func joinSymbols(current *entities.Symbols, new *entities.Symbols) {
	for path, file := range new.Files {
		current.Files[path] = file
	}

	for symbol, gotoSymbol := range new.Symbols {
		current.Symbols[symbol] = gotoSymbol
	}
}

func GetSymbolsFromFileInBranch(e aladino.Env, commitFile *codehost.File, branch *pbe.Branch) (*entities.Symbols, string, error) {
	fp := commitFile.Repr.GetFilename()

	blob, err := e.GetGithubClient().DownloadContents(e.GetCtx(), fp, branch)
	if err != nil {
		return nil, "", err
	}

	// Convert the bytes array to string, remove the escape and convert again to bytes
	str, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(blob)), `\\u`, `\u`, -1))
	blob = []byte(str)

	blocks := getBlocks(commitFile)

	service, ok := e.GetBuiltIns().Services[plugins_aladino_services.SEMANTIC_SERVICE_KEY]
	if !ok {
		return nil, "", fmt.Errorf("semantic service not found")
	}

	semanticClient := service.(api.SemanticClient)
	req := &api.GetSymbolsRequest{
		Uri:      branch.Repo.Uri,
		CommitId: branch.Sha,
		Filepath: fp,
		Blob:     blob,
		BlobId:   commitFile.Repr.Sha,
		Diff:     &entities.ResolveFileDiff{Blocks: blocks},
	}
	reply, err := semanticClient.GetSymbols(e.GetCtx(), req)
	if err != nil {
		return nil, "", err
	}

	return reply.GetSymbols(), str, nil
}

func getBlocks(commitFile *codehost.File) []*entities.ResolveBlockDiff {
	res := make([]*entities.ResolveBlockDiff, len(commitFile.Diff))
	for i, block := range commitFile.Diff {
		if block.New != nil {
			res[i] = &entities.ResolveBlockDiff{
				Start: block.New.Start,
				End:   block.New.End,
			}
		} else {
			// invariant: block.Old != nil
			res[i] = &entities.ResolveBlockDiff{
				Start: block.Old.Start,
				End:   block.Old.End,
			}
		}
	}

	return res
}
