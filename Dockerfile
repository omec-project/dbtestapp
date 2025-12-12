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

FROM alpine:3.23@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375 AS dbtestapp
RUN apk add --no-cache gcompat vim nano strace net-tools curl netcat-openbsd bind-tools bash && rm -rf /var/cache/apk/*

RUN mkdir -p /dbtestapp/bin
COPY --from=test /go/bin/* /dbtestapp/bin/
WORKDIR /dbtestapp
