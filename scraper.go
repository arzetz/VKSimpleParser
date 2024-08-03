package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type User struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Birthdate string `json:"birthdate,omitempty"`
	Status    string `json:"status,omitempty"`
}

func main() {
	var link string
	var filename string
	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}
	flag.StringVar(&link, "u", "", "Enter the URL please")
	flag.StringVar(&filename, "f", "", "Enter the filename, where you want info to be in in .json")
	flag.Parse()
	// configure the browser options
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"no-sandbox",
		"disable-dev-shm-usage",
		"disable-extensions",
		"user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
	}})

	// create a new remote client with the specified options
	driver, _ := selenium.NewRemote(caps, "")
	// maximize the current window to avoid responsive rendering
	driver.MaximizeWindow("")
	driver.Get(link)
	page, _ := driver.PageSource()

	re := regexp.MustCompile(`"domain":"([^"]*)`)
	user_username := re.FindStringSubmatch(page)
	re = regexp.MustCompile(`title>([^|]*)`)
	user_info := re.FindStringSubmatch(page) // Learning user name and surname
	re = regexp.MustCompile(`"bdate":"([^"]*)`)
	user_age := re.FindStringSubmatch(page)
	re = regexp.MustCompile(`"activity":"(.*?)"`)
	user_status := re.FindStringSubmatch(page)
	if len(user_username) > 1 {

		parts := strings.Split(user_info[1], " ")
		user := User{
			Username: user_username[1], // That's weird FindStringSubmatch structure, that's why we take like [1]
			Name:     parts[0],
			Surname:  parts[1],
		}

		if len(user_age) > 1 {
			user.Birthdate = user_age[1]
		}
		if len(user_status) > 1 {
			user.Status = user_status[1]
		}
		jsonData, _ := json.Marshal(user)
		os.WriteFile(filename, jsonData, 0666)
		fmt.Printf("\nEverything is cool! Check %s", filename)
	}

	defer func() {
		fmt.Println("\nHave a good BROOT") //Sorry but i don't give a foc
		os.Exit(0)
	}()
	defer service.Stop()
	defer driver.Quit()
}
