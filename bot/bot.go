package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/algren123/gordle/config"
	"github.com/bwmarrin/discordgo"
)

var (
	BotId string
	goBot *discordgo.Session

	g           []guess_data
	wonArray    []string
	botResponse string
)

type guess_data struct {
	Slot   int    `json:"slot"`
	Guess  string `json:"guess"`
	Result string `json:"result"`
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running :)")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	command := strings.Fields(m.Content)

	// Returns Pong
	if command[0] == config.BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}

	if command[0] == config.BotPrefix+"guess" {
		_, err := isValidGuess(command[1])
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Guess is incorrect format. It needs to be 5 letters long.")
			return
		}

		resp, err := http.Get("https://api-slack-wordle.herokuapp.com/daily?guess=" + command[1])
		if err != nil {
			fmt.Println("error accessing API")
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error reading response body")
		}

		err = json.Unmarshal(body, &g)
		if err != nil {
			fmt.Println("error decoding guess data")
		}

		for _, letter := range g {
			// Add all the correct letters in an array that checks if the user has got the correct answer
			if letter.Result == "correct" {
				wonArray = append(wonArray, letter.Result)
			}

			if len(wonArray) == 5 && wonArray[0] == "correct" {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "You have completed today's Wordle! 🎉\n🟩 🟩 🟩 🟩 🟩", m.Reference())
				botResponse = ""
				wonArray = nil
				return
			}

			// Format the answer
			botResponse += formatResponse(letter.Result) + " "
		}

		_, err = s.ChannelMessageSendReply(m.ChannelID, botResponse, m.Reference())
		if err != nil {
			fmt.Println(err)
		}

		// Reset the response
		botResponse = ""
	}
}

func isValidGuess(g string) (bool, error) {
	r, _ := regexp.Compile(`^[A-Za-z]+$`)

	if len(g) != 5 || !r.MatchString(g) {
		return false, fmt.Errorf("incorrect format")
	}

	return true, nil
}

func formatResponse(result string) string {
	if result == "correct" {
		return "🟩"
	}

	if result == "present" {
		return "🟧"
	}

	if result == "absent" {
		return "⬛"
	}

	return result
}
