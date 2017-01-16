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

type Member struct {
	Id       string `json:"id"`
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
	Deleted  bool   `json:"deleted"`
	Presence string `json:"presence,omitempty"`
}

type Members []Member

func (slice Members) Len() int {
	return len(slice)
}

func (slice Members) Less(i, j int) bool {
	return slice[i].IsAdmin && !slice[j].IsAdmin
}

func (slice Members) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type UserListData struct {
	Ok      bool    `json:"ok"`
	Members Members `json:"members"`
}

type FilteredMember struct {
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Tz       string `json:"tz"`
	Profile  struct {
		Image192 string `json:"image_192"`
		Image512 string `json:"image_512"`
	} `json:"profile"`
	Presence string `json:"presence,omitempty"`
}

type FilteredMembers []FilteredMember

type FilteredUserListData struct {
	Ok      bool            `json:"ok"`
	Members FilteredMembers `json:"members"`
}
