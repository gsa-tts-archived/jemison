package pages

import (
	"fmt"

	"github.com/rivo/tview"
)

func (g *GUI) MainPage(name string) {

	list := tview.NewList().
		AddItem("query", "query the search server", 'a', func() {

		}).
		AddItem("crawl", "crawl a site", 'b', func() {
			g.Pages.SwitchToPage(fmt.Sprintf("page-%d", g.PagesByName["crawl"]))
		}).
		AddItem("monitor", "monitor the job queue", 'm', func() {
			is_refreshing = true
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
