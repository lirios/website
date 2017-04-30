/****************************************************************************
 * This file is part of Liri.
 *
 * Copyright (C) 2017 Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
 * Copyright (C) 2017 Ziga Patacko Koderman <ziga.patacko@gmail.com>
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

// Member represents a Slack team member.
type Member struct {
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

// Members is a list of Slack team members.
type Members []Member

// UserListData is the content of Slack team members response.
type UserListData struct {
	Ok      bool    `json:"ok"`
	Members Members `json:"members"`
}

// FilteredMember is a Member with filtered out information.
type FilteredMember struct {
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Tz       string `json:"tz"`
	Image    string `json:"image"`
	Presence string `json:"presence,omitempty"`
}

// FilteredMembers is a list of filtered out members.
type FilteredMembers []FilteredMember

// FilteredUserListData is the response of our API service.
type FilteredUserListData struct {
	Ok      bool            `json:"ok"`
	Members FilteredMembers `json:"members"`
}

// Len returns the length of the slice.
func (slice Members) Len() int {
	return len(slice)
}

// Less compares two slice items and returns true if index i should go before index j.
func (slice Members) Less(i, j int) bool {
	return slice[i].IsAdmin && !slice[j].IsAdmin
}

// Swap swaps two slice items.
func (slice Members) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
