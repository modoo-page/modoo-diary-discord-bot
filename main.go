package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}
	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentsGuildMessages
	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer discord.Close()
	discord.AddHandler(messageCreate)
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	ch, _ := s.Channel(m.ChannelID)

	if m.Content == "~로그인" {
		if ch.Type == discordgo.ChannelTypeDM {
			s.ChannelMessageSend(m.ChannelID, "DM")
			return
		}
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println(err)
		}
		s.ChannelMessageSend(channel.ID, "1. ~토큰요청 {EMAIL}\n2.~토큰로그인 {EMAIL} {TOKEN}\n\n 시도해주세요")

	}
	if strings.HasPrefix(m.Content, "~일기 ") {
		text := m.Content[8:] // bytes 계산법
		fmt.Println(text)
		s.ChannelMessageSend(m.ChannelID, text)
	}
}
