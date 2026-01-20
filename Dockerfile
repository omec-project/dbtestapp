# SPDX-FileCopyrightText: 2022-present Intel Corporation
# Copyright 2019-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.25.6-bookworm@sha256:2f768d462dbffbb0f0b3a5171009f162945b086f326e0b2a8fd5d29c3219ff14 AS test

LABEL maintainer="Aether SD-Core <dev@lists.aetherproject.org>"

RUN apt-get update && apt-get -y install --no-install-recommends vim

WORKDIR $GOPATH/src/dbtestapp
COPY . .
RUN go install

FROM alpine:3.23@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62 AS dbtestapp
RUN apk add --no-cache gcompat vim nano strace net-tools curl netcat-openbsd bind-tools bash && rm -rf /var/cache/apk/*

RUN mkdir -p /dbtestapp/bin
COPY --from=test /go/bin/* /dbtestapp/bin/
WORKDIR /dbtestapp
