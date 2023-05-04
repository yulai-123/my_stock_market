package tushare

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"my_stock_market/config"
	"net/http"
	"strings"
	"time"
)

type RequestBody struct {
	APIName string                 `json:"api_name"`
	Token   string                 `json:"token"`
	Params  map[string]interface{} `json:"params"`
	Fields  string                 `json:"fields"`
}

type ResponseBody struct {
	Code int64                  `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	Timeout: 10 * time.Second,
}

func getParamsAndFields(params map[string]interface{}, fieldList []string) (map[string]interface{}, string, error) {
	fields := ""
	if len(fieldList) > 0 {
		fields += fieldList[0]
		for i := 1; i < len(fieldList); i++ {
			fields += "," + fieldList[i]
		}
	}

	return params, fields, nil
}

func getRequest(ctx context.Context, apiName string, params map[string]interface{}, fieldList []string) (*http.Request, error) {
	conf := config.GetTuShareConf(ctx)
	url := fmt.Sprintf("%v", conf.Host)

	reParams, fields, err := getParamsAndFields(params, fieldList)
	if err != nil {
		return nil, err
	}

	requestBody := RequestBody{
		APIName: apiName,
		Token:   conf.Token,
		Params:  reParams,
		Fields:  fields,
	}
	requestBodyStr, err := json.Marshal(requestBody)
	if err != nil {
		logrus.Errorf("[getRequest] marshal requestBody error: %v", err)
		return nil, err
	}
	s := string(requestBodyStr)
	reqBody := strings.NewReader(s)

	request, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		logrus.Errorf("[getRequest] NewRequest error: %v", err)
	}

	return request, nil
}

func getResponseData(ctx context.Context, response *http.Response) (map[string]interface{}, error) {
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err := fmt.Errorf("response status is not 200")
		logrus.Errorf("[getResponseData] reponse status is not 200")
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("[getResponseData] read response body error: %v", err)
		return nil, err
	}

	respBody := ResponseBody{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		logrus.Errorf("[getResponseData] unmarshal response body error: %v", err)
		return nil, err
	}
	if respBody.Code != 0 {
		logrus.Errorf("[getResponseData] response body code is not 0, code: %v, msg: %v", respBody.Code, respBody.Msg)
		err := fmt.Errorf("response body code is not 0")
		return nil, err
	}

	return respBody.Data, nil
}

func Post(ctx context.Context, apiName string, params map[string]interface{}, fieldList []string) (map[string]interface{}, error) {
	var response *http.Response
	var err error
	for i := 0; i < 5; i++ {
		request, err := getRequest(ctx, apiName, params, fieldList)
		if err != nil {
			return nil, err
		}
		response, err = client.Do(request)
		if err != nil {
			logrus.Errorf("[Post] client post error: %v, try: %v", err, i)
			time.Sleep(10 * time.Second)
		}
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	data, err := getResponseData(ctx, response)
	if err != nil {
		return nil, err
	}

	return data, nil
}
