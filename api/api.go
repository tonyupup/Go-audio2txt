package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var token string
// var client http.Client

// func DefaultClient() *http.Client {
// 	return &client
// }

func init() {
	rsp, err := http.Get("https://openapi.baidu.com/oauth/2.0/token?grant_type=client_credentials&client_id=rgAdP6ltviEeB67uvB0yaPjf&client_secret=YVk80xf4QkfV9y4RVeA77Tv1maFm0Ye2")
	defer rsp.Body.Close()
	if err != nil {
		panic("init error.")
	}
	if body, err := ioutil.ReadAll(rsp.Body); err == nil {
		m := make(map[string]interface{})
		json.Unmarshal([]byte(body), &m)
		token = m["access_token"].(string)
	} else {
		panic("Unmarshal err")
	}
}
func Audio2StrHelper(id int) func([]byte) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://vop.baidu.com/server_api?dev_pid=%d&cuid=rgAdP6ltviEeB67uvB0yaPjf&token=%s", id, token)
	return func(data []byte) (map[string]interface{}, error) {
		req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
		req.Header.Set("Content-Type", "audio/pcm;rate=16000")
		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if err != nil {
			panic(err)
		}

		if 200 != resp.StatusCode {
			return nil, fmt.Errorf("HTTP Error:code :%d", resp.StatusCode)
		} else {
			data, _ = ioutil.ReadAll(resp.Body)
			m := make(map[string]interface{})
			err = json.Unmarshal(data, &m)
			if err != nil {
				panic(err)
			}
			return m, nil
		}
	}
}
