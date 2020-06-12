package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const starturl1 = "https://search.yahoo.co.jp/image/search?p="
const starturl2 = "&op=&ei=UTF-8&b="

func main() {
	send, word := CreateURL()
	url := GetPage(send)
	EncodingCSV(url, word)
}

func test() {
	send, word := CreateURL()
	url := GetPage(send)
	EncodingCSV(url, word)
}

func CreateURL() (<-chan string, string) {
	var word string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word = s.Text()
		break
	}
	enword := url.QueryEscape(word)

	send := make(chan string, 10)
	go func() {
		defer close(send)
		for i := 1; i < 700; i += 20 {
			a := strconv.Itoa(i)
			urls := starturl1 + enword + starturl2 + a
			send <- urls
		}
	}()
	return send, word
}

func GetPage(reci <-chan string) <-chan string {
	send := make(chan string, 10)
	go func() {
		defer close(send)
		for i := range reci {
			doc, _ := goquery.NewDocument(i)
			doc.Find("img").Each(func(_ int, s *goquery.Selection) {
				url, _ := s.Attr("src")
				send <- url
			})
		}
	}()
	return send
}

func EncodingCSV(reci <-chan string, word string) {
	filename := word + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("cannot create file", err)
	}
	defer file.Close()

	var record = []string{}
	i := 1

	csvfile := csv.NewWriter(file)
	for j := range reci {
		id := strconv.Itoa(i)
		record = []string{
			id,
			j,
		}
		err := csvfile.Write(record)
		if err != nil {
			log.Fatal("cannot write record:", err)
		}
		i++
	}
}
