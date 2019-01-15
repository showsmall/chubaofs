// Copyright 2018 The Container File System Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package wrapper

/* TODO why we need this separate package? */

import (
	"fmt"
	"github.com/tiglabs/containerfs/proto"
	"strings"
)

// DataPartition defines the wrapper of the data partition.
type DataPartition struct {
	// Will not be changed
	proto.DataPartitionResponse
	RandomWrite   bool
	PartitionType string
	Metrics *DataPartitionMetrics
}

// DataPartitionMetrics defines the wrapper of the metrics related to the data partition.
type DataPartitionMetrics struct {
	WriteCnt        uint64
	ReadCnt         uint64
	SumWriteLatency uint64
	SumReadLatency  uint64
	WriteLatency    float64
	ReadLatency     float64
}

// TODO do we have any naming convention for sorted items? DataPartitionSorter?
type DataPartitionSlice []*DataPartition

func (ds DataPartitionSlice) Len() int {
	return len(ds)
}
func (ds DataPartitionSlice) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}
func (ds DataPartitionSlice) Less(i, j int) bool {
	return ds[i].Metrics.WriteLatency < ds[j].Metrics.WriteLatency
}

// NewDataPartitionMetrics returns a new DataPartitionMetrics instance.
func NewDataPartitionMetrics() *DataPartitionMetrics {
	metrics := new(DataPartitionMetrics)
	metrics.WriteCnt = 1
	metrics.ReadCnt = 1
	return metrics
}

// String returns the string format of the data partition.
func (dp *DataPartition) String() string {
	return fmt.Sprintf("PartitionID(%v) Status(%v) ReplicaNum(%v) PartitionType(%v) Hosts(%v)",
		dp.PartitionID, dp.Status, dp.ReplicaNum, dp.PartitionType, dp.Hosts)
}

// GetAllAddrs returns the addresses of all the replicas of the data partition.
// TODO remove m?
func (dp *DataPartition) GetAllAddrs() (m string) {
	return strings.Join(dp.Hosts[1:], proto.AddrSplit) + proto.AddrSplit
}

func isExcluded(partitionId uint64, excludes []uint64) bool {
	for _, id := range excludes {
		if id == partitionId {
			return true
		}
	}
	return false
}

// TODO unused? remove?
func NewGetDataPartitionMetricsPacket(partitionid uint64) (p *proto.Packet) {
	p = new(proto.Packet)
	p.PartitionID = partitionid
	p.Magic = proto.ProtoMagic
	p.ExtentType = proto.NormalExtentType
	p.ReqID = proto.GenerateRequestID()
	p.Opcode = proto.OpGetDataPartitionMetrics

	return
}
