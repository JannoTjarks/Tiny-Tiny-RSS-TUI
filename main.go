package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rivo/tview"
)

var sessionId string
var config Config

func main() {
	readConfig()

	sessionId = login(config.Username, config.Password)

	apiLevel := getApiLevel()
	categories := getCategories()
	fmt.Println(apiLevel)
	fmt.Println(categories)

	InitTvViewList(categories)
}

func InitTvView() {
	box := tview.NewBox().SetBorder(true).SetTitle("Tiny Tiny RSS TUI")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}

func InitTvViewList(categories []Category) {
	list := tview.NewList()
	var runeForList = 'a'

	for _, category := range categories {
		list.AddItem(category.Title, "", runeForList, nil)
		runeForList++
		// Makes sure, that q is not used as rune for a categoriy
		if runeForList == 'q' {
			runeForList++
		}
	}
	app := tview.NewApplication()

	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func readConfig() {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	d := json.NewDecoder(file)
	err = d.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
}

func requestApi(values map[string]string) (responseBody []byte) {
	request_data, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(config.Ttrss_Api_Endpoint, "application/json", bytes.NewBuffer(request_data))
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Println(config.Ttrss_Api_Endpoint + " was not found!")
			fmt.Println("There are two things to check:")
			fmt.Println("\t- Check if the endpoint url is correct")
			fmt.Println("\t- Check if you are using the correct dns servers")
			fmt.Println()
		} else if strings.Contains(err.Error(), "dial tcp: lookup") {
			fmt.Println(config.Ttrss_Api_Endpoint + " was not resolveable!")
			fmt.Println("There are two things to check:")
			fmt.Println("\t- Check the dns configuration on your client")
			fmt.Println("\t- Check if you have access to the internet")
			fmt.Println()
		}

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

func isLoggedIn() (isLoggedIn bool) {
	values := map[string]string{"op": "isLoggedIn", "sid": sessionId}
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

func getApiLevel() (currentApiLevel int) {
	if !isLoggedIn() {
		login(config.Username, config.Password)
	}

	values := map[string]string{"op": "getApiLevel", "sid": sessionId}

	body := requestApi(values)

	apiLevel := ApiLevel{}
	err := json.Unmarshal(body, &apiLevel)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentApiLevel = apiLevel.Content.Level
	return
}

func getCategories() (currentCategories []Category) {
	if !isLoggedIn() {
		login(config.Username, config.Password)
	}

	values := map[string]string{"op": "getCategories", "sid": sessionId}

	body := requestApi(values)

	categories := Categories{}
	err := json.Unmarshal(body, &categories)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentCategories = categories.Content
	return
}

type Config struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	Ttrss_Api_Endpoint string `json:"ttrss_api_endpoint"`
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

type ApiLevel struct {
	Seq     int `json:"seq"`
	Status  int `json:"status"`
	Content struct {
		Level int `json:"level"`
	} `json:"content"`
}

type Categories struct {
	Seq     int        `json:"seq"`
	Status  int        `json:"status"`
	Content []Category `json:"content"`
}
type Category struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Unread  int    `json:"unread"`
	OrderID int    `json:"order_id,omitempty"`
}
