package main

import (
	"context"
	"fmt"
	"github.com/Xrefullx/YanDip/client/model"
	"github.com/Xrefullx/YanDip/client/pkg"
	"github.com/Xrefullx/YanDip/client/provider/http"
	"github.com/Xrefullx/YanDip/client/services"
	"github.com/Xrefullx/YanDip/client/storage/sqllte"
	"github.com/Xrefullx/YanDip/client/tui"
	"github.com/rivo/tview"
	"log"
	"os"
	"os/signal"
	"time"
)

var app = tview.NewApplication()

func main() {
	fmt.Printf("Build version:%v\n", model.BuildVersion)
	fmt.Printf("Build date:%v\n", model.BuildDate)
	fmt.Printf("Build commit:%v\n", model.BuildCommit)
	cfg, err := pkg.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	provCfg := http.HTTPConfig{
		AuthURL:     "/api/user/login",
		RegisterURL: "/api/user/register",
		SecretURL:   "/api/secret",
		SyncListURL: "/api/sync",
		PingURL:     "/api/ping",
		BaseURL:     cfg.ServerURL,
		Timeout:     time.Millisecond * 500,
	}
	db, err := sqllte.NewStorage(cfg.StorageFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	provider := http.NewHTTPProvider(provCfg)
	secretService := services.NewSecret(cfg, db)
	svcSecret := &secretService
	svcSync := services.NewSyncService(db, provider, cfg)
	if err := svcSync.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
	//  run tui
	tui.SetQ(app, svcSecret)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc
}
