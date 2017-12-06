package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	for {
		lot := fmt.Sprintf("%s-_-%s", time.Now().String(), uuid.New().String())
		for i := 0; i < 1000; i++ {
			i := i
			go func() {
				err := Post(lot, i)
				if err != nil {
					fmt.Println(err.Error())
				}
			}()
		}
		time.Sleep(1 * time.Minute)
	}
}

func Post(lot string, index int) error {
	log := slog.Start(time.Now())
	defer log.Flush()

	contents := []string{
		lot,
		fmt.Sprintf("%d", index),
		"hello client",
	}
	body := ItemAPIPostRequest{
		Contents: contents,
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
		return errors.Wrap(err, "client.Do err")
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("request.Body %s", err.Error())
		return errors.Wrap(err, "read request.Body")
	}

	if res.StatusCode != http.StatusOK {
		log.Errorf("response code = %d, body = %s", res.StatusCode, resBody)
	}

	lm := struct {
		Lot                string   `json:"lot"`
		Index              int      `json:"index"`
		Contents           []string `json:"contents"`
		ResponseStatusCode int      `json:"responseStatusCode"`
		ResponseBody       string   `json:"responseBody"`
	}{
		Lot:                lot,
		Index:              index,
		Contents:           contents,
		ResponseStatusCode: res.StatusCode,
		ResponseBody:       string(resBody),
	}
	logJson, err := json.Marshal(lm)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	log.Info(string(logJson))

	return nil
}
