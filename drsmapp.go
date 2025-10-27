// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"os"
	"time"

	"github.com/omec-project/util/drsm"
	"github.com/omec-project/util/logger"
)

type drsmInterface struct {
	initDrsm bool
	Mode     drsm.DrsmMode
	d        *drsm.Drsm
	poolName string
}

var drsmIntf drsmInterface

func scanChunk(i int32) bool {
	logger.DrsmLog.Debugf("received callback from module to scan Chunk resource %+v", i)
	return false
}

func initDrsm(resName string) {
	if drsmIntf.initDrsm {
		return
	}
	drsmIntf.initDrsm = true
	drsmIntf.poolName = resName

	podn := os.Getenv("HOSTNAME") // pod-name
	podi := os.Getenv("POD_IP")
	podId := drsm.PodId{PodName: podn, PodIp: podi}
	db := drsm.DbInfo{Url: "mongodb://mongodb-arbiter-headless", Name: "sdcore"}

	t := time.Now().UnixNano()
	opt := &drsm.Options{}
	if t%2 == 0 {
		logger.DrsmLog.Debugln("running in Demux Mode")
		drsmIntf.Mode = drsm.ResourceDemux
	} else {
		opt.ResourceValidCb = scanChunk
		opt.IpPool = make(map[string]string)
		opt.IpPool["pool1"] = "192.168.1.0/24"
		opt.IpPool["pool2"] = "192.168.2.0/24"
	}
	drsmInitialize, err := drsm.InitDRSM(resName, podId, db, opt)
	if err != nil {
		logger.DrsmLog.Fatalf("DRSM initialization failed: %+v", err)
	}
	drsmIntf.d = drsmInitialize.(*drsm.Drsm)
}

func AllocateInt32One(resName string) int32 {
	id, err := drsmIntf.d.AllocateInt32ID()
	if err != nil {
		logger.DrsmLog.Debugf("id allocation error %+v", err)
		return 0
	}
	logger.DrsmLog.Infof("received id %d", id)
	return id
}

func AllocateInt32Many(resName string, number int32) []int32 {
	var resIds []int32
	var count int32 = 0

	ticker := time.NewTicker(50 * time.Millisecond)
	for range ticker.C {
		id, err := drsmIntf.d.AllocateInt32ID()
		if err != nil {
			logger.DrsmLog.Debugf("id allocation error %+v", err)
			continue
		}
		if id != 0 {
			resIds = append(resIds, id)
		}
		logger.DrsmLog.Infof("received id %d", id)
		count++
		if count >= number {
			ticker.Stop()
			return resIds
		}
	}
	return resIds
}

func ReleaseInt32One(resName string, resId int32) error {
	err := drsmIntf.d.ReleaseInt32ID(resId)
	if err != nil {
		logger.DrsmLog.Debugf("id release error %+v", err)
		return err
	}
	return nil
}
