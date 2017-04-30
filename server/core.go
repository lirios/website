/****************************************************************************
 * This file is part of Liri.
 *
 * Copyright (C) 2017 Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
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

package server

import (
	cla "github.com/lirios/website/cla"
)

// Settings contains settings from a configuration file.
type Settings struct {
	Server struct {
		Port    string
		BaseURL string
		SiteURL string
	}
	Slack struct {
		Token string
	}
	CLA struct {
		DatabasePath string
		Token        string
		HookSecret   string
		ClientID     string
		ClientSecret string
	}
}

// Context interface.
type Context interface {
	Settings() *Settings
	CLA() *cla.CLA
}