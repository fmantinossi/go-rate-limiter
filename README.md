# Go Rate Limiter

A configurable rate limiter implementation in Go that supports both IP-based and token-based rate limiting.

## Features

- IP-based rate limiting
- Token-based rate limiting
- Configurable request limits and time windows
- Configurable blocking duration for exceeded limits
- Redis-based storage with extensible storage interface
- Docker and Docker Compose support

## Configuration

The rate limiter can be configured using environment variables:

```env
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# IP-based Rate Limiting
RATE_LIMIT_IP_REQUESTS=5      # Maximum requests per window
RATE_LIMIT_IP_WINDOW=1s       # Time window (e.g., 1s, 1m, 1h)
RATE_LIMIT_IP_BLOCK_DURATION=5m  # Duration to block after limit exceeded

# Token-based Rate Limiting
RATE_LIMIT_TOKEN_REQUESTS=10
RATE_LIMIT_TOKEN_WINDOW=1s
RATE_LIMIT_TOKEN_BLOCK_DURATION=5m

# Server Configuration
SERVER_PORT=8080
```

## Running with Docker Compose

1. Clone the repository
2. Run the application:
   ```bash
   docker-compose up --build
   ```

The server will start on port 8080.

## API Usage

### Making Requests

- Without a token (IP-based limiting):
  ```bash
  curl http://localhost:8080/
  ```

- With a token (token-based limiting):
  ```bash
  curl -H "API_KEY: your-token" http://localhost:8080/
  ```

### Rate Limit Response

When the rate limit is exceeded, the API will respond with:
- Status Code: 429
- Message: "you have reached the maximum number of requests or actions allowed within a certain time frame"

## Examples

### Testing IP-based Rate Limiting

1. **Basic Test**
```bash
# Make 5 requests (within limit)
for i in {1..5}; do curl http://localhost:8080/; done

# Make 6th request (should be blocked)
curl http://localhost:8080/
```

2. **Testing with Different IPs**
```bash
# Using different IPs
curl --interface eth0 http://localhost:8080/
curl --interface eth1 http://localhost:8080/
```

### Testing Token-based Rate Limiting

1. **Basic Test**
```bash
# Make 10 requests (within limit)
for i in {1..10}; do curl -H "API_KEY: test-token" http://localhost:8080/; done

# Make 11th request (should be blocked)
curl -H "API_KEY: test-token" http://localhost:8080/
```

2. **Testing with Different Tokens**
```bash
# Using different tokens
curl -H "API_KEY: token1" http://localhost:8080/
curl -H "API_KEY: token2" http://localhost:8080/
```

### Testing Block Duration

1. **IP Block Test**
```bash
# Exceed limit
for i in {1..6}; do curl http://localhost:8080/; done

# Wait for block duration (5 minutes)
sleep 300

# Try again (should work)
curl http://localhost:8080/
```

2. **Token Block Test**
```bash
# Exceed limit
for i in {1..11}; do curl -H "API_KEY: test-token" http://localhost:8080/; done

# Wait for block duration (5 minutes)
sleep 300

# Try again (should work)
curl -H "API_KEY: test-token" http://localhost:8080/
```

### Advanced Testing

1. **Concurrent Requests**
```bash
# Test with 10 concurrent requests
for i in {1..10}; do
  curl -H "API_KEY: test-token" http://localhost:8080/ &
done
wait
```

2. **Monitoring Redis**
```bash
# Access Redis CLI
docker-compose exec redis redis-cli

# Monitor keys
redis-cli monitor
```

## Development

### Running Tests
```bash
docker-compose exec app go test ./...
```

### Modifying Configuration
Edit the `docker-compose.yml` file to change environment variables:
```yaml
environment:
  - RATE_LIMIT_IP_REQUESTS=5
  - RATE_LIMIT_IP_WINDOW=1s
  - RATE_LIMIT_TOKEN_REQUESTS=10
  - RATE_LIMIT_TOKEN_WINDOW=1s
```

### Debugging
```bash
# Access application container
docker-compose exec app sh

# View logs
docker-compose logs -f
```

## Architecture

The rate limiter is designed with the following components:

1. **Storage Interface**: Abstracts the storage backend, making it easy to switch between different storage solutions.
2. **Redis Implementation**: Default storage implementation using Redis.
3. **Rate Limiter**: Core logic for rate limiting, supporting both IP and token-based limiting.
4. **Middleware**: HTTP middleware that can be easily integrated with any Go web application.

## Cleanup

To stop and remove all containers and volumes:
```bash
docker-compose down -v
```