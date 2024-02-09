package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
	"github.com/fedesog/webdriver"
	"github.com/labstack/echo/v4"
	cron2 "github.com/robfig/cron"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type NoSolveList struct {
	People []string `json:"people"`
}

type Request struct {
	Selector string `json:"selector"`
}

func webhook(isToday bool) string {
	selector := new(Request)
	if isToday {
		selector.Selector = ".css-fpwzir > svg > rect:nth-child(13)"
	} else {
		selector.Selector = ".css-fpwzir > svg > rect:nth-child(14)"
	}
	bodyBytes, _ := json.Marshal(selector)
	req, _ := http.NewRequest("POST", os.Getenv("BASE_URL")+"/webhook", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] : " + err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("[ERROR] : " + err.Error())
	}
	rs := NoSolveList{}
	json.Unmarshal(body, &rs)
	returnString := ""
	if len(rs.People) == 0 {
		switch isToday {
		case true:
			return "오늘 모든 친구들이 백준을 풀었습니다."
		case false:
			return "어제 모든 친구들이 백준을 풀었습니다."
		}
	}
	for _, v := range rs.People {
		returnString += v + " "
	}
	if isToday {
		return "아직 안푼사람 : " + returnString
	}
	return "어제 안푼사람 : " + returnString
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
	cron.AddFunc("0 0 22 * * *", func() {
		fmt.Println("저녁 10시 웹훅을 실행합니다.")
		dg.ChannelMessageSend(os.Getenv("COMMIT_POT_DISCORD_CHANNEL_ID"), webhook(true))
	})
	cron.AddFunc("0 0 9 * * *", func() {
		fmt.Println("오전 9시 웹훅을 실행합니다.")
		dg.ChannelMessageSend(os.Getenv("COMMIT_POT_DISCORD_CHANNEL_ID"), webhook(false))
	})
	cron.Start()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "김연우 스트릭좀 관리해봅시다")
	})

	e.POST("/webhook", func(c echo.Context) error {
		var req Request
		_ = c.Bind(&req)
		commitPotList := []string{"jeongho1209", "chanhong1206", "bjcho1503", "jikwan12", "yeon8747", "chul8747"}
		n := NoSolveList{}
		chromeDriver := webdriver.NewChromeDriver("./chromedriver-mac-arm64/chromedriver")
		err := chromeDriver.Start()
		checkErr(err)
		desired := webdriver.Capabilities{"Platform": "linux"}
		required := webdriver.Capabilities{}
		for _, v := range commitPotList {
			session, err := chromeDriver.NewSession(desired, required)
			checkErr(err)
			err = session.Url("https://solved.ac/profile/" + v)
			checkErr(err)
			resp, err := session.Source()
			checkErr(err)
			htmlNode, err := html.Parse(strings.NewReader(resp))
			checkErr(err)
			doc := goquery.NewDocumentFromNode(htmlNode)
			val, exist := doc.Find(req.Selector).Attr("fill")
			if !exist {
				n.People = append(n.People, v+"(스트릭 프리즈)")
			}
			if val == "#dddfe0" {
				n.People = append(n.People, v)
			}
			session.Delete()
		}
		chromeDriver.Stop()
		return c.JSON(200, n)
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("[ERROR] : 연결에 실패했습니다, ", err)
		return
	}

	fmt.Println("돌아가는중...")
	e.Logger.Fatal(e.Start(":8080"))
}
