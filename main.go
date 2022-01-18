package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln("Error creating discord session", err)
		return
	}

}
