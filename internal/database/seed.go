package database

import "log"

func SeedData() {
	// Some sample data for testing
	seedRequests := []struct {
		method       string
		path         string
		responseCode int
		responseTime int
		responseBody string
	}{
		{"GET", "/posts", 200, 145, `[{"id":1,"title":"Sample"}]`},
		{"GET", "/users/1", 200, 89, `{"id":1,"name":"John"}`},
		{"POST", "/posts", 201, 234, `{"id":101,"title":"New Post"}`},
		{"GET", "/posts/999", 404, 67, `{}`},
		{"GET", "/comments", 500, 3421, `Internal Server Error`},
		{"DELETE", "/posts/1", 200, 112, `{}`},
		{"GET", "/users", 200, 2567, `[{"id":1},{"id":2}]`},
	}

	for _, req := range seedRequests {
		query := `INSERT INTO api_requests (method, path, response_code, response_time, response_body) 
                  VALUES (?, ?, ?, ?, ?)`
		_, err := DB.Exec(query, req.method, req.path, req.responseCode, req.responseTime, req.responseBody)
		if err != nil {
			log.Println("Error seeding data:", err)
		}
	}

	log.Println("Seed data inserted successfully")
}
