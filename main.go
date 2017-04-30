/****************************************************************************
 * This file is part of Liri.
 *
 * Copyright (C) 2017 Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
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
	"net/url"
	"os"
	"sort"
	"strings"

	"gopkg.in/gcfg.v1"
)

// Settings contains settings from a configuration file.
type Settings struct {
	Server struct {
		Port string
	}
	Slack struct {
		Token string
	}
}

// Config is a global configuration object.
var Config Settings

func main() {
	fillConfig()

	http.HandleFunc("/team", teamHandler)
	http.ListenAndServe(Config.Server.Port, nil)
}

func fillConfig() {
	// Default ini path
	var path = "./config.ini"

	// Check for config path agrument
	if len(os.Args) > 1 {
		// Get the args and ignore program path
		args := os.Args[1:]
		// Replace default path
		path = args[0]
	}

	err := gcfg.ReadFileInto(&Config, path)
	check(err)
}

func teamHandler(w http.ResponseWriter, r *http.Request) {
	// For easier debugging - JavaScript won't accept json from another domain otherwise
	if strings.Contains(r.Host, "localhost") || strings.Contains(r.Host, "127.0.0.1") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	resp := textFromURL("https://slack.com/api/users.list?presence=1&token=" + Config.Slack.Token)

	// Parse to go object (all unnecessary info is ignored)
	data := UserListData{}
	err := json.Unmarshal([]byte(resp), &data)
	check(err)

	// Put administrators first
	sort.Sort(data.Members)

	// Parse the object back to json and print it
	result := FilteredUserListData{Ok: data.Ok}
	for _, v := range data.Members {
		// Exclude deleted members and slackbot and filter out some information
		if v.Id != "USLACKBOT" && !v.Deleted {
			member := FilteredMember{}
			member.Name = v.Name
			member.RealName = v.RealName
			member.Tz = v.Tz
			u, err := url.QueryUnescape(v.Profile.Image512)
			if err == nil {
				member.Image = u
			} else {
				member.Image = v.Profile.Image512
			}
			member.Presence = v.Presence
			result.Members = append(result.Members, member)
		}
	}
	finalJSON, err := json.Marshal(result)
	check(err)
	fmt.Fprintf(w, string(finalJSON))
}

func textFromURL(url string) string {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return string(body)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
