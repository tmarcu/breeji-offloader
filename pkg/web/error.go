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
	"io"
	"net/http"
)

// WriteHTTPMessage sets the header and returns a message for the failure
func WriteHTTPMessage(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_, _ = io.WriteString(w, msg)
}
