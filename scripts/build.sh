#!/bin/bash

# Copyright (c) 2023 Schubert Anselme <schubert@anselm.es>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program. If not, see <https://www.gnu.org/licenses/>.
set -e

DATE="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"

case "${1}" in
"daemon")
  APP_NAME="${APP_NAME}d"
  ;;
"client")
  APP_NAME="${APP_NAME}ctl"
  ;;
"operator")
  APP_NAME="${APP_NAME}-operator"
  ;;
"plugin")
  APP_NAME="kubectl-${APP_NAME}"
  ;;
*)
  echo "error: missing argument!!!"
  echo "usage: build.sh <daemon|client|operator|plugin>"
  exit 1
  ;;
esac

# debug
echo "Building ${2:-${APP_NAME}}..."
echo "Version: ${GIT_BRANCH}"
echo "Commit: ${GIT_COMMIT}"
echo "Date: ${DATE}"

# build
go build \
  -ldflags "-X  github.com/sanselme/helloworld/pkg/version.Branch=${GIT_BRANCH} \
    -X  github.com/sanselme/helloworld/pkg/version.Commit=${GIT_COMMIT} \
    -X github.com/sanselme/helloworld/pkg/version.Date=${DATE}" \
  -o "${2:-bin/${APP_NAME}}" \
  "./cmd/${1}"
