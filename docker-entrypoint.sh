#!/bin/bash
set -e

code-concierge

exec dumb-init code-server "$@"