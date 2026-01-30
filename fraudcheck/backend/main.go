package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"fraudcheck-backend/figchain"

	"github.com/figchain/go-client/pkg/client"
	"github.com/figchain/go-client/pkg/evaluation"
	"github.com/gin-gonic/gin"
)

// PollManager handles long polling connections
type PollManager struct {
	mu        sync.RWMutex
	listeners []chan bool
}

var pollManager = &PollManager{
	listeners: make([]chan bool, 0),
}

var (
	currentFailThreshold = 70 // default: 70% risk threshold
	failThresholdMutex   sync.RWMutex
)

// AddListener registers a new long polling listener
func (pm *PollManager) AddListener() chan bool {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	ch := make(chan bool, 1)
	pm.listeners = append(pm.listeners, ch)
	return ch
}

// RemoveListener unregisters a long polling listener
func (pm *PollManager) RemoveListener(ch chan bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for i, listener := range pm.listeners {
		if listener == ch {
			pm.listeners = append(pm.listeners[:i], pm.listeners[i+1:]...)
			close(ch)
			break
		}
	}
}

// NotifyAll sends a refresh signal to all listeners
func (pm *PollManager) NotifyAll() {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, ch := range pm.listeners {
		select {
		case ch <- true:
		default:
		}
	}
}

// SetFailThreshold updates the fail threshold
func SetFailThreshold(threshold int) {
	failThresholdMutex.Lock()
	defer failThresholdMutex.Unlock()
	currentFailThreshold = threshold
}

// GetFailThreshold returns the current fail threshold
func GetFailThreshold() int {
	failThresholdMutex.RLock()
	defer failThresholdMutex.RUnlock()
	return currentFailThreshold
}

// RecalculateStatuses recalculates application statuses based on current threshold
func RecalculateStatuses() error {
	threshold := GetFailThreshold()
	log.Printf("Recalculating statuses with fail threshold: %d%%", threshold)

	applications, err := GetAllApplications()
	if err != nil {
		return err
	}

	for _, app := range applications {
		riskPercentage := app.RiskFactor * 100
		newStatus := determineStatusByThreshold(riskPercentage, float64(threshold))

		if newStatus != app.Status {
			err := UpdateApplicationStatus(app.ID, newStatus)
			if err != nil {
				log.Printf("Failed to update application %d status: %v", app.ID, err)
			}
		}
	}

	return nil
}

// determineStatusByThreshold assigns status based on risk percentage and threshold
func determineStatusByThreshold(riskPercentage, threshold float64) string {
	if riskPercentage > threshold {
		return "rejected"
	} else if riskPercentage > threshold*0.7 { // 70% of threshold = review required
		return "review_required"
	}
	return "approved"
}

func main() {
	// Initialize database
	err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Seed database with fake data
	err = SeedDatabase()
	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize FigChain client
	configPath := "client-config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Warning: Config file %s not found. Running without FigChain integration.", configPath)
	} else {
		go initFigChain(configPath)
	}

	// Create Gin router
	router := gin.Default()

	// Enable CORS for frontend
	router.Use(corsMiddleware())

	// API routes
	router.GET("/api/applications", getApplications)
	router.GET("/api/health", healthCheck)
	router.GET("/api/poll", longPoll)
	router.POST("/api/trigger-refresh", triggerRefresh)

	// Start server
	port := ":8080"
	log.Printf("Starting Fraud Check API server on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initFigChain initializes the FigChain client and registers listeners
func initFigChain(configPath string) {
	c, err := client.NewClientFromConfig(configPath)
	if err != nil {
		log.Printf("Failed to create FigChain client: %v", err)
		return
	}
	defer c.Close()

	log.Println("FigChain client initialized successfully")

	figKey := "parameters"
	log.Printf("Listening for FigChain updates on key: %s", figKey)

	// Register listener for configuration updates
	c.RegisterListener(figKey, &figchain.FraudCheckParameters{}, func(r client.AvroRecord) {
		params, ok := r.(*figchain.FraudCheckParameters)
		if !ok {
			log.Printf("Received record is not of type *figchain.FraudCheckParameters")
			return
		}

		log.Printf(">>> FigChain UPDATE RECEIVED for %s <<<", figKey)
		log.Printf("New FailThreshold: %d%%", params.FailThreshold)

		// Update the fail threshold
		SetFailThreshold(params.FailThreshold)

		// Recalculate application statuses based on new threshold
		if err := RecalculateStatuses(); err != nil {
			log.Printf("Error recalculating statuses: %v", err)
			return
		}

		// Notify all connected clients to refresh
		log.Println("Notifying all clients to refresh due to config change")
		pollManager.NotifyAll()
	})

	// Fetch initial value
	evalContext := evaluation.NewEvaluationContext(nil)
	var initialVal figchain.FraudCheckParameters
	if err := c.GetFig(figKey, &initialVal, evalContext); err != nil {
		log.Printf("Initial GetFig failed: %v (using default threshold)", err)
	} else {
		log.Printf("Initial FigChain value fetched. FailThreshold=%d%%", initialVal.FailThreshold)
		SetFailThreshold(initialVal.FailThreshold)

		// Recalculate statuses with initial threshold
		if err := RecalculateStatuses(); err != nil {
			log.Printf("Error recalculating statuses with initial threshold: %v", err)
		}
	}

	// Keep FigChain client running
	select {}
}

// getApplications handles GET /api/applications
func getApplications(c *gin.Context) {
	applications, err := GetAllApplications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch applications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": applications,
		"count":        len(applications),
	})
}

// healthCheck handles GET /api/health
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// longPoll handles GET /api/poll - long polling endpoint
func longPoll(c *gin.Context) {
	// Register this connection as a listener
	ch := pollManager.AddListener()
	defer pollManager.RemoveListener(ch)

	// Wait for either a notification or timeout (30 seconds)
	select {
	case <-ch:
		// Received a refresh signal
		c.JSON(http.StatusOK, gin.H{
			"action": "refresh",
		})
	case <-time.After(30 * time.Second):
		// Timeout - no updates
		c.JSON(http.StatusOK, gin.H{
			"action": "none",
		})
	case <-c.Request.Context().Done():
		// Client disconnected
		return
	}
}

// triggerRefresh handles POST /api/trigger-refresh - manually trigger a refresh for testing
func triggerRefresh(c *gin.Context) {
	log.Println("Manually triggering refresh for all connected clients")
	pollManager.NotifyAll()
	c.JSON(http.StatusOK, gin.H{
		"message": "Refresh triggered",
	})
}

// corsMiddleware enables CORS for all routes
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
