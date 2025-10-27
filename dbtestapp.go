// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"github.com/omec-project/util/logger"
	"github.com/omec-project/util/mongoapi"
)

var mongoHndl *mongoapi.MongoClient

// TODO : take DB name from helm chart
// TODO : inbuild shell commands to

func main() {
	logger.MongoapiLog.Infoln("dbtestapp started")

	// connect to mongoDB
	var err error
	mongoHndl, err = mongoapi.NewMongoClient("mongodb://mongodb-arbiter-headless", "sdcore")
	if err != nil {
		logger.MongoapiLog.Fatalf("failed to connect to MongoDB: %+v", err)
	}

	initDrsm("resourceids")

	// blocking
	http_server()
}
