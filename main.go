package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Websites struct {
	Name []string `json:"websites"`
}

var mp = make(map[string]string)

func dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the dashboard!")
}

func CheckStatus(w http.ResponseWriter, r *http.Request) {
	go Status()
	io.WriteString(w, "Status has been checked.\n")
}

func Status() {
	for {
		time.Sleep(5 * time.Second)
		var wg sync.WaitGroup

		for key := range mp {
			wg.Add(1)
			key := key

			go func() {
				defer wg.Done()
				resp, err := http.Get("https://" + key)

				if err != nil {
					mp[key] = "DOWN"
					fmt.Println("Url :", key, "status : DOWN")
					return
				}

				statusCode := resp.StatusCode
				if statusCode == 200 {
					mp[key] = "UP"
					fmt.Println("Url :", key, "status : 200 OK")
				} else {
					mp[key] = "DOWN"
					fmt.Println("Url :", key, "status : DOWN")
				}
			}()
		}
		wg.Wait()
	}
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Has("name") {
		w.WriteHeader(http.StatusOK)

		var m = make(map[string]string)
		m[r.URL.Query().Get("name")] = mp[r.URL.Query().Get("name")]

		err := json.NewEncoder(w).Encode(m)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(mp)
		if err != nil {
			fmt.Println(err)
		}
	}
	io.WriteString(w, "Status returned!\n")
}

func PostWebsites(w http.ResponseWriter, r *http.Request) {
	ws := Websites{}
	err := json.NewDecoder(r.Body).Decode(&ws)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ws)

	for _, website := range ws.Name {
		mp[website] = "DOWN"
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ws)
	if err != nil {
		fmt.Println(err)
	}

	io.WriteString(w, "Welcome to status checker website!\n")
}

func main() {
	fmt.Println("Server has started...")
	http.HandleFunc("/", dashboard)

	http.HandleFunc("/post", PostWebsites)
	http.HandleFunc("/Check", CheckStatus)
	http.HandleFunc("/websites", getStatus)
	http.ListenAndServe("127.0.0.1:3000", nil)
}

// curl -X POST -H "Content-Type: application/json" -d '{"websites":["www.google.com","www.bing.com","www.reallyfaketaxi.com","www.notarealwebsite.com","www.yahoo.com"]}' http://localhost:3000/post
// curl http://localhost:3000/Check
// curl "http://localhost:3000/websites?name=www.google.com"
