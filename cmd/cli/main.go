package main

import (
	"flag"
	"fmt"
	"github.com/sarcaustech/go-telegram-awb/pkg/awb"
	"log"
)

func main() {

	buildingNo := flag.Int("buildingNo", -1, "Building number")
	streetCode := flag.Int("streetCode", -1, "Code for target street")
	startMonth := flag.Int("startMonth", -1, "Start month of fetch period")
	startYear := flag.Int("startYear", -1, "Start year of fetch period")
	endMonth := flag.Int("endMonth", -1, "End month of fetch period")
	endYear := flag.Int("endYear", -1, "End year of fetch period")
	flag.Parse()

	for index, value := range []int{*buildingNo, *streetCode, *startMonth, *startYear, *endMonth, *endYear} {
		if value < 0 {
			log.Println(index, value)
			log.Fatalln("Insufficient data provided")
		}
	}

	fetcher := awb.Fetcher{
		BuildingNo: *buildingNo,
		StreetCode: *streetCode,
		StartMonth: *startMonth,
		StartYear:  *startYear,
		EndMonth:   *endMonth,
		EndYear:    *endYear,
	}

	dates, err := fetcher.Fetch()
	if err != nil {
		log.Println(err)
	}

	for _, date := range dates {
		log.Println(fmt.Sprintf("[%02d.%02d.%4d] %s", date.Day, date.Month, date.Year, date.Type))
	}
}
