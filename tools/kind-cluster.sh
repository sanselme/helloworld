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
CLUSTER: ${CLUSTER_NAME}
IMAGE: ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
CHART: ${CHART_NAME}
CHART_VERSION: ${CHART_VERSION}
TEMP_DIR: ${TEMP_DIR}
PATH: ${TEMP_DIR}/${CHART_NAME}-${CHART_VERSION}.tgz
"""

# create kind cluster
CLUSTERS="$(kind get clusters)"
[[ ${CLUSTERS} != *"${CLUSTER_NAME}"* ]] &&
  kind create cluster --name="${CLUSTER_NAME}"

# load images into cluster
kind load docker-image --name="${CLUSTER_NAME}" "${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"

# deploy hello-daemon chart
helm package -u "./deployment/chart/${CHART_NAME}" -d "/${TEMP_DIR}/"
helm uninstall "${CHART_NAME}" || true
[[ -f "./config/samples/values.yaml" ]] &&
  helm upgrade \
    --install "${CHART_NAME}" \
    "/${TEMP_DIR}/${CHART_NAME}-${CHART_VERSION}.tgz" \
    --values ./config/samples/values.yaml ||
  helm upgrade \
    --install "${CHART_NAME}" \
    "/${TEMP_DIR}/${CHART_NAME}-${CHART_VERSION}.tgz"
