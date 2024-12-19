# rate-limiter-go

This project implements a configurable rate limiter in Go, adhering to the Clean Architecture principles. It restricts the number of HTTP requests based on IP addresses or access tokens, leveraging Redis for data persistence.

---

## Features
- **Rate Limiting by IP**: Configurable limits for requests based on the client's IP address.
- **Rate Limiting by Token**: Custom limits for requests based on access tokens.
- **Persistence**: Uses Redis to store rate limiting data.
- **HTTP Responses**: Handles request limits with appropriate HTTP 429 responses and `Retry-After` headers.

---

## How to Execute

### 1. Clone the Repository
```bash
git clone https://github.com/willychavez/rate-limiter-go.git
cd rate-limiter-go
```

### 2. Configure the Environment Variables
Navigate to the `cmd/server` directory and copy the sample `.env.dist` to `.env`:
```bash
cd cmd/server
cp .env.dist .env
```
Edit the `.env` file as needed.

### 3. Start the Application
Return to the root directory and run:
```bash
docker compose up --build
```

---

## Configuration

### Limits by IP
Set these variables in the `.env` file to configure limits for IP addresses:
```env
REQUEST_LIMIT=10
BLOCK_TIME=60
```
- **`REQUEST_LIMIT`**: Maximum number of requests allowed per IP.
- **`BLOCK_TIME`**: Time (in seconds) the IP will be blocked after exceeding the limit.

### Limits by Token
Use the Redis CLI to configure limits for specific tokens.

#### Access Redis CLI:
```bash
docker exec -it rate_limiter_redis redis-cli
```

#### Set Limits for a Token:
For token `abc123`:
```bash
SET token:limit:abc123 5
SET token:block_time:abc123 120
```
For token `xyz789`:
```bash
SET token:limit:xyz789 10
SET token:block_time:xyz789 60
```

#### Validate Configuration:
```bash
GET token:limit:abc123
GET token:block_time:abc123
```

---

## Testing

### Test by IP
Run the following script to test the limits by IP:
```bash
for i in {1..20}; do
  response=$(curl -s -i -w "\n%{http_code}" http://localhost:8080/)
  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | sed -n '/^\r$/,$p' | tail -n +2 | head -n -1)
  retry_after=$(echo "$response" | grep -i "Retry-After" | awk '{print $2}' | tr -d '\r')
  if [ -n "$retry_after" ]; then
    echo "$http_code | $body | Retry-After: $retry_after"
  else
    echo "$http_code | $body"
  fi
done
```

### Test by Token
Replace `abc123` or `xyz789` with the desired token in the following script:
```bash
for i in {1..20}; do
  response=$(curl -s -i -w "\n%{http_code}" http://localhost:8080/ -H "API_KEY: abc123")
  http_code=$(echo "$response" | tail -n1)
  body=$(echo "$response" | sed -n '/^\r$/,$p' | tail -n +2 | head -n -1)
  retry_after=$(echo "$response" | grep -i "Retry-After" | awk '{print $2}' | tr -d '\r')
  if [ -n "$retry_after" ]; then
    echo "$http_code | $body | Retry-After: $retry_after"
  else
    echo "$http_code | $body"
  fi
done
```

### Example Requests
Additional example requests can be found in the `api/request.http` file.

---

## Endpoints
- **`GET /`**: Test the rate limiting functionality.

---

## Dependencies
- **Go**: Application language.
- **Redis**: Data persistence.
- **Docker Compose**: Container orchestration.

---

## Contributions
Feel free to submit issues or pull requests to improve this project!

