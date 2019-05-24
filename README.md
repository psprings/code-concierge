# code-concierge

This is a supplemental binary for a code-server environment which will automate fetching dependencies for a GitHub repo to set up a remote text editor/workstation.

## Concept

Given a code-server Docker image and a GitHub repo URL, identify the language(s) and download/install useful extensions for VS Code,
install binaries for development, and clone the repository.

## Usage

### Docker

The easiest way to get up and running

```shell
# This needs to not include the `.git` for now
# will sanitize that URL later
export GITHUB_REPO_URL=<YOUR GITHUB REPO URL HERE>
# OPTIONAL - unless your repo is private or you are on enterprise with private mode
export GITHUB_API_TOKEN=""
docker run -it -e GITHUB_REPO_URL -e GITHUB_API_TOKEN -p 8443:8443 psprings/code-concierge --allow-http --no-auth
```

This will be availabe at http://localhost:8443 after pre-requisites are installed and `code-server` starts

#### Environment variables

The following environment variables can be set to modify the arguments for `code-concierge`

| Name                               | Description                                                          | Default          |
|------------------------------------|----------------------------------------------------------------------|------------------|
| `ADDITIONAL_EXTENSIONS`            | Code Server extension ids to install in addition to defaults (comma separated) |        |
| `ADDITIONAL_PACKAGES`              | Additional system packages to install (comma separated)              |                  |
| `CODE_CONCIERGE_ARGS`              | Can be used to pass arguments to `code-concierge`                    |                  |
| `GITHUB_REPO_URL`                  | The GitHub repository URL to clone code from and setup               |                  |
| `INSTALL_DOCKER_CLI`               | Determines whether the Docker CLI is installed (will auto-install if repo uses Docker) | `false`          |

### CLI Options

```shell
Usage of code-concierge:
  -additional-extensions string
        Comma separated list of extension IDs to install
  -additional-packages string
        Comma separated list of packages to install
  -api-token string
        The token to use for authentication to GitHub
  -api-url string
        (optional) The base URL for the GitHub API
  -install-docker
        Whether to install the Docker CLI
  -repo-url string
        The (https) URL of the GitHub repo to use
```

## TODO

* Add table of environment variable options to documentation
* Return logs on `8443` until process is ready to hand off to `code-server`