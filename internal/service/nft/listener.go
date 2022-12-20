package nft

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"videown-server/global"
	"videown-server/internal/dto"
)

var CmpBaseUrl = "/"

const LISTEN_TIME_INTERVAL_SECOND = 1 * 60

type Listener interface {
	AddListenItem(data any)
}

// ListenedData is the object that be listened
type ListenedData interface {
	GetHash() string
	Handler(item *ListenItem) *ListenItem
}

// Listen group for file status updata
type StatusListener struct {
	B        int
	QueueLen int
	Queues   []chan ListenItem
}

type ListenItem struct {
	Data  any
	Count int
	Timer *time.Timer
}

var listener *StatusListener

func GetStatusListener() Listener {
	return listener
}

func InitDefaultStatusListener(b, qLen int) {
	listener = new(StatusListener)
	InitStatusListener(b, qLen, listener)
}

func InitStatusListener(b int, qLen int, listener *StatusListener) error {
	if b < 0 || qLen < 0 {
		return errors.New("invalid args")
	}
	listener.B = b
	listener.QueueLen = qLen
	listener.Queues = make([]chan ListenItem, 1<<b)
	for i := 0; i < len(listener.Queues); i++ {
		listener.Queues[i] = make(chan ListenItem, qLen)
		go listener.listenAndUpdate(i)
	}
	return nil
}

func (t *StatusListener) AddListenItem(data any) {
	ld, ok := data.(ListenedData)
	if !ok {
		global.Logger.Error("[File status update] can not prase data to ListenedData interface.")
		return
	}
	selector, err := strconv.ParseInt(ld.GetHash(), 16, 64)
	if err != nil {
		global.Logger.Errorf("[File status update] data hash: %s parse selector error,%v.", ld.GetHash(), err)
		return
	}
	selector = selector % (1 << t.B)
	item := ListenItem{
		Data:  ld,
		Count: 0,
		Timer: time.NewTimer(time.Millisecond),
	}
	t.Queues[selector] <- item
	global.Logger.Infof("[File status update] data hash: %s entry async listening", ld.GetHash())
}

func (t *StatusListener) listenAndUpdate(selecter int) {
	global.Logger.Infof("Listen routine %d start listen...", selecter)
	for {
		queueLen := len(t.Queues[selecter])
		for i := 0; i < queueLen; i++ {
			item := <-t.Queues[selecter]
			select {
			case <-item.Timer.C:
				go func() {
					itemPtr := item.Data.(ListenedData).Handler(&item)
					if itemPtr != nil {
						t.Queues[selecter] <- *itemPtr
					}
				}()
			default:
				t.Queues[selecter] <- item
			}
		}
		time.Sleep(10 * time.Second)
	}
}

type FileHash string

func (t FileHash) GetHash() string {
	return string(t)[len(t)-8:]
}

func (t FileHash) Handler(item *ListenItem) *ListenItem {
	fileMeta, err := getFileStatus(string(t))
	interval := (item.Count/10 + 1) * LISTEN_TIME_INTERVAL_SECOND
	item.Timer = time.NewTimer(time.Duration(interval) * time.Second)
	item.Count++
	if err != nil {
		global.Logger.Errorf("[Listen file status] get file status error:%v.", err)
		return item
	}
	UpdateVideoStatus(string(t), fileMeta.State)
	if fileMeta.State != FILE_STATUS_ACTIVE &&
		fileMeta.State != FILE_STATUS_CANCEL {
		return item
	}
	global.Logger.Infof("[File status update] data hash: %s exit async listening", t.GetHash())
	return nil
}

func getFileStatus(filehash string) (dto.FileMeta, error) {
	var resp dto.QueryResponse
	url, _ := url.JoinPath(CmpBaseUrl, filehash)
	response, err := http.Get(url)
	if err != nil {
		return resp.Ok, err
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return resp.Ok, err
	}
	err = json.Unmarshal(bytes, &resp)
	return resp.Ok, err
}
