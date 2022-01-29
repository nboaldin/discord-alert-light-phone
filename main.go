package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln("Error creating discord session", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	//TODO get most of these in a config file silly balls!
	// testOwnersChannelId := "936728072315613294"
	ownersChannelId := os.Getenv("OWNERS_CHANNEL_ID")
	nathanUsername := os.Getenv("NATHAN_USER_NAME")
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_FROM_PHONE_NUMBER")
	to := os.Getenv("TWILIO_TO_PHONE_NUMBER")

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println("Mentions", m.Mentions)

	if DiscordUserSliceHasUsername(m.Mentions, nathanUsername) || m.ChannelID == ownersChannelId {
		//TODO get that client out of this func you NOOBer!
		twilioClient := twilio.NewRestClientWithParams(twilio.RestClientParams{
			Username: accountSid,
			Password: authToken,
		})
		params := &openapi.CreateMessageParams{}
		params.SetTo(to)
		params.SetFrom(from)
		params.SetBody(fmt.Sprintf("Discord Message ::: %s", m.ContentWithMentionsReplaced()))

		resp, err := twilioClient.ApiV2010.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			response, _ := json.Marshal(*resp)
			fmt.Println("Response: " + string(response))
		}
	}

}

func DiscordUserSliceHasUsername(users []*discordgo.User, mentionedUsername string) bool {
	for _, user := range users {
		if user.Username == mentionedUsername {
			return true
		}
	}
	return false
}
