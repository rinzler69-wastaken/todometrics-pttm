package main

import (
	"html/template"
	"log"
	"time"

	"github.com/rinzler69-wastaken/todometrics-pttm/config"
	"github.com/rinzler69-wastaken/todometrics-pttm/controllers"

	"github.com/gin-gonic/gin"
)

func seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

func main() {
	// Initialize database
	config.InitDatabase()

	// Wipe and reset the tasks table (only during dev/testing)
	// config.DB.Migrator().DropTable(&models.Task{})
	// config.DB.AutoMigrate(&models.Task{})

	// Setup Gin router
	router := gin.Default()

	// Serve static assets
	router.Static("/static", "./static")

	router.SetFuncMap(template.FuncMap{
		"time": func(value string, layout string) time.Time {
			t, _ := time.Parse(layout, value)
			return t
		},
	})

	// Load HTML templates
	router.LoadHTMLGlob("templates/*.html")

	// Routes
	router.GET("/", controllers.ShowOverview) // 👈 new root page
	router.GET("/overview", controllers.ShowOverview)
	router.GET("/tasks", controllers.ShowTasks)              // Read
	router.POST("/tasks", controllers.CreateTask)            // Create
	router.POST("/tasks/:id/update", controllers.UpdateTask) // Update (only description + completed)
	router.GET("/completed", controllers.ShowCompletedTasks)
	router.POST("/tasks/:id/delete", controllers.DeleteTask) // Delete
	router.GET("/about", controllers.ShowAboutPage)
	router.GET("/export/pdf", controllers.ExportTasksPDF)
	router.GET("/completed/filter", controllers.FilterCompletedTasks)
	router.GET("/api/overview/chart-data", controllers.GetChartData)
	router.GET("/api/overview/stats", controllers.GetOverviewStats)

	//router.GET("/report/pdf", controllers.ExportPDF)         // PDF Export
	//router.GET("/report/excel", controllers.ExportExcel)     // Excel Export

	// Start server
	port := "0.0.0.0:8080"
	log.Printf("Listening on http://localhost%s", port)
	err := router.Run(port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
