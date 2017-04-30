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

package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	api "github.com/lirios/website/api"
	server "github.com/lirios/website/server"
	"gopkg.in/gcfg.v1"
)

// Context of the application.
type ctx struct {
	settings *server.Settings
}

func (c ctx) Settings() *server.Settings {
	return c.settings
}

// Application handler.
type appHandler struct {
	*ctx
	handler func(server.Context, http.ResponseWriter, *http.Request) (int, []byte)
}

func (t appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code, data := t.handler(t.ctx, w, r)
	if code != http.StatusOK {
		http.Error(w, string(data), code)
		return
	}
	w.Write(data)
}

// Routes.
var routes = []struct {
	method  string
	route   string
	handler func(server.Context, http.ResponseWriter, *http.Request) (int, []byte)
}{
	{"GET", "/api/team", api.TeamHandler},
}

func main() {
	// Load settings
	var settingsFileName = "./config.ini"
	if len(os.Args) > 1 {
		settingsFileName = os.Args[1:][0]
	}
	var settings server.Settings
	err := gcfg.ReadFileInto(&settings, settingsFileName)
	if err != nil {
		panic(err)
	}

	// Create context
	appContext := &ctx{&settings}

	// Create router
	r := mux.NewRouter()

	// Add routes
	for _, detail := range routes {
		r.Handle(detail.route, appHandler{appContext, detail.handler}).Methods(detail.method)
	}

	// Serve
	http.ListenAndServe(settings.Server.Port, r)
}
