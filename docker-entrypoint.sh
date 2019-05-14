#!/bin/bash
set -e

if [ -z "${GITHUB_REPO_URL}" ]; then
    exit 1
fi

code_concierge_args=""

if [ "${INSTALL_DOCKER_CLI}" = "true" ]; then
    code_concierge_args="${code_concierge_args} --install-docker"
fi

if [ ! -z "${ADDITIONAL_EXTENSIONS}" ]; then
    code_concierge_args="${code_concierge_args} --additional-extensions=${ADDITIONAL_EXTENSIONS}"
fi

if [ ! -z "${ADDITIONAL_PACKAGES}" ]; then
    code_concierge_args="${code_concierge_args} --additional-packages=${ADDITIONAL_PACKAGES}"
fi

code-concierge "${code_concierge_args}"

exec dumb-init code-server "$@"