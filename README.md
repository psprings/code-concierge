# code-concierge

This is a supplemental binary for a code-server environment which will automate fetching dependencies for a GitHub repo to set up a remote text editor/workstation.

## Concept

Given a code-server Docker image and a GitHub repo URL, identify the language(s) and download/install useful extensions for VS Code,
install binaries for development, and clone the repository.

## Usage

The easiest way to get up and running

```shell
# This needs to not include the `.git` for now
# will sanitize that URL later
export GITHUB_REPO_URL=<YOUR GITHUB REPO URL HERE>
# OPTIONAL - unless your repo is private or you are on enterprise with private mode
export GITHUB_API_TOKEN=""
docker run -it -e GITHUB_REPO_URL -e GITHUB_API_TOKEN -p 8443:8443 psprings/code-concierge --allow-http --no-auth
```

This will be availabe at <http://localhost:8443> after pre-requisites are installed and `code-server` starts

## TODO

* Add table of environment variable options to documentation
* Return logs on `8443` until process is ready to hand off to `code-server`