package tui

import "github.com/charmbracelet/lipgloss"

// Colors based on Codeforces rank system
var (
	ColorNewbie    = lipgloss.Color("#808080") // Gray
	ColorPupil     = lipgloss.Color("#008000") // Green
	ColorSpecialist = lipgloss.Color("#03A89E") // Cyan
	ColorExpert    = lipgloss.Color("#0000FF") // Blue
	ColorCM        = lipgloss.Color("#AA00AA") // Violet
	ColorMaster    = lipgloss.Color("#FF8C00") // Orange
	ColorGM        = lipgloss.Color("#FF0000") // Red

	ColorPrimary   = lipgloss.Color("#7C3AED") // Purple
	ColorSecondary = lipgloss.Color("#06B6D4") // Cyan
	ColorSuccess   = lipgloss.Color("#10B981") // Green
	ColorWarning   = lipgloss.Color("#F59E0B") // Amber
	ColorError     = lipgloss.Color("#EF4444") // Red
	ColorMuted     = lipgloss.Color("#6B7280") // Gray
	ColorBorder    = lipgloss.Color("#374151") // Dark gray
	ColorBg        = lipgloss.Color("#1F2937") // Dark background
	ColorBgLight   = lipgloss.Color("#374151") // Lighter background
)

// Base styles
var (
	// Text styles
	TextBold = lipgloss.NewStyle().Bold(true)
	TextDim  = lipgloss.NewStyle().Foreground(ColorMuted)
	TextMuted = lipgloss.NewStyle().Foreground(ColorMuted)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			MarginBottom(1)

	// Header style for page headers
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 2).
			MarginBottom(1)

	// Box/Card styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	FocusedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess)
	WarningStyle = lipgloss.NewStyle().Foreground(ColorWarning)
	ErrorStyle   = lipgloss.NewStyle().Foreground(ColorError)

	// Help bar style
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)

	// Active/Selected item
	ActiveItemStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	// Tab styles
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 2)

	// Badge styles
	BadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 1)

	SolvedBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorSuccess).
			Padding(0, 1)

	// Tag style
	TagStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Background(ColorBgLight).
			Padding(0, 1)

	// Progress bar colors
	ProgressFilled = lipgloss.NewStyle().Foreground(ColorSuccess)
	ProgressEmpty  = lipgloss.NewStyle().Foreground(ColorBorder)

	// Spinner style
	SpinnerStyle = lipgloss.NewStyle().Foreground(ColorPrimary)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(ColorBorder)

	TableRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	TableSelectedRowStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Background(ColorBgLight).
				Bold(true)
)

// RatingStyle returns a style based on the Codeforces rating
func RatingStyle(rating int) lipgloss.Style {
	var color lipgloss.Color
	switch {
	case rating == 0:
		color = ColorNewbie
	case rating < 1200:
		color = ColorNewbie
	case rating < 1400:
		color = ColorPupil
	case rating < 1600:
		color = ColorSpecialist
	case rating < 1900:
		color = ColorExpert
	case rating < 2100:
		color = ColorCM
	case rating < 2400:
		color = ColorMaster
	default:
		color = ColorGM
	}
	return lipgloss.NewStyle().Foreground(color).Bold(true)
}

// RankStyle returns a style based on the Codeforces rank string
func RankStyle(rank string) lipgloss.Style {
	var color lipgloss.Color
	switch rank {
	case "newbie":
		color = ColorNewbie
	case "pupil":
		color = ColorPupil
	case "specialist":
		color = ColorSpecialist
	case "expert":
		color = ColorExpert
	case "candidate master":
		color = ColorCM
	case "master", "international master":
		color = ColorMaster
	case "grandmaster", "international grandmaster", "legendary grandmaster":
		color = ColorGM
	default:
		color = ColorNewbie
	}
	return lipgloss.NewStyle().Foreground(color).Bold(true)
}

// Logo returns the styled app logo
func Logo() string {
	logo := `
 ____  ____    _    ____
|  _ \/ ___|  / \  |  _ \ _ __ ___ _ __
| | | \___ \ / _ \ | |_) | '__/ _ \ '_ \
| |_| |___) / ___ \|  __/| | |  __/ |_) |
|____/|____/_/   \_\_|   |_|  \___| .__/
                                  |_|
`
	return lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render(logo)
}

// StatusBar returns a styled status bar
func StatusBar(content string) string {
	return lipgloss.NewStyle().
		Foreground(ColorMuted).
		Background(ColorBgLight).
		Padding(0, 1).
		Width(80).
		Render(content)
}

// ProgressBar renders a simple ASCII progress bar
func ProgressBar(current, total, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))
	empty := width - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += ProgressFilled.Render("█")
	}
	for i := 0; i < empty; i++ {
		bar += ProgressEmpty.Render("░")
	}

	return bar
}
