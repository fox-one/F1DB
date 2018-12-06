package main

import (
	"context"
	"flag"
	"log"

	account "github.com/fox-one/f1db/account"
	config "github.com/fox-one/f1db/config"
	storage "github.com/fox-one/f1db/storage"
)

func main() {
	mode := flag.String("m", "serve", "Run Mode")
	flag.Parse()

	ctx := context.Background()

	config.Init()
	account.InitSession()
	account.SetupBroker(config.GetConfig().Broker)
	pk, err := account.GetFoxPublicKey(ctx)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("use pk: %s\n", pk)
	log.Printf("Run Mode: %s\n", *mode)

	if *mode == "serve" {
		storage.InitIpfs()
		serve(ctx, pk)
	} else if *mode == "register" {
		if newUserID, err := account.Register(ctx, pk); err == nil {
			log.Printf("New User:\n- ID: %s\n", newUserID)
		} else {
			log.Panic(err)
		}
	} else {
		log.Println("Bye!")
	}

}
