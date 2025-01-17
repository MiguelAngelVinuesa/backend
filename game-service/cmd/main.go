package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/net/kafka"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/handlers"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/events"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/hashes"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/utils"
)

func main() {
	appStart := time.Now().UTC()

	// fix-up version constants
	if consts.SemBuild == "" {
		consts.SemBuild = consts.ValueDev
		consts.SemVerFull += consts.SemBuild
	}
	if consts.SemDate == "" {
		consts.SemDate = time.Now().UTC().Format(time.RFC3339)
	}

	// init logging.
	log.Init()
	log.Logger.Info(consts.MsgBinaryHashes, hashes.MainFile, hashes.MainHash, hashes.RngLib, hashes.RngLibHash, hashes.RngInclude, hashes.RngIncludeHash)
	log.Logger.Info(consts.MsgBinaryModDate, hashes.MainFile, hashes.MainModDate.Format(time.RFC3339), hashes.RngLib, hashes.RngLibModDate.Format(time.RFC3339))
	log.Logger.Info(consts.MsgBinaryFileSize, hashes.MainFile, hashes.MainFileSize, hashes.RngLib, hashes.RngLibFileSize)
	log.Logger.Info(consts.MsgGoModules, consts.FieldModules, hashes.Modules)

	// init global metrics.
	metrics.InitMetrics()

	// init message consumer.
	if config.EventsTopic != "" || config.MessagesTopic != "" {
		topics := make([]string, 0, 2)
		if config.EventsTopic != "" {
			topics = append(topics, config.EventsTopic)
		}
		if config.MessagesTopic != "" {
			topics = append(topics, config.MessagesTopic)
		}

		kafka.NewConsumer(
			kafka.WithBrokers(config.MqBrokers...),
			kafka.WithTopics(topics...),
			kafka.FromTimestamp(appStart.Add(-15*time.Minute)),
			kafka.WithClient(events.Consumer),
			kafka.WithLogger(log.Logger),
			kafka.WithDone(utils.Final().Done()),
		)
	}

	// calculate game config hashes.
	hashes.InitGameHashes()
	go hashes.SendReport(appStart)
	log.Logger.Info(consts.MsgGameConfigs, consts.FieldGames, hashes.GameHashes)

	// start the http server.
	app := initFiber()
	go func() {
		log.Logger.Info(consts.MsgServerStarting, consts.FieldHost, config.Server, consts.FieldPort, config.Port)
		log.Logger.Info(consts.MsgHttp, consts.FieldRoutes, app.Stack())
		if err := app.Listen(fmt.Sprintf("%s:%d", config.Server, config.Port)); err != nil {
			log.Logger.Panic(consts.MsgServiceFailed, consts.FieldError, err)
		}
	}()

	time.AfterFunc(250*time.Millisecond, func() { log.Logger.Info(consts.MsgServerStarted) })

	// shutdown on SIGINT.
	intChan := make(chan os.Signal)
	signal.Notify(intChan, syscall.SIGINT, syscall.SIGTERM)
	<-intChan

	utils.Final().Shutdown()

	go func() {
		log.Logger.Info(consts.MsgServerShuttingDown)
		app.Shutdown()
		// final metrics
		metrics.Metrics.Print()
		time.Sleep(100 * time.Millisecond)
		log.Logger.Info(consts.MsgServerShutdown)
		os.Exit(0)
	}()

	// Another signal will force process termination.
	signal.Notify(intChan, syscall.SIGINT, syscall.SIGTERM)
	<-intChan
	os.Exit(0)
}

func initFiber() *fiber.App {
	// initialize the fast http server.
	app := fiber.New(fiber.Config{
		CaseSensitive:         true,
		DisableDefaultDate:    true,
		DisableStartupMessage: true,
		BodyLimit:             64 * 1024,
		ReadTimeout:           10 * time.Second,
		WriteTimeout:          10 * time.Second,
		IdleTimeout:           config.ConnCleanup,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	app.Server().Logger = &log.InfoPrinter{}

	if !config.NoDefaultHeaders {
		app.Use(func(ctx *fiber.Ctx) error {
			ctx.Set("X-XSS-Protection", "1; mode=block")
			ctx.Set("X-Content-Type-Options", "nosniff")
			ctx.Set("X-Download-Options", "noopen")
			ctx.Set("Strict-Transport-Security", "max-age=5184000")
			ctx.Set("X-Frame-Options", "SAMEORIGIN")
			ctx.Set("X-DNS-Prefetch-Control", "off")
			return ctx.Next()
		})
	}
	if !config.NoCors {
		app.Use(cors.New())
	}
	if !config.NoCompression {
		app.Use(compress.New())
	}

	app.Get(consts.PathPing, handlers.Ping)
	app.Get(consts.PathBinHashes, handlers.GetBinHashes)
	app.Get(consts.PathGameHash, handlers.GameHash)
	app.Get(consts.PathGameInfo, handlers.GetGameInfo)
	app.Post(consts.PathPreferences, handlers.PostPreferences)
	app.Get(consts.PathCcbFlags, handlers.GetCcbFlags)
	app.Post(consts.PathStrings, handlers.PostStrings)
	app.Get(consts.PathPlural, handlers.GetPluralString)
	app.Post(consts.PathPlurals, handlers.PostPluralStrings)
	app.Get(consts.PathMessages, handlers.GetMessages)

	app.Post(consts.PathRound, handlers.PostRound)
	app.Post(consts.PathRoundPaid, handlers.PostRoundPaid)
	app.Post(consts.PathRoundSecond, handlers.PostRoundSecond)
	app.Post(consts.PathRoundResume, handlers.PostRoundResume)
	app.Post(consts.PathRoundNext, handlers.PostRoundNext)
	app.Post(consts.PathRoundFinish, handlers.PostRoundFinish)

	app.Get(consts.PathSessionInfo, handlers.GetSessionInfo)

	// SUPERVISED-BUILD-REMOVE-START
	if config.DebugMode {
		app.Post(consts.PathRoundDebug, handlers.PostRoundDebug)
		app.Post(consts.PathRoundDebugSecond, handlers.PostRoundDebugSecond)
		app.Post(consts.PathRoundDebugResume, handlers.PostRoundDebugResume)
	}

	if config.Environment == consts.ValueDev {
		app.Get(consts.PathRngConditionsLU, handlers.GetRngConditionsLU)
		app.Post(consts.PathRngMagicTest, handlers.PostRngMagicTest)
		app.Post(consts.PathRngMagic, handlers.PostRngMagic)
	}
	// SUPERVISED-BUILD-REMOVE-END

	go func() {
		ticker := time.NewTicker(metrics.PrintMetrics)
		fh := app.Server()
		for {
			select {
			case <-ticker.C:
				log.Logger.Info(consts.MsgHttp,
					consts.FieldFiberHandlers, app.HandlersCount(),
					consts.FieldConnections, fh.GetOpenConnectionsCount(),
					consts.FieldConcurrency, fh.GetCurrentConcurrency())
			}
		}
	}()

	return app
}
