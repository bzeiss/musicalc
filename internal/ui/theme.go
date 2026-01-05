package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct {
	isMobile bool
}

func (m customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m customTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInlineIcon {
		if m.isMobile {
			return 16 // Smaller tab icons on mobile
		}
		return 32
	}

	// Adjust sizes for mobile
	if m.isMobile {
		switch name {
		case theme.SizeNamePadding:
			return 8 // Larger padding for better touch targets in menus
		case theme.SizeNameInnerPadding:
			return 6 // Larger inner padding for menus
		case theme.SizeNameText:
			return 16 // Larger text for menus
		case theme.SizeNameScrollBarSmall:
			return 2
		}
	}

	return theme.DefaultTheme().Size(name)
}

func NewCustomTheme() fyne.Theme {
	return customTheme{isMobile: fyne.CurrentDevice().IsMobile()}
}
