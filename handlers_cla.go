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
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
	server "github.com/lirios/website/server"
	oauth2 "golang.org/x/oauth2"
	oauth2_github "golang.org/x/oauth2/github"
)

func claApp(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/cla/app.html")
}

func claGithubAuth(c server.Context, w http.ResponseWriter, r *http.Request) (int, []byte) {
	// Configure OAuth2
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     c.Settings().CLA.ClientID,
		ClientSecret: c.Settings().CLA.ClientSecret,
		Scopes:       []string{"user", "repo"},
		Endpoint:     oauth2_github.Endpoint,
	}

	// Redirect URL
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	client.Get("...")

	return http.StatusTemporaryRedirect, []byte(url)
}

// CLAWebhookHandler is the GitHub webooks implementation.
func CLAWebhookHandler(c server.Context, w http.ResponseWriter, r *http.Request) (int, []byte) {
	// Parse request
	payload, err := github.ValidatePayload(r, []byte(c.Settings().CLA.HookSecret))
	if err != nil {
		return http.StatusBadRequest, []byte(err.Error())
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Handle event
	switch event := event.(type) {
	case *github.PingEvent:
		break
	case *github.PullRequestEvent:
		err = processPullRequest(c, event)
		if err != nil {
			return http.StatusBadRequest, []byte(err.Error())
		}
	default:
		return http.StatusBadRequest, nil
	}

	return http.StatusOK, []byte("OK")
}

// CLAListHandler is an API that lists all CLA available.
func CLAListHandler(c server.Context, w http.ResponseWriter, r *http.Request) (int, []byte) {
	return http.StatusOK, []byte("OK")
}

// CLAAgreeHandler is an API to agree with a CLA.
func CLAAgreeHandler(c server.Context, w http.ResponseWriter, r *http.Request) (int, []byte) {
	/*
		vars := mux.Vars(r)
			if err := c.CLA().StoreUserAgreement(ghUser, vars["name"]); err != nil {
				return http.StatusInternalServerError, nil
			}
	*/
	return http.StatusOK, []byte("OK")
}

// CLADisagreeHandler is an API to disagree with a CLA.
func CLADisagreeHandler(c server.Context, w http.ResponseWriter, r *http.Request) (int, []byte) {
	/*
		vars := mux.Vars(r)
			if err := c.CLA().DiscardUserAgreement(ghUser, vars["name"]); err != nil {
				return http.StatusInternalServerError, nil
			}
	*/
	return http.StatusOK, []byte("OK")
}

func processPullRequest(c server.Context, p *github.PullRequestEvent) error {
	// Log pull request
	var sender string
	if p.Sender != nil {
		sender = p.Sender.GetLogin()
	} else {
		sender = "unknown"
	}
	log.Printf("Pull request #%d %s from %s", p.GetNumber(), p.GetAction(), sender)

	// Authenticate
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Settings().CLA.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Prepare a status
	status := &github.RepoStatus{}

	// List status to find an existing one
	statuses, _, err := client.Repositories.ListStatuses(ctx, p.Repo.Owner.GetLogin(), p.Repo.GetName(), p.PullRequest.Head.GetSHA(), nil)
	if err != nil {
		log.Printf("Failed to list statuses for %v: %s", p.Repo.GetFullName(), err.Error())
		return err
	}
	for _, s := range statuses {
		if s.GetContext() == "license/cla" {
			status = s
			break
		}
	}

	// Determine whether the user has agreed at least to one CLA
	agreed := false
	claNames, err := c.CLA().ListAgreements()
	if err != nil {
		log.Printf("Failed to list agreements: %s", err.Error())
		return err
	}
	for _, claName := range claNames {
		if result, _ := c.CLA().HasAgreed(p.Sender.GetID(), claName); result {
			agreed = true
			break
		}
	}

	// Create or update status check
	targetURL := fmt.Sprintf("%s/cla/agree?utm_source=github_status&utm_medium=notification&repo=%s&pullrequest=%d", c.Settings().Server.SiteURL, p.Repo.GetFullName(), p.PullRequest.GetID())
	if len(claNames) == 0 {
		status.State = github.String("failure")
		status.Description = github.String("Organization is not properly configured")
	} else {
		if agreed {
			status.State = github.String("success")
			status.Description = github.String("Contributor License Agreement signed")
		} else {
			status.State = github.String("pending")
			status.Description = github.String("Contributor License Agreement is not signed yet")
		}
	}
	status.TargetURL = github.String(targetURL)
	status.Context = github.String("license/cla")
	status.Creator = p.Repo.Owner
	_, _, err = client.Repositories.CreateStatus(ctx, p.Repo.Owner.GetLogin(), p.Repo.GetName(), p.PullRequest.Head.GetSHA(), status)
	if err != nil {
		log.Printf("Failed to create status for %v: %s", p.Repo.GetFullName(), err.Error())
		return err
	}

	return nil
}
