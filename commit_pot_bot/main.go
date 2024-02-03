package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	cron2 "github.com/robfig/cron"
	"os"
	"os/signal"
	"syscall"
)

var cron cron2.Cron

func main() {
	e := echo.New()
	cron = *cron2.New()
	dg, err := discordgo.New("Bot ${process.env.DISCORD_BOT_TOKEN}")
	if err != nil {
		fmt.Println("[ERROR] : 디스코드 세션 생성에 실패했습니다, ", err)
		return
	}
	defer dg.Close()
	cron.Start()

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("[ERROR] : 연결에 실패했습니다, ", err)
		return
	}

	fmt.Println("돌아가는중...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "김연우 스트릭좀 관리해봅시다")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
