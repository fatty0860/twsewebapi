package main

import (
	"fatty0860/twsewebapi"
	"fmt"
	"strings"
	"time"
)

func work(data interface{}) {
	o, _ := data.(twsewebapi.StockInfoResponse)
	fmt.Printf("---------------------------- \n")
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

func main() {
	var src twsewebapi.TwseMktSource
	var hdl twsewebapi.TwseStkHdl
	var cmd string

	hdl.Init(5)
	src.Init()
	src.OnDataArrived = work
	/*
		src.AddSymbol("tse_2330.tw")
		src.AddSymbol("tse_2062.tw")
		src.AddSymbol("tse_1101.tw")
		src.AddSymbol("tse_2603.tw")
		src.AddSymbol("tse_2610.tw")
		src.AddSymbol("tse_2882.tw")
		src.AddSymbol("tse_3008.tw")
		src.AddSymbol("tse_2311.tw")
		src.AddSymbol("tse_2049.tw")
		src.AddSymbol("tse_1301.tw")
		src.AddSymbol("otc_2736.tw")
	*/

	src.Start(&hdl, 5)
	defer src.Stop()
	for {
		fmt.Scanln(&cmd)
		switch cmd {
		case "quit":
			return
		case "list":
			fmt.Printf("%s\n", cmd)
			ary := src.GetRegisterSymbol()
			for _, stk := range ary {
				fmt.Printf("Reg : %s\n", stk)
			}
		default:
			fmt.Printf("DEBUG - [%s]\n", cmd)
			src.AddSymbol(cmd)
		}
		//time.Sleep(time.Second * 10)
	}

}
