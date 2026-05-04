//go:build !release && cgo

package gui

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/services"
)

const timeFormat = "2-1-2006"

type Gui struct{}

type guiState struct {
	app    fyne.App
	window fyne.Window

	mu     sync.Mutex
	logs   []string
	report services.TotalsReport

	navigation  *widget.List
	screenNames []string
	logOutput   *widget.TextGrid
	darkMode    bool
}

func (g *Gui) Description() string {
	return "Experimental desktop GUI"
}

func (g *Gui) Name() string {
	return "gui"
}

func (g *Gui) Usage() string {
	return "gui"
}

func (g *Gui) Run(args []string) error {
	fs := flag.NewFlagSet("gui", flag.ExitOnError)
	fs.Parse(args)

	a := app.NewWithID("github.com.ahmedsat.ebda-cli.gui")
	w := a.NewWindow("EBDA CLI GUI")
	w.Resize(fyne.NewSize(1100, 700))

	state := &guiState{
		app:         a,
		window:      w,
		screenNames: []string{"Dashboard", "Totals", "Logs"},
	}
	screens := map[string]fyne.CanvasObject{}
	content := container.NewMax()

	showScreen := func(name string) {
		content.Objects = []fyne.CanvasObject{screens[name]}
		content.Refresh()
		state.logf("opened %s screen", name)
	}

	screens["Dashboard"] = state.newDashboardScreen()
	screens["Totals"] = state.newTotalsScreen()
	screens["Logs"] = state.newLogsScreen()

	navigation := widget.NewList(
		func() int { return len(state.screenNames) },
		func() fyne.CanvasObject { return widget.NewLabel("screen") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(state.screenNames[id])
		},
	)
	state.navigation = navigation
	navigation.OnSelected = func(id widget.ListItemID) {
		showScreen(state.screenNames[id])
	}

	split := container.NewHSplit(
		container.NewBorder(
			widget.NewLabel("Screens"),
			nil,
			nil,
			nil,
			navigation,
		),
		content,
	)
	split.Offset = 0.18

	w.SetContent(split)
	state.logf("GUI started")
	state.selectScreen("Dashboard")
	w.ShowAndRun()
	return nil
}

func (s *guiState) newDashboardScreen() fyne.CanvasObject {
	info := widget.NewForm(
		widget.NewFormItem("ERP base URL", widget.NewLabel(config.ErpBaseUrl)),
		widget.NewFormItem("Kobo base URL", widget.NewLabel(config.KoboBaseURL)),
		widget.NewFormItem("DB path", widget.NewLabel(config.DBFilePath)),
		widget.NewFormItem("Notifications", widget.NewLabel(boolLabel(!config.DisableNotify, "enabled", "disabled"))),
		widget.NewFormItem("Settings source", widget.NewLabel("Environment variables")),
	)

	darkMode := widget.NewCheck("Dark mode", func(enabled bool) {
		s.setDarkMode(enabled)
	})
	darkMode.SetChecked(s.darkMode)

	actions := container.NewVBox(
		widget.NewButton("Open Totals", func() { s.selectScreen("Totals") }),
		widget.NewButton("Open Logs", func() { s.selectScreen("Logs") }),
	)

	return container.NewBorder(
		widget.NewLabel("Dashboard"),
		nil,
		nil,
		nil,
		container.NewVBox(
			widget.NewCard("Runtime", "Current environment-backed configuration", info),
			widget.NewCard("Appearance", "Prototype theme controls", darkMode),
			widget.NewCard("Quick Actions", "Prototype navigation", actions),
		),
	)
}

