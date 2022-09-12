# How to contribute

Reviewpad is an open source project to empower teams to automate their development process - issues and pull requests. We believe that teams should be able to contribute to a project the the most fast and efficient way. We are very excited to see your interest on how to help us make Reviewpad the best automation tool.

### Code of Conduct

Reviewpad projects have a [Code of Conduct](https://github.com/reviewpad/.github/blob/main/docs/CODE_OF_CONDUCT.md) to which all contributors must adhere.
Please read it before interacting with the repository or the community in any way.

### Open Development

All work on Reviewpad happens directly on [GitHub](https://github.com/reviewpad/reviewpad). Both core team members and external contributors send pull requests which go through the same review process.

### Semantic Versioning

Reviewpad follows [semantic versioning](https://semver.org/). We release patch versions for critical bugfixes, minor versions for new features or non-essential changes, and major versions for any breaking changes. When we make breaking changes, we also introduce deprecation warnings in a minor version so that our users learn about the upcoming changes and migrate their code in advance.

Every significant change is documented in each [release](https://github.com/reviewpad/reviewpad/releases).

### Branch Organization

Submit all changes directly to the [`main branch`](https://github.com/reviewpad/reviewpad/tree/main). We don't use separate branches for development or for upcoming releases. We do our best to keep `main` in good shape, with all tests passing.

### Bugs

#### Where to Find Known Issues

We are using [GitHub Issues](https://github.com/reviewpad/reviewpad/issues) for our bugs. We keep a close eye on this and try to make it clear when we have an internal fix in progress. Before filing a new task, try to make sure your problem doesn't already exist.

#### Reporting New Issues

The best way to get your bug fixed is fill our [bug report template](https://github.com/reviewpad/reviewad/issues/new?assignees=&labels=bug&template=bug_report.md). Please follow the instructions in the template and provide as much information as possible. The more information you provide, the more likely someone will be successful reproducing the issue and finding a fix.

### How to Get in Touch

* [Discord](https://reviewpad.com/discord)
* [GitHub Discussions](https://github.com/reviewpad/reviewpad/discussions)

Please note that the Reviewpad [Discord](https://reviewpad.com/discord) is the same used by the team. So you'll be talking directly with the team 💪

### Proposing a Change

If you intend to change a piece of Reviewpad that is not yet being addressed on any issue we first invite you to [fill an issue](https://github.com/reviewpad/reviewpad/issues/new?assignees=&labels=enhancement&template=feature_request.md). This lets us reach an agreement on your proposal before you put significant effort into it.

If you're only fixing a bug, it's fine to submit a pull request right away but we still recommend to file an issue detailing what you're fixing. This is helpful in case we don't accept that specific fix but want to keep track of the issue.

### Your First Pull Request

Working on your first Pull Request? You can learn how from this free video series:

**[How to contribute to open source within 10 minutes - start now!](https://www.youtube.com/watch?v=8B_JWf7pG20)** by our friend [Eddie Jaoude](https://twitter.com/eddiejaoude)

To help you get your feet wet and get you familiar with our contribution process, we have a list of **[good first issues](https://github.com/reviewpad/reviewpad/issues?q=is:open+is:issue+label:"good+first+issue")** that contain bugs that have a relatively limited scope. This is a great place to get started.

If you decide to fix an issue, please be sure to check the comment thread in case somebody is already working on a fix. If nobody is working on it at the moment, please leave a comment stating that you intend to work on it so other people don't accidentally duplicate your effort.

### Sending a Pull Request

The [Reviewpad Action](https://github.com/reviewpad/action) is monitoring the pull requests and notifying the team. We will review your pull request and either merge it, request changes to it, or close it with an explanation. We'll do our best to provide updates and feedback throughout the process.

**Before submitting a pull request,** please make sure the following is done:

1. Fork [the repository](https://github.com/reviewpad/reviewpad) and create your branch from `main`.
2. Run `go get` in the repository root.
3. If you've fixed a bug or added code that should be tested, add tests!
4. Ensure the test suite passes (`task test`).
5. If you need a debugger, checkout how to run our [Reviewpad CLI](https://github.com/reviewpad/reviewpad#compilation).
6. Format you format your code code (`task format`).
7. Do a final check (`task check -f`).

<!-- Add Contributor License Agreement (CLA) -->

### Contribution Prerequisites

* You have [GO](https://go.dev) installed with minimum version of 1.18.2.
* You have [Taskfile](https://taskfile.dev/installation/) installed.
* You have `gcc` installed or are comfortable installing a compiler if needed. Some of our dependencies may require a compilation step. On OS X, the Xcode Command Line Tools will cover this. On Ubuntu, `apt-get install build-essential` will install the required packages. Similar commands should work on other Linux distros. Windows will require some additional steps, see the [`node-gyp` installation instructions](https://github.com/nodejs/node-gyp#installation) for details.
* You are familiar with Git.

### Development Workflow

After cloning Reviewpad, run `go get` to fetch its dependencies.
Then, you can run several commands:

* `task format` formats your code.
* `task test` runs the complete test suite.
* `task build` build the project (no output is generated).
* `task build-ci` build Reviewpad CLI and creates the executable `reviewpad-cli`.
* `task check` runs the code formatter, the complete test suite, and builds the project.

We recommend running `task check -f` to make sure you don't introduce any issue.

The easiest way to try your changes is to use [Reviewpad CLI](https://github.com/reviewpad/reviewpad#compilation) `./reviewpad-cli run -u [PULL_REQUEST_URL] -t [GITHUB_TOKEN] -e [GITHUB_EVENT] -f [REVIEWPAD_FILE]`.

Another great way to test your changes it to use `VSCode` debug tool. You can do that by creating a `.vscode` folder in the root of the repository, add the file `launch.json` to it, and copy the following content into the file:

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch cli",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "args": [
                "run",
                // Flag to run on dry run
                "-d=true",
                // Absolute path to reviewpad.yml file to run
                "-f=REVIEWPAD_FILE",
                // GitHub url to run the reviewpad.yml against to
                // e.g. https://github.com/reviewpad/action-demo/pull/1
                // e.g.https://github.com/reviewpad/action-demo/issue/1
                "-u=PULL_REQUEST_URL",
                // Absolute path to JSON file with GitHub event payload (OPTIONAL)
                "-e=GITHUB_EVENT",
                // GiHub PAT (Personal Access Token)
                // https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
                "-t=GITHUB_TOKEN",
            ],
            "program": "${workspaceFolder}/cli",
        },
    ]
}
```

We still require that your pull request contains unit tests for any new functionality. This way we can ensure that we don't break your code in the future.

### Commit convention

As commit convention, we adopt [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

This is not the case for all the commits in the history but any new commit should follow it.

If you want to enforce it as a pre-hook commit you can use [tiger](https://github.com/marcelosousa/tiger).


### Style Guide

We use an automatic code formatter called [gofmt](https://pkg.go.dev/cmd/gofmt).
Run `task format` after making any changes to the code.

Then, our formatter will catch most issues that may exist in your code.

<!-- Add Request for Comments (RFC) -->

### License

By contributing to Reviewpad, you agree that your contributions will be licensed under its LGPL-3.0 license.

### Attribution

This documentation is based on [reactjs.org](https://reactjs.org/) contributing guidelines.

<!-- Add What Next? -->
