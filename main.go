package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/mdp/qrterminal/v3"
	"io"
	"net/http"
	"os"
	"strings"
)

type PostResponse struct {
	Id string `json:"id"`
}

type GetRequest struct {
	Code string `json:"code"`
}

type GetResponse struct {
	Msg string `json:"msg"`
}

func main() {
	host := flag.String("host", "https://webclip2.mytools.express/", "host to operation with webclip2")
	post := flag.Bool("post", false, "post message")
	fromClipboard := flag.Bool("from-clipboard", false, "get from the clipboard")
	toClipboard := flag.Bool("to-clipboard", false, "set result to the clipboard")
	get := flag.Bool("get", false, "get message")
	code := flag.String("code", "", "retrieve code")
	showId := flag.Bool("show-id", false, "show the code instead of url")
	showQr := flag.Bool("show-qr", false, "show the qr code")
	flag.Parse()

	if *post {
		var msg string
		if *fromClipboard {
			b, err := clipboard.ReadAll()
			if err != nil {
				panic(err)
			}
			msg = b
		} else {
			b, err := io.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			msg = string(b)
		}
		msg = strings.ReplaceAll(msg, "\r", "\\r")
		msg = strings.ReplaceAll(msg, "\n", "\\n")

		url := *host + "api/msg/create"
		responseObj := PostResponse{}
		response, err := http.Post(url, "application/json", bytes.NewBufferString("{\"msg\":\""+msg+"\"}"))
		if err != nil {
			panic(err)
		} else if msg, err := io.ReadAll(response.Body); err != nil {
			panic(err)
		} else if err := json.Unmarshal(msg, &responseObj); err != nil {
			panic(err)
		} else {
			url := "https://webclip2.mytools.express/#/get?id=" + responseObj.Id
			if *showQr {
				config := qrterminal.Config{
					Level:     qrterminal.M,
					Writer:    os.Stdout,
					BlackChar: qrterminal.WHITE,
					WhiteChar: qrterminal.BLACK,
					QuietZone: 1,
				}
				qrterminal.GenerateWithConfig(url, config)
			}
			fmt.Println(url)
			if *toClipboard {
				if *showId {
					err := clipboard.WriteAll(responseObj.Id)
					if err != nil {
						panic(err)
					}
				} else {
					err := clipboard.WriteAll(url)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	} else if *get && *code != "" {
		url := *host + "api/msg/retrieve"
		codeObj := GetRequest{}
		codeObj.Code = *code
		responseObj := GetResponse{}
		response, err := http.Post(url, "application/json", bytes.NewBufferString("{\"code\":\""+*code+"\"}"))
		if err != nil {
			panic(err)
		} else if msg, err := io.ReadAll(response.Body); err != nil {
			panic(err)
		} else if err := json.Unmarshal(msg, &responseObj); err != nil {
			panic(err)
		} else {
			responseObj.Msg = strings.ReplaceAll(responseObj.Msg, "\\r", "\r")
			responseObj.Msg = strings.ReplaceAll(responseObj.Msg, "\\n", "\n")
			if *toClipboard {
				err := clipboard.WriteAll(responseObj.Msg)
				if err != nil {
					return
				}
			} else {
				fmt.Println(responseObj.Msg)
			}
		}
	} else {
		flag.PrintDefaults()
	}
}
