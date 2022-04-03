package bot

import (
	"fmt"
	"strings"

	"github.com/algren123/gordle/api"
	"github.com/algren123/gordle/config"
	"github.com/algren123/gordle/store"
	"github.com/algren123/gordle/tools"
	"github.com/bwmarrin/discordgo"
)

var (
	BotId string
	goBot *discordgo.Session

	wonArray    []string
	botResponse string
)

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

	store.GetServerByID = make(map[string]*store.ServerConfig)
	store.GetUserByID = make(map[string]*store.User)

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

	// Returns amount of points user has
	if command[0] == config.BotPrefix+"points" {
		if !store.UserExists(m.Author.ID) {
			_, _ = s.ChannelMessageSend(m.ChannelID, "You need to at least have a guess first using the command '!guess' to be added to the database")
			return
		}

		_, _ = s.ChannelMessageSendReply(m.ChannelID, "You have "+store.GetPoints(m.Author.ID)+" point(s)", m.Reference())
	}

	// Guess handler
	if command[0] == config.BotPrefix+"guess" {
		// Handle registering a new server to local memory
		if !store.ServerExists(m.GuildID) {
			newServer := store.ServerConfig{ID: m.GuildID, CanPlay: true, Tries: 6}
			store.StoreServer(newServer)
		}

		// Handle registering new users to local memory
		if !store.UserExists(m.Author.ID) {
			newUser := store.User{ID: m.Author.ID, Author: m.Author.Username, Points: 0}
			store.StoreUser(newUser)
		}

		if !store.AbleToPlay(m.GuildID) {
			_, _ = s.ChannelMessageSend(m.ChannelID, "You cannot play yet. You have to wait until 00:00 GMT Time.")
			return
		}

		_, err := tools.IsValidGuess(command[1])
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Guess is incorrect format. It needs to be 5 letters long.")
			return
		}

		guessData, err := api.GetGuessData(command[1])
		if err != nil {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
			return
		}

		for _, letter := range guessData {
			// Add all the correct letters in an array that checks if the user has got the correct answer
			if letter.Result == "correct" {
				wonArray = append(wonArray, letter.Result)
			}

			// Handle win
			if len(wonArray) == 5 && wonArray[0] == "correct" {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, "You have completed today's Wordle! 🎉\n🟩 🟩 🟩 🟩 🟩", m.Reference())
				store.GivePoint(m.Author.ID)
				store.SetCanPlay(m.GuildID, false)
				botResponse = ""
				wonArray = nil
				return
			}

			// Format the answer
			botResponse += tools.FormatResponse(letter.Result) + " "
		}

		_, err = s.ChannelMessageSendReply(m.ChannelID, botResponse, m.Reference())
		if err != nil {
			fmt.Println(err)
		}

		store.SubtractServerTry(m.GuildID)

		if store.GetServerTries(m.GuildID) < 1 {
			store.SetCanPlay(m.GuildID, false)
		}

		// Reset params
		botResponse = ""
		wonArray = nil
	}
}
