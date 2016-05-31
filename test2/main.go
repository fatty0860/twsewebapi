package main

import (
	"fatty0860/twsewebapi"
	"fmt"
	"strings"
	"time"
)

func thdwork(sec int, mktdata *chan interface{}) {
	var syms []string
	var hdl twsewebapi.TwseStkHdl
	hdl.Init(5)

	syms = append(syms, "tse_2002.tw")
	syms = append(syms, "tse_2062.tw")
	syms = append(syms, "tse_2303.tw")

	for {
		//rep, err := hdl.QryStkInfo("tse_2330.tw")
		rep, err := hdl.QryStkInfoBatch(syms)
		if err != nil {
			if hdl.IsTimeout(err) {
				fmt.Printf("connection timeout wait 1 sec")
				time.Sleep(time.Second * 1)
				continue
			}

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
					bidpx := stk.Best5BidPx[:strings.IndexAny(stk.Best5BidPx, "_")]
					askpx := stk.Best5AskPx[:strings.IndexAny(stk.Best5AskPx, "_")]
					bidqty := stk.Best5BidQty[:strings.IndexAny(stk.Best5BidQty, "_")]
					askqty := stk.Best5AskQty[:strings.IndexAny(stk.Best5AskQty, "_")]

					fmt.Printf("%s [%-8s] Bid[%-8s, %-4s] Ask[%-8s, %-4s] M[%-8s, %-4s] [%d]\n",
						time.Now().Format("15:04:05"),
						stk.Channel,
						bidpx, bidqty,
						askpx, askqty,
						stk.MatchPx, stk.MatchQty,
						stk.Tlong)
				}
			}
		}
	}

	//fmt.Printf("test\n")

}
