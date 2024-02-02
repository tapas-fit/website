package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/meyskens/go-turnstile"
)

// TODO: referrer

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "%s <waitlist-path> <addr> <secretKey>\n", os.Args[0])
		os.Exit(2)
	}
	path := strings.TrimSuffix(os.Args[1], "/")
	addr := os.Args[2]
	secretKey := os.Args[3]

	// Secret key	Description
	// 1x0000000000000000000000000000000AA	success: true, ErrorCodes: []
	// 2x0000000000000000000000000000000AA	success: false, ErrorCodes: ["invalid-input-response"]
	// 3x0000000000000000000000000000000AA	success: false, ErrorCodes: ["timeout-or-duplicate"]
	ts := turnstile.New(secretKey)

	fmt.Println(time.Now(), "Starting waitlist")
	fmt.Println("will use path", path)
	fmt.Println("will listen on", addr)

	db, err := OpenDB(path)
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				handlePost(w, r, db, ts)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		},
	),
	)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, db *DB, ts *turnstile.Turnstile) {
	err := r.ParseMultipartForm(1024 * 1024)
	if err != nil {
		fmt.Println("POST", r.URL.Path, "unparseable form:", err)
		http.Error(w, "Couldn't parse form: "+err.Error(), http.StatusBadRequest)
		return
	}
	tsResponse := r.Form["cf-turnstile-response"]
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	resp, err := ts.Verify(tsResponse[0], ip)
	if err != nil {
		fmt.Println("POST", r.URL.Path, "error while verifying captcha: ", err)
		http.Error(w, "Couldn't verify captcha: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.Success {
		fmt.Println("POST", r.URL.Path, "verification success")
		db.add(r.Form)
		return
	}
	if len(resp.ErrorCodes) == 0 {
		fmt.Println("POST", r.URL.Path, "internal server error: verification failed without errorcode", resp)
		http.Error(w, "verification failed", http.StatusInternalServerError)
		return
	}
	if resp.ErrorCodes[0] == "invalid-input-response" {
		fmt.Println("POST", r.URL.Path, "verification failed with message: invalid-input-response")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// should only be timeout-or-duplicate
	fmt.Println("POST", r.URL.Path, "verification failed with message:", resp.ErrorCodes[0])
	http.Error(w, "verification failed", http.StatusInternalServerError)
}
