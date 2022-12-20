package chain

import (
	"errors"
	"time"
	"videown-server/global"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/goccy/go-json"
)

const rpcConnNum = 16 * 1024

var rpcWorkPool chan struct{}

func InitRpcWorkPool() {
	if rpcWorkPool == nil {
		rpcWorkPool = make(chan struct{}, rpcConnNum)
	}
}

func (c *chainClient) SendTx(signtx string) (string, error) {
	var ext types.Extrinsic
	var txhash string
	if ext.UnmarshalJSON([]byte(signtx)) != nil {
		bytes, _ := json.Marshal(signtx)
		if err := ext.UnmarshalJSON(bytes); err != nil {
			global.Logger.Errorf("[Send tx error] %v.", err)
			return txhash, err
		}
	}
	rpcWorkPool <- struct{}{}
	defer func() {
		<-rpcWorkPool
	}()
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		global.Logger.Errorf("[Send tx error] %v.", err)
		return txhash, err
	}
	timeout := time.After(c.timeForBlockOut)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				events := types.EventRecords{}
				h, err := c.api.RPC.State.GetStorageRaw(c.keyEvents, status.AsInBlock)
				if err != nil {
					global.Logger.Errorf("[Send tx error] get event raw error, %v.", err)
					return txhash, err
				}
				types.EventRecordsRaw(*h).DecodeEventRecords(c.metadata, &events)
				if len(events.System_ExtrinsicFailed) > 0 {
					err := errors.New("system.ExtrinsicFailed")
					global.Logger.Errorf("[Send tx error] %v.", err)
					return txhash, err
				}
				if len(events.System_ExtrinsicSuccess) > 0 {
					global.Logger.Errorf("[Send tx success] tx hash is %v.", txhash)
					return txhash, err
				}
			}
		case err := <-sub.Err():
			global.Logger.Errorf("[Send tx error] channel error, %v.", err)
			return txhash, err
		case <-timeout:
			err := errors.New("send tx timeout")
			global.Logger.Errorf("[Send tx error] %v.", err)
			return txhash, err
		}
	}
}
