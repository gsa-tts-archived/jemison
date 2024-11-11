package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/rivo/tview"
	"github.com/tidwall/sjson"
)

const refreshInterval = 1000 * time.Millisecond

var (
	view *tview.Modal
	app  *tview.Application
)

func currentTimeString() string {
	t := time.Now()
	return fmt.Sprintf(t.Format("Current time is 15:04:05"))
}

func updateTime() {
	for {
		time.Sleep(refreshInterval)
		app.QueueUpdateDraw(func() {
			view.SetText(currentTimeString())
		})
	}
}

type GUI struct {
	Pages        *tview.Pages
	PagesByName  map[string]int
	page_counter int
}

func NewGui() *GUI {
	return &GUI{
		Pages:        tview.NewPages(),
		PagesByName:  make(map[string]int),
		page_counter: 0,
	}
}

func PutData(base_url string, data string) {
	log.Println(base_url, data)
	data += "x"

	req, err := http.NewRequest(http.MethodPut,
		base_url,
		bytes.NewBuffer([]byte(data)))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
}

func (g *GUI) CrawlPage(name string) {
	s, _ := env.Env.GetUserService("admin")
	base := s.GetParamString("api_base_url")
	schemeDd := tview.NewDropDown().
		SetOptions([]string{"http", "https"}, nil).
		SetLabel("Scheme")
	hostInput := tview.NewInputField().SetLabel("Host")
	pathInput := tview.NewInputField().SetLabel("Path")

	form := tview.NewForm().
		AddFormItem(schemeDd).
		AddFormItem(hostInput).
		AddFormItem(pathInput).
		// AddInputField("Host", "", 20, nil, nil).
		// AddInputField("Path", "", 20, nil, nil).
		AddButton("Go", func() {
			req_url := base + "/crawl"

			_, scheme := schemeDd.GetCurrentOption()
			host := hostInput.GetText()
			path := pathInput.GetText()

			data := `{"scheme": "", "host": "", "path": ""}`
			data, _ = sjson.Set(data, "scheme", scheme)
			data, _ = sjson.Set(data, "host", host)
			data, _ = sjson.Set(data, "path", path)

			PutData(req_url, data)
		}).
		AddButton("Back", func() {})
	form.SetBorder(true).SetTitle("Crawl a site").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)
	flex.AddItem(form, 0, 1, true)

	g.Pages.AddPage(
		fmt.Sprintf("page-%d", g.page_counter),
		flex,
		true,
		g.page_counter == 0)

	g.PagesByName[name] = g.page_counter
	g.page_counter += 1
}

func (g *GUI) MainPage(name string) {

	list := tview.NewList().
		AddItem("query", "query the search server", 'a', func() {

		}).
		AddItem("crawl", "crawl a site", 'b', func() {
			g.Pages.SwitchToPage(fmt.Sprintf("page-%d", g.PagesByName["crawl"]))
		}).
		AddItem("download", "download the db", 'c', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	list.Box.SetBorder(true).SetTitle("jemison admin")

	flex := tview.NewFlex().SetFullScreen(true)
	flex.AddItem(list, 0, 1, true)

	g.Pages.AddPage(
		fmt.Sprintf("page-%d", g.page_counter),
		flex,
		true,
		g.page_counter == 0)

	g.PagesByName[name] = g.page_counter
	g.page_counter += 1
}

var ThisServiceName = "admin_cli"

func main() {
	os.Setenv("ENV", "LOCALHOST")

	env.InitGlobalEnv(ThisServiceName)
	// InitializeQueues()

	app = tview.NewApplication()
	g := NewGui()

	g.MainPage("main")
	g.CrawlPage("crawl")

	if err := app.SetRoot(g.Pages, false).Run(); err != nil {
		panic(err)
	}
}
