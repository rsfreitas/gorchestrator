//A chat bot orchestrator ;-)
package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/processone/gox/xmpp"
	"menteslibres.net/gosexy/yaml"
)

const AppName string = "orchestrator"
const Version string = "0.1"

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

type BotConfig struct {
	ConfigServer
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

	return b, nil
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

func main() {
	var options CLIOptions

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

	if err := loadCertificate(*options.certFilename); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	xmppOptions := xmpp.Options{
		Address:      fmt.Sprintf("%s:%d", botConfig.hostname, botConfig.port),
		Jid:          botConfig.username,
		Password:     botConfig.password,
		PacketLogger: os.Stdout,
		Retry:        3,
	}

	client, err := xmpp.NewClient(xmppOptions)

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	_, err = client.Connect()

	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Iterator to receive packets coming from our XMPP connection
	for packet := range client.Recv() {
		switch packet := packet.(type) {
		case *xmpp.ClientMessage:
			fmt.Fprintf(os.Stdout, "Body = %s - from = %s\n", packet.Body, packet.From)
			//			reply := xmpp.ClientMessage{Packet: xmpp.Packet{To: packet.From}, Body: packet.Body}
			//			client.Send(reply.XMPPFormat())
		default:
			fmt.Fprintf(os.Stdout, "Ignoring packet: %T\n", packet)
		}
	}
}
