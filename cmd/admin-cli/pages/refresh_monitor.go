package pages

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func getJobData() string {
	base := viper.GetString("local_base_url")
	data := string(GetData(base, "/jobs"))
	return data
}

var is_refreshing = false

func Ping(wait chan int) {
	for {
		if is_refreshing {
			wait <- 1
		}
		time.Sleep(1 * time.Second)
	}
}

func RefreshMonitor(wait chan int, mainflex *tview.Flex, g *GUI) {
	first := true
	for {

		if first || is_refreshing {
			first = false
			mainflex.Clear()
			upperflex := tview.NewFlex()
			upperflex.SetDirection(tview.FlexColumn)
			mainflex.AddItem(upperflex, 0, 1, true)

			views := make([]*tview.TextView, 0)

			for range 3 {
				t := tview.NewTextView()
				//t.SetChangedFunc(func() { g.App.Draw() })
				t.SetBorder(true)
				upperflex.AddItem(t, 0, 1, false)
				views = append(views, t)
			}

			count := 0
			gjson.Get(getJobData(), "counts").ForEach(func(k gjson.Result, v gjson.Result) bool {
				t := views[count]
				fmt.Fprintf(t, "%s\n-----------\n", k.String())
				v.ForEach(func(k, v gjson.Result) bool {
					fmt.Fprintf(t, "%s\t%s\n", k.String(), v.String())
					return true
				})
				fmt.Fprintf(t, "\n\n")
				count = (count + 1) % 3

				return true
			})
		}

		t := tview.NewTextView()
		t.SetBorder(true)
		t.SetDynamicColors(true)

		fmt.Fprintf(t, "Press [yellow]Enter[white] to toggle refresh, [yellow](m)[white] for main menu.\n")
		onoff := ""

		if is_refreshing {
			onoff = "on"
		} else {
			onoff = "off"
		}
		fmt.Fprintf(t, "Refresh is [green]"+onoff+"[white]\n")
		mainflex.AddItem(t, 5, 1, false)

		g.App.Draw()
		<-wait
	}
}

func (g *GUI) MonitorJobQueues(name string) {
	var wait = make(chan int, 1)

	flex := tview.NewFlex().SetFullScreen(true)
	flex.SetDirection(tview.FlexRow)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch e := event.Key(); true {
		case e == tcell.KeyEnter:
			is_refreshing = !is_refreshing
		case tcell.Key(event.Rune()) == 'm':
			g.Pages.SwitchToPage(fmt.Sprintf("page-%d", g.PagesByName["main"]))
		}
		return nil
	})

	go Ping(wait)
	go RefreshMonitor(wait, flex, g)

	g.Pages.AddPage(
		fmt.Sprintf("page-%d", g.page_counter),
		flex,
		true,
		g.page_counter == 0)

	g.PagesByName[name] = g.page_counter
	g.page_counter += 1
}
