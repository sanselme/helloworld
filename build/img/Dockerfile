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
ARG BUILD_IMAGE=golang:1.19
ARG RUNTIME_IMAGE=gcr.io/distroless/static:nonroot

FROM ${BUILD_IMAGE} AS builder

ARG GIT_COMMIT
ARG GIT_BRANCH
ARG TARGET

ENV GIT_BRANCH=${GIT_BRANCH}
ENV CGO_ENABLED=0
ENV GIT_COMMIT=${GIT_COMMIT}
ENV GOOS=linux

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN ./scripts/build.sh ${TARGET} /hellod

FROM ${RUNTIME_IMAGE}
COPY --from=builder /hellod /hellod
USER nonroot:nonroot
