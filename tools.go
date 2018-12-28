package tools

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const DATE_FORMAT = "2006-01-02 15:04:05"
const DATE_FORMAT_DAY = "2006-01-02"

var _httpClient = &http.Client{

	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func TypeOf(obj interface{}) {
	fmt.Println(obj, reflect.TypeOf(obj))
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("[E] ", err.Error())
	}
}

func GET(url string) (int, string, error) {

	targetReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 500, err.Error(), err
	}

	targetRes, err := _httpClient.Do(targetReq) //提交
	if err != nil {
		return 500, err.Error(), err
	}
	defer targetRes.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(targetRes.Body)
	if err != nil {
		return 500, err.Error(), err
	}

	return targetRes.StatusCode, string(body), err

}

func POST(url string, params string) (int, string, error) {

	targetReq, err := http.NewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		return 500, err.Error(), err
	}

	targetRes, err := _httpClient.Do(targetReq)
	if err != nil {
		return 500, err.Error(), err
	}
	defer targetRes.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(targetRes.Body)
	if err != nil {
		return 500, err.Error(), err
	}

	return targetRes.StatusCode, string(body), err
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
