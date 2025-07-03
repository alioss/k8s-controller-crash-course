package formatter

import (
	"os"
	"sort"

	"atlasctl/pkg/models"

	"github.com/olekukonko/tablewriter"
)

// PrintAtlasAppsTable prints Atlas apps in a formatted table
func PrintAtlasAppsTable(apps models.AtlasAppList) {
	if len(apps) == 0 {
		println("No Atlas applications found.")
		return
	}

	// Sort apps by namespace for consistent output
	sort.Slice(apps, func(i, j int) bool {
		if apps[i].Namespace != apps[j].Namespace {
			return apps[i].Namespace < apps[j].Namespace
		}
		return apps[i].Name < apps[j].Name
	})

	table := tablewriter.NewWriter(os.Stdout)
	
	// Set table style
	table.SetHeader([]string{
		"Namespace",
		"App", 
		"Version",
		"Migration ID",
		"Status",
		"Replicas",
		"Last Update",
		"Age",
	})
	
	// Configure table appearance
	table.SetBorders(tablewriter.Border{
		Left:   true,
		Top:    true,
		Right:  true,
		Bottom: true,
	})
	table.SetCenterSeparator("┼")
	table.SetColumnSeparator("│")
	table.SetRowSeparator("─")
	table.SetHeaderLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)

	// Add color coding for status
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)

	// Add rows
	for _, app := range apps {
		row := []string{
			app.Namespace,
			app.Name,
			app.Version,
			app.MigrationID,
			app.Status,
			app.Replicas,
			app.LastUpdate.Format("2006-01-02 15:04:05"),
			app.Age,
		}

		// Color code status
		colors := []tablewriter.Colors{
			{}, // Namespace
			{}, // App
			{}, // Version
			{}, // Migration ID
			getStatusColor(app.Status), // Status
			{}, // Replicas
			{}, // Last Update
			{}, // Age
		}

		table.Rich(row, colors)
	}

	table.Render()
}

// getStatusColor returns appropriate color for status
func getStatusColor(status string) tablewriter.Colors {
	switch status {
	case "Running":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor}
	case "Pending":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgYellowColor}
	case "Failed":
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgRedColor}
	default:
		return tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}
	}
}
