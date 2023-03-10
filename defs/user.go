// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package defs

import "github.com/google/go-github/v49/github"

type User struct {
	ID    string
	Login string
}

type Users []User

func FromGithubUsers(githubUsers []*github.User) Users {
	users := make(Users, len(githubUsers))
	for i, githubUser := range githubUsers {
		users[i] = FromGithubUser(githubUser)
	}
	return users
}

func FromGithubUser(githubUser *github.User) User {
	return User{
		ID:    githubUser.GetNodeID(),
		Login: githubUser.GetLogin(),
	}
}
