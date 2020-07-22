package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token         string
	fizzbuzzRegex = "!fizzbuzz\\s(\\d+)"
	validCommand  = regexp.MustCompile(fizzbuzzRegex)
)

func init() {
	flag.StringVar(&token, "token", "", "bot token")
}

func main() {
	flag.Parse()
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error while creating discord bot session: %v", err)
	}

	// Register handlers for fizzbuzz
	session.AddHandler(onMessageCreate)

	log.Println("connecting to discord servers")

	if err := session.Open(); err != nil {
		log.Fatalf("error while connecting to the discord servers: %v", err)
	}

	log.Println("fizzbuzz is now running...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := session.Close(); err != nil {
		log.Fatalf("error while closing session: %v", err)
	}
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("new message: %v", m)
	// Ignore message if its from the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!fizzbuzz") {
		number, err := getFizzbuzzInput(m.Content)
		if err != nil {
			log.Println("error while getting fizzbuzz input: ", err)
		}
		ans := fizzbuzz(number)
		s.ChannelMessageSend(m.ChannelID, formatResponse(ans))
	}

}

func getFizzbuzzInput(message string) (int, error) {
	parts := validCommand.FindStringSubmatch(message)
	if len(parts) > 2 {
		return 0, errors.New("the input must match the following regex: " + fizzbuzzRegex)
	}
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("the input must match the following regex: " + fizzbuzzRegex)
	}
	return num, nil
}

func fizzbuzz(i int) string {
	if i%3 == 0 && i%5 == 0 {
		return "fizzbuzz"
	} else if i%3 == 0 {
		return "fizz"
	} else if i%5 == 0 {
		return "buzz"
	}
	return strconv.Itoa(i)
}

func formatResponse(ans string) string {
	return fmt.Sprintf("```%s```", ans)
}
