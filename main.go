package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type endPoint struct {
	Url string `json:"url"`
}

var hostMap map[string]endPoint

func main() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &hostMap)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(":8085", nil)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathSplits := strings.Split(path, "/")
	realPath := "/" + strings.Join(pathSplits[2:], "/")

	prefix := pathSplits[1]
	serviceConfig := hostMap[prefix]

	if serviceConfig.Url != "" {
		fullUrl := strings.TrimRight(serviceConfig.Url, "/") + realPath
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " - " + fullUrl)
		ForwardHandler(w, r, fullUrl)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404"))
	}
}

func ForwardHandler(writer http.ResponseWriter, request *http.Request, path string) {
	u, err := url.Parse(path)
	if nil != err {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " + error")
		fmt.Println(err)
		return
	}

	proxy := httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.URL = u
			request.Host = u.Host
		},
	}

	proxy.ServeHTTP(writer, request)
	defer setCORS(writer, request)
}

func setCORS(rw http.ResponseWriter, r *http.Request) {
	if rw.Header().Get("Access-Control-Allow-Origin") == "" {
		rw.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	}
}
