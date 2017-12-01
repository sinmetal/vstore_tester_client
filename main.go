package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sinmetal/slog"
)

type ItemAPIPostRequest struct {
	Contents []string
}

type ItemAPIPostResponse struct {
	Key       string
	Contents  []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const vtServerURL = "http://vt-server-service.default.svc.cluster.local:8080"

func main() {
	for i := 0; i < 1000; i++ {
		Post()
	}
}

func Post() {
	log := slog.Start(time.Now())
	defer log.Flush()

	body := ItemAPIPostRequest{
		Contents: []string{"hello client"},
	}
	b, err := json.Marshal(body)
	if err != nil {
		log.Errorf("json.Marshal err = %s", err.Error())
	}

	client := new(http.Client)
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/item", vtServerURL),
		strings.NewReader(string(b)),
	)
	log.Infof("%s", string(b))

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("client.Do err = %s", err.Error())
		return
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("request.Body %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Errorf("response code = %d, body = %s", res.StatusCode, resBody)
		return
	}
}
