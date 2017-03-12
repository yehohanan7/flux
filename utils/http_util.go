package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

func GetAbsoluteUrl(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host + r.URL.Path
}

func HttpGetJson(url string, container interface{}) error {
	var body []byte
	fmt.Println("fetching ", url)
	res, err := http.Get(url)
	if err != nil {
		glog.Error("Error while getting ", err)
		return err
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Error("Error while reading the data ", err)
		return err
	}

	fmt.Println("unmarshalling ... ", url)
	err = json.Unmarshal(body, container)
	fmt.Println("unmarshalled ", url)

	if err != nil {
		glog.Error("Error while decoding data ", err)
		return err
	}

	return nil
}
