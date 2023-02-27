package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func GetResponseBody(year int, rawUrl string) ([]byte, error) {
	if year < 0 || year > 2049 {
		return nil, errors.New("unsupported year")
	}
	_, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, errors.New("invalid url")
	}
	responseBody, err := http.Get(rawUrl)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(responseBody.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(responseBody.Body)
	return body, nil
}

type HolidayResponse struct {
	Response struct {
		Holidays []struct {
			Date struct {
				ISO string `json:"iso"`
			} `json:"date"`
		} `json:"holidays"`
	} `json:"response"`
}

func ConvertTOJSONWithDate(body []byte) ([]time.Time, error) {
	var publicHoliday HolidayResponse

	err := json.Unmarshal(body, &publicHoliday)
	if err != nil {
		log.Print(err)
		return []time.Time{}, err
	}

	var holidays []time.Time

	for _, holiday := range publicHoliday.Response.Holidays {
		date, err := time.Parse("2006-01-02T15:04:05-07:00", holiday.Date.ISO)
		if err != nil {
			date, err = time.Parse("2006-01-02", holiday.Date.ISO[:10])
			if err != nil {
				log.Print(err)
				continue
			}
		}
		holidays = append(holidays, date)
	}
	return holidays, nil
}

func AddHolidaysToDB(holidays []time.Time, urlDB string) error {

	//pgx is the driver for the Postgres db
	conn, err := pgx.Connect(context.Background(), urlDB)
	if err != nil {
		log.Fatal(err)
	}

	// Insert the JSON data into the database
	for _, holiday := range holidays {
		_, err = conn.Exec(context.Background(), "INSERT INTO Holidays(Date) VALUES($1) ON CONFLICT (Date) DO UPDATE SET Date = $1", holiday)
		if err != nil {
			return err
		}
	}

	fmt.Println("Successfully inserted JSON data into database")
	return nil
}