func (s *guiState) newTotalsScreen() fyne.CanvasObject {
	fromEntry := widget.NewEntry()
	fromEntry.SetText("1-1-2022")

	toEntry := widget.NewEntry()
	toEntry.SetText(time.Now().Format(timeFormat))

	status := widget.NewLabel("Idle")
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	table := widget.NewTable(
		func() (int, int) {
			s.mu.Lock()
			defer s.mu.Unlock()
			return len(s.report.Rows) + 2, 4
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapWord
			return label
		},
		func(id widget.TableCellID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(s.cellText(id))
		},
	)
	table.SetColumnWidth(0, 220)
	table.SetColumnWidth(1, 90)
	table.SetColumnWidth(2, 110)
	table.SetColumnWidth(3, 120)

	run := func() {
		from, err := time.Parse(timeFormat, strings.TrimSpace(fromEntry.Text))
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid from date: %w", err), s.window)
			return
		}
		to, err := time.Parse(timeFormat, strings.TrimSpace(toEntry.Text))
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid to date: %w", err), s.window)
			return
		}

		progress.Show()
		status.SetText("Loading totals...")
		s.logf("loading totals report from %s to %s", from.Format(timeFormat), to.Format(timeFormat))

		go func() {
			report, err := services.LoadTotalsReport(from, to)
			if err != nil {
				s.logf("totals failed: %v", err)
				fyne.Do(func() {
					progress.Hide()
					status.SetText("Failed")
					dialog.ShowError(err, s.window)
				})
				return
			}

			s.mu.Lock()
			s.report = report
			s.mu.Unlock()

			s.logf("totals loaded: %d regions, %d farms", len(report.Rows), report.TotalFarms)
			fyne.Do(func() {
				table.Refresh()
				progress.Hide()
				status.SetText(fmt.Sprintf("Loaded %d regions", len(report.Rows)))
			})
		}()
	}

	controls := widget.NewForm(
		widget.NewFormItem("From", fromEntry),
		widget.NewFormItem("To", toEntry),
	)

	header := container.NewVBox(
		widget.NewLabel("Totals"),
		controls,
		container.NewHBox(
			widget.NewButton("Run", run),
			progress,
			status,
		),
	)

	return container.NewBorder(
		header,
		nil,
		nil,
		nil,
		table,
	)
}

func (s *guiState) newLogsScreen() fyne.CanvasObject {
	logOutput := widget.NewTextGrid()
	s.logOutput = logOutput
	s.refreshLogs()

	return container.NewBorder(
		widget.NewLabel("Logs"),
		nil,
		nil,
		nil,
		container.NewScroll(logOutput),
	)
}

func (s *guiState) cellText(id widget.TableCellID) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch id.Row {
	case 0:
		switch id.Col {
		case 0:
			return "Region"
		case 1:
			return "Farms"
		case 2:
			return "Farmers"
		case 3:
			return "Area"
		}
	case len(s.report.Rows) + 1:
		switch id.Col {
		case 0:
			return "Total"
		case 1:
			return fmt.Sprintf("%d", s.report.TotalFarms)
		case 2:
			return fmt.Sprintf("%d", s.report.TotalFarmers)
		case 3:
			return fmt.Sprintf("%.2f", s.report.TotalArea)
		}
	default:
		row := s.report.Rows[id.Row-1]
		switch id.Col {
		case 0:
			return row.Region
		case 1:
			return fmt.Sprintf("%d", row.Farms)
		case 2:
			return fmt.Sprintf("%d", row.Farmers)
		case 3:
			return fmt.Sprintf("%.2f", row.Area)
		}
	}

	return ""
}

func (s *guiState) logf(format string, args ...any) {
	s.mu.Lock()
	s.logs = append(s.logs, fmt.Sprintf("%s  %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, args...)))
	s.mu.Unlock()
	s.refreshLogs()
}

func (s *guiState) refreshLogs() {
	if s.logOutput == nil {
		return
	}

	s.mu.Lock()
	text := strings.Join(s.logs, "\n")
	s.mu.Unlock()

	fyne.Do(func() {
		s.logOutput.SetText(text)
	})
}

func (s *guiState) selectScreen(name string) {
	if s.navigation == nil {
		return
	}

	for i, screenName := range s.screenNames {
		if screenName == name {
			s.navigation.Select(i)
			return
		}
	}
}

func (s *guiState) setDarkMode(enabled bool) {
	s.darkMode = enabled
	if enabled {
		s.app.Settings().SetTheme(theme.DarkTheme())
		s.logf("switched to dark mode")
		return
	}

	s.app.Settings().SetTheme(theme.LightTheme())
	s.logf("switched to light mode")
}

func boolLabel(value bool, trueLabel, falseLabel string) string {
	if value {
		return trueLabel
	}
	return falseLabel
}
