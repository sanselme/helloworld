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

source ./scripts/load-env.sh

# debug
echo """
Image:      ${IMAGE_NAME}
Platforms:  ${IMAGE_PLATFORMS}
Registry:   ${REGISTRY}
Tag:        ${IMAGE_TAG}
Target:     ${TARGET}
"""

# build
docker buildx build \
  --platform "${IMAGE_PLATFORMS}" \
  --tag "${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}" \
  --tag "${REGISTRY}/${IMAGE_NAME}:latest" \
  --file "${DOCKER_FILE}" \
  --build-arg "GIT_BRANCH=${GIT_BRANCH}" \
  --build-arg "GIT_COMMIT=${GIT_COMMIT}" \
  --build-arg "TARGET=${TARGET}" \
  --build-arg "IMAGE_NAME=${IMAGE_NAME}" \
  .
