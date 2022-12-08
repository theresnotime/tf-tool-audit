package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/user"
	"regexp"
	"time"
)

type tool struct {
	Name         string
	Title        string
	Description  string
	URL          string
	Repository   string
	Deprecated   bool
	Experimental bool
	License      string
}

func main() {
	rTools := regexp.MustCompile(`(?:tools\.)(.*)`)

	thisUser, userErr := user.Current()
	if userErr != nil {
		log.Fatal(userErr)
	}

	groups, grpsErr := thisUser.GroupIds()
	if grpsErr != nil {
		log.Fatal(grpsErr)
	}

	for _, element := range groups {
		group, grpErr := user.LookupGroupId(element)
		if grpErr != nil {
			log.Fatal(grpErr)
		}

		if rTools.MatchString(group.Name) {
			toolName := rTools.FindStringSubmatch(group.Name)[1]
			url := "https://toolhub.wikimedia.org/api/tools/toolforge-" + toolName + "/"

			spaceClient := http.Client{
				Timeout: time.Second * 2,
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Set("User-Agent", "User:TheresNoTime")

			res, getErr := spaceClient.Do(req)
			if getErr != nil {
				log.Fatal(getErr)
			}

			if res.Body != nil {
				defer res.Body.Close()
			}

			body, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {
				log.Fatal(readErr)
			}

			tool := tool{}
			jsonErr := json.Unmarshal(body, &tool)
			if jsonErr != nil {
				log.Fatal(jsonErr)
			}

			if res.Status == "404 Not Found" {
				// Didn't return a status code, assume OK
				fmt.Println(toolName + ": Missing")
			} else {
				// Returned a status_code
				fmt.Println(toolName + ": OK")
			}

		}
	}
}
