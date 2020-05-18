/*
 * Copyright 2020 zpxio (Jeff Sharpe)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package site

import (
	"bytes"
	"errors"
	"github.com/zpxio/mdsite/pkg/config"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const DefaultWeight = 10

type PageEntry struct {
	Id         uint64
	Path       string
	Extension  string
	Url        string
	Label      string
	ListWeight float64
	Modified   time.Time
}

var nextId uint64 = 1

func nextPageId() uint64 {
	i := nextId
	nextId++

	return i
}

const (
	WroteChar = iota
	WroteSpace
)

func generateLabel(path string) string {
	buf := bytes.Buffer{}

	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// State flags
	lastState := WroteSpace
	breakOnUpper := false

	for _, r := range name {
		if unicode.IsSpace(r) || strings.ContainsRune("-_", r) {
			if lastState != WroteSpace {
				buf.WriteRune(' ')
				lastState = WroteSpace
			}
		} else {
			rx := r
			if lastState == WroteSpace {
				rx = unicode.ToUpper(r)
				breakOnUpper = true
			} else if breakOnUpper && unicode.IsUpper(r) {
				buf.WriteRune(' ')
				rx = unicode.ToUpper(r)
				breakOnUpper = false
			} else if !breakOnUpper && unicode.IsLower(r) {
				breakOnUpper = true
			}

			buf.WriteRune(rx)
			lastState = WroteChar
		}
	}

	return buf.String()
}

func LoadPageEntry(path string) (*PageEntry, error) {
	id := nextPageId()
	fullPath := filepath.Join(config.Global().SitePath, path)
	fs, statErr := os.Stat(fullPath)
	if statErr != nil {
		return nil, statErr
	}

	// Check that nothing weird is happening with the file
	if fs.IsDir() {
		return nil, errors.New("file entry is actually a director")
	}

	// Generate a Url from the file path
	ext := filepath.Ext(path)
	url := "/" + strings.TrimSuffix(path, ext)
	// Fix the extension
	ext = strings.TrimPrefix(ext, ".")

	pe := PageEntry{
		Id:         id,
		Path:       path,
		Extension:  ext,
		Url:        url,
		Label:      generateLabel(path),
		Modified:   fs.ModTime(),
		ListWeight: DefaultWeight,
	}

	return &pe, nil
}
