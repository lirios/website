package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
    "strings"
)

var token string

func main() {
	var path string = "./token.txt"

	if len(os.Args) > 1{
		//get the args and ignore program path
		args := os.Args[1:]

		//replace default path
		path = args[0]
	}

	//read token from file
	file := readFile(path)
	//extract first line to token
	token = file[0:strings.Index(file, "\n")]

    http.HandleFunc("/team", teamHandler)
    http.ListenAndServe(":8887", nil)
}

func teamHandler(w http.ResponseWriter, r *http.Request) {
	resp := textFromUrl("https://slack.com/api/users.list?presence=1&token="+token)

	//parse to go object (all unnecessary info is ignored)
	result := UserListData{}
	err := json.Unmarshal([]byte(resp), &result)
	check(err)

	//parse the object back to json and print it
	final_json, err := json.Marshal(result)
	check(err)
    fmt.Fprintf(w, string(final_json))
}

func textFromUrl(url string) string {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
    return string(body)
}

func readFile(path string) string{
	dat, err := ioutil.ReadFile(path)
    check(err)
    return string(dat)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type UserListData struct {
	Ok bool `json:"ok"`
	Members []struct {
		Name string `json:"name"`
		RealName string `json:"real_name"`
		TzLabel string `json:"tz_label"`
		Profile struct {
			Image192 string `json:"image_192"`
			Image512 string `json:"image_512"`
		} `json:"profile"`
		IsBot bool `json:"is_bot"`
		Presence string `json:"presence,omitempty"`
	} `json:"members"`
}