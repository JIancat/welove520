package welove

import (
	bb "bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const KEY = "8b5b6eca8a9d1d1f"

func TreePost(accessToken, appKey string, op int) (*http.Response, error) {
	u := "http://api.welove520.com/v1/game/tree/op"
	sigEncoder := NewSig([]byte(KEY))
	d1 := Data{"access_token", accessToken}
	d2 := Data{"app_key", appKey}
	d3 := Data{"op", strconv.Itoa(op)}
	sig := sigEncoder.Encode("POST", u, d1, d2, d3)

	data := make(url.Values)
	data.Add("access_token", accessToken)
	data.Add("app_key", appKey)
	data.Add("op", strconv.Itoa(op))
	data.Add("sig", sig)
	res, err := http.PostForm(u, data)
	return res, err
}

func HomePost(accessToken string, taskType int, loveSpaceId string) (*http.Response, error) {
	u := "http://api.welove520.com/v1/game/house/task"
	sigEncoder := NewSig([]byte(KEY))
	d1 := Data{"access_token", accessToken}
	d2 := Data{"love_space_id", loveSpaceId}
	d3 := Data{"task_type", strconv.Itoa(taskType)}
	sig := sigEncoder.Encode("POST", u, d1, d2, d3)

	data := make(url.Values)
	data.Add("access_token", accessToken)
	data.Add("task_type", strconv.Itoa(taskType))
	data.Add("love_space_id", loveSpaceId)
	data.Add("sig", sig)
	res, err := http.PostForm(u, data)
	return res, err
}

func RandomHouse(accessToken string) (string, bool) {
	var u = "http://api.welove520.com/v1/game/house/info"
	sigEncoder := NewSig([]byte(KEY))
	d1 := Data{"access_token", accessToken}
	d2 := Data{"love_space_id", "random"}
	sig := sigEncoder.Encode("POST", u, d1, d2)

	values := make(url.Values)
	values.Add("access_token", accessToken)
	values.Add("love_space_id", "random")
	values.Add("sig", sig)
	res, err := http.PostForm(u, values)
	if err != nil {
		log.Fatal(err)
	}
	bytes, _ := ioutil.ReadAll(res.Body)
	js, err := simplejson.NewJson(bytes)
	if err != nil {
		log.Fatal(err)
	}
	arr, err := js.Get("messages").Array()
	house, ok := arr[0].(map[string]interface{})["house"].(map[string]interface{})
	if ok {
		id, ok := house["love_space_id"].(string)
		return id, ok
	} else {
		return "", ok
	}
}

func Visit(accessToken, loveSpaceId string) (*http.Response, error) {
	u := "http://api.welove520.com/v1/game/house/task"

	d1 := Data{"task_type", "8"}
	d2 := Data{"house_num", "0"}
	d3 := Data{"access_token", accessToken}
	d4 := Data{"love_space_id", loveSpaceId}
	sigEncoder := NewSig([]byte(KEY))
	sig := sigEncoder.Encode("POST", u, d3, d2, d4, d1)

	values := make(url.Values)
	values.Add("task_type", "8")
	values.Add("house_num", "0")
	values.Add("access_token", accessToken)
	values.Add("love_space_id", loveSpaceId)
	values.Add("sig", sig)
	res, err := http.PostForm(u, values)
	return res, err
}

type QueryItem struct {
	Result   int `json:"result"`
	Messages []struct {
		OpTime  int64 `json:"op_time"`
		MsgType int `json:"msg_type"`
		AdItems []struct {
			ItemID        int `json:"item_id"`
			Count         int `json:"count"`
			OpTime        int64 `json:"op_time"`
			NeedHelp      int `json:"need_help"`
			SellerFarmID  string `json:"seller_farm_id"`
			HeadURLFamale string `json:"head_url_famale"`
			HeadURLMale   string `json:"head_url_male"`
			ID            int `json:"id"`
			FarmName      string `json:"farm_name"`
			Coin          int `json:"coin"`
		} `json:"ad_items"`
	} `json:"messages"`
}

func QueryItems(accessToken string) QueryItem {
	u := "http://api.welove520.com/v1/game/farm/ad/query"
	d1 := Data{"access_token", accessToken}
	sigEncoder := NewSig([]byte(KEY))
	sig := sigEncoder.Encode("POST", u, d1)
	data := make(url.Values)
	data.Add("access_token", accessToken)
	data.Add("sig", sig)
	res, err := http.PostForm(u, data)
	if err != nil {
		log.Fatal(err)
	}
	bytes, _ := ioutil.ReadAll(res.Body)
	queryItem := QueryItem{}
	json.Unmarshal(bytes, &queryItem)
	return queryItem
}

type BuyItemStatus struct {
	Result   int `json:"result"`
	Messages []struct {
		StallItem  struct {
				   BuyerHeadURL  string `json:"buyer_head_url"`
				   BuyerFarmName string `json:"buyer_farm_name"`
				   ID            int `json:"id"`
			   } `json:"stall_item,omitempty"`
		OpTime     int64 `json:"op_time"`
		MsgType    int `json:"msg_type"`
		Warehouses []struct {
			Category int `json:"category"`
			ItemsInc []struct {
				ItemID int `json:"item_id"`
				Count  int `json:"count"`
			} `json:"items_inc"`
		} `json:"warehouses,omitempty"`
		FarmID     string `json:"farm_id,omitempty"`
		GoldCost   int `json:"gold_cost,omitempty"`
	} `json:"messages"`
}

func BuyItem(accessToken, sellerFarmId string, stallSaleId int) BuyItemStatus {
	u := "http://api.welove520.com/v1/game/farm/stall/buy"
	d1 := Data{"access_token", accessToken}
	d2 := Data{"seller_farm_id", sellerFarmId}
	d3 := Data{"stall_sale_id", strconv.Itoa(stallSaleId)}
	sigEncoder := NewSig([]byte(KEY))
	sig := sigEncoder.Encode("POST", u, d1, d2, d3)

	data := make(url.Values)
	data.Add("access_token", accessToken)
	data.Add("seller_farm_id", sellerFarmId)
	data.Add("stall_sale_id", strconv.Itoa(stallSaleId))
	data.Add("sig", sig)
	res, err := http.PostForm(u, data)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)
	buyItemStatus := BuyItemStatus{}
	json.Unmarshal(bytes, &buyItemStatus)
	return buyItemStatus
}

