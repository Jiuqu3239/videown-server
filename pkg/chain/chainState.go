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
	"github.com/CESSProject/cess-oss/pkg/utils"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// GetPublicKey returns your own public key
func (c *chainClient) GetPublicKey() []byte {
	return c.keyring.PublicKey
}

func (c *chainClient) GetMnemonicSeed() string {
	return c.keyring.URI
}

func (c *chainClient) GetSyncStatus() (bool, error) {
	if !c.IsChainClientOk() {
		return false, ERR_RPC_CONNECTION
	}
	h, err := c.api.RPC.System.Health()
	if err != nil {
		return false, err
	}
	return h.IsSyncing, nil
}

func (c *chainClient) GetChainStatus() bool {
	return c.GetChainState()
}

func (c *chainClient) GetChainMethodList() []string {
	methods := make([]string, 0)
	for _, v := range c.metadata.AsMetadataV8.Modules {
		methods = append(methods, string(v.Name))
	}
	return methods
}

func (c *chainClient) GetCessAccount() (string, error) {
	return utils.EncodePublicKeyAsCessAccount(c.keyring.PublicKey)
}

func (c *chainClient) GetAccountInfo(pkey []byte) (types.AccountInfo, error) {
	var data types.AccountInfo

	if !c.IsChainClientOk() {
		c.SetChainState(false)
		return data, ERR_RPC_CONNECTION
	}
	c.SetChainState(true)

	b, err := types.Encode(types.NewAccountID(pkey))
	if err != nil {
		return data, errors.Wrap(err, "[EncodeToBytes]")
	}

	key, err := types.CreateStorageKey(
		c.metadata,
		pallet_System,
		account,
		b,
	)
	if err != nil {
		return data, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return data, ERR_RPC_EMPTY_VALUE
	}
	return data, nil
}
