package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var defaultStyle = lipgloss.NewStyle()
var firstCellStyle = lipgloss.NewStyle().Padding(0, 0).Width(14).Bold(true)
var headerStyle = lipgloss.NewStyle().Bold(true)

func PrintInitialTUI(userName, serverUrl, localUrl string) {
	t := table.New().Border(lipgloss.HiddenBorder()).StyleFunc(func(row, col int) lipgloss.Style {
		switch col {
		case 0:
			return firstCellStyle
		default:
			return defaultStyle
		}
	})
	t.Headers()
	t.Row("Account", userName)
	t.Row("Receiving", serverUrl)
	t.Row("Forwarding", localUrl)

	fmt.Println(headerStyle.Render("whrelay") + " by " + headerStyle.Render("@yagnikpt"))
	fmt.Println(t.Render())
	fmt.Println("HTTP Requests")
	fmt.Println("-------------")
	fmt.Println()
}
