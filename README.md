# Reviewpad

Reviewpad is an open-source tool to automate pull requests workflow.

## Development

### Debugging on VSCode

Add the following to your `.vscode/launch.json`.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch GH Action",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "env": {
                // Repository with the format "organization/repository"
                // e.g. "nodejs/node"
                "INPUT_REPOSITORY": "_REPOSITORY_",
                // GitHub personal access token
                // https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
                "INPUT_TOKEN": "_GITHUB_TOKEN_",
                // Pull request number
                "INPUT_PRNUMBER": "_PR_NUMBER_",
                // "true" to run on dry run, "false" otherwise
                // By default is "false"
                "INPUT_DRYRUN": "_DRY_RUN_"
            },
            "program": "${workspaceFolder}/cmd/github-action/main.go"
        },
        {
            "name": "Launch CLI",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "args": [
                // Flag to run on dry run
                "-dry-run",
                // Absolute path to reviewpad.yml file to run
                "-reviewpad=_PATH_TO_REVIEWPAD_FILE_",
                // GitHub url to run the reviewpad.yml against to
                // e.g. https://github.com/reviewpad/action-demo/pull/1
                "-pull-request=_PULL_REQUEST_URL_",
                // GiHub personal access token
                // https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
                "-token=_GIT_HUB_TOKEN_",
            ],
            "program": "${workspaceFolder}/cmd/cli/main.go"
        }
    ]
}
```

## License

Reviewpad is available under the GNU Lesser General Public License v3.0 license. See LICENSE for the full license text.