![Reviewpad-Background-Logo-Shorter@1 5x](https://user-images.githubusercontent.com/38539/185980242-be2dfd87-2e0c-4bb1-87ce-462f629cedf6.png)

# Reviewpad

[![x-ray-badge](https://img.shields.io/badge/Time%20to%20Merge-Fair%20team-bb3e03?link=https://xray.reviewpad.com/analysis?repository=https%3A%2F%2Fgithub.com%2Freviewpad%2Freviewpad&style=plastic.svg)](https://xray.reviewpad.com/analysis?repository=https%3A%2F%2Fgithub.com%2Freviewpad%2Freviewpad) [![Build](https://github.com/reviewpad/reviewpad/actions/workflows/build.yml/badge.svg)](https://github.com/reviewpad/reviewpad/actions/workflows/build.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/reviewpad/reviewpad)](https://goreportcard.com/report/github.com/reviewpad/reviewpad)

## Welcome to Reviewpad!

For **questions**, check out the GitHub [discussions](https://github.com/reviewpad/reviewpad/discussions).

For **documentation**, check out this document and the [official documentation](https://docs.reviewpad.com).

**Join our Community on [Discord](https://reviewpad.com/discord)!**

To start using Reviewpad, check out our [website](https://reviewpad.com).

## What is Reviewpad?

Reviewpad is a service to automate pull requests and issues workflows.

The workflows are specified in a YML-based configuration language described in the [official documentation](https://docs.reviewpad.com/guides/syntax).

In Reviewpad, you can automate actions over the pull requests and issues such as:

1. Automated [comments](https://docs.reviewpad.com/use-cases/comment-on-pull-requests);
2. Add or remove [labels](https://docs.reviewpad.com/use-cases/automated-labelling/);
3. Specify [reviewer assignments](https://docs.reviewpad.com/use-cases/reviewer-assignment);
4. Automate [close/merge actions](https://docs.reviewpad.com/use-cases/auto-merge);
5. [Block the merge action](https://docs.reviewpad.com/use-cases/blocking-merge);
6. Automatically [summarize pull requests](https://docs.reviewpad.com/guides/built-ins#summarize);

As an example, the following workflow:

```yml
labels:
    ship:
        description: Ship mode
        color: "#76dbbe"

workflows:
    - name: ship
      description: Ship process - bypass the review and merge with rebase
      run:
          - if: $hasFileExtensions([".md"])
            then:
                - $addLabel("ship")
                - $merge()
```

Automatically adds the label `ship` and merges all pull requests that only change `.md` files.

You can execute Reviewpad through the CLI or install Reviewpad [GitHub App](https://github.com/marketplace/reviewpad).

## Architecture

This repository generates two artifacts:

1. CLI [cli](cli/main.go) that runs Reviewpad open source edition.
2. Reviewpad library packages:
    - github.com/reviewpad/reviewpad/collector
    - github.com/reviewpad/reviewpad/engine
    - github.com/reviewpad/reviewpad/lang/aladino
    - github.com/reviewpad/reviewpad/plugins/aladino
    - github.com/reviewpad/reviewpad/utils

Conceptually, the packages are divided into four categories:

1. Engine: The engine is the package responsible for processing the YML file. This process is divided into two stages:
    - Process the YML file to determine which workflows are enabled. The outcome of this phase is a program with the actions that will be executed over the pull request.
    - Execution of the synthesised program.
2. [Aladino Language](https://docs.reviewpad.com/guides/aladino/specification): This is the language that is used in the `spec` property of the rules and also the actions of the workflows. The engine of Reviewpad is not specific to Aladino - this means that it is possible to add support for a new language such as `Javascript` or `Golang` in these specifications.
3. Plugins: The plugin package contains the built-in functions and actions that act as an abstraction to the 3rd party services such as GitHub, Jira, etc. This package is specific to each supported specification language. In the case of `plugins/aladino`, it contains the implementations of the [built-ins](https://docs.reviewpad.com/guides/built-ins).
4. Utilities: packages, such as the collector, that provide utilities that are used in multiple places.

## Development

### Prerequisites

Before you begin, ensure you have met the following requirements:

-   [Go](https://golang.org/doc/install) with the minimum version of 1.16.
-   [goyacc](https://pkg.go.dev/golang.org/x/tools/cmd/goyacc) used to generate Reviewpad Aladino parser (`go install golang.org/x/tools/cmd/goyacc@master`).
-   [libgit2](https://github.com/libgit2/libgit2) with version v1.2.
-   To run the tests, Reviewpad requires the following environment variables:
     - `INPUT_SEMANTIC_SERVICE`. You can do this by running the following command in your terminal: `export INPUT_SEMANTIC_SERVICE="0.0.0.0:3006"`.
     - `INPUT_ROBIN_SERVICE`. You can do this by running the following command in your terminal: `export INPUT_ROBIN_SERVICE="0.0.0.0:3011"`.
     - `INPUT_CODEHOST_SERVICE`. You can do this by running the following command in your terminal: `export INPUT_CODEHOST_SERVICE="0.0.0.0:3004"`.

### Compilation

We use [Taskfile](https://taskfile.dev/). To compile the packages simply run:

```sh
task build
```

To generate the CLI run:

```sh
task build-cli
```

This command generates the binary `reviewpad-cli` which you can run to try Reviewpad directly.

The CLI has the following argument list:

```sh
./reviewpad-cli
reviewpad-cli is command line interface to run reviewpad commands.

Usage:
  reviewpad-cli [command]

Available Commands:
  check       Check if input reviewpad file is valid
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Runs reviewpad

Flags:
  -f, --file string   input reviewpad file
  -h, --help          help for reviewpad-cli

Use "reviewpad-cli [command] --help" for more information about a command.
```

### Running unit tests

Run the tests with:

```
task test
```

If you get the error:

```
panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: too many open files [recovered]
        panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: socket: too many open files
```

You can solve with:

```
ulimit -Sn 500
```

#### Coverage

To generate the coverage report run:

```
task test
```

To display the code coverage for every package run:

```
go tool cover -func coverage.out
```

To display the total code coverage percentage run:

```
go tool cover -func coverage.out | grep total:
```

### Running integration tests

The integration tests run Reviewpad on an actual repository and pull request. The repository where the integration tests run requires the following setup:

- At least one milestone;
- At least 3 labels named `bug`, `documentation`, `wontfix` (GitHub adds these labels to every new repository by default);
- A team called `integration-test` with at least 3 members;
- A project called `[INTEGRATION TESTS] Reviewpad` with `Todo` and `In Progress` status
- A GitHub status check called `log event`.

#### Required Environment Variables

- `GITHUB_INTEGRATION_TEST_TOKEN` : GitHub access token used to setup tests and run Reviewpad
- `GITHUB_INTEGRATION_TEST_REPO_OWNER` : The owner of the repository used to run integration tests on
- `GITHUB_INTEGRATION_TEST_REPO_NAME` : The name of the repository used to run integration tests on

After setting these variables, you can run the integration tests with:

```
task integration-test
```

### VSCode

We strongly recommend using [VSCode](https://code.visualstudio.com/) with the following extensions:

-   [Go](https://marketplace.visualstudio.com/items?itemName=golang.go) for Go language support.
-   [EditorConfig](https://marketplace.visualstudio.com/items?itemName=EditorConfig.EditorConfig) to helps to maintain consistent coding styles.
-   [licenser](https://marketplace.visualstudio.com/items?itemName=ymotongpoo.licenser) for adding license headers.
-   [YAML](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml) for enabling `reviewpad.yml` JSON schema.
-   [JSON Parse & Stringify](https://marketplace.visualstudio.com/items?itemName=nextfaze.json-parse-stringify) for formatting JSON files.

Open the project in VSCode, open the command palette (Ctrl+Shift+P) and search for `Preferences: Open Workspace Settings (JSON)`.

Paste the following configuration:

```json
{
    // Licenser configuration
    "licenser.license": "Custom",
    "licenser.author": "Explore.dev, Unipessoal Lda",
    "licenser.customHeader": "Copyright (C) @YEAR@ @AUTHOR@ - All Rights Reserved\nUse of this source code is governed by a license that can be\nfound in the LICENSE file.",
    // Associate Reviewpad schema to reviewpad.(yml|yaml) files
    "yaml.schemas": {
        "https://raw.githubusercontent.com/reviewpad/schemas/main/latest/schema.json": [
            "reviewpad.yml",
        ]
    },
    // Go configuration
    "go.testFlags": [
        "integration"
    ],
    "go.buildFlags": [
        "-tags=integration"
    ],
}
```

### Debugging with VSCode

Add the following to your `.vscode/launch.json`.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch CLI",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "args": [
                "run",
                // Flag to run on dry run (optional. Default is false)
                "-d",
                // Flag to run on safe mode (optional. Default is false)
                "-s",
                // Flag to define the log level (optional. Default is debug)
                "-l=debug",
                // GiHub personal access token.
                // To generate a token, go to https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token.
                "-t=<GIT_HUB_TOKEN>",
                // Absolute path to reviewpad.yml file to run.
                "-f=<PATH_TO_REVIEWPAD_FILE>",
                // Absolute path to JSON file with GitHub event.
                // This GitHub event defines the action that will run on the pull request / issue.
                "-e=<PATH_TO_EVENT_JSON>"
                // GitHub URL of the pull request / issue.
                // The reviewpad.yml file provided with the -f flag will run against this pull request / issue.
                // This URL should be same the URL of the pull request / issues in the GitHub event provided with the -e flag.
                // For instance, if you are using a GitHub event of a comment on the pull request X, the -u flag should be the URL of the pull request X.
                "-u=<GITHUB_URL>",
            ],
            "env": {
                "INPUT_CODEHOST_SERVICE": "<CODEHOST_SERVICE_ENDPOINT>",
                "INPUT_SEMANTIC_SERVICE": "<SEMANTIC_SERVICE_ENDPOINT>",
                "INPUT_ROBIN_SERVICE": "<ROBIN_SERVICE_ENDPOINT>",
            }
            "program": "${workspaceFolder}/cli/main.go"
        }
    ]
}
```

#### How to get the GitHub event (internal use only)

The `-e` flag is mandatory to run the debugger.

It represents the GitHub event that you wish to run the reviewpad.yml file against.

To extract a GitHub event, you can use the following steps:

1. Navigate to the logs of the Reviewpad GitHub App run where you wish to copy the event from. You can use the following query to filter the events payload `{$.level=debug && $.msg="request received" && $.request_received.delivery_id=<DELIVERY_ID> }`.
2. Copy the content inside the property `body`.
3. Paste the content inside a file (e.g. `my_event.json`) and save it under `cli > debugdata`.
4. This content is an escape JSON string. Use the [JSON Parse & Stringify](https://marketplace.visualstudio.com/items?itemName=nextfaze.json-parse-stringify) extension to parse the content by pressing `Ctrl+Shift+P` and searching for `JSON: Parse Stringified JSON`.
5. Rename the root property `eventType` to `event_name`.
6. Rename the root property `eventPayload` to `event`.
7. Update the argument `-e` to point to the full path of the file you just created.

You can then run the debugger by pressing F5.

### Contributing

We welcome contributions to Reviewpad from the community!

See the [Contributing Guide](CONTRIBUTING.md).

If you need any assistance, please join [discord](https://reviewpad.com/discord) to reach the core contributors.

## License

Reviewpad is available under the GNU Lesser General Public License v3.0 license.

See [LICENSE](./LICENSE.md) for the full license text.

## Compliance

<a href="https://drata.com/soc-2?utm_source=soc2badge" target="_blank"><img src="https://cdn.drata.com/badge/soc2-light.png" height="120" /></a>

