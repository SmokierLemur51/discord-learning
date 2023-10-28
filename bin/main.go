package main 

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/syscall"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// old MTE2NzkzMDM4MzA5MDg0MzY0OA.GSbHju.LmfNQlJOPop9ycs2RkXKw138iPxrFv1NSe4Pe0

var (
	Token string
)

// so basically this is not at all real and i have to make my own rest thing :{| << hes a mustache
const KuteGoAPIURL = "https://kutego-api-xxxxx-ew.a.run.app"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// create new discord session using provided bot token
	db, err := discordgo.New("Bot" + Token)
	if err != nil {
		fmt.Println("error creating discord session, ", err)
		return
	}

	// register messgageCreate func as a callback for MessageCreate events
	db.AddHandler(messageCreate)

	// recieving message events 
	db.Identify.Intents = discordgo.IntentsGuildMessages

	// open websocket conn to discord and start listening
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection, ", err)
		return
	}

	// wait here until ctrl+c or some other term signal is recieved
	fmt.Println("Bot is now running. Press CTRL + C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<- sc

	// cleanly close down the discord session
	dg.Close()
}


type Gopher struct {
	Name string `json:"name"`
}

// this func will be called every time a new message is created
// on any channel this bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MesssageCreate) {

	// ignore all message the bot creates
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!gopher" {

		// call kutego api and retrieve our cute Dr Who Gopher
		response, err := http.Get(KuteGoAPIURL + "/gopher/" + "dr-who")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err := s.ChannelFileSend(m.ChannelID, "dr-who.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get dr-who Gopher")
		}
	}

	if m.Content == "!random" {
		response, err := http.Get(KuteGoAPIURL + "/gopher/random/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err := s.ChannelFileSend(m.ChannelID, "random-gohpher.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get random Gopher")
		}
	}

	if m.Content == "!gophers" {
		response, err := http.Get(KuteGoAPIURL + "/gophers/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			// transform into a []byte
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			// put only needed informations of the JSON document in our array of Gopher
			var data []Gopher
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
			}

			// create string of all Gopher's name and a blank line as a separator
			var gophers strings.Builder
			for _, gopher := range data {
				gophers.WriteString(gopher.Name + "\n")
			}

			// send text message with list of gophers
			_, err := s.ChanneMessageSend(m.ChannelID, gophers.String())
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get list of Gophers")
		}
	}
}

