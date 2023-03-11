package difiber

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"log"
	"os"
	"os/signal"
	"sync"
)

func initFiber() *fiber.App {
	app := fiber.New(fiber.Config{Prefork: false})
	return app
}

func handleFiber(lc fx.Lifecycle, app *fiber.App) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	var fiberShutdownWG sync.WaitGroup

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(":7500"); err != nil {
					log.Fatal("[!] fiber not started", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			go func() {
				_ = <-ch
				log.Println("Gracefully shutting down...")
				fiberShutdownWG.Add(1)
				defer fiberShutdownWG.Done()
				_ = app.Shutdown()
			}()
			fmt.Println("Running cleanup tasks...")
			return nil
		},
	})

	fiberShutdownWG.Wait()
}

// DI
// sample usage
// fx.New(DI).Run()
var DI = fx.Options(
	fx.Provide(initFiber),
	fx.Invoke(handleFiber),
)

