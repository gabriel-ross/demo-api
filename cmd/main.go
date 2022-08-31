package main

import (
	"airtonomy-assessment/server"
	"fmt"
	"net/http"
)

const PORT = "8080"

func main() {
	s := server.NewServer()
	fmt.Printf("Running server on port %s\n", PORT)
	http.ListenAndServe(":"+PORT, s)
	// resp, err := http.Get("https://api.hatchways.io/assessment/blog/posts?tag=tech")
	// if err != nil {
	// 	log.Fatal("error: ", err)
	// }
	// defer resp.Body.Close()
	// data, err := io.ReadAll(resp.Body)
	// fmt.Printf("raw data: %s\n", data)
	// if err != nil {
	// 	log.Fatal("error2: ", err)
	// }
	// respArr := server.APIResponse{}
	// err = json.Unmarshal(data, &respArr)
	// fmt.Printf("%v", respArr)
	// fmt.Printf("%s", reflect.TypeOf(respArr.Posts[0]))
}
