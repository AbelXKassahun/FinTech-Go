Project Structure i will ideally be following for this project
<pre>
<code>

├── cmd/
│   └── api/                     # Main application entrypoint
│       └── main.go
├── internal/
│   ├── api/                     # HTTP handlers & routers
│   │   ├── middleware/          # Middlewares: auth, rate limit, session
│   │   ├── handlers/            # Handlers: simulate, report, webhook, etc.
│   │   └── routes.go
│   ├── config/                  # Config loading (env, flags)
│   ├── service/                 # Business logic layer
│   │   ├── session/
│   │   ├── feeengine/
│   │   ├── simulation/
│   │   ├── webhook/
│   │   ├── report/
│   │   └── tier/
│   ├── model/                   # DTOs and schema definitions
│   ├── storage/                 # DB and Redis access
│   │   ├── redis/
│   │   └── postgres/
│   ├── jobs/                    # Async workers: write-behind, webhook retries
│   ├── events/                  # Event bus, pub-sub or stream logic
│   ├── utils/                   # Helpers: CSV, logging, validation
│   └── observability/           # Logging, metrics, tracing setup
├── pkg/                         # Public reusable libraries
│   ├── auth/                    # JWT generation and validation
│   └── ratelimit/               # Redis-based rate limiter
├── test/
│   ├── integration/             # Integration tests
│   └── mocks/                   # Testify mocks
├── api/
│   └── openapi.yaml             # Swagger/OpenAPI spec
├── deployments/
│   ├── docker/                  # Dockerfile, docker-compose.yml
│   └── k6/                      # Load testing scripts
├── scripts/                     # Seed scripts, CSV generators
├── .env                         # Environment variables
├── Makefile                     # Build/test commands
├── go.mod
├── go.sum
├── README.md
└── postman_collection.json      # Postman collection

</code>
</pre>
