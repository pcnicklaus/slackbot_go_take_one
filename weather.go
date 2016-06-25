package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
)


// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getWeather(city, state string) string {
   // need to parse it into different words but for now just run with one
	city = strings.ToUpper(city)
   state = strings.ToUpper(state)

   gUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s,+%s&key=AIzaSyDdxHZhiR83jJa1-mpFRTMIt28YhtpcytY", city, state)

   resp, err := http.Get(gUrl)
   if err != nil {
      return fmt.Sprintf("error: %v", err)
   }
   return fmt.Sprintf("response %v", resp)


	// url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
	// resp, err := http.Get(url)
	// if err != nil {
	// 	return fmt.Sprintf("error: %v", err)
	// }
	// rows, err := csv.NewReader(resp.Body).ReadAll()
	// if err != nil {
	// 	return fmt.Sprintf("error: %v", err)
	// }
	// if len(rows) >= 1 && len(rows[0]) == 5 {
	// 	return fmt.Sprintf("%s (%s) is trading at $%s", rows[0][0], rows[0][1], rows[0][2])
	// }
	// return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}

//https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&key=YOUR_API_KEY
// google key = AIzaSyDdxHZhiR83jJa1-mpFRTMIt28YhtpcytY
// dark star key = f1ba36ed1b17183a9a37dc37367e1b38
