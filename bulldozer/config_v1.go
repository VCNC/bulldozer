// Copyright 2018 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bulldozer

import (
	"strings"
	"regexp"
)

type MessageStrategy string
type TitleStrategy string
type MergeMethod string

const (
	PullRequestBody  MessageStrategy = "pull_request_body"
	SummarizeCommits MessageStrategy = "summarize_commits"
	EmptyBody        MessageStrategy = "empty_body"

	PullRequestTitle   TitleStrategy = "pull_request_title"
	FirstCommitTitle   TitleStrategy = "first_commit_title"
	GithubDefaultTitle TitleStrategy = "github_default"

	MergeCommit    MergeMethod = "merge"
	SquashAndMerge MergeMethod = "squash"
	RebaseAndMerge MergeMethod = "rebase"
)

type MergeConfig struct {
	Whitelist Signals `yaml:"whitelist"`
	Blacklist Signals `yaml:"blacklist"`

	DeleteAfterMerge bool `yaml:"delete_after_merge"`

	Method  MergeMethod  `yaml:"method"`
	Options MergeOptions `yaml:"options"`

	BranchMethod map[string]MergeMethod `yaml:"branch_method"`

	// Additional status checks that bulldozer should require
	// (even if the branch protection settings doesn't require it)
	RequiredStatuses []string `yaml:"required_statuses"`
}

func (m MergeConfig) branchMethod(branch string) (method MergeMethod, ok bool) {
	for key, val := range m.BranchMethod {
		regexBranch := strings.ReplaceAll(key, "*", "(.*)")
		matcher, err := regexp.Compile(regexBranch)
		if err != nil {
			return "", false
		}
		
		if len(matcher.FindString(branch)) > 0 {
			return val, true
		}
	}
	return "", false
}

type MergeOptions struct {
	Squash *SquashOptions `yaml:"squash"`
}

type SquashOptions struct {
	Title            TitleStrategy   `yaml:"title"`
	Body             MessageStrategy `yaml:"body"`
	MessageDelimiter string          `yaml:"message_delimiter"`
}

type UpdateConfig struct {
	Whitelist Signals `yaml:"whitelist"`
	Blacklist Signals `yaml:"blacklist"`
}

type Config struct {
	Version int `yaml:"version"`

	Merge  MergeConfig  `yaml:"merge"`
	Update UpdateConfig `yaml:"update"`
}
