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
                "-github-token=_GIT_HUB_TOKEN_",
            ],
            "program": "${workspaceFolder}/cmd/cli/main.go"
        }
    ]
}
```

## License

Reviewpad is available under the GNU Lesser General Public License v3.0 license. See LICENSE for the full license text.