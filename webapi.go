package twsewebapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// TwseMISWeb : Twse mis website
	twseMISWeb = "http://mis.twse.com.tw/"
	//TwseAPIBase : Twse api url
	twseAPIBase = "http://mis.twse.com.tw/stock/api"
	//TwseAPIGetStock : 查詢商品基本資料
	twseAPIGetStock = "getStock.jsp"
	//TwseAPIGetStockInfo : 查詢商品基本資訊
	twseAPIGetStockInfo = "getStockInfo.jsp"
)

// TwseStkHdl : 證交所資訊網站連線物件
type TwseStkHdl struct {
	Client  *http.Client
	Timeout time.Duration
	cookies []*http.Cookie
}

//Stock : TWSE 商品基本資料結構
type Stock struct {
	Exchange   string `json:"ex"`
	TradeDate  string `json:"d"`
	It         string `json:"it"`
	Name       string `json:"\n"`
	I          string `json:"i"`
	IP         string `json:"ip"`
	LowerPrice string `json:"w"`
	UpperPrice string `json:"u"`
	T          string `json:"t"`
	P          string `json:"p"`
	StkCode    string `json:"ch"`
	StkKey     string `json:"key"`
	Y          string `json:"y"`
	StopEnd    string `json:"rt"`
	StopBegin  string `json:"st"`
}

//StockRespose : TWSE 商品基本資料擋查詢結果
type StockRespose struct {
	Info       []Stock `json:"msgArray"`
	RtnMessage string  `json:"rtmessage"`
	RtnCode    string  `json:"rtcode"`
}

//StockInfo : TWSE 商品報價資料結構
type StockInfo struct {
	Simulator    string `json:"ts"`
	Fv           string `json:"fv"`
	Tk0          string `json:"tk0"`
	Tk1          string `json:"tk1"`
	Oa           string `json:"oa"`
	Ob           string `json:"ob"`
	Tlong        int64  `json:"tlong"`
	Ot           string `json:"ot"`
	Best5AskQty  string `json:"f"`
	Exchange     string `json:"ex"`
	Best5BidQty  string `json:"g"`
	Ov           string `json:"ov"`
	TradeDate    string `json:"d"`
	It           string `json:"it"`
	Best5BidPx   string `json:"b"`
	Symbol       string `json:"c"`
	Mt           string `json:"mt"`
	Best5AskPx   string `json:"a"`
	Name         string `json:"n"`
	OpenPx       string `json:"o"`
	DayLowerPx   string `json:"l"`
	Oz           string `json:"oz"`
	DayUpperPx   string `json:"h"`
	IP           string `json:"ip"` // 4 : 暫緩收盤 , 5 : 暫緩開盤
	I            string `json:"i"`
	LowerLimitPx string `json:"w"`
	TotalVolumn  string `json:"v"`
	UpperLimitPx string `json:"u"`
	DateTime     string `json:"t"`
	S            string `json:"s"`
	Pz           string `json:"pz"`
	MatchQty     string `json:"tv"`
	P            string `json:"p"`
	FullName     string `json:"nf"`
	Channel      string `json:"ch"`
	MatchPx      string `json:"z"`
	LastDayPx    string `json:"y"`
	Ps           string `json:"ps"`
}

// StockInfoResponse :
type StockInfoResponse struct {
	Info       []StockInfo `json:"msgArray"`
	RtnMessage string      `json:"rtmessage"`
	RtnCode    string      `json:"rtcode"`
	UserDelay  string      `json:"userDelay"`
}

// Init : 初始化 http Client
func (hdl *TwseStkHdl) Init(TimeoutSec time.Duration) {
	hdl.Client = &http.Client{
		Timeout: time.Second * TimeoutSec,
	}
}

// Private function
// --------------------------------------------------------------------
func (hdl *TwseStkHdl) getCookie() (err error) {
	var url string
	url = twseMISWeb
	var request *http.Request
	var response *http.Response

	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	response, err = hdl.Client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	hdl.cookies = response.Cookies()

	return
}

// Public function
// --------------------------------------------------------------------

/*QryStock : 查詢 商品資料
  http://mis.twse.com.tw/stock/api/getStock.jsp?ch=2330.tw&json=1
*/
func (hdl *TwseStkHdl) QryStock(sym string) (rep StockRespose, err error) {
	var url string
	var request *http.Request
	var response *http.Response
	var data []byte

	url = fmt.Sprintf("%s/%s?ch=%s.tw&json=1", twseAPIBase, twseAPIGetStock, sym)
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	response, err = hdl.Client.Do(request)

	if err != nil {
		return
	}
	defer response.Body.Close()

	data, err = ioutil.ReadAll(response.Body)
	json.Unmarshal([]byte(data), &rep)

	return
}

//QryStkInfo : 查詢商品報價資料
// http://mis.twse.com.tw/stock/api/getStockInfo.jsp?ex_ch=tse_2330.tw&json=1&delay=0
func (hdl *TwseStkHdl) QryStkInfo(sym string) (rep StockInfoResponse, err error) {
	var url string
	var data []byte
	var request *http.Request
	var response *http.Response

	hdl.Client = &http.Client{
		Timeout: time.Second * 5,
	}
	// get cookie
	hdl.getCookie()

	url = fmt.Sprintf("%s/%s?ex_ch=%s&json=1", twseAPIBase, twseAPIGetStockInfo, sym)
	fmt.Printf("url : %s\n", url)
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	for _, c := range hdl.cookies {
		request.AddCookie(c)
	}

	response, err = hdl.Client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	data, err = ioutil.ReadAll(response.Body)
	json.Unmarshal(data, &rep)

	return
}

//QryStkInfoBatch : 批次查尋商品資料
func (hdl *TwseStkHdl) QryStkInfoBatch(syms []string) (rep StockInfoResponse, err error) {
	var sym string
	for _, s := range syms {
		sym += s + "|"
	}
	rep, err = hdl.QryStkInfo(sym)
	return
}
