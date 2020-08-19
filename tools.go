package tools

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
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
		log.Println("[E] ", err.Error())
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

type ResponseSeccess struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

type ResponsePage struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Count   int64       `json:"count"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


func PrintStruct(i interface{}) {
	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)
	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i)
	}
	fmt.Println("获取到数据:")
	for k, v := range kv {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func PrintBody(req *http.Request) {
	s, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(s))
}

// 内网地址
func GetNetworkIP() (string, error) {
	addrs, err := net.InterfaceAddrs();
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("valid local IP not found")
}

// 外网地址
func GetExternalIPAddr() (exip string, err error) {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	exip = string(bytes.TrimSpace(b))
	return
}
func GetContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}
