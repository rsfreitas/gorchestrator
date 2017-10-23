//A chat bot orchestrator ;-)
package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	//	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/processone/gox/xmpp"
	"github.com/spacemonkeygo/spacelog"
	"menteslibres.net/gosexy/yaml"
)

const AppName string = "orchestrator"
const Version string = "0.1.1"

//The command line options.
type CLIOptions struct {
	version      *bool
	botFilename  *string
	certFilename *string
}

//The content of a YAML bot configuration file.
type ConfigServer struct {
	hostname string
	port     int
	username string
	password string
}

type BotInfo struct {
	name string
}

type BotConfig struct {
	ConfigServer
	BotInfo
	logFilename string
	level       string
}

//buildCLIOptions configures the application supported command line options.
func buildCLIOptions(options *CLIOptions) {
	options.version = flag.Bool("version", false,
		"Shows the current application version.")

	flag.BoolVar(options.version, "v", false,
		"Shows the current application version.")

	options.botFilename = flag.String("B", "",
		"Indicate a YAML config file with the bot-XMPP server informations.")

	options.certFilename = flag.String("C", "",
		"Indicate the certificate to establish connection with XMPP server.")
}

//loadBotConfiguration loads all required configuration from the YAML file.
func loadBotConfiguration(filename string) (BotConfig, error) {
	cfg, err := yaml.Open(filename)

	if err != nil {
		return BotConfig{}, err
	}

	var b BotConfig

	// The server's information
	b.hostname = cfg.Get("server", "host").(string)
	b.port = cfg.Get("server", "port").(int)
	b.username = cfg.Get("server", "username").(string)
	b.password = cfg.Get("server", "password").(string)

	// Log configuration
	b.logFilename = cfg.Get("log", "path").(string)
	b.level = cfg.Get("log", "level").(string)

	// Bot configuration
	b.name = cfg.Get("bot", "name").(string)

	return b, nil
}

//startLogging creates the application log file.
func startLogging(botConfig BotConfig) {
	spacelog.Setup(AppName, spacelog.SetupConfig{
		Output: botConfig.logFilename,
		Level:  botConfig.level,
	})
}

//loadCertificate loads the server public certificate so we can login.
func loadCertificate(filename string) error {
	cf, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	cpb, _ := pem.Decode(cf)
	cert, err := x509.ParseCertificate(cpb.Bytes)

	if err != nil {
		return err
	}

	roots := x509.NewCertPool()
	roots.AddCert(cert)
	xmpp.DefaultTlsConfig.RootCAs = roots

	return nil
}

//establishChatConnection establishes the XMPP connection with the server.
func establishChatConnection(certFilename string, botConfig BotConfig) (*xmpp.Client, error) {
	if err := loadCertificate(certFilename); err != nil {
		fmt.Println(err)
		return nil, err
	}

	/*	writer := func() io.Writer {
		if botConfig.level == "debug" {
			logger := spacelog.GetLogger()
			return logger.Writer(spacelog.Debug)
		}

		return nil
	}()*/

	xmppOptions := xmpp.Options{
		Address:  fmt.Sprintf("%s:%d", botConfig.hostname, botConfig.port),
		Jid:      botConfig.username,
		Password: botConfig.password,
		//PacketLogger: os.Stdout,
		PacketLogger: nil,
		Retry:        3,
	}

	client, err := xmpp.NewClient(xmppOptions)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = client.Connect()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return client, nil
}

//handleMessage does the job of calling the bot message handling function and give
//back an answer.
func handleMessage(client *xmpp.Client, bot BotModel, packet *xmpp.ClientMessage, session ChatSession) {
	// First, we let the bot do its job
	answers := bot.HandleMessage(*packet, session)

	// And now we can send the answer
	for _, answer := range answers {
		reply := xmpp.ClientMessage{Packet: xmpp.Packet{To: packet.From}, Body: answer.Message}
		client.Send(reply.XMPPFormat())
	}
}

//receiveMessages is the function to receive all incoming messages to us.
//It also notifies through the chanSession channel when a new chat message just arrived.
func receiveMessages(client *xmpp.Client, chanSession chan ChatSession, bot BotModel, sessions map[string]ChatSession) {
	for packet := range client.Recv() {
		switch packet := packet.(type) {
		case *xmpp.ClientMessage:
			if packet.Type != "chat" {
				/* Ignores unsupported messages */
				break
			}

			chanSession <- ChatSession{
				lastActivity: time.Now(),
				Name:         packet.From,
				ChatMessage: ChatMessage{
					Type: Message,
					From: packet.From,
					To:   packet.To,
				},
			}

			fmt.Fprintf(os.Stdout, "Body = %s - from = %s\n", packet.Body, packet.From)
			go handleMessage(client, bot, packet, sessions[packet.From])
		}
	}
}

func main() {
	var options CLIOptions
	sessions := make(map[string]ChatSession)

	buildCLIOptions(&options)
	flag.Parse()

	if *options.version {
		fmt.Printf("%s - version %s\n", AppName, Version)
		os.Exit(0)
	}

	if len(*options.botFilename) == 0 {
		fmt.Printf("Error: Must inform the bot configuration file name\n")
		os.Exit(-1)
	}

	botConfig, err := loadBotConfiguration(*options.botFilename)

	if err != nil {
		fmt.Printf("Error: While loading the bot configuration file.\n")
		os.Exit(-1)
	}

	// create the Bot
	startLogging(botConfig)
	bot, err := CreateBot(botConfig.name)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// make the XMPP connection
	client, err := establishChatConnection(*options.certFilename, botConfig)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	chanSession := make(chan ChatSession)
	go receiveMessages(client, chanSession, bot, sessions)

	for {
		select {
		case session := <-chanSession:
			fmt.Println("Adding session for user:", session.Name)
			sessions[session.Name] = session

		case <-time.After(5 * time.Second):
			DeleteInactives(sessions, 1)
			fmt.Println(len(sessions))
		}
	}
}
