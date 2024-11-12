package pages

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
)

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
