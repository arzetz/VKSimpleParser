package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type User struct {
	username, name, surname, birthdate, quote string
}

func main() {
	var users []User
	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}

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
	driver.Get("https://vk.com/shcherbakovv27")
	page, _ := driver.PageSource()

	file, _ := os.Create("helloween.txt")

	file.WriteString(page)
	re := regexp.MustCompile(`"domain":"([^"]*)`)
	user_username := re.FindStringSubmatch(page)
	re = regexp.MustCompile(`title>([^|]*)`)
	user_info := re.FindStringSubmatch(page) // Learning user name and surname
	re = regexp.MustCompile(`"bdate":"([^"]*)`)
	user_age := re.FindStringSubmatch(page)
	fmt.Println(user_info)
	if len(user_username) > 1 {
		user := User{username: user_username[1]}
		users = append(users, user)
	}
	if len(user_info) > 1 {
		parts := strings.Split(user_info[1], " ")
		users[0].name = parts[0]
		users[0].surname = parts[1]

	}
	if len(user_age) > 1 {
		users[0].birthdate = user_age[1] // That's weird FindStringSubmatch structure, that's why we take like [1]
	}
	fmt.Println(users, "Done.")
	defer func() {
		fmt.Println("Bb") //Sorry but i don't give a foc
		os.Exit(0)
	}()
	defer service.Stop()
	defer driver.Quit()
	defer file.Close()
}
