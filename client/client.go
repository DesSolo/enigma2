package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func sendEnigma(addr, msg string, due int) (string, error) {
	data := url.Values{
		"msg":  {msg},
		"due":  {strconv.Itoa(due)},
		"send": {"send"},
		"type": {"text"},
	}
	resp, err := http.PostForm(addr+"/post/", data)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(body), nil
}

func main() {
	var (
		serverAddress string
		dues, copies  int
	)
	flag.StringVar(&serverAddress, "s", "http://127.0.0.1:9000", "Server address")
	flag.IntVar(&dues, "d", 1, "How many days to keep the message 1..4")
	flag.IntVar(&copies, "c", 1, "How many times to copy messages 1...")
	flag.Parse()

	if !(dues >= 1 && dues <= 4) {
		log.Fatalln("Due must be within 1...4")
	}
	if copies <= 0 {
		log.Fatalln("Copies must be great 0")
	}

	for _, message := range flag.Args() {
		for i := 0; i < copies; i++ {
			resp, err := sendEnigma(serverAddress, message, dues)
			if err != nil {
				log.Println("fault send data err:", err)
			}
			fmt.Println(resp)
		}
	}
}
