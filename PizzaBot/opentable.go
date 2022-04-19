package main

import (
/*	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"*/
	"time"
)

const (
	NO_AVAILABLE string = "NO_AVAILABLE"
	IN_ADVANCE          = "IN_ADVANCE"
)

type OpenTableResult struct {
	Message string
	Results []time.Time
}

func Query(id string, datetime time.Time, partySize string) (OpenTableResult, error) {
	result := OpenTableResult{}

	//	url := fmt.Sprintf("https://www.opentable.com/restaurant/profile/%s/reserve?restref=%s&datetime=%s&covers=%s", id, id, datetime, partySize)

<<<<<<< HEAD
/*	searchYear := datetime.Year()
=======
	/*searchYear := datetime.Year()
>>>>>>> cb3275b8d1b581e240b085cc24b8b124d636901d
	searchMonth := datetime.Month()
	searchDay := datetime.Day()
	searchHour := datetime.Hour()
	searchMinute := datetime.Minute()

	baseUrl := "https://www.opentable.com/restref/client/"
	paramUrl := fmt.Sprintf("%s?rid=%s&restref=%s&datetime=%d-%02d-%02dT%02d:%02d&partysize=%s",
		baseUrl, id, id, searchYear, searchMonth, searchDay, searchHour, searchMinute, partySize)
	res, err := http.Get(paramUrl)

	fmt.Println(paramUrl)
	if err != nil {
		return result, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return result, err
	}

	//fmt.Println(doc.Text())
	doc.Find(".dtp-results").Each(func(i int, s *goquery.Selection) {
		raw := s.Text()
		reg, _ := regexp.Compile("(.*)See affiliated restaurants")

		base := reg.FindStringSubmatch(raw)[1]

		matched, _ := regexp.MatchString("no online availability", base)

		if matched {
			result.Message = NO_AVAILABLE
		}

		matched, _ = regexp.MatchString("minutes in advance of the time", base)

		if matched {
			result.Message = IN_ADVANCE
		}

		reg, _ = regexp.Compile(`(\d+):(\d\d) (AM|PM)`)

		if reg.MatchString(base) {
			for _, v := range reg.FindAllStringSubmatch(base, -1) {
				hour, _ := strconv.Atoi(v[1])

				minutes, _ := strconv.Atoi(v[2])
				period := v[3]
				if period == "PM" {
					hour += 12
				}
				parsedTime := time.Date(searchYear, searchMonth, searchDay, hour, minutes, 0, 0, time.Local)
				result.Results = append(result.Results, parsedTime)
			}
		}

	})*/

	result.Results = append(result.Results, datetime)
	return result, nil
}
