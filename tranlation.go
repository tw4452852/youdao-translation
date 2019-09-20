package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	showHelp = flag.Bool("h", false, "show help")
)

func main() {
	flag.Parse()

	//must have at least one word
	if flag.NArg() == 0 || *showHelp {
		usage()
		os.Exit(1)
	}

	translateWords(flag.Args())
}

func usage() {
	fmt.Printf("%s word1 [word2...wordN]\n", os.Args[0])
	flag.PrintDefaults()
}

type Results struct {
	Entries []Entry `xml:"custom-translation>translation>content"`
}

type Entry struct {
	Content string `xml:",chardata"`
}

func translateWords(words []string) {
	const urlPrefix = "http://dict.youdao.com/fsearch?q="
	for _, word := range words {
		resp, err := http.Get(urlPrefix + word)
		if err != nil {
			log.Println(err)
			continue
		}
		resultXml, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println(err)
			continue
		}

		var result Results
		if err := xml.Unmarshal(resultXml, &result); err != nil {
			log.Println(err)
			log.Printf("resp:\n%s\n", string(resultXml))
			continue
		}

		resultS := ""
		for i, ent := range result.Entries {
			if i != 0 {
				resultS += " "
			}
			resultS += ent.Content
		}
		fmt.Printf("%s: %s", word, resultS)
	}
}
