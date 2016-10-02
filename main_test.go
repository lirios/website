/****************************************************************************
 * This file is part of Liri.
 *
 * Copyright (C) 2016 Pier Luigi Fiorini <pierluigi.fiorini@gmail.com>
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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTeam(t *testing.T) {
	// Configure the Web server.
	Config.Server.Port = ":8887"
	Config.Slack.Token = os.Getenv("SLACK_TOKEN")

	// New request to the API.
	req, err := http.NewRequest("GET", "/team", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(teamHandler)

	// Serve.
	handler.ServeHTTP(rr, req)

	// Verify the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	result := UserListData{}
	err = json.Unmarshal([]byte(rr.Body.String()), &result)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Ok {
		t.Errorf("handler returned unexpected Ok: got %v want %v",
			result.Ok, true)
	}
}
