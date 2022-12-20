/*
   Copyright 2022 CESS scheduler authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

        http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// Pallets
const (
	pallet_FileBank    = "FileBank"
	pallet_FileMap     = "FileMap"
	pallet_Sminer      = "Sminer"
	pallet_SegmentBook = "SegmentBook"
	pallet_System      = "System"
	pallet_Oss         = "Oss"
)

// Pallet's method
const (
	// System
	account = "Account"
	events  = "Events"
)

// Extrinsics
const ()

const (
	FILE_STATE_ACTIVE  = "active"
	FILE_STATE_PENDING = "pending"
)

const (
	ERR_Failed  = "failed"
	ERR_Timeout = "timeout"
	ERR_Empty   = "empty"
)

// error type
var (
	ERR_RPC_CONNECTION  = errors.New("rpc connection failed")
	ERR_RPC_IP_FORMAT   = errors.New("unsupported ip format")
	ERR_RPC_TIMEOUT     = errors.New("timeout")
	ERR_RPC_EMPTY_VALUE = errors.New("empty")
)

type EventEmitted struct {
	From types.AccountID
	To   types.AccountID
	Id   types.U8
}