func GetLoveSpaceIdRaw(accessToken, appKey string) (*http.Response, error) {
	u := "http://api.welove520.com/v5/useremotion/getone"
	d1 := Data{"access_token", accessToken}
	d2 := Data{"app_key", appKey}
	d3 := Data{"user_id", "0"}
	sigEncoder := NewSig([]byte(KEY))
	sig := sigEncoder.Encode("POST", u, d1, d2, d3)

	data := make(url.Values)
	data.Add("access_token", accessToken)
	data.Add("app_key", appKey)
	data.Add("user_id", "0")
	data.Add("sig", sig)
	res, err := http.PostForm(u, data)
	return res, err
}

func GetLoveSpaceId(body string) string {
	js, _ := simplejson.NewJson([]byte(body))
	loveSpaceId, _ := js.Get("love_space_id").Float64()
	return strconv.Itoa(int(loveSpaceId))
}

type PetStatus struct {
	Result   int `json:"result"`
	Messages []struct {
		MsgType int `json:"msg_type"`
		Pets    []struct {
			PetID    int `json:"pet_id"`
			PetTasks []struct {
				Count      int `json:"count"`
				TaskType   int `json:"task_type"`
				RemainTime int `json:"remain_time"`
			} `json:"pet_tasks"`
		} `json:"pets,omitempty"`
		Count   int `json:"count,omitempty"`
	} `json:"messages"`
}

func GetPetStatus(accessToken string) PetStatus {
	u := "http://api.welove520.com/v1/game/house/pet/task/list"
	sigEncoder := NewSig([]byte(KEY))
	d1 := Data{"access_token", accessToken}
	sig := sigEncoder.Encode("POST", u, d1)

	data := make(url.Values)
	data.Add("access_token", accessToken)
	data.Add("sig", sig)
	res, err := http.PostForm(u, data)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	bytes, _ := ioutil.ReadAll(res.Body)
	pet := PetStatus{}
	err = json.Unmarshal(bytes, &pet)
	if err != nil {
		log.Fatal(err)
	}
	return pet
}

type PetTaskResult struct {
	Messages []struct {
		Count      int `json:"count"`
		MsgType    int `json:"msg_type"`
		PetID      int `json:"pet_id"`
		RemainTime int `json:"remain_time"`
		TaskType   int `json:"task_type"`
	} `json:"messages"`
	Result   int    `json:"result"`
	ErrorMsg string `json:"error_msg"`
}

func DoPetTask(accessToken, petId, taskType string) PetTaskResult {
	u := "http://api.welove520.com/v1/game/house/pet/task/do"
	sigEncoder := NewSig([]byte(KEY))
	d1 := Data{"access_token", accessToken}
	d2 := Data{"pet_id", petId}
	d3 := Data{"task_type", taskType}
	sig := sigEncoder.Encode("POST", u, d1, d2, d3)

	data := make(url.Values)
	data.Add("access_token", accessToken)
	data.Add("pet_id", petId)
	data.Add("sig", sig)
	data.Add("task_type", taskType)

	res, err := http.PostForm(u, data)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)
	result := PetTaskResult{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func ServerRun(path, port string) {
	log.Printf("请将手机Http代理设置为[本机IP%s]\n", port)
	go contentHandler(path)
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.OnRequest().DoFunc(httpHandler)
	log.Fatal(http.ListenAndServe(port, proxy))
}

var sChan = make(chan string)

type Love struct {
	AccessToken string `json:"access_token"`
	AppKey      string `json:"app_key"`
	TaskType    []int  `json:"task_type"`
}

func contentHandler(path string) {
	var f, _ = os.OpenFile(path, os.O_CREATE | os.O_RDWR, os.ModeAppend)
	defer f.Close()
	for v := range sChan {
		accessToken, _ := getValue(v, "access_token")
		appKey, _ := getValue(v, "app_key")
		love := Love{}
		love.AccessToken = accessToken
		love.AppKey = appKey
		love.TaskType = []int{1, 4, 5, 6, 7, 11}
		bytes, _ := json.Marshal(love)
		f.Write(bytes)
		fmt.Println("生成配置文件完毕：" + path)
		os.Exit(0)
	}
}

func getValue(content, key string) (string, error) {
	r := "(?:" + key + ")=(.+?)(&|$)"
	reg, err := regexp.Compile(r)
	return reg.FindAllStringSubmatch(content, -1)[0][1], err
}

func httpHandler(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	if r.Method != "POST" {
		return r, nil
	}
	if r.Host != "api.welove520.com" {
		return r, nil
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	ori := ioutil.NopCloser(bb.NewBuffer(buf))
	r.Body = ori
	bytesString := bb.Buffer{}
	bytesString.Write(buf)
	content := bytesString.String()
	if strings.Contains(content, "access_token") && strings.Contains(content, "app_key") {
		dContent, err := url.QueryUnescape(content)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Decode [%s] to [%s]\n", content, dContent)
		sChan <- dContent
	}
	return r, nil
}
