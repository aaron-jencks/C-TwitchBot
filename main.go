package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/aaronjencks/gitchbot/storage"
	"github.com/oriser/regroup"
)

const MAX_MSG_LEN = 512

// Regex for parsing user commands, from already parsed PRIVMSG strings.
//
// First matched group is the command name and the second matched group is the argument for the
// command.
var CmdRegex *regroup.ReGroup = regroup.MustCompile(`^!(\w+)\s?(\w+)?`)

type Credentials struct {
	Username     string `json:"username"`
	Token        string `json:"oauth_token"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var (
	irc_addr    string = "irc.chat.twitch.tv:6667"
	credentials string = "./config.json"
	channel     string = "cheezitthehedgehog"
	backing     string = "./data.db"
)

func main() {
	flag.StringVar(&irc_addr, "address", irc_addr, "the address to use for twitch connection")
	flag.StringVar(&credentials, "credentials", credentials, "the location of the credentials json file")
	flag.StringVar(&channel, "channel", channel, "the channel for the bot to join")
	flag.StringVar(&backing, "db", backing, "the location of the sql database for data backing")
	flag.Parse()

	fp, err := os.Open(credentials)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	decode := json.NewDecoder(fp)
	account := Credentials{}
	err = decode.Decode(&account)
	if err != nil {
		panic(err)
	}

	sqlBacking, err := storage.CreateSqliteBacker(backing)
	if err != nil {
		panic(err)
	}

	bot := CreateBasicTwitchBot(account.Username, account.Token, sqlBacking)
	CreateCounterHandler(bot, "oops", 0, "Whoopsie, I made a mistake")
	bot.Join(channel)
	bot.Loop()
}
