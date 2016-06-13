package twsewebapi

import (
	"container/list"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

//TwseMktSource : 報價源
type TwseMktSource struct {
	//OnDataArrived : 收到報價資料
	OnDataArrived func(data interface{})
	hdl           *TwseStkHdl
	syms          *list.List
	isEnable      bool
	mtx           sync.Mutex
	wait          sync.WaitGroup
}

//Init : 物件初始化, 請優先執行此方法
func (src *TwseMktSource) Init() {
	src.syms = list.New()
}

//AddSymbol : 新增註冊商品
func (src *TwseMktSource) AddSymbol(sym string) bool {
	var stkKey string
	rep, err := src.hdl.QryStock(sym)
	fmt.Printf("DEBUG : %v\n", rep)
	for err != nil {
		if src.hdl.IsTimeout(err) {
			rep, err = src.hdl.QryStock(sym)
			continue
		} else {
			return false
		}
	}
	if len(rep.Info) < 1 {
		return false
	}
	stkKey = rep.Info[0].StkKey

	// add to symbol list
	src.mtx.Lock()
	defer src.mtx.Unlock()

	for v := src.syms.Front(); v != nil; v = v.Next() {
		if v.Value.(string) == stkKey {
			return false
		}
	}
	src.syms.PushFront(stkKey)
	return true
}

//RemoveSymbol : 刪除註冊商品
func (src *TwseMktSource) RemoveSymbol(sym string) bool {
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

//GetRegisterSymbol : 取得註冊商品清單
func (src *TwseMktSource) GetRegisterSymbol() (ary []string) {
	// copy all register symbol
	src.mtx.Lock()
	for v := src.syms.Front(); v != nil; v = v.Next() {
		ary = append(ary, v.Value.(string))
	}
	src.mtx.Unlock()
	return
}

//Start : 開始定時查詢商品資料
func (src *TwseMktSource) Start(hdl *TwseStkHdl, interval int) {
	if src.isEnable == false {
		src.isEnable = true
		src.hdl = hdl

		src.wait.Add(1)

		go func() {
			var syncnt, i int
			var rep StockInfoResponse
			var err error
		NEXT:
			for src.isEnable {
				syncnt = src.syms.Len()

				if src.OnDataArrived != nil && syncnt > 0 {

					ary := src.GetRegisterSymbol()

					// send to query
					if syncnt > 10 {
						for i = 0; i < (syncnt - 10); i += 10 {
							rep, err = src.hdl.QryStkInfoBatch(ary[i:(i + 10)])
							if err != nil {
								if src.hdl.IsTimeout(err) {
									i -= 10
									continue
								}
							}
							src.OnDataArrived(rep)
						}
					}

					if i < syncnt {
						rep, err = src.hdl.QryStkInfoBatch(ary[i:])
						for err != nil {
							rep, err = src.hdl.QryStkInfoBatch(ary[i:])
							if src.hdl.IsTimeout(err) {
								fmt.Printf("1\n")
								continue
							} else {
								fmt.Printf("2\n")
								break NEXT
							}
						}
						src.OnDataArrived(rep)
					}
				}
				time.Sleep(time.Second * time.Duration(interval))
			}
			src.wait.Done()
		}()
	}
	return
}

//Stop : 停止資料查詢
func (src *TwseMktSource) Stop() {
	src.isEnable = false
	src.wait.Wait()
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
