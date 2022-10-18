package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var wg sync.WaitGroup

func main() {
	//Введите ключевое слово для поиска

	var search_word = "кошка"

	start_time := time.Now()
	writer := create_csv()
	wg.Add(5)

	go parse_rt(writer, search_word)
	go parse_aif(writer, search_word)
	go parse_m24(writer, search_word)
	go parse_vz(writer, search_word)
	go parse_ria(writer, search_word)

	wg.Wait()

	fmt.Println("Execution time", time.Since(start_time).Seconds())

}

func parse_ria(writer *csv.Writer, search_word string) {
	link := "https://ria.ru/search/?query=" + url.QueryEscape(search_word)
	body := Run_user(link)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	print_error(err)
	doc.Find("div.list-item").Find("div.list-item__content").Each(func(id int, item *goquery.Selection) {
		link, _ := item.Find("a").Attr("href")
		title := item.Find("a").Text()
		data := []string{title, link}
		writer.Write(data)
	})
	writer.Flush()
	defer wg.Done()
}

func Run_user(link string) (body []byte) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", link, nil)
	print_error(err)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.160 YaBrowser/22.5.4.904 Yowser/2.5 Safari/537.36")
	resp, err := client.Do(req)
	print_error(err)
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	print_error(err)
	return body
}

func print_error(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func create_csv() *csv.Writer {
	file, err := os.Create("links.csv")
	print_error(err)
	writer := csv.NewWriter(file)
	return writer
}

func parse_rt(writer *csv.Writer, search_word string) {
	link := "https://russian.rt.com/search?q=" + url.QueryEscape(search_word)
	body := Run_user(link)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	print_error(err)
	doc.Find("div.listing__card ").Find("div").Find("div").Each(func(i int, item *goquery.Selection) {
		link_, _ := item.Find("a").Attr("href")
		title := item.Find("a").Text()
		strings.Split(link_, "\n")
		data := []string{strings.TrimSpace(title), "https://russian.rt.com" + link_}
		writer.Write(data)
	})
	writer.Flush()
	defer wg.Done()
}

func parse_m24(writer *csv.Writer, search_word string) {
	link := "https://www.m24.ru/sphinx?criteria=" + url.QueryEscape(search_word)
	body := Run_user(link)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	print_error(err)
	doc.Find("div").Find("main").Find("div.columns").Find("div.columns-right").Find("div").Find("div").Find("div").Find("section").Find("#SearchBody").Find("div").Find("ul").Each(func(i int, item *goquery.Selection) {
		link_, _ := item.Find("a").Attr("href")
		title := item.Find("a").Text()
		data := []string{strings.TrimSpace(title), "https://www.m24.ru" + link_}
		writer.Write(data)
	})
	writer.Flush()

	defer wg.Done()
}

func parse_aif(writer *csv.Writer, search_word string) {
	link := "https://aif.ru/search?text=" + url.QueryEscape(search_word)
	body := Run_user(link)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	print_error(err)
	doc.Find("div.search_list_js").Find("section").Find("div.list_item").Find("div.text_box").Each(func(i int, item *goquery.Selection) {
		link_, _ := item.Find("a").Attr("href")
		title := item.Find("a").Text()
		data := []string{strings.TrimSpace(title), "https://www.m24.ru" + link_}
		writer.Write(data)
	})
	writer.Flush()
	defer wg.Done()
}

func parse_vz(writer *csv.Writer, search_word string) {
	link := "https://vz.ru/search/?q=" + url.QueryEscape(search_word)
	body := Run_user(link)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	print_error(err)
	doc.Find("#main").Find("#ajax-search-res").Find("li").Each(func(i int, item *goquery.Selection) {
		link_, _ := item.Find("a").Attr("href")
		title := item.Find("a").Text()
		data := []string{strings.TrimSpace(title), "https://www.m24.ru" + link_}
		writer.Write(data)
	})
	writer.Flush()
	defer wg.Done()
}
