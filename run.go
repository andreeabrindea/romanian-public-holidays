package main

import (
	"fmt"
	"log"
	"publicHolidays/handler"
)

func main() {
	year := 2024
	apiKey := "6ff1c210e443df1122a57a52aa383388be119213"
	rawUrl := fmt.Sprintf("https://calendarific.com/api/v2/holidays?api_key=%s&country=RO&year=%v", apiKey, year)

	response, err := handler.GetResponseBody(year, rawUrl)
	if err != nil {
		log.Println(err)
		return
	}
	holidays, err := handler.ConvertTOJSONWithDate(response)
	err = handler.AddHolidaysToDB(holidays, "postgres://xvyctfje:5yGXTCPQKkKJe0rjuvsJtFOQF7BiOBJp@mouse.db.elephantsql.com/xvyctfje")
	if err != nil {
		log.Print(err)
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}
