package main

import (
	_ "embed"
	"musicalc/internal/ui"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//go:embed VERSION
var version string

// CategoryInfo represents a category with tab indices
type CategoryInfo struct {
	Name       string
	TabIndices []int
}

// createCompactHeader creates a compact header with: hamburger + "category / tab heading"
func createCompactHeader(window fyne.Window, categories []CategoryInfo, switchCategory func(int), headerLabel *widget.Label, toolbar *widget.Toolbar) *fyne.Container {
	menuButton := widget.NewButtonWithIcon("", theme.MenuIcon(), func() {
		showCategoryMenu(window, categories, switchCategory)
	})
	menuButton.Importance = widget.LowImportance
	return container.NewBorder(nil, nil, menuButton, toolbar, container.NewCenter(headerLabel))
}

func main() {
	myApp := app.NewWithID("com.musicalc")
	myApp.Settings().SetTheme(ui.NewCustomTheme())

	ver := strings.TrimSpace(version)
	windowTitle := "MusiCalc"
	if ver != "" {
		windowTitle = "MusiCalc v" + ver
	}
	window := myApp.NewWindow(windowTitle)
	window.Resize(fyne.NewSize(450, 650))
	window.SetFixedSize(true)

	if icon, err := fyne.LoadResourceFromPath("icons/appicon.png"); err == nil {
		window.SetIcon(icon)
	}

	tabHeadings := map[string]string{
		"timecode":     "Timecode Calculator",
		"tempo":        "Tempo to Delay",
		"tempochange":  "Tempo Change",
		"note2freq":    "Note to Frequency",
		"freq2note":    "Frequency to Note",
		"samplelength": "Sample Length",
		"alignment":    "Alignment Delay",
	}

	var switchCategory func(int)

	// Determine tab text based on device type
	isMobile := fyne.CurrentDevice().IsMobile()
	var timecodeText, tempoText, tempoChangeText, note2freqText, freq2noteText, sampleLengthText, alignmentText string
	if !isMobile {
		timecodeText = "Timecode"
		tempoText = "Delay"
		tempoChangeText = "Tempo Chg"
		note2freqText = "Note→Freq"
		freq2noteText = "Freq→Note"
		sampleLengthText = "Sample Len"
		alignmentText = "Align Dly"
	}

	// Create all tab items
	timecodeTab := container.NewTabItem(timecodeText, ui.NewTimecodeTab())
	timecodeTab.Icon = ui.ResourceTimecodeSvg

	tempoTab := container.NewTabItem(tempoText, ui.NewTempoTab())
	tempoTab.Icon = ui.ResourceDelaySvg

	tempoChangeTab := container.NewTabItem(tempoChangeText, ui.NewTempoChangeTab())
	tempoChangeTab.Icon = ui.ResourceTempochangeSvg

	note2freqTab := container.NewTabItem(note2freqText, ui.NewDiapasonTab())
	note2freqTab.Icon = ui.ResourceNote2freqSvg

	freq2noteTab := container.NewTabItem(freq2noteText, ui.NewFrequencyToNoteTab())
	freq2noteTab.Icon = ui.ResourceFreq2noteSvg

	sampleLengthTab := container.NewTabItem(sampleLengthText, ui.NewSampleLengthTab())
	sampleLengthTab.Icon = ui.ResourceSamplelengthSvg

	alignmentContent, alignmentExportCSV := ui.NewAlignmentDelayTabWithExport()
	alignmentTab := container.NewTabItem(alignmentText, alignmentContent)
	alignmentTab.Icon = ui.ResourceAlignmentdelaySvg

	// Create single AppTabs with ALL tabs (maintains left alignment)
	allTabs := []*container.TabItem{
		timecodeTab, tempoTab, tempoChangeTab,
		note2freqTab, freq2noteTab,
		sampleLengthTab,
		alignmentTab,
	}
	tabs := container.NewAppTabs(allTabs...)
	tabs.SetTabLocation(container.TabLocationBottom)

	// Define categories with their tab indices
	categories := []CategoryInfo{
		{Name: "Time & Tempo", TabIndices: []int{0, 1, 2}},
		{Name: "Frequency & Pitch", TabIndices: []int{3, 4}},
		{Name: "Analysis", TabIndices: []int{5}},
		{Name: "Multi-Mic", TabIndices: []int{6}},
	}

	// Tab heading keys for each global tab index
	tabHeadingKeys := []string{
		"timecode", "tempo", "tempochange",
		"note2freq", "freq2note",
		"samplelength",
		"alignment",
	}

	// Always start with Time & Tempo category (Timecode Calculator)
	currentCategoryIndex := 0

	headerLabel := widget.NewLabel("")
	headerToolbar := widget.NewToolbar()
	headerToolbar.Hide()

	// Per-tab toolbar items (empty for all calculators for now)
	perTabToolbarItems := make([][]widget.ToolbarItem, len(allTabs))
	perTabToolbarItems[6] = []widget.ToolbarItem{
		widget.NewToolbarAction(theme.DownloadIcon(), func() {
			d := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				if uc == nil {
					return
				}
				defer uc.Close()

				data, err := alignmentExportCSV()
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				if _, err := uc.Write([]byte(data)); err != nil {
					dialog.ShowError(err, window)
					return
				}
			}, window)
			d.SetFileName("alignment_delay.csv")
			d.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			d.Show()
		}),
	}
	setHeaderToolbar := func(items []widget.ToolbarItem) {
		if len(items) == 0 {
			headerToolbar.Hide()
			headerToolbar.Items = nil
			headerToolbar.Refresh()
			return
		}
		headerToolbar.Items = items
		headerToolbar.Show()
		headerToolbar.Refresh()
	}

	// Helper to show only tabs for a specific category
	showCategoryTabs := func(categoryIdx int) {
		if categoryIdx < 0 || categoryIdx >= len(categories) {
			return
		}

		cat := categories[categoryIdx]

		// Build list of visible tabs for this category
		var visibleTabs []*container.TabItem
		for _, tabIdx := range cat.TabIndices {
			if tabIdx >= 0 && tabIdx < len(allTabs) {
				visibleTabs = append(visibleTabs, allTabs[tabIdx])
			}
		}

		// Replace tabs.Items with visible tabs
		tabs.Items = visibleTabs

		// Select first tab in category
		if len(visibleTabs) > 0 {
			tabs.Select(visibleTabs[0])
		}

		tabs.Refresh()
	}

	// Helper to update header text based on selected tab
	updateHeader := func() {
		selectedTab := tabs.Selected()
		if selectedTab == nil {
			return
		}

		// Find which tab is selected
		for i, tab := range allTabs {
			if tab == selectedTab {
				// Find category for this tab
				for _, cat := range categories {
					for _, tabIdx := range cat.TabIndices {
						if tabIdx == i {
							headerText := cat.Name + " / " + tabHeadings[tabHeadingKeys[i]]
							headerLabel.SetText(headerText)
							setHeaderToolbar(perTabToolbarItems[i])
							return
						}
					}
				}
			}
		}
	}

	// Set initial category tabs visibility
	showCategoryTabs(currentCategoryIndex)
	updateHeader()

	// Update header when tab changes
	tabs.OnSelected = func(item *container.TabItem) {
		updateHeader()
	}

	// Function to switch categories
	switchCategory = func(index int) {
		if index >= 0 && index < len(categories) {
			if index == currentCategoryIndex {
				return
			}

			currentCategoryIndex = index
			showCategoryTabs(index)
			updateHeader()

			myApp.Preferences().SetInt("lastCategory", index)
		}
	}

	// Create initial header and content
	header := createCompactHeader(window, categories, switchCategory, headerLabel, headerToolbar)

	content := container.NewBorder(
		container.NewVBox(header, widget.NewSeparator()),
		nil, nil, nil,
		tabs,
	)
	window.SetContent(content)
	window.ShowAndRun()
}

// showCategoryMenu displays a popup menu for category selection
func showCategoryMenu(window fyne.Window, categories []CategoryInfo, onSelect func(int)) {
	var menuItems []*fyne.MenuItem

	for i, category := range categories {
		index := i
		menuItems = append(menuItems, fyne.NewMenuItem(category.Name, func() {
			onSelect(index)
		}))
	}

	menu := fyne.NewMenu("Categories", menuItems...)
	popupMenu := widget.NewPopUpMenu(menu, window.Canvas())
	popupMenu.ShowAtPosition(fyne.NewPos(10, 50))
}
