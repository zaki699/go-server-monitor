package server

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"visualon.com/go-server-monitor/config"
)

// PostHandler Writes a file
func PostMonitorHandler(r *http.Request) {

	var session_id int64 = -1
	var byteSlice []string = make([]string, 0)
	if r.URL.Query().Get("hash") == config.CONFIG.Secure.Hash {
		var err error
		if n, err := CheckifSessionAlreadyExist(*r.URL); err == nil {
			if n != -1 {
				session_id = n
			} else {
				// Session does not exist yet !
				if n, err = CreateSession(*r.URL); err == nil {
					session_id = n
				}
			}
		}

		reader := bufio.NewReader(r.Body)

		if err != nil {
			fmt.Println(err)
			r.Body.Close()
		}

		doThis := true
		for doThis {
			// Read the next chunk
			buf, _, err := reader.ReadLine()
			if err != nil {
				EndSession(session_id)
				fmt.Println("Session End")
				// Got an error back (e.g. EOF), so exit the loop
				doThis = false

			} else {
				if strings.HasSuffix(string(buf[:]), "continue") {
					byteSlice = append(byteSlice, string(buf[:]))
					Ingest(byteSlice, session_id)
					byteSlice = nil
				} else {
					byteSlice = append(byteSlice, string(buf[:]))

				}
			}
		}
	}
	r.Body.Close()
}

// PostHandler Writes a file
func PostLogHandler(rw http.ResponseWriter, r *http.Request) {

	var session_id int64

	if r.URL.Query().Get("hash") == config.CONFIG.Secure.Hash {

		reader := bufio.NewReader(r.Body)
		var err error
		if n, err := CheckifSessionAlreadyExist(*r.URL); err == nil {
			if n != -1 {
				session_id = n
			} else {
				// Session does not exist yet !
				if n, err = CreateSession(*r.URL); err == nil {
					session_id = n
				}
			}
		}

		if err != nil {
			fmt.Println(err)
			r.Body.Close()
		}

		doThis := true
		for doThis {
			// Read the next chunk
			buf, _, err := reader.ReadLine()
			if err != nil {
				doThis = false
			} else {
				value := string(buf[:])
				IngestLog(value, session_id)
			}
		}
	}
	r.Body.Close()
}
