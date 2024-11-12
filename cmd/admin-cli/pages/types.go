package pages

import "github.com/rivo/tview"

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
