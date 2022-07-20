# How to contribute

It is awesome that you are reading this, because we need volunteer developers to help this project!

- [How to contribute](#how-to-contribute)
  - [Code of Conduct](#code-of-conduct)
  - [Community](#community)
  - [Contributing](#contributing)
    - [Issues](#issues)
    - [Development Policy](#development-policy)
    - [Pull Requests](#pull-requests)
    - [Reviews](#reviews)
    - [Commit convention](#commit-convention)
  - [Attribution](#attribution)

## Code of Conduct

Reviewpad projects have a [Code of Conduct](CODE_OF_CONDUCT.md) to which all contributors must adhere.
Please read it before interacting with the repository or the community in any way.

## Community

Discussion and **support requests** should go through [Discord](http://reviewpad.com/discord).

We also publish FAQs and announcements in [GitHub discussions](https://github.com/reviewpad/reviewpad/discussions).

## Contributing

Not surprisingly, all reviewpad projects use the [Reviewpad action](https://github.com/reviewpad/action) for contributions.

The workflows are specified at [Reviewpad.yml](reviewpad.yml). 

Take the time to get familiar with the configuration in the official [docs](https://docs.reviewpad.com).

### Issues

We plan to extensively use GitHub issues to communicate the ongoing and near future development. 

There are mainly three kinds of issues you can open:

* Bug report: you believe you found a problem in a project and you want to discuss and get it fixed,
  creating an issue with the **bug report template** is the best way to do so.
* Feature request: would you like a new feature/rule to be added to Reviewpad? This is the kind of issue you'll need then.
  Do your best at explaining your intent, it is always important that others can understand what you mean in order to discuss.
  Be open and collaborative in letting others help you get things done!

The best way to **get involved** in the project is through issues. You can help in many ways:

* Issues triaging: participating in the discussion and adding details to open issues is always a good thing;
sometimes issues need to be verified, you could be the one writing a test case to fix a bug!
* Helping to resolve the issue: you can help in getting it fixed in many ways, more often by opening a pull request.

### Development Policy

Our development policies are explicitly formalized in the [Reviewpad.yml](https://github.com/reviewpad/reviewpad/blob/main/reviewpad.yml) file.

We use short-lived feature branches and pull requests to introduce changes to the codebase.

Reviewpad will take care of most of the automation so that 90% of the pull requests do not stay open for longer than a few hours (even for outside contributors).

Because of the limitation in described in the [official GitHub documentation](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#using-the-github_token-in-a-workflow):

> When you use the repository's GITHUB_TOKEN to perform tasks, events triggered by the GITHUB_TOKEN will not create a new workflow run.

We run the Reviewpad action with an access token from the [reviewpad-team](https://github.com/reviewpad-team).
As soon as GitHub resolves this [issue](https://github.community/t/triggering-a-new-workflow-from-another-workflow/16250),
the automation actions should be done through the `github-actions (bot)`.

### Pull Requests

Thanks for taking time to make a [pull request](https://help.github.com/articles/about-pull-requests) (PR).

The PR template is there to guide you through the process of opening it.

### Reviews

Reviewing a pull request is also a very good way of contributing.

### Commit convention

As commit convention, we adopt [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/).

This is not the case for all the commits in the history but any new commit should follow it.

If you want to enforce it as a pre-hook commit you can use [tiger](https://github.com/marcelosousa/tiger).

## Attribution

This contributing guide is inspired from the [Falco project](https://github.com/falcosecurity/.github/blob/master/CONTRIBUTING.md).
