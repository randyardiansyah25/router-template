package main

import (
	"fmt"
	"os"
	"os/signal"
	"router-template/delivery"
	"router-template/delivery/broker"
	"router-template/delivery/broker/consumer"
	"router-template/delivery/http/router"
	"router-template/entities/app"
	"router-template/entities/common"
	"router-template/repository/built_in/databasefactory"
	"router-template/repository/built_in/keyvaluefactory"
	"runtime"
	"syscall"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/kpango/glg"
	"github.com/randyardiansyah25/libpkg/util/env"
)

func main() {
	go delivery.StartPrintoutObserver()
	isUse := env.GetBool("rabbit.use", false)
	if isUse {
		broker.ConnectToRabbit()
		go consumer.Start()
		go broker.BrokerClosedChannelObserver()
	}
	router.Start()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if loc, er := time.LoadLocation("Asia/Jakarta"); er == nil {
		glg.Get().SetTimeLocation(loc)
	}
	LoadConfiguration(false)

	SetLoggerFileWriter()
	go StartLoggerRotation()

	if os.Getenv("app.database_driver") != "" {
		PrepareRepo()
	}

	go ReloadObserver()
}

func LoadConfiguration(isReload bool) {
	var er error

	if isReload {
		_ = glg.Log("Reloading configuration file...")
		er = godotenv.Overload(".env")
	} else {
		_ = glg.Log("Loading configuration file...")
		er = godotenv.Load(".env")
	}

	if er != nil {
		_ = glg.Error("Configuration file not found...")
		os.Exit(1)
	}

}

func SetLoggerFileWriter() {
	//Opsi agar log utk level LOG, DEBUG, INFO dicatat atau tidak
	//Jika menggunakan docker atau dibuatkan service, log sudah dibuatkan, sehingga direkomendasikan
	//app log di set false
	useAppLog := env.GetBool("app.log", false)
	logPath := fmt.Sprintf("%s%c%s", app.LogDir, os.PathSeparator, app.LogFileName)
	if useAppLog {
		log := glg.FileWriter(fmt.Sprintf("%s.log", logPath), 0766)
		glg.Get().
			SetMode(glg.BOTH).
			AddLevelWriter(glg.LOG, log).
			AddLevelWriter(glg.DEBG, log).
			AddLevelWriter(glg.INFO, log).
			AddLevelWriter(glg.ERR, log).
			AddLevelWriter(glg.WARN, log)
	}

	//Untuk error akan selalu dicatat dalam file
	logEr := glg.FileWriter(fmt.Sprintf("%s.err", logPath), 0766)
	glg.Get().
		SetMode(glg.BOTH).
		AddLevelWriter(glg.ERR, logEr).
		AddLevelWriter(glg.WARN, logEr)

	glg.Info("App logging started..")
}

func PrepareRepo() {
	var er error
	if os.Getenv("app.database_driver") != "" {
		databasefactory.AppDb, er = databasefactory.GetDatabase()
		if er != nil {
			glg.Fatal(er.Error())
		}

		_ = glg.Log("Connecting to database...")
		if er = databasefactory.AppDb.Connect(); er != nil {
			_ = glg.Error("Connection to database failed : ", er.Error())
			os.Exit(1)
		}

		if er = databasefactory.AppDb.Ping(); er != nil {
			_ = glg.Error("Cannot ping database : ", er.Error())
			os.Exit(1)
		}

		_ = glg.Log("Database Connected")
	}

	if os.Getenv("app.keyvalue_driver") != "" {
		keyvaluefactory.AppStore, er = keyvaluefactory.GetStore()
		if er != nil {
			glg.Fatal(er.Error())
		}
		_ = glg.Log("Connecting to keyvalue store...")
		if er = keyvaluefactory.AppStore.Open(); er != nil {
			_ = glg.Error("failed to open keyvalue store : ", er.Error())
			os.Exit(1)
		}

		if er = keyvaluefactory.AppStore.Echo(); er != nil {
			_ = glg.Error("failed to echo to keyvalue store : ", er.Error())
			os.Exit(1)
		}

		_ = glg.Log("Key value storage ready...")

	}
}

func ReloadObserver() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGHUP)

	func() {
		for {
			<-sign
			glg.Info("Received SIGHUP, reloading configuration...")
			LoadConfiguration(true)
		}
	}()
}

func StartLoggerRotation() {
	at := env.GetString("app.log_rotation_at")
	if at != "" {
		gocron.Every(1).Day().At(at).Do(func() {
			glg.Info("Stopping logging temporarily for backup..")
			glg.Reset()
			suffixFile := time.Now().Format(app.DATE_FORMAT_DATEONLY)
			logPath := fmt.Sprintf("%s%c%s", app.LogDir, os.PathSeparator, app.LogFileName)
			backupFile := fmt.Sprintf("%s_%s", logPath, suffixFile)

			srcLog := fmt.Sprintf("%s.log", logPath)
			destLog := fmt.Sprintf("%s.log", backupFile)
			if env.GetBool("app.log", false) {
				if er := common.CopyAndDelete(srcLog, destLog); er != nil {
					glg.Info(er)
				}
			}
			erSrcLog := fmt.Sprintf("%s.err", logPath)
			erDestLog := fmt.Sprintf("%s.err", backupFile)
			if er := common.CopyAndDelete(erSrcLog, erDestLog); er != nil {
				glg.Info(er)
			}

			SetLoggerFileWriter()
			glg.Infof("The log has been truncated. Check the backup files: %s, %s ", destLog, erDestLog)
		})
		<-gocron.Start()
	}
}
