package ljtop

// {"jsonrpc":"2.0","method":"homepage.get_rating","params":{"country":"cyr","full_text":"1","page":0,"pagesize":25},"id":4777887904}
// ratingURL = "http://l-api.livejournal.com/__api/?request=" + Uri.encode(jsonObject.toString());

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	RatingRequest struct {
		RPCVersion string `json:"jsonrpc"`
		Method     string `json:"method"`
		Id         int64  `json:"id"`
		Params     struct {
			Country    string `json:"country"`
			IsFullText string `json:"full_text"`
			Page       int    `json:"page"`
			PageSize   int    `json:"pagesize"`
		} `json:"params"`
	}

	RatingResponse struct {
    Result  struct {
      Rating  []struct {
        PostUrl string `json:"post_url"`
      } `json:"rating"`
    } `json:"result"`
	}
)

func GetLJTop(country string) RatingResponse {

	rating_req := RatingRequest{}
	rating_req.RPCVersion = "2.0"
	rating_req.Method = "homepage.get_rating"
	rating_req.Params.Country = country
	rating_req.Params.IsFullText = "0"
	rating_req.Params.Page = 0
	rating_req.Params.PageSize = 3
	rating_req.Id = time.Now().Unix()

	b, err := json.Marshal(rating_req)
	if err != nil {
		fmt.Println(err)
		return RatingResponse{}
	}

	// fmt.Println(url.QueryEscape(string(b)))
	res, err := http.Get("http://l-api.livejournal.com/__api/?request=" + url.QueryEscape(string(b)))
	if err != nil {
		fmt.Println(err)
		return RatingResponse{}
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return RatingResponse{}
	}

  // fmt.Println(string(body))

	rating_res := RatingResponse{}
	err = json.Unmarshal(body, &rating_res)
	if err != nil {
		fmt.Println(err)
		return RatingResponse{}
	}

	return rating_res
}
