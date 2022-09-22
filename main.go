// Author Gera Jahja
// Last update : 21/09/2022
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var links = []string{}        //global var
var visitedlinks = []string{} //global var
// variable declaring the website
var webpgName = "https://www.monzo.com/" //:= declares and assigns the variable , whereas = is simply assignment
var domainName = strings.TrimPrefix(webpgName, "https://www")

func main() {
	links = append(links, webpgName)
	getLinks()
}

func getLinks() { //recursive function... very slow ;/ passed on 21/09/2022
	links = removeDuplicateStr(links)
	visitedlinks = removeDuplicateStr(visitedlinks)
	i := 0
	for _, link := range links { //infinite loop
		response, e := http.Get(link)
		GetCheck(e)
		defer response.Body.Close()                                 // whenthe body of a website is retrieved it must be closed.
		document, e := goquery.NewDocumentFromReader(response.Body) // Create a goquery document from the HTTP response
		GetCheck(e)
		// Find all links and verify/print them with the function elementIsPresent    // defined earlier
		document.Find("a").Each(hrefCheck) //displays all links on a single website , by filtering tags by <a>

		visitedlinks = append(visitedlinks, link)
		linkNo := len(visitedlinks)
		linktoseeNo := len(links)

		fmt.Println(link) //prints links during runtime
		fmt.Println("Visited: ", linkNo)
		fmt.Println("Links to Crawl: ", linktoseeNo)
		fmt.Println("")

		// Remove the element
		copy(links[i:], links[i+1:]) // Shift a[i+1:] left one index.
		links[len(links)-1] = ""     // Erase last element (write zero value).
		links = links[:len(links)-1] // Truncate slice.
		i += 1

		if len(links) > 0 {
			getLinks()
		} else {
			fmt.Println("")
			fmt.Println("Overall Web Crawler Result:")
			fmt.Println("Visited: ", linkNo)
			fmt.Println("Links discovered:") //should be empty!
			for _, link := range links {
				fmt.Println(link)
			}
			fmt.Println("Links discovered and visited:")
			for _, link := range visitedlinks {
				fmt.Println(link)
			}
		}
	}
}

// Looks through the element that has been passed, and sees whether 'href' is present
// Tested : passed on 20/09/2022
func hrefCheck(index int, element *goquery.Selection) {

	// See if the href attribute exists on the element
	href, exists := element.Attr("href") //accesses the attribute that has the tag <a href... tags (which store links in HTML)
	if exists {
		if strings.Contains(href, domainName) {
			if !(strings.Contains(href, "/u/")) && !(strings.Contains(href, "/t/")) {
				links = append(links, href)
			}
		}
	}
}

func GetCheck(err error) {
	//http.Get returns a response and an error so we have to handle the error with an if statement
	if err != nil {
		log.Fatal(err)
	}
}

// https://stackoverflow.com/questions/66643946/how-to-remove-duplicates-strings-or-int-from-slice-in-go
func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

