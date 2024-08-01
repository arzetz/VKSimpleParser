package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

type User struct {
	Name    string
	Surname string
}

func main() {

	var user_info []User
	c := colly.NewCollector(colly.AllowedDomains("vk.com", "www.vk.com", "https://vk.com", "m.vk.com"))
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	scrapeUrl := "https://vk.com/"

	c.OnHTML("script", func(e *colly.HTMLElement) {

		re := regexp.MustCompile(`"first_name_nom":"([^"]*)`)
		match_name := re.FindStringSubmatch(e.Text)

		reZ := regexp.MustCompile(`"last_name_gen":"([^"]*)`)
		match_surname := reZ.FindStringSubmatch(e.Text)

		if len(match_name) > 1 || len(match_surname) > 1 {
			fmt.Println(match_name[1], (match_surname[1]))
			user_info = append(user_info, User{match_name[1], match_surname[1]})
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {

		fmt.Printf("Error while shittin %s\n", e.Error())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("response ect' %d\n", r.StatusCode)
	})

	c.Visit(scrapeUrl)
	jsonFile, jsonErr := os.Create("user_info.json")
	if jsonErr != nil {
		log.Fatalln("Failed to create the output JSON file", jsonErr)
	}
	defer jsonFile.Close()

	jsonString, _ := json.MarshalIndent(user_info, " ", " ")

	jsonFile.Write(jsonString)
}
