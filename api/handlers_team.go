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

package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	server "github.com/lirios/website/server"
)

// member represents a Slack team member.
type member struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	TzLabel  string `json:"tz_label"`
	Tz       string `json:"tz"`
	TzOffset int    `json:"tz_offset"`
	Profile  struct {
		Image24   string `json:"image_24"`
		Image32   string `json:"image_32"`
		Image48   string `json:"image_48"`
		Image72   string `json:"image_72"`
		Image192  string `json:"image_192"`
		Image512  string `json:"image_512"`
		Image1024 string `json:"image_1024"`
	} `json:"profile"`
	IsBot    bool   `json:"is_bot"`
	IsAdmin  bool   `json:"is_admin"`
	Deleted  bool   `json:"deleted"`
	Presence string `json:"presence,omitempty"`
}

// members is a list of Slack team members.
type members []member

// userListData is the content of Slack team members response.
type userListData struct {
	Ok      bool    `json:"ok"`
	Members members `json:"members"`
}

// filteredMember is a Member with filtered out information.
type filteredMember struct {
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Tz       string `json:"tz"`
	Image    string `json:"image"`
	Presence string `json:"presence,omitempty"`
}

// filteredMembers is a list of filtered out members.
type filteredMembers []filteredMember

// filteredUserListData is the response of our API service.
type filteredUserListData struct {
	Ok      bool            `json:"ok"`
	Members filteredMembers `json:"members"`
}

// Len returns the length of the slice.
func (slice members) Len() int {
	return len(slice)
}

// Less compares two slice items and returns true if index i should go before index j.
func (slice members) Less(i, j int) bool {
	return slice[i].IsAdmin && !slice[j].IsAdmin
}

// Swap swaps two slice items.
func (slice members) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Read text from URL.
func textFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// TeamHandler is a http handler for the team API.
func TeamHandler(c server.Context, w http.ResponseWriter, r *http.Request) (int, []byte) {
	w.Header().Set("Content-Type", "application/json")

	// For easier debugging - JavaScript won't accept json from another domain otherwise
	if strings.Contains(r.Host, "localhost") || strings.Contains(r.Host, "127.0.0.1") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}

	resp, err := textFromURL("https://slack.com/api/users.list?presence=1&token=" + c.Settings().Slack.Token)
	if err != nil {
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Parse to go object (all unnecessary info is ignored)
	data := userListData{}
	err = json.Unmarshal([]byte(resp), &data)
	if err != nil {
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Put administrators first
	sort.Sort(data.Members)

	// Parse the object back to json and print it
	result := filteredUserListData{Ok: data.Ok}
	for _, v := range data.Members {
		// Exclude deleted members and slackbot and filter out some information
		if v.ID != "USLACKBOT" && !v.Deleted {
			member := filteredMember{}
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
	if err != nil {
		return http.StatusBadRequest, []byte(err.Error())
	}
	return http.StatusOK, finalJSON
}
