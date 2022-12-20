package main

import (
	"time"
	g "videown-server/global"

	"videown-server/internal/db"
	ginlet "videown-server/internal/ginlet"
	"videown-server/internal/model"
	"videown-server/internal/service"
	"videown-server/internal/service/nft"

	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"videown-server/pkg/chain"
	"videown-server/pkg/logger"
	"videown-server/pkg/setting"

	"github.com/CESSProject/cess-oss/configs"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/urfave/cli"
)

var Version = "v1.1.0"

var timeoutSec = 20

func main() {
	if err := setupApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var configDir string

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "Videown Server"
	app.Action = doSetup
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Value:       "./configs/",
			Usage:       "config director for the application",
			Destination: &configDir,
		}}
	app.Commands = []cli.Command{}
	app.Before = func(_ *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func doSetup(_ *cli.Context) {
	if err := setupSetting(); err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
		return
	}

	if err := setupLogger(); err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
		return
	}

	if err := setupDbEngine(); err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
		return
	}
	if err := buildChain(); err != nil {
		log.Fatalf("init.buildChainClient err: %v", err)
		return
	}
	service.SetupService(g.Settings, g.GormDb, g.Logger)
	initCmpConfig()
	setupGin()
	signalHandle()
}

func setupGin() {
	gin.SetMode(g.Settings.ServerSetting.RunMode)
	routerHandler := ginlet.NewRouter()
	httpServer := &http.Server{
		Addr:           ":" + g.Settings.ServerSetting.HttpPort,
		Handler:        routerHandler,
		ReadTimeout:    g.Settings.ServerSetting.ReadTimeout,
		WriteTimeout:   g.Settings.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	os.MkdirAll(g.COVER_IMAGE_PATH, os.ModeDir)
	nft.InitDefaultStatusListener(3, 1024)
}

func setupSetting() error {
	s, err := setting.NewSettingsWithDirectory(configDir)
	if err != nil {
		return err
	}
	g.Settings = s
	return nil
}

func setupLogger() error {
	g.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  g.Settings.AppSetting.FullLogFilePath(),
		MaxSize:   500, // Size of a single file
		MaxAge:    10,  // Number of backup files
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	os.MkdirAll("storage/logs", os.ModeDir)
	return nil
}

func initCmpConfig() {
	nft.CmpBaseUrl = g.Settings.AppSetting.CmpHttpUrl
	fmt.Println("init cmp base url:", nft.CmpBaseUrl)
}

func setupDbEngine() error {
	var err error
	g.GormDb, err = db.NewGormDbForMySql(g.Settings.DatabaseSetting, g.Settings.ServerSetting.IsDebugMode())
	if err != nil {
		panic("database init err! " + err.Error())
	}
	model.AutoMigrate(g.GormDb)
	return err
}

func signalHandle() {
	fmt.Println("videown server startup success!")
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Printf("get a signal: %s, stop the videown server process\n", si.String())
			//txresult.Shutdown()
			fmt.Println("videown server shutdown success!")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func buildChain() error {
	// connecting chain
	var err error
	chain.ChainCli, err = chain.NewChainClient(
		g.Settings.Web3Setting.RpcEndpoints[0],
		g.Settings.Web3Setting.Mnemonic,
		time.Duration(timeoutSec*int(time.Second)),
	)
	if err != nil {
		return err
	}

	// judge the balance
	// accountinfo, err := chain.ChainCli.GetAccountInfo(chain.ChainCli.GetPublicKey())
	// if err != nil {
	// 	return err
	// }

	// if accountinfo.Data.Free.CmpAbs(new(big.Int).SetUint64(configs.MinimumBalance)) == -1 {
	// 	return fmt.Errorf("account balance is less than %v pico", configs.MinimumBalance)
	// }

	// sync block
	for {
		ok, err := chain.ChainCli.GetSyncStatus()
		if err != nil {
			return err
		}
		if !ok {
			break
		}
		fmt.Println("In sync block...")
		time.Sleep(time.Second * configs.BlockInterval)
	}
	chain.InitRpcWorkPool()
	fmt.Println("Complete synchronization of primary network block data")
	fmt.Println("building chain success!")
	return nil
}
