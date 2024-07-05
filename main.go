package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"

	"github.com/oriser/regroup"
)

const PSTFormat = "Jan 2 15:04:05 PST"
const MAX_MSG_LEN = 512

// Regex for parsing PRIVMSG strings.
//
// First matched group is the user's name and the second matched group is the content of the
// user's message.
var MsgRegex *regroup.ReGroup = regroup.MustCompile(`^\w*\s*:(?P<username>\w+)!\w+@\w+.tmi.twitch.tv PRIVMSG (#(?P<channel>\w+))?\s*:(?P<message>.*)$`)

// Regex for parsing user commands, from already parsed PRIVMSG strings.
//
// First matched group is the command name and the second matched group is the argument for the
// command.
var CmdRegex *regroup.ReGroup = regroup.MustCompile(`^!(\w+)\s?(\w+)?`)

type Bot interface {
	Connect(addr string) error
	Disconnect()
	Pong() error
	ReadMsg() (string, error)
	JoinChannel(chl, username, password string) error
	LeaveChannel(chl string) error
	Say(msg string) error
	Whisper(user, msg string) error
	HandleCommand(cmd, msg string) error
}

type CommandHandler func(bot Bot, msg string) error

type BasicBot struct {
	// The address of the server the bot is currently connected to
	addr string

	// A reference to the bot's connection to the server
	conn net.Conn

	// list of channels that the bot is a part of
	channels []string

	// Contains a map of command names and their handlers
	commandHandlers map[string]CommandHandler

	reader *textproto.Reader
}

func (bb *BasicBot) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("Bot connected to %s\n", addr)
	bb.addr = addr
	bb.conn = conn
	return nil
}

func (bb *BasicBot) Disconnect() {
	bb.conn.Close()
	log.Printf("Bot disconnected from %s\n", bb.addr)
	bb.addr = ""
	bb.conn = nil
}

func (bb *BasicBot) Pong() error {
	_, err := fmt.Fprint(bb.conn, "PONG :tmi.twitch.tv\r\n")
	return err
}

func (bb *BasicBot) ReadMsg() (msg string, err error) {
	if bb.reader == nil {
		bb.reader = textproto.NewReader(bufio.NewReader(bb.conn))
	}
	msg, err = bb.reader.ReadLine()
	return
}

func (bb *BasicBot) isInChannel(chl string) (int, bool) {
	for ci, c := range bb.channels {
		if c == chl {
			return ci, true
		}
	}
	return -1, false
}

func (bb *BasicBot) JoinChannel(chl, username, password string) error {
	if _, in := bb.isInChannel(chl); in {
		return fmt.Errorf("bot already in channel %s", chl)
	}

	log.Printf("Bot joining channel %s as %s\n", chl, username)
	_, err := fmt.Fprintf(bb.conn, "PASS %s\r\nNICK %s\r\nJOIN #%s\r\n", password, username, chl)
	if err != nil {
		return err
	}
	log.Printf("Bot Successfully joined channel %s as %s\n", chl, username)
	bb.channels = append(bb.channels, chl)
	return nil
}

func (bb *BasicBot) LeaveChannel(chl string) error {
	ci, in := bb.isInChannel(chl)
	if !in {
		return fmt.Errorf("bot not in channel %s", chl)
	}
	_, err := fmt.Fprintf(bb.conn, "PART #%s\r\n", bb.channels[ci])
	bb.channels = append(bb.channels[:ci], bb.channels[ci+1:]...)
	log.Printf("Bot left chanel %s\n", bb.channels[ci])
	return err
}

func (bb *BasicBot) Say(msg string) error {
	if msg == "" {
		return fmt.Errorf("bot attempted to send empty message!")
	}

	// check if message is too large for IRC
	if len(msg) > MAX_MSG_LEN {
		return fmt.Errorf("bot attempted to send message that was too large, %d/%d", len(msg), MAX_MSG_LEN)
	}

	for _, chl := range bb.channels {
		_, err := fmt.Fprintf(bb.conn, "PRIVMSG #%s :%s\r\n", chl, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bb *BasicBot) Whisper(user, msg string) error {
	return fmt.Errorf("whisper is not implemented due to twitch chat restrictions")
}

func (bb *BasicBot) HandleCommand(cmd, msg string) error {
	handler, ok := bb.commandHandlers[cmd]
	if !ok {
		return fmt.Errorf("unrecognized command %s", cmd)
	}
	return handler(bb, msg)
}

func ParseIRCMessage(line string) (user, channel, msg string, err error) {
	matches, err := MsgRegex.Groups(line)
	if err != nil {
		return
	}
	muser, ok := matches["username"]
	if ok {
		user = muser
	}
	mchan, ok := matches["channel"]
	if ok {
		channel = mchan
	}
	mmsg, ok := matches["message"]
	if ok {
		msg = mmsg
	}
	return
}

func RunBotLoop(b Bot, addr, username, password string, channels []string) {
	err := b.Connect(addr)
	if err != nil {
		panic(err)
	}
	defer b.Disconnect()

	for _, chl := range channels {
		err = b.JoinChannel(chl, username, password)
		if err != nil {
			log.Printf("unable to join channel %s: %v\n", chl, err.Error())
		}
	}

	log.Println("bot entering command handler...")
	for {
		msg, err := b.ReadMsg()
		if err != nil {
			log.Printf("unexpected error occured while reading msg: %v\n", err.Error())
			continue
		}
		if "PING" == msg[:4] {
			b.Pong()
			continue
		}
		u, c, m, err := ParseIRCMessage(msg)
		if err != nil {
			log.Printf("failed to parse line '%s': %s\n", msg, err.Error())
			continue
		}
		log.Printf("Message: %s@%s '%s'\n", u, c, m)

	}
}

type Credentials struct {
	Username string
	Password string
}

var (
	irc_addr    string = "irc.chat.twitch.tv:6667"
	credentials string = "./config.json"
	channel     string = "cheezitthehedgehog"
)

func main() {
	flag.StringVar(&irc_addr, "address", irc_addr, "the address to use for twitch connection")
	flag.StringVar(&credentials, "credentials", credentials, "the location of the credentials json file")
	flag.StringVar(&channel, "channel", channel, "the channel for the bot to join")
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

	tok, err := GetAccessToken(account.Username, account.Password)
	if err != nil {
		panic(err)
	}

	RunBotLoop(&BasicBot{}, irc_addr, account.Username, fmt.Sprintf("oauth:%s", tok.AccessToken), []string{channel})
}
