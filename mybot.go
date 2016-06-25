/*

mybot - Illustrative Slack bot in Go

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
		"io/ioutil"
		"encoding/json"
	"strings"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
	// 	os.Exit(1)
	// }

	// start a websocket-based Real Time API session
	ws, id := slackConnect("xoxb-54265211282-nOq9gSAiPWz7S8lPP3YSThwt")
	fmt.Println("mybot ready, ^C exits")

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 3 && parts[1] == "stock" {
				// looks good, get the quote and reply with the result
				go func(m Message) {
					m.Text = getQuote(parts[2])
					postMessage(ws, m)
				}(m)
				// NOTE: the Message object is copied, this is intentional
			}
			if len(parts) == 3 && parts[1] == "weather" {
				go func(m Message) {
					m.Text = getWeather(parts[2])
					postMessage(ws, m)
				}(m)
			}
			// else {
			// 	// huh?
			// 	m.Text = fmt.Sprintf("sorry, that does not compute\n")
			// 	postMessage(ws, m)
			// }
		}

	}
}

// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getQuote(sym string) string {
	sym = strings.ToUpper(sym)
	url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if len(rows) >= 1 && len(rows[0]) == 5 {
		return fmt.Sprintf("%s (%s) is trading at $%s", rows[0][0], rows[0][1], rows[0][2])
	}
	return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}

func getWeather(city string) string {
   // need to parse it into different words but for now just run with one
	city = strings.ToUpper(city)
   // state = strings.ToUpper(state)

   gUrl := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s,+CO&key=AIzaSyDdxHZhiR83jJa1-mpFRTMIt28YhtpcytY", city)

   resp, err := http.Get(gUrl)
   if err != nil {
      return fmt.Sprintf("error: %v", err)
   }
   return fmt.Sprintf("response %v", resp.body)


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


type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	Id string `json:"id"`
}

// slackStart does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func slackStart(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	var respObj responseRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.Url
	id = respObj.Self.Id
	return
}

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func getMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

var counter uint64

var hi Message

func postMessage(ws *websocket.Conn, m Message) error {
	m.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}

// Starts a websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func slackConnect(token string) (*websocket.Conn, string) {
	wsurl, id, err := slackStart(token)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		log.Fatal(err)
	}

	return ws, id
}
