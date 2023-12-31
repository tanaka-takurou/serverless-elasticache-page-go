package main

import (
	"os"
	"fmt"
	"log"
	"time"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/bradfitz/gomemcache/memcache"
)

type APIResponse struct {
	Message string `json:"message"`
}

type Response events.APIGatewayProxyResponse

var memcacheClient *memcache.Client

const itemKey string = "memcache_item_key"
const layout  string = "2006-01-02 15:04"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "get" :
			res, e := getItem(itemKey)
			if e != nil {
				err = e
			} else {
				jsonBytes, _ = json.Marshal(APIResponse{Message: res})
			}
		case "update" :
			res, e := updateItem(itemKey)
			if e != nil {
				err = e
			} else {
				jsonBytes, _ = json.Marshal(APIResponse{Message: res})
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: 500,
			Body: string(jsonBytes),
		}, nil
	}
	return Response {
		StatusCode: 200,
		Body: string(jsonBytes),
	}, nil
}

func getMemcacheClient() *memcache.Client {
	return memcache.New(os.Getenv("ADDRESS")  + ":" + os.Getenv("PORT"))
}

func getItem(key string)(string, error) {
	if memcacheClient == nil {
		memcacheClient = getMemcacheClient()
	}
	it, err := memcacheClient.Get(key)
	if err != nil {
		return "", err
	}
	return string(it.Value), nil
}

func setItem(key string, value string) {
	if memcacheClient == nil {
		memcacheClient = getMemcacheClient()
	}
	memcacheClient.Set(&memcache.Item{Key: itemKey, Value: []byte(value)})
	return
}

func updateItem(key string)(string, error) {
	t := time.Now()
	setItem(key, t.Format(layout))
	return getItem(key)
}

func main() {
	lambda.Start(HandleRequest)
}
