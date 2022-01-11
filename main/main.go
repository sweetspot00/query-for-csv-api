package main

import (
	"fmt"
	"log"
	"net/http"

	"net/url"

	"../server_utils"
)

func handleQuery(w http.ResponseWriter, r *http.Request) {
	urlStr := r.RequestURI
	static1 := "/?query="
	static2 := "/favicon.ico"
	if urlStr[:8] != static1 && urlStr != static2 {
		fmt.Println(urlStr[:8])
		fmt.Print("bad URL\n")
	} else if urlStr == static2 {
		// skip
	} else {
		s := urlStr[8:]
		s, _ = url.QueryUnescape(s)
		// urlStr, _ = url.QueryUnescape(urlStr)
		// s, _ := ioutil.ReadAll(r.Body) //
		query := string(s)
		fmt.Printf("%s\n", query)
		if len(query) > 0 {
			data, col_name, err := server_utils.ReadCsv()
			if err != nil {
				fmt.Printf("read failed\n")
			} else {
				err := server_utils.Select_data(query, col_name, data)
				if err != nil {
					fmt.Printf("bad query statement\n")
					fmt.Fprintln(w, err.Error())
				} else {
					fmt.Printf("select succesfully!\n")
					fmt.Fprintln(w, "select successfully!")
				}
			}
		}
	}

}

func main() {
	http.HandleFunc("/", handleQuery)
	err := http.ListenAndServe(":9527", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
