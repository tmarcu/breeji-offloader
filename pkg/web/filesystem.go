// Copyright 2022 Tudor Marcu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"net/http"
	"strings"
)

// NeuteredFileSystem sets up a FileSystem to use with http.FilerServer
type NeuteredFileSystem struct {
	FileSystem http.FileSystem
}

// Open satisfies the filesystem interface
// Directory/filesystem "neutering" technique from https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.FileSystem.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.FileSystem.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
