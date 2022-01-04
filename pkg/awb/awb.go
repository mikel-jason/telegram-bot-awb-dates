package awb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
)

type Date struct {
	Day   int    `json:"day"`
	Month int    `json:"month"`
	Year  int    `json:"year"`
	Type  string `json:"type"`
}

type dto struct {
	Data []Date `json:"data"`
}

type Fetcher struct {
	BuildingNo int
	StreetCode int
	StartMonth int
	StartYear  int
	EndMonth   int
	EndYear    int
}

func (f *Fetcher) Fetch() ([]Date, error) {
	url := fmt.Sprintf("https://www.awbkoeln.de/api/calendar?building_number=%d&street_code=%d&start_month=%d&start_year=%d&end_month=%d&end_year=%d&form=json", f.BuildingNo, f.StreetCode, f.StartMonth, f.StartYear, f.EndMonth, f.EndYear)
	response, err := http.Get(url)
	if err != nil {
		return []Date{}, err
	}

	if response.StatusCode != http.StatusOK {
		return []Date{}, errors.New(fmt.Sprintf("AWB Response with error code %d", response.StatusCode))
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return []Date{}, err
	}
	defer response.Body.Close()

	var dto dto
	err = json.Unmarshal(bodyBytes, &dto)
	if err != nil {
		return []Date{}, err
	}

	sort.Slice(dto.Data, func(p, q int) bool {
		pp := dto.Data[p]
		qq := dto.Data[q]
		if pp.Year != qq.Year {
			return pp.Year < qq.Year
		}
		if pp.Month != qq.Month {
			return pp.Month < qq.Month
		}
		return pp.Day < qq.Day
	})

	return dto.Data, nil
}
