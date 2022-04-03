package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GuessData struct {
	Slot   int    `json:"slot"`
	Guess  string `json:"guess"`
	Result string `json:"result"`
}

var guessData []GuessData

func GetGuessData(guess string) ([]GuessData, error) {
	resp, err := http.Get("https://api-slack-wordle.herokuapp.com/daily?guess=" + guess)
	if err != nil {
		return nil, fmt.Errorf("error accessing API")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body")
	}

	err = json.Unmarshal(body, &guessData)
	if err != nil {
		return nil, fmt.Errorf("word doesn't exist in database")
	}

	return guessData, err
}
