// Author Gera Jahja
// Last update : 30/09/2022
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)
//map of all links on website
var alllinks = make(map[string]bool)

// variable declaring the website
var webpgName = "https://monzo.com" //:= declares and assigns the variable , whereas = is simply assignment

//devlares the waitgroup used in func main() and func getLinks()
var wg sync.WaitGroup

/*
   Tracks time taken to run the program
   Initial call to getLinks()
   adds the web domain as the first value in the map alllinks to be crawled 
*/
func main() {
	start := time.Now()
	alllinks[webpgName+"/"] = true

	getLinks()
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Time taken: %s", elapsed)
}

/*  Creates a goquery document from the HTTP response
	gets all links on a single website , by filtering tags by <a>
	then calls itself so that the updated map of links is crawled also
	Uses waitgroups for concurrency
*/
func getLinks() { 

	for link, inMap := range alllinks {
		wg.Add(1)
		if inMap {
			response, e := http.Get(link)
			if e != nil {
				break
				//log.Fatal(e)
			}
			defer response.Body.Close()                                 
			document, e := goquery.NewDocumentFromReader(response.Body) 
			if e != nil {
				break
				//log.Fatal(e)
			}
			document.Find("a").Each(hrefCheck) 
		} else {
			break
		}
	}
	go getLinks()
	wg.Done()
}

/* 
   Looks through the element that has been passed, and sees whether 'href' is present
   accesses the attribute that has the tag <a href... tags (which store links in HTML)
   writes the links and relative links in an document (see func writeLinkToFile)
   adds link to map allLinks
   uses locks for concurrency
   ensures no duplicate values are present in the map
*/
func hrefCheck(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")
	var m sync.RWMutex
	if exists {
		if !(strings.Contains(href, "email-protection")) {
			if (!alllinks[href]) && (!alllinks[webpgName+href]) {
				if strings.HasPrefix(href, webpgName) {
					m.RLock()
					alllinks[href] = true
					go writeLinkToFile(href)
					fmt.Println(href)
					m.RUnlock()
					return
				}
				if strings.HasPrefix(href, "/") {
					m.RLock()
					alllinks[webpgName+href] = true
					go writeLinkToFile(webpgName + href)
					fmt.Println(webpgName + href)
					m.RUnlock()
					return
				}
			} else {
				return
			}
		}
	}
}

// 	Writes a string (passes as a parameter) on a new line in a txt file named Crawled
func writeLinkToFile(data string) {
	file, err := os.OpenFile("Crawled.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)
	_, _ = datawriter.WriteString(data + "\n")
	datawriter.Flush()
	file.Close()
}
