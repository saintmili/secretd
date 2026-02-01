package main

import (
	"fmt"
	"log"
	"os"

	"github.com/saintmili/secretd/internal/app"
	"github.com/saintmili/secretd/internal/config"
	"github.com/saintmili/secretd/internal/lock"
)

var version = "v1.0.0"

func main() {
	// if err := execute(); err != nil {
	// 	log.Fatalln(err)
	// }
	cfg, warnings, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	for _, w := range warnings {
		fmt.Printf("⚠️  config.%s: %s\n", w.Field, w.Message)
	}

	appInstance, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	defer appInstance.Logger.Close()

	appInstance.Logger.Info("secretd starting")

	// only one instance at a time
	lock, err := lock.Aquire()
	if err != nil {
		log.Fatalln(err)
	}
	defer lock.Release()

	if len(os.Args) < 2 {
		if err := app.PrintUsage(); err != nil {
			log.Fatalln(err)
		}
	}

	switch os.Args[1] {
	case "init":
		if err := app.Init(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "unlock":
		if err := app.Unlock(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "add":
		if err := app.Add(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "list":
		if err := app.List(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "show":
		if err := app.Show(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "generate":
		if err := app.Generate(); err != nil {
			log.Fatalln(err)
		}
	case "update":
		if err := app.Update(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "delete":
		if err := app.Delete(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "change-master-password":
		if err := app.ChangeMasterPassword(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "doctor":
		if err := app.Doctor(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "help":
		if err := app.PrintUsage(); err != nil {
			log.Fatalln(err)
		}
	case "export":
		if err := app.Export(appInstance); err != nil {
			log.Fatalln(err)
		}
	case "config":
		if err := app.Config(*appInstance.Config); err != nil {
			log.Fatalln(err)
		}
	case "version":
		fmt.Printf("secretd %s\n", version)
	default:
		fmt.Println("Unknown command")
		if err := app.PrintUsage(); err != nil {
			log.Fatalln(err)
		}
	}
	os.Exit(1)
}
