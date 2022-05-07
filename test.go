package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var ttrss_api_endpoint string
var session_id string

func main() {
	user := os.Args[1]
	password := os.Args[2]
	ttrss_api_endpoint = os.Args[3]
	session_id = login(user, password)
	isLoggedIn := isLoggedIn(session_id)
	if isLoggedIn {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

func requestApi(values map[string]string) (responseBody []byte) {
	request_data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(ttrss_api_endpoint, "application/json", bytes.NewBuffer(request_data))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	responseBody, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func login(user string, password string) (sessionId string) {
	values := map[string]string{"op": "login", "user": user, "password": password}
	body := requestApi(values)

	loginResponse := LoginResponse{}
	err := json.Unmarshal(body, &loginResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	sessionId = loginResponse.Content.SessionID
	return
}

func isLoggedIn(sid string) (isLoggedIn bool) {
	values := map[string]string{"op": "isLoggedIn", "sid": session_id}
	body := requestApi(values)

	logInfo := LogInfo{}
	err := json.Unmarshal(body, &logInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	isLoggedIn = logInfo.Content.Status
	return
}

type LoginResponse struct {
	Seq     int `json:"seq"`
	Status  int `json:"status"`
	Content struct {
		SessionID string `json:"session_id"`
		Config    struct {
			IconsDir        string   `json:"icons_dir"`
			IconsURL        string   `json:"icons_url"`
			DaemonIsRunning bool     `json:"daemon_is_running"`
			CustomSortTypes []string `json:"custom_sort_types"`
			NumFeeds        int      `json:"num_feeds"`
		} `json:"config"`
		APILevel int `json:"api_level"`
	} `json:"content"`
}

type LogInfo struct {
	Seq     int `json:"seq"`
	Status  int `json:"status"`
	Content struct {
		Status bool `json:"status"`
	} `json:"content"`
}
