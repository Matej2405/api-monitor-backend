package handlers

import (
	"api-monitor-backend/internal/database"
	"api-monitor-backend/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const JSONPlaceholderURL = "https://jsonplaceholder.typicode.com"

// GetRequests godoc
// @Summary Get all API requests
// @Description Get list of API requests with optional filtering and sorting
// @Tags requests
// @Accept json
// @Produce json
// @Param method query string false "Filter by HTTP method (GET, POST, PUT, DELETE)"
// @Param response_code query int false "Filter by exact response code"
// @Param min_response_code query int false "Minimum response code"
// @Param max_response_code query int false "Maximum response code"
// @Param min_response_time query int false "Minimum response time in ms"
// @Param max_response_time query int false "Maximum response time in ms"
// @Param start_date query string false "Filter by start date (RFC3339 format)"
// @Param end_date query string false "Filter by end date (RFC3339 format)"
// @Param search query string false "Search in path"
// @Param sort_by query string false "Sort by field (created_at, response_time)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {array} models.APIRequest
// @Router /api/requests [get]
func GetRequests(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, method, path, response_code, response_time, response_body, created_at FROM api_requests WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	// Filter by method
	if method := r.URL.Query().Get("method"); method != "" {
		query += fmt.Sprintf(" AND method = $%d", argCount)
		args = append(args, method)
		argCount++
	}
	// Filter by response code
	if responseCode := r.URL.Query().Get("response_code"); responseCode != "" {
		query += fmt.Sprintf(" AND response_code = $%d", argCount)
		args = append(args, responseCode)
		argCount++
	}

	// Filter by response code range (for 2xx, 4xx, 5xx filtering)
	if minCode := r.URL.Query().Get("min_response_code"); minCode != "" {
		query += fmt.Sprintf(" AND response_code >= $%d", argCount)
		args = append(args, minCode)
		argCount++
	}
	if maxCode := r.URL.Query().Get("max_response_code"); maxCode != "" {
		query += fmt.Sprintf(" AND response_code <= $%d", argCount)
		args = append(args, maxCode)
		argCount++
	}

	// Filter by response time range
	if minTime := r.URL.Query().Get("min_response_time"); minTime != "" {
		query += fmt.Sprintf(" AND response_time >= $%d", argCount)
		args = append(args, minTime)
		argCount++
	}
	if maxTime := r.URL.Query().Get("max_response_time"); maxTime != "" {
		query += fmt.Sprintf(" AND response_time <= $%d", argCount)
		args = append(args, maxTime)
		argCount++
	}

	// Filter by date range
	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, startDate)
		argCount++
	}
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, endDate)
		argCount++
	}

	// Search in path
	if search := r.URL.Query().Get("search"); search != "" {
		query += fmt.Sprintf(" AND path LIKE $%d", argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Sorting
	sortBy := r.URL.Query().Get("sort_by")
	order := r.URL.Query().Get("order")

	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	// Validate sort fields
	validSortFields := map[string]bool{
		"created_at":    true,
		"response_time": true,
	}

	if validSortFields[sortBy] {
		query += fmt.Sprintf(" ORDER BY %s %s", sortBy, strings.ToUpper(order))
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	requests := []models.APIRequest{}
	for rows.Next() {
		var req models.APIRequest
		err := rows.Scan(&req.ID, &req.Method, &req.Path, &req.ResponseCode,
			&req.ResponseTime, &req.ResponseBody, &req.CreatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		requests = append(requests, req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetProblems godoc
// @Summary Get all API problems
// @Description Get list of detected API problems with optional filtering
// @Tags problems
// @Accept json
// @Produce json
// @Param problem_type query string false "Filter by problem type"
// @Param severity query string false "Filter by severity"
// @Param sort_by query string false "Sort by field"
// @Param order query string false "Sort order"
// @Success 200 {array} models.Problem
// @Router /api/problems [get]
func GetProblems(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, request_id, problem_type, severity, description, created_at FROM problems WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	// Filter by problem type
	if problemType := r.URL.Query().Get("problem_type"); problemType != "" {
		query += fmt.Sprintf(" AND problem_type = $%d", argCount)
		args = append(args, problemType)
		argCount++
	}

	// Filter by severity
	if severity := r.URL.Query().Get("severity"); severity != "" {
		query += fmt.Sprintf(" AND severity = $%d", argCount)
		args = append(args, severity)
		argCount++
	}

	// Sorting
	sortBy := r.URL.Query().Get("sort_by")
	order := r.URL.Query().Get("order")

	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	validSortFields := map[string]bool{
		"created_at": true,
		"severity":   true,
	}

	if validSortFields[sortBy] {
		query += fmt.Sprintf(" ORDER BY %s %s", sortBy, strings.ToUpper(order))
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	problems := []models.Problem{}
	for rows.Next() {
		var prob models.Problem
		err := rows.Scan(&prob.ID, &prob.RequestID, &prob.ProblemType,
			&prob.Severity, &prob.Description, &prob.CreatedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		problems = append(problems, prob)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problems)
}

// ProxyRequest godoc
// @Summary Proxy request to JSONPlaceholder
// @Description Proxy API request through backend and log it
// @Tags proxy
// @Accept json
// @Produce json
// @Param endpoint path string true "API endpoint to proxy"
// @Success 200 {object} map[string]interface{}
// @Router /api/proxy/{endpoint} [get]
// @Router /api/proxy/{endpoint} [post]
// @Router /api/proxy/{endpoint} [put]
// @Router /api/proxy/{endpoint} [delete]
func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	endpoint := vars["endpoint"]

	// Build the full URL
	targetURL := fmt.Sprintf("%s/%s", JSONPlaceholderURL, endpoint)

	// Record start time
	startTime := time.Now()

	// Create the proxy request
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Make the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(proxyReq)

	responseTime := int(time.Since(startTime).Milliseconds())
	responseCode := 0
	responseBody := ""

	if err != nil {
		// Handle timeout or connection errors
		responseCode = 0 // Indicates failure
		responseBody = err.Error()

		// Log the request
		logRequest(r.Method, "/"+endpoint, responseCode, responseTime, responseBody)

		// Create problem for timeout
		createProblem(r.Method, "/"+endpoint, responseCode, responseTime)

		http.Error(w, "Request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	responseCode = resp.StatusCode

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseBody = string(bodyBytes)

	// Log the request to database
	logRequest(r.Method, "/"+endpoint, responseCode, responseTime, responseBody)

	// Check for problems and create entries
	createProblem(r.Method, "/"+endpoint, responseCode, responseTime)

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(bodyBytes)
}

// Helper function to log request to database
func logRequest(method, path string, responseCode, responseTime int, responseBody string) {
	query := `INSERT INTO api_requests (method, path, response_code, response_time, response_body) 
              VALUES (?, ?, ?, ?, ?)`

	_, err := database.DB.Exec(query, method, path, responseCode, responseTime, responseBody)
	if err != nil {
		log.Println("Error logging request:", err)
	}
}

// Helper function to create problem if detected
func createProblem(method, path string, responseCode, responseTime int) {
	var problemType, severity, description string
	shouldCreateProblem := false

	// Check for 5xx errors
	if responseCode >= 500 && responseCode < 600 {
		problemType = "error_5xx"
		severity = "critical"
		description = fmt.Sprintf("Server error %d on %s %s", responseCode, method, path)
		shouldCreateProblem = true
	} else if responseCode >= 400 && responseCode < 500 {
		// Check for 4xx errors
		problemType = "error_4xx"
		severity = "medium"
		description = fmt.Sprintf("Client error %d on %s %s", responseCode, method, path)
		shouldCreateProblem = true
	} else if responseCode == 429 {
		// Rate limit
		problemType = "rate_limit"
		severity = "high"
		description = fmt.Sprintf("Rate limit exceeded on %s %s", method, path)
		shouldCreateProblem = true
	} else if responseTime > 2000 {
		// Slow response
		problemType = "slow_response"
		severity = "medium"
		description = fmt.Sprintf("Slow response (%dms) on %s %s", responseTime, method, path)
		shouldCreateProblem = true
	} else if responseCode == 0 {
		// Timeout or connection failure
		problemType = "timeout"
		severity = "critical"
		description = fmt.Sprintf("Request timeout on %s %s", method, path)
		shouldCreateProblem = true
	}

	if shouldCreateProblem {
		// Get the last inserted request ID
		var requestID int
		err := database.DB.QueryRow("SELECT id FROM api_requests ORDER BY id DESC LIMIT 1").Scan(&requestID)
		if err != nil {
			log.Println("Error getting last request ID:", err)
			return
		}

		query := `INSERT INTO problems (request_id, problem_type, severity, description) 
                  VALUES (?, ?, ?, ?)`

		_, err = database.DB.Exec(query, requestID, problemType, severity, description)
		if err != nil {
			log.Println("Error creating problem:", err)
		}
	}
}
