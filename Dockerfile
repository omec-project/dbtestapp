# SPDX-FileCopyrightText: 2022-present Intel Corporation
# Copyright 2019-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.25.5-bookworm@sha256:09f53deea14d4019922334afe6258b7b776afc1d57952be2012f2c8c4076db05 AS test

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
