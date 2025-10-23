CREATE TABLE IF NOT EXISTS api_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    response_code INTEGER NOT NULL,
    response_time INTEGER NOT NULL, -- in milliseconds
    response_body TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS problems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id INTEGER NOT NULL,
    problem_type TEXT NOT NULL, -- 'error_5xx', 'error_4xx', 'slow_response', 'timeout', 'rate_limit'
    severity TEXT NOT NULL, -- 'low', 'medium', 'high', 'critical'
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (request_id) REFERENCES api_requests(id)
);

CREATE INDEX IF NOT EXISTS idx_requests_created_at ON api_requests(created_at);
CREATE INDEX IF NOT EXISTS idx_requests_method ON api_requests(method);
CREATE INDEX IF NOT EXISTS idx_requests_response_time ON api_requests(response_time);
CREATE INDEX IF NOT EXISTS idx_problems_request_id ON problems(request_id);