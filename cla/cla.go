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

package cla

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
)

/*
 * Database schema example:
 *
 * [users]
 *   12345678
 *     id: 12345678
 *     name: "Mario Rossi"
 *     login: "mrossi"
 *     email: "mrossi@mrossi.it"
 * [users]
 *   234567890
 *     id: 234567890
 *     name: "Giuseppe Verdi"
 *     login: "gverdi"
 *     email: "gverdi@gverdi.it"
 * [companies]
 *   [cool-company]
 *     name: "Cool Company"
 *     url: "https://www.coolcompany.com"
 *     address: "Via Tal dei Tali, 1\n40100 Bologna\nItaly"
 *     code: "111333444666"
 *     representative: "Pinco Pallino"
 *     email: "pincopallino@coolcompany.com"
 * [usercompanies]
 *   [cool-company]
 *     - 234567890
 * [agreements]
 *   "individual-1.0": {
 *     title: "Individual Contributor License v1.0"
 *     url: "https://www.site.org/cla/individual-1.0.pdf"
 *     is_entity: false
 *   },
 *   "entity-1.0": {
 *     title: "Entity Contributor License v1.0"
 *     url: "https://www.site.org/cla/entity-1.0.pdf"
 *     is_entity: true
 *   }
 * [useragreements]
 *   [12345678]
 *     "individual-1.0": {
 *       when: "2009-11-10T23:00:00"
 *     }
 *   [234567890]
 *     "individual-1.0": {
 *       when: "2009-11-10T23:00:00"
 *     }
 *  [companyagreements]
 *    [cool-company]
 *      - "entity-1.0"
 */

// CLA is a structure that contains internal data for the CLA module.
type CLA struct {
	db *bolt.DB
}

// A GitHub user.
type user struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
	Email string `json:"email"`
}

// A company.
type company struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	Address        string `json:"address"`
	Code           string `json:"code"`
	Representative string `json:"representative"`
	Email          string `json:"email"`
}

// CLA text.
type agreement struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	IsEntity bool   `json:"is_entity"`
}

// User agreement or disagreement with a CLA.
type userAgreement struct {
	When time.Time `json:"when"`
}

// Open opens the CLA database.
func Open(path string) (*CLA, error) {
	// Open database
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	// Create object
	c := &CLA{db: db}

	// Create buckets
	c.createBuckets()

	// Save agreements
	a := &agreement{
		Title:    "Liri Individual Contributor License Agreement v1.0",
		URL:      "https://github.com/lirios/lirios/raw/documents/cla/Liri-Individual-1.0.pdf",
		IsEntity: false,
	}
	if err := c.SaveAgreement("individual-1.0", a); err != nil {
		return nil, err
	}
	a = &agreement{
		Title:    "Liri Entity Contributor License Agreement v1.0",
		URL:      "https://github.com/lirios/lirios/raw/documents/cla/Liri-Entity-1.0.pdf",
		IsEntity: true,
	}
	if err := c.SaveAgreement("entity-1.0", a); err != nil {
		return nil, err
	}

	return c, nil
}

// Close closes the CLA database.
func (c *CLA) Close() error {
	return c.db.Close()
}

// HasAgreement returns whether the agreement exists.
func (c *CLA) HasAgreement(name string) (bool, error) {
	found := false
	err := c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("agreements"))
		if bucket == nil {
			return errors.New("No \"agreements\" bucket found")
		}
		bucket.ForEach(func(k, v []byte) error {
			if bytes.Compare(k, []byte(name)) == 0 {
				found = true
				return nil
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return false, err
	}
	return found, nil
}

// ListAgreements returns a list of agreements names.
func (c *CLA) ListAgreements() ([]string, error) {
	names := []string{}
	err := c.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("agreements"))
		if bucket == nil {
			return errors.New("No \"agreements\" bucket found")
		}
		bucket.ForEach(func(k, v []byte) error {
			names = append(names, string(k))
			return nil
		})
		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return names, nil
}

// SaveAgreement saves agreement into the "agreements" bucket.
func (c *CLA) SaveAgreement(name string, a *agreement) error {
	// Start the transaction
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get bucket
	bucket := tx.Bucket([]byte("agreements"))
	if bucket == nil {
		return errors.New("No \"agreements\" bucket found")
	}

	// Save agreement
	if buf, err := json.Marshal(a); err != nil {
		return err
	} else if err := bucket.Put([]byte(name), buf); err != nil {
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// StoreUserAgreement saves user agreement to a CLA.
func (c *CLA) StoreUserAgreement(ghUser *github.User, name string) error {
	// Start the transaction
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save user
	if err = saveUser(ghUser, tx); err != nil {
		return err
	}

	// Save user agreement
	if err = saveUserAgreement(ghUser, name, tx); err != nil {
		return nil
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DiscardUserAgreement saves user agreement to a CLA.
func (c *CLA) DiscardUserAgreement(ghUser *github.User, name string) error {
	// Start the transaction
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save user
	if err = saveUser(ghUser, tx); err != nil {
		return err
	}

	// Delete user agreement
	if err = deleteUserAgreement(ghUser, name, tx); err != nil {
		return nil
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// HasAgreed returns whether the user has agreed to the CLA.
func (c *CLA) HasAgreed(userID int, name string) (bool, error) {
	var agreed = false
	err := c.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte("useragreements"))
		if root == nil {
			return errors.New("No \"useragreements\" bucket found")
		}

		keyName := []byte(strconv.FormatInt(int64(userID), 10))
		bucket := root.Bucket(keyName)
		if bucket == nil {
			return fmt.Errorf("No \"%s\" bucket found under \"useragreements\"", keyName)
		}

		bucket.ForEach(func(k, v []byte) error {
			if bytes.Compare(k, []byte(name)) == 0 {
				agreed = true
				return nil
			}
			return nil
		})

		return nil
	})
	return agreed, err
}

// Creates all buckets.
func (c *CLA) createBuckets() error {
	// Start the transaction
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create buckets
	if _, err := tx.CreateBucketIfNotExists([]byte("agreements")); err != nil {
		return err
	}
	if _, err := tx.CreateBucketIfNotExists([]byte("users")); err != nil {
		return err
	}
	if _, err := tx.CreateBucketIfNotExists([]byte("useragreements")); err != nil {
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Saves user into the "users" bucket.
func saveUser(ghUser *github.User, tx *bolt.Tx) error {
	// Get or create users bucket
	bucket := tx.Bucket([]byte("users"))
	if bucket == nil {
		return errors.New("No \"users\" bucket found")
	}

	// Save user
	user := &user{
		ID:    ghUser.GetID(),
		Name:  ghUser.GetName(),
		Login: ghUser.GetLogin(),
		Email: ghUser.GetEmail(),
	}
	keyName := []byte(strconv.FormatInt(int64(user.ID), 10))
	if buf, err := json.Marshal(user); err != nil {
		return err
	} else if err := bucket.Put(keyName, buf); err != nil {
		return err
	}

	return nil
}

// Saves user agreement into the "agreements" bucket.
func saveUserAgreement(ghUser *github.User, name string, tx *bolt.Tx) error {
	// Get or create root bucket
	root := tx.Bucket([]byte("useragreements"))
	if root == nil {
		return errors.New("No \"useragreements\" bucket found")
	}

	// Get or create bucket for the user
	keyName := []byte(strconv.FormatInt(int64(ghUser.GetID()), 10))
	bucket, err := root.CreateBucketIfNotExists(keyName)
	if err != nil {
		return err
	}
	a := &userAgreement{When: time.Now()}
	if buf, err := json.Marshal(a); err != nil {
		return err
	} else if err := bucket.Put([]byte(name), buf); err != nil {
		return err
	}

	return nil
}

// Deletes user agreement into the "agreements" bucket.
func deleteUserAgreement(ghUser *github.User, name string, tx *bolt.Tx) error {
	// Get or create root bucket
	root := tx.Bucket([]byte("useragreements"))
	if root == nil {
		return errors.New("No \"useragreements\" bucket found")
	}

	// Get or create bucket for the user
	keyName := []byte(strconv.FormatInt(int64(ghUser.GetID()), 10))
	bucket, err := root.CreateBucketIfNotExists(keyName)
	if err != nil {
		return err
	}
	bucket.Delete([]byte(name))

	return nil
}
