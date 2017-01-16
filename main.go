/****************************************************************************
 * This file is part of Liri.
 *
 * Copyright (C) 2016 Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
 * Copyright (C) 2016 Ziga Patacko Koderman <ziga.patacko@gmail.com>
 *
 * $BEGIN_LICENSE:AGPL3+$
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * $END_LICENSE$
 ***************************************************************************/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"gopkg.in/gcfg.v1"
)

// Global configuration object.
var Config Settings

func main() {
	fillConfig()

	http.HandleFunc("/team", teamHandler)
	http.ListenAndServe(Config.Server.Port, nil)
}

func fillConfig() {
	//default ini path
	var path string = "./config.ini"

	//check for config path agrument
	if len(os.Args) > 1 {
		//get the args and ignore program path
		args := os.Args[1:]
		//replace default path
		path = args[0]
	}

	err := gcfg.ReadFileInto(&Config, path)
	check(err)

}

func teamHandler(w http.ResponseWriter, r *http.Request) {
	//for easyer debugging - js won't accept json from another domain otherwise
	if strings.Contains(r.Host, "localhost") || strings.Contains(r.Host, "127.0.0.1") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	resp := textFromUrl("https://slack.com/api/users.list?presence=1&token=" + Config.Slack.Token)

	//parse to go object (all unnecessary info is ignored)
	result := UserListData{}
	err := json.Unmarshal([]byte(resp), &result)
	check(err)

	// Put administrators first
	sort.Sort(result.Members)

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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Member struct {
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	TzLabel  string `json:"tz_label"`
	Tz       string `json:"tz"`
	TzOffset int    `json:"tz_offset"`
	Profile  struct {
		Image192 string `json:"image_192"`
		Image512 string `json:"image_512"`
	} `json:"profile"`
	IsBot    bool   `json:"is_bot"`
	IsAdmin  bool   `json:"is_admin"`
	Presence string `json:"presence,omitempty"`
}

type Members []Member

func (slice Members) Len() int {
	return len(slice)
}

func (slice Members) Less(i, j int) bool {
	return (slice[i].IsAdmin && !slice[j].IsAdmin) || slice[i].Name < slice[j].Name
}

func (slice Members) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type UserListData struct {
	Ok      bool    `json:"ok"`
	Members Members `json:"members"`
}

// Represents settings file.
type Settings struct {
	Server struct {
		Port string
	}
	Slack struct {
		Token string
	}
}
