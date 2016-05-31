package main

import (
	"fatty0860/twsewebapi"
	"fmt"
)

func main() {
	twsewebapi.Test()
	fmt.Printf("done\n")
	/*
		var hdl twsewebapi.TwseStkHdl

		hdl.Init(5)

		rep, err := hdl.QryStock("2330")

		if err != nil {
			fmt.Printf("err %s\n", err.Error())
			return
		}
		fmt.Printf("%v\n", rep)

		rep2, err2 := hdl.QryStkInfo(rep.Info[0].StkKey)
		if err2 != nil {
			fmt.Printf("err %s\n", err.Error())
			return
		}

		fmt.Printf("MatchPx : %s, %s\n", rep2.Info[0].MatchPx, rep2.Info[0].MatchQty)

		vec := []string{"tse_2330.tw", "tse_2062.tw", "tse_1101.tw"}
		rep2, err2 = hdl.QryStkInfoBatch(vec)
		if err2 != nil {
			fmt.Printf("err %s\n", err.Error())
			return
		}
		fmt.Printf("%v\n", rep2)
	*/
}
