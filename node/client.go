package node

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-oracle/config"
	"github.com/ontio/ontology-oracle/core"
	"github.com/ontio/ontology-oracle/log"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/password"
	"github.com/urfave/cli"
)

// Client is the shell for the node. It has fields for the
// AppFactory (the services application).
type Client struct {
	AppFactory AppFactory
}

// RunNode starts the oracle node.
func (client *Client) RunNode(c *cli.Context) {

	log.Info("Starting Oracle Node... ")
	log.Info("Open the account")
	if !common.FileExisted(config.Configuration.WalletFile) {
		log.Fatal(fmt.Sprintf("No %s detected, please create a wallet.", config.Configuration.WalletFile))
		os.Exit(1)
	}
	ontSdk := sdk.NewOntologySdk()
	wallet, err := ontSdk.OpenWallet(config.Configuration.WalletFile)
	if err != nil {
		log.Fatalf("Can't open local wallet: %s", err)
		os.Exit(1)
	}
	pwd, err := password.GetPassword()
	if err != nil {
		log.Fatalf("password.GetPassword erro: %sr", err)
		os.Exit(1)
	}
	acct, err := wallet.GetDefaultAccount(pwd)
	if err != nil {
		log.Fatalf("Can't get default account: %s", err)
		os.Exit(1)
	}

	app := client.AppFactory.NewApplication(acct)

	app.Start()
	defer app.Stop()

	waitToExit()
}

// AppFactory implements the NewApplication method.
type AppFactory interface {
	NewApplication(*sdk.Account) core.Application
}

// AppFactory is used to create a new Application.
type OracleAppFactory struct{}

// NewApplication returns a new instance of the node with the given config.
func (n OracleAppFactory) NewApplication(account *sdk.Account) core.Application {
	return core.NewApplication(account)
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("Ontology Oracle received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
