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

# ref: https://github.com/johanbrandhorst/grpc-gateway-boilerplate/blob/main/scripts/generate-swagger-ui.sh
set -e

CACHE_DIR=".cache/swagger-ui"
OUTPUT_DIR="./docs/openapi"
SWAGGER_UI_REPO="https://github.com/swagger-api/swagger-ui.git"
SWAGGER_UI_VERSION="${SWAGGER_UI_VERSION:-$1}"

escape_str() {
  echo "$1" | sed -e 's/[]\/$*.^[]/\\&/g'
}

[[ -z ${SWAGGER_UI_VERSION} ]] &&
  echo 'missing ${SWAGGER_UI_VERSION}' &&
  exit 1

[[ ! -d ${OUTPUT_DIR} ]] && mkdir -p "${OUTPUT_DIR}"

# cache swagger-ui dist
if [[ ! -d ${CACHE_DIR} ]]; then
  mkdir -p "${CACHE_DIR}"
  tmp="$(mktemp -d)"
  git clone --depth 1 --branch "${SWAGGER_UI_VERSION}" "${SWAGGER_UI_REPO}" "${tmp}"
  cp -r "${tmp}/dist/"* "${CACHE_DIR}"
  cp -r "${tmp}/LICENSE" "${CACHE_DIR}"
  rm -rf "${tmp}"
fi

# populate swagger.json
buf generate
tmp="    urls: ["
for i in $(find "${OUTPUT_DIR}" -name "*.swagger.json"); do
  escaped_gen_dir="$(escape_str "${OUTPUT_DIR}/")"
  path="${i//${escaped_gen_dir}/}"
  tmp="${tmp}{\"url\":\"${path}\",\"name\":\"${path}\"},"
done
# delete last characters from $tmp
tmp="${tmp//.$/}"
tmp="${tmp}],"

# generate swagger-ui
find "${OUTPUT_DIR}" -type f -name "*.swagger.json" -delete
mkdir -p "${OUTPUT_DIR}"
cp -r "${CACHE_DIR}/"* "${OUTPUT_DIR}"

# replace the default URL
line="$(cat "${OUTPUT_DIR}/swagger-initializer.js" | grep -n "url" | cut -f1 -d:)"
escaped_tmp="$(escape_str "${tmp}")"
sed -i'' -e "${line} s/^.*$/${escaped_tmp}/" "${OUTPUT_DIR}/swagger-initializer.js"
rm -f "${OUTPUT_DIR}/swagger-initializer.js-e"

# trunk.io
trunk fmt --all
