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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

// Uploader is our handle to the server instance
type Uploader struct {
	CurrentFile int    `json:"currentFile"`
	PathPrefix  string `json:"pathPrefix"`
}

// NewUploader returns an uploader object with a specific prefix to save files to
func NewUploader(pathPrefix string) Uploader {
	basedir, _ := os.Open(filepath.Join(pathPrefix, "uploads/", "mobile"))
	files, err := basedir.Readdirnames(0)
	latestfile := 0
	if err == nil && len(files) != 0 {
		sort.SliceStable(files, func(i, j int) bool {
			base := files[i][:strings.LastIndex(files[i], ".")]

			base2 := files[j][:strings.LastIndex(files[j], ".")]

			return base < base2
		})
		files = strings.Split(files[len(files)-1], "_")
		cleanedString := files[len(files)-1]
		base := cleanedString[:strings.LastIndex(cleanedString, ".")]
		latestfile, err = strconv.Atoi(base)
		if err != nil {
			latestfile = 0
		}
		latestfile++
	}

	return Uploader{PathPrefix: pathPrefix, CurrentFile: latestfile}
}

// RemoveSpaces efficiently trims spaces out of strings
func RemoveSpaces(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// ReceiveHandler takes a filepond upload POST and copies the file to the
// filesystem
func (u *Uploader) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the file POST from our FilePond plugin
	if r.Method != http.MethodPost {
		return
	}
	file, header, err := r.FormFile("filepond")
	if err != nil {
		_, _ = fmt.Fprintln(w, err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	cleanedString := RemoveSpaces(strings.ToLower(header.Filename))
	extension := cleanedString[strings.LastIndex(cleanedString, "."):]
	// Write to path, create folder(s) if it does not exist already
	// TODO: Make this a configurable path based on the user logged in
	userFolder := "mobile"
	filename := fmt.Sprintf("%04d", u.CurrentFile)
	basedir := filepath.Join(u.PathPrefix, "uploads/", userFolder)
	writepath := filepath.Join(basedir, filename+extension)
	if _, err = os.Stat(basedir); err != nil {
		if os.IsNotExist(err) {
			if mkerr := os.MkdirAll(basedir, 0755); mkerr != nil {
				WriteHTTPMessage(w, http.StatusInternalServerError, mkerr.Error())
			}
		}
	}
	out, err := os.Create(writepath)
	if err != nil {
		WriteHTTPMessage(w, http.StatusForbidden, "Unable to create the file for writing. Check your write access privilege")
		return
	}
	defer func() {
		_ = out.Close()
	}()

	// Write te file contents out to filesystem
	_, err = io.Copy(out, file)
	if err != nil {
		_, err = fmt.Fprintln(w, err)
		if err != nil {
			log.Error(err)
		}
	}
	log.Infof("Wrote %s", writepath)
	_, _ = fmt.Fprintf(w, "File uploaded successfully: ")
	_, _ = fmt.Fprint(w, filename)
	u.CurrentFile++
}
