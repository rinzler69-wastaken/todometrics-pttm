package controllers

import (
	// Import the JSON package
	"fmt"
	"go-pttm/config"
	"go-pttm/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

// StatBox represents the data needed for a single box in the overview grid.
type StatBox struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// func GetOverviewStats(c *gin.Context) {
// 	view := c.DefaultQuery("view", "ongoing")
// 	sort := c.DefaultQuery("sort", "urgency")
// 	now := time.Now()

// 	var statBoxes [4]StatBox

// 	if view == "completed" {
// 		var completedTasks []models.Task
// 		config.DB.Where("completed = ?", true).Find(&completedTasks)

// 		completionCounts := map[string]int{"ahead": 0, "justintime": 0, "overdue": 0}
// 		for _, task := range completedTasks {
// 			sameDay := task.CompletedAt.Year() == task.DueDate.Year() && task.CompletedAt.YearDay() == task.DueDate.YearDay()
// 			if task.CompletedAt.After(task.DueDate) {
// 				completionCounts["overdue"]++
// 			} else if sameDay {
// 				completionCounts["justintime"]++
// 			} else {
// 				completionCounts["ahead"]++
// 			}
// 		}
// 		statBoxes[0] = StatBox{"TOTAL COMPLETED", len(completedTasks), "blue", "🏆"}
// 		statBoxes[1] = StatBox{"AHEAD OF TIME", completionCounts["ahead"], "purple", "🚀"}
// 		statBoxes[2] = StatBox{"JUST IN TIME", completionCounts["justintime"], "green", "🎯"}
// 		statBoxes[3] = StatBox{"OVERDUE", completionCounts["overdue"], "red", "🚨"}
// 	} else { // "ongoing" view
// 		var ongoingTasks []models.Task
// 		config.DB.Where("completed = ?", false).Find(&ongoingTasks)

// 		statBoxes[0] = StatBox{"TOTAL TASKS", len(ongoingTasks), "purple", "📊"}

// 		if sort == "priority" {
// 			priorityCounts := map[string]int{"high": 0, "medium": 0, "low": 0}
// 			for _, task := range ongoingTasks {
// 				switch task.Priority {
// 				case "High":
// 					priorityCounts["high"]++
// 				case "Medium":
// 					priorityCounts["medium"]++
// 				case "Low":
// 					priorityCounts["low"]++
// 				}
// 			}
// 			statBoxes[1] = StatBox{"HIGH PRIORITY", priorityCounts["high"], "red", "🔥"}
// 			statBoxes[2] = StatBox{"MEDIUM PRIORITY", priorityCounts["medium"], "yellow", "⚠️"}
// 			statBoxes[3] = StatBox{"LOW PRIORITY", priorityCounts["low"], "blue", "🧊"}
// 		} else { // "urgency" sort
// 			urgencyCounts := map[string]int{"dueToday": 0, "overdue": 0, "ongoing": 0}
// 			for _, task := range ongoingTasks {
// 				// FIX: Check if the task is overdue FIRST, using the full timestamp.
// 				if task.DueDate.Before(now) {
// 					urgencyCounts["overdue"]++
// 					// This will now only catch tasks due LATER today.
// 				} else if task.DueDate.Format("2006-01-02") == now.Format("2006-01-02") {
// 					urgencyCounts["dueToday"]++
// 				} else {
// 					urgencyCounts["ongoing"]++
// 				}
// 			}
// 			statBoxes[1] = StatBox{"OVERDUE", urgencyCounts["overdue"], "red", "⏰"}
// 			statBoxes[2] = StatBox{"DUE TODAY", urgencyCounts["dueToday"], "yellow", "📅"}
// 			statBoxes[3] = StatBox{"ONGOING", urgencyCounts["ongoing"], "blue", "🕒"}
// 		}
// 	}

// 	c.JSON(http.StatusOK, statBoxes)
// }

func GetOverviewStats(c *gin.Context) {
	view := c.DefaultQuery("view", "ongoing")
	sort := c.DefaultQuery("sort", "urgency")
	now := time.Now()

	var statBoxes [4]StatBox

	if view == "completed" {
		// ... (no changes needed in this part) ...
	} else { // "ongoing" view
		var ongoingTasks []models.Task
		config.DB.Where("completed = ?", false).Find(&ongoingTasks)

		statBoxes[0] = StatBox{"TOTAL TASKS", len(ongoingTasks), "purple", "📊"}

		if sort == "priority" {
			// ... (no changes needed in this part) ...
		} else { // "urgency" sort
			urgencyCounts := map[string]int{"dueToday": 0, "overdue": 0, "ongoing": 0}

			for _, task := range ongoingTasks {

				// --- ROBUST FIX IS HERE ---
				// The task.DueDate.Before(now) check is correct IF the timezone is right.
				// This new logic is safer.
				if task.DueDate.Before(now) {
					urgencyCounts["overdue"]++
					// Compare Year, Month, and Day to be immune to timezone differences.
				} else if task.DueDate.Year() == now.Year() && task.DueDate.Month() == now.Month() && task.DueDate.Day() == now.Day() {
					urgencyCounts["dueToday"]++
				} else {
					urgencyCounts["ongoing"]++
				}
			}
			statBoxes[1] = StatBox{"OVERDUE", urgencyCounts["overdue"], "red", "⏰"}
			statBoxes[2] = StatBox{"DUE TODAY", urgencyCounts["dueToday"], "yellow", "📅"}
			statBoxes[3] = StatBox{"ONGOING", urgencyCounts["ongoing"], "blue", "🕒"}
		}
	}

	c.JSON(http.StatusOK, statBoxes)
}

func ShowOverview(c *gin.Context) {
	view := c.DefaultQuery("view", "ongoing")
	sort := c.DefaultQuery("sort", "urgency")
	chartView := c.DefaultQuery("chart", "monthly")

	// The initial data passed is minimal, as JS will fetch the rest.
	c.HTML(http.StatusOK, "overview.html", gin.H{
		"title":  "Overview",
		"active": "overview",
		"view":   view,
		"sort":   sort,
		"chart":  chartView,
	})
}

// GetChartData is a new API endpoint that only returns chart data as JSON
func GetChartData(c *gin.Context) {
	chartView := c.DefaultQuery("chart", "monthly")
	now := time.Now()

	var completedTasks []models.Task
	if err := config.DB.Where("completed = ?", true).Order("completed_at asc").Find(&completedTasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	var chartLabels []string
	var chartScores []float64

	if chartView == "yearly" {
		yearlyScores := make(map[int]float64)
		yearlyCounts := make(map[int]int)
		for _, task := range completedTasks {
			year := task.CompletedAt.Year()
			score := 0.0
			sameDay := task.CompletedAt.Year() == task.DueDate.Year() && task.CompletedAt.YearDay() == task.DueDate.YearDay()
			if task.CompletedAt.After(task.DueDate) {
				score = -1
			} else if sameDay {
				score = 1
			} else {
				score = 2
			}
			yearlyScores[year] += score
			yearlyCounts[year]++
		}

		currentYear := now.Year()
		for i := 4; i >= 0; i-- {
			year := currentYear - i
			chartLabels = append(chartLabels, fmt.Sprintf("%d", year))
			if count, ok := yearlyCounts[year]; ok && count > 0 {
				avgScore := yearlyScores[year] / float64(count)
				chartScores = append(chartScores, avgScore)
			} else {
				chartScores = append(chartScores, 0)
			}
		}
	} else { // Default to monthly
		monthlyScores := make(map[string]float64)
		monthlyCounts := make(map[string]int)
		for _, task := range completedTasks {
			monthKey := task.CompletedAt.Format("Jan '06")
			score := 0.0
			sameDay := task.CompletedAt.Year() == task.DueDate.Year() && task.CompletedAt.YearDay() == task.DueDate.YearDay()
			if task.CompletedAt.After(task.DueDate) {
				score = -1
			} else if sameDay {
				score = 1
			} else {
				score = 2
			}
			monthlyScores[monthKey] += score
			monthlyCounts[monthKey]++
		}
		for i := 11; i >= 0; i-- {
			month := now.AddDate(0, -i, 0)
			monthKey := month.Format("Jan '06")
			chartLabels = append(chartLabels, monthKey)
			if count, ok := monthlyCounts[monthKey]; ok && count > 0 {
				avgScore := monthlyScores[monthKey] / float64(count)
				chartScores = append(chartScores, avgScore)
			} else {
				chartScores = append(chartScores, 0)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": chartLabels,
		"scores": chartScores,
	})
}

// ShowTasks handles GET /
func ShowTasks(c *gin.Context) {
	var tasks []models.Task

	if err := config.DB.
		Where("completed = ?", false).
		Order("CASE priority WHEN 'High' THEN 1 WHEN 'Medium' THEN 2 ELSE 3 END").
		Order("due_date ASC").
		Find(&tasks).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching tasks: %v", err)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":  "Your Tasks",
		"tasks":  tasks,
		"active": "tasks",
		"year":   time.Now().Year(),
	})
}

func ShowCompletedTasks(c *gin.Context) {
	var completedTasks []models.Task

	month := c.Query("month")
	year := c.Query("year")

	db := config.DB

	if month != "" && year != "" {
		// Construct start and end time for the selected month
		startDate, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%s-%s-01", year, month), time.Local)
		if err == nil {
			endDate := startDate.AddDate(0, 1, 0) // next month
			db = db.Where("completed_at >= ? AND completed_at < ?", startDate, endDate)
		}
	}

	if err := db.
		Where("completed = ?", true).
		Order("completed_at DESC").
		Find(&completedTasks).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching completed tasks")
		return
	}

	c.HTML(http.StatusOK, "completed.html", gin.H{
		"title":         "Completed Tasks",
		"tasks":         completedTasks,
		"active":        "completed",
		"selectedMonth": month,
		"selectedYear":  year,
		"filterMonth":   c.Query("month"),
	})
}

func ShowAboutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"title":  "About",
		"active": "about",
	})
}

func FilterCompletedTasks(c *gin.Context) {
	monthStr := c.Query("month") // format: YYYY-MM

	parsedTime, err := time.ParseInLocation("2006-01", monthStr, time.Local)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid month format")
		return
	}

	startDate := time.Date(parsedTime.Year(), parsedTime.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	var filteredTasks []models.Task
	if err := config.DB.
		Where("completed = ? AND completed_at >= ? AND completed_at < ?", true, startDate, endDate).
		Order("completed_at DESC").
		Find(&filteredTasks).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error filtering tasks")
		return
	}

	c.HTML(http.StatusOK, "completed.html", gin.H{
		"title":       "Completed Tasks",
		"tasks":       filteredTasks,
		"active":      "completed",
		"filterMonth": c.Query("month"),
	})
}

// CreateTask handles POST /tasks
func CreateTask(c *gin.Context) {
	var form struct {
		Name        string `form:"name" binding:"required"`
		Description string `form:"description"`
		DueDate     string `form:"due_date" binding:"required"`
		Priority    string `form:"priority" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form input: %v", err)
		return
	}

	parsedDate, err := time.ParseInLocation("2006-01-02T15:04", form.DueDate, time.Local)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid datetime format. Use YYYY-MM-DDTHH:MM")
		return
	}

	task := models.Task{
		Name:        form.Name,
		Description: form.Description,
		DueDate:     parsedDate,
		Priority:    models.PriorityLevel(form.Priority),
		Completed:   false,
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to create task: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/tasks?toast=added")

}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		c.String(http.StatusNotFound, "Task not found")
		return
	}

	action := c.PostForm("action")
	var toast string

	if action == "edit_description" {
		newDesc := c.PostForm("description")
		task.Description = newDesc
		toast = "edited"
	} else if action == "complete" {
		task.Completed = true
		task.CompletedAt = time.Now()
		toast = "completed"
	}

	if err := config.DB.Save(&task).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to update task")
		return
	}

	c.Redirect(http.StatusSeeOther, "/tasks?toast="+toast)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		c.String(http.StatusNotFound, "Task not found")
		return
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete task: %v", err)
		return
	}

	if task.Completed {
		c.Redirect(http.StatusSeeOther, "/completed?toast=deleted")
	} else {
		c.Redirect(http.StatusSeeOther, "/tasks?toast=deleted")
	}
}

func ExportTasksPDF(c *gin.Context) {
	var tasks []models.Task
	monthFilter := c.Query("month")

	var reportDate time.Time
	parsedMonth, err := time.ParseInLocation("2006-01", monthFilter, time.Local)
	if err == nil {
		reportDate = parsedMonth
	} else {
		reportDate = time.Now()
	}

	db := config.DB.Model(&models.Task{}).Where("completed = ?", true)
	if monthFilter != "" && err == nil {
		startOfMonth := time.Date(reportDate.Year(), reportDate.Month(), 1, 0, 0, 0, 0, time.Local)
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
		db = db.Where("completed_at BETWEEN ? AND ?", startOfMonth, endOfMonth)
	}

	if err := db.Order("completed_at DESC").Find(&tasks).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching completed tasks")
		return
	}

	fileName := fmt.Sprintf("GoPTTM-%s-Report.pdf", reportDate.Format("January-2006"))
	reportHeader := fmt.Sprintf("GoPTTM - Completed Tasks Report of %s", reportDate.Format("January 2006"))

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, reportHeader)
	pdf.Ln(12)

	colorRed := []int{220, 53, 69}
	colorGreen := []int{25, 135, 84}
	colorPurple := []int{102, 51, 153}

	for _, task := range tasks {
		var statusLabel string
		var statusColor []int

		sameDay := task.CompletedAt.Year() == task.DueDate.Year() && task.CompletedAt.YearDay() == task.DueDate.YearDay()

		if task.CompletedAt.After(task.DueDate) {
			statusLabel = "Overdue"
			statusColor = colorRed
		} else if sameDay {
			statusLabel = "Just In Time"
			statusColor = colorGreen
		} else {
			statusLabel = "Ahead of Time"
			statusColor = colorPurple
		}

		pdf.SetFont("Arial", "B", 14)
		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(0, 10, task.Name)
		pdf.Ln(8)

		if task.Description != "" {
			pdf.SetFont("Arial", "", 12)
			pdf.MultiCell(0, 7, "Description: "+task.Description, "", "", false)
			pdf.Ln(1)
		}

		pdf.SetFont("Arial", "", 12)
		pdf.SetTextColor(100, 100, 100)
		pdf.Cell(0, 7, "Due: "+task.DueDate.Format("02 Jan 2006, 15:04"))
		pdf.Ln(6)

		pdf.SetTextColor(statusColor[0], statusColor[1], statusColor[2])
		pdf.Cell(0, 7, "Completed: "+task.CompletedAt.Format("02 Jan 2006, 15:04")+" ("+statusLabel+")")
		pdf.Ln(8)

		pdf.SetTextColor(0, 0, 0)
		pdf.Cell(0, 6, "Priority: "+string(task.Priority))
		pdf.Ln(10)

		pdf.SetDrawColor(220, 220, 220)
		pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
		pdf.Ln(10)
	}

	if err := pdf.Output(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate PDF: %v", err)
	}
}
