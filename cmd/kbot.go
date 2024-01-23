package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
)

var (
	tracer trace.Tracer
)

func myFunction() {
	ctx, span := tracer.Start(context.Background(), "myFunction")
	defer span.End()

	startTime := time.Now()

	// Your code

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	log.Printf("Operation took %s", elapsed)
}

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A simple bot created for demonstration purposes",
	Long: `Kbot is a simple bot created for demonstration purposes. It can respond to the following commands:
- /start: Start the bot
- /time: Display the current time`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started\n", appVersion)
		bot, err := tgbotapi.NewBotAPI(TeleToken)

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		bot.Debug = true

		log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

			updates, err := bot.GetUpdatesChan(u)

			for update := range updates {
				if update.Message == nil { // ignore any non-Message Updates
					continue
				}

				switch update.Message.Text {
				case "/start":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! You have started the bot.")
					bot.Send(msg)
				case "/time":
					currentTime := time.Now().Format("15:04:05")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("The current time is: %s", currentTime))
					bot.Send(msg)
				default:
					// do nothing
				}

				// Call function to measure execution time
				myFunction()
			}
		},
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Set up OpenTelemetry
	ctx := context.Background()

	tp := sdktrace.NewTracerProvider()

	otel.SetTracerProvider(tp)

	tracer = otel.Tracer("kbot")
}
