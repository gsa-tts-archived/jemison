package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

type GUI struct {
	App          *tview.Application
	Pages        *tview.Pages
	PagesByName  map[string]int
	page_counter int
}

func NewGui() *GUI {
	return &GUI{
		Pages:        tview.NewPages(),
		PagesByName:  make(map[string]int),
		page_counter: 0,
		App:          tview.NewApplication(),
	}
}

func (g *GUI) Run() {
	if err := g.App.SetRoot(g.Pages, false).Run(); err != nil {
		panic(err)
	}
}

func PutData(base_url string, data string) {
	log.Println(base_url, data)

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

func GetData(base_url string, path string) []byte {
	resp, err := http.Get(base_url + path)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return b
}

func (g *GUI) CrawlPage(name string) {

	base := viper.GetString("local_base_url")

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

var monitoring = false

func (g *GUI) MonitorJobQueues(name string) {
	if monitoring {
		base := viper.GetString("local_base_url")
		data := string(GetData(base, "/jobs"))
		log.Println(data)
	}
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
		AddItem("monitor", "monitor the job queue", 'm', func() {
			monitoring = true
			g.Pages.SwitchToPage(fmt.Sprintf("page-%d", g.PagesByName["monitor"]))
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			g.App.Stop()
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
