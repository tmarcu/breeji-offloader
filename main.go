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

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"

	"github.com/tmarcu/breeji-offloader/pkg/web"
)

var (
	Version    = "0.0.1"
	pathPrefix = "./"
	//go:embed assets
	content embed.FS
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "breeji-offloader version %s\n Copyright © 2022 Tudor Marcu\n", Version)
	if len(os.Args) > 1 {
		pathPrefix = os.Args[1]
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	file, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()
	log.SetOutput(file)

	// Page routes
	fsys := fs.FS(content)
	html, _ := fs.Sub(fsys, "assets")
	http.Handle("/", http.FileServer(http.FS(html)))
	uploader := web.NewUploader(pathPrefix)

	http.HandleFunc("/upload", uploader.ReceiveHandler)

	if err != nil {
		log.Fatalf(err.Error())
	}
	ip := web.MachineIP()
	err = qrcode.WriteFile("http://"+ip, qrcode.Medium, 256, "qr.png")
	if err != nil {
		log.Fatal(err)
	}

	_ = open("qr.png")

	if err = http.ListenAndServe(ip, nil); err != nil {
		_ = fmt.Errorf("server crashed with error: %w", err)
	}
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
