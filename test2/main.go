package main

import (
	"fatty0860/twsewebapi"
	"fmt"
	"time"
)

func thdwork(sec int, mktdata *chan interface{}) {
	var syms []string
	var hdl twsewebapi.TwseStkHdl
	hdl.Init(5)

	syms = append(syms, "tse_2330.tw")
	syms = append(syms, "tse_2062.tw")

	for {
		//rep, err := hdl.QryStkInfo("tse_2330.tw")
		rep, err := hdl.QryStkInfoBatch(syms)
		if err != nil {
			fmt.Printf("err : %s\n", err.Error())
			return
		}
		*mktdata <- rep
		time.Sleep(time.Second * 5)
	}
	//return
}

func main() {
	mktdata := make(chan interface{})

	go thdwork(5, &mktdata)

	for {
		select {
		case data := <-mktdata:
			fmt.Printf("----------------------------\n")
			if o, ok := data.(twsewebapi.StockInfoResponse); ok {
				for _, stk := range o.Info {
					fmt.Printf("%s Sym %s = %s, %s [%d]\n", time.Now().Format("2006-01-02 15:04:05"),
						stk.Channel, stk.MatchPx, stk.MatchQty,
						stk.Tlong)
				}
			}
		}
	}

	//fmt.Printf("test\n")

}
