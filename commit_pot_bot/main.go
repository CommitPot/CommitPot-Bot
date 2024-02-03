package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	cron2 "github.com/robfig/cron"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var cron cron2.Cron

func test(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(m.Content, "연우") && strings.Contains(m.Content, "벌금") {
		s.ChannelMessageSend(m.ChannelID, "김연우님의 벌금은 2500원입니다.\n또한 다음부터 스트릭이 깨질 시 2500원씩 벌금이 추가됩니다.")
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func main() {
	e := echo.New()
	cron = *cron2.New()
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("[ERROR] : 디스코드 세션 생성에 실패했습니다, ", err)
		return
	}
	dg.AddHandler(test)
	defer dg.Close()
	cron.Start()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "김연우 스트릭좀 관리해봅시다")
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("[ERROR] : 연결에 실패했습니다, ", err)
		return
	}

	fmt.Println("돌아가는중...")
	e.Logger.Fatal(e.Start(":8080"))
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
