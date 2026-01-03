# Backend Optimizer - AI Chat with Memory

Fast AI chat backend with conversation memory using Redis caching.

## Features
- JWT authentication
- Conversation history (30 min sessions)
- Multi-user support
- Groq AI integration

## Prerequisites
- Go 1.21+
- Redis 7.0+
- MongoDB 6.0+
- Auth0 account
- Groq API key

## Auth0 Setup
1. Create account at auth0.com
2. Create new application (Regular Web Application)
3. Note your Domain and configure in environment variables
4. Add `http://localhost:8080` to Allowed Callback URLs

## Installation

1. Clone repository:
```bash
git clone https://github.com/yourusername/backend-optimizer
cd backend-optimizer
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables (create `.env` file or export):
```bash
export REDIS_URL="localhost:6379"
export MONGO_URL="mongodb://localhost:27017"
export AUTH0_DOMAIN="your-domain.auth0.com"  # From Auth0 dashboard
export AUTH0_AUDIENCE="your-audience"         # From Auth0 dashboard
export GROQ_API_KEY="gsk_..."                 # From console.groq.com
```

4. Run the server:
```bash
go run cmd/server/main.go
```

5. Server runs at `http://localhost:8080`

## Usage

### Register/Login
POST `/Register` or `/Login` with user credentials

Body:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

Returns JWT token for authenticated requests.

### Chat with AI
POST `/askAI`

Headers:
```json
{
  "Authorization": "Bearer YOUR_JWT_TOKEN",
  "Content-Type": "application/json"
}
```

Body:
```json
{
  "role": "user",
  "content": "Your message here"
}
```

## Quick Test

After starting the server:

1. Register a user:
```bash
curl -X POST http://localhost:8080/Register \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"Test123!"}'
```

2. Login to get JWT token:
```bash
curl -X POST http://localhost:8080/Login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"Test123!"}'
```

3. Chat with AI (use token from step 2):
```bash
curl -X POST http://localhost:8080/askAI \
    -H "Authorization: Bearer YOUR_TOKEN_HERE" \
    -H "Content-Type: application/json" \
    -d '{"role":"user","content":"Hello!"}'
```

## Project Structure
```
.
├── README.md
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
└── internal
    ├── ai
    │   ├── groq_dto.go
    │   └── model.go
    ├── api
    │   ├── client_auth.go
    │   ├── handlers.go
    │   ├── handlers_ai.go
    │   ├── handlers_todo.go
    │   ├── handlers_user.go
    │   └── routes.go
    ├── cache
    │   └── cache.go
    ├── db
    │   └── db.go
    ├── models
    │   ├── cachedb.go
    │   ├── todos.go
    │   └── user.go
    ├── service
    │   └── risk_service.go
    ├── tests
    │   └── benchmark_test.go
    └── utils
        ├── config.go
        └── logger.go

12 directories, 21 files
```

## Roadmap
- [ ] Cosine similarity for smart context selection
- [ ] Two-tier caching (hot/cold)
- [ ] Embedding-based message retrieval

## License
MIT