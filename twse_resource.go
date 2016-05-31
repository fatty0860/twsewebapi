package twsewebapi

import (
	"container/list"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type twseMktSource struct {
	//OnDataArrived : 收到報價資料
	OnDataArrived func(data interface{})
	hdl           *TwseStkHdl
	syms          *list.List
	quit          chan int
	isEnable      bool
	mtx           sync.Mutex
}

func (src *twseMktSource) Init() {
	src.syms = list.New()
}

func (src *twseMktSource) AddSymbol(sym string) bool {
	src.mtx.Lock()
	defer src.mtx.Unlock()

	for v := src.syms.Front(); v != nil; v = v.Next() {
		if v.Value.(string) == sym {
			return false
		}
	}
	src.syms.PushBack(sym)
	return true
}

func (src *twseMktSource) RemoveSymbol(sym string) bool {
	src.mtx.Lock()
	defer src.mtx.Unlock()

	for v := src.syms.Front(); v != nil; v = v.Next() {
		if v.Value.(string) == sym {
			src.syms.Remove(v)
			return true
		}
	}
	return false
}

func (src *twseMktSource) Start(hdl *TwseStkHdl, interval int) {
	if src.isEnable == false {
		src.isEnable = true
		src.hdl = hdl
		go func() {
			var syncnt, i int
			var rep StockInfoResponse
			var err error

			for src.isEnable {
				syncnt = src.syms.Len()
				if src.OnDataArrived != nil && syncnt > 0 {
					// copy all register symbol
					src.mtx.Lock()
					ary := make([]string, syncnt, syncnt)
					for v := src.syms.Front(); v != nil; v = v.Next() {
						ary = append(ary, v.Value.(string))
					}
					src.mtx.Unlock()

					// send to query
					for i = 0; i < syncnt; i += 10 {
						rep, err = src.hdl.QryStkInfoBatch(ary[i:(i + 10)])
						if err != nil {
							if src.hdl.IsTimeout(err) {
								i -= 10
								continue
							}
						}
						src.OnDataArrived(rep)
					}

					if i < (syncnt - 1) {
						rep, err = src.hdl.QryStkInfoBatch(ary[i:])
						for err != nil {
							rep, err = src.hdl.QryStkInfoBatch(ary[i:])
						}
						src.OnDataArrived(rep)
					}
				}
				time.Sleep(time.Second * time.Duration(interval))
			}
		}()
	}
	return
}

func (src *twseMktSource) Stop() {

}

//Test : for test
func Test() {
	fmt.Printf("test\n")
	c := &http.Client{
		Timeout: time.Second * 2}

	_, err := c.Get("http://1.1.1.1")
	if err != nil {
		fmt.Printf("err : %s\n", err.Error())
		nerr := err.(net.Error)
		if nerr.Timeout() {
			fmt.Printf("timeout\n")
		}
		return
	}

}
