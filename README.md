# ghrepos

[![Build Status](https://travis-ci.org/seiffert/ghrepos.svg?branch=master)](https://travis-ci.org/seiffert/ghrepos)

**ghrepos** is a simple CLI tool that allows you to find GitHub repositories with a specific
[topic](https://help.github.com/articles/about-topics/). This is particularly useful when you automate the 
configuration of build tools or the management of your GitHub repositories.

## Example Usage

```bash
$ ghrepos --owner seiffert example-topic
seiffert/example-repo-1
seiffert/example-repo-2
seiffert/example-repo-3
```

To authenticate against the GitHub API, **ghrepos** optionally takes a GitHub access token. This is required if you want
to list private repositories.
Generate a token one in [your account settings](https://github.com/settings/tokens) and pass it either as environment 
variable `GITHUB_TOKEN` or via the `--token` option:

```bash
$ ghrepos --owner seiffert --token <GITHUB_TOKEN> example-topic
```

## Installation

To install **ghrepos**, download a binary from the provided
[GitHub releases](https://github.com/seiffert/ghrepos/releases) and put it into a folder that is part of your 
system's `$PATH`.

## Contribution

If you have ideas for improving this little tools or just a question, please don't hesitate to open an
[issue](https://github.com/seiffert/ghrepos/issues/new) or even fork this repository and create a pull-request!
