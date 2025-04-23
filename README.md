# Starting up the project
- This project (that runs inside a container) and the other containers all use environment variables, so you must define them in a .env file.
- Create an environment file and create these variables
    - JWT_SECRET="super_secret_seed" // seed for creating jwt tokens
    - DB_URL="postgres://user:password@postgres:5432/finTechDB?sslmode=disable" // connection string for postgres
        - this connection string is read as intended by sql.Open() as the host part (postgres:5432) is being translated to (127.0.0.1:5432)
        - so i have harcoded the connection string insie /internal/storage/postgres.go, you can change the connection string there
        - the same translation is happening in redis client as well, so i have harcoded that aswell 
        - you can change that in /internal/storage/redis.go  
    - APP_PORT="8080" // the port where this project (the golang app) server listens through
    - APP_HOST_PORT="8081" // the port where the golang app contianer listens through

- You can configure the following service variables in docker-compose.yaml under the environment section of each service
- postgres service
    - POSTGRES_USER: user
    - POSTGRES_PASSWORD: password
    - POSTGRES_DB: finTechDB
- pgAdmin service
    - PGADMIN_DEFAULT_EMAIL: admin@example.com
    - PGADMIN_DEFAULT_PASSWORD: admin123
---
- After configuring the variables run the container for the first time using `docker compose up --build`
- After that you anytime you want to run the contianer use `docker compose up`
- If you make a change to the code you have to restart the container but not all of the containers just restart finTech_app container 
---
- When running the container for the first time, docker will create a folder called `data` which is a bind mount, configured in the volume section of postgres, it is the docker volume for our postgres server.
- And also `init` is another bind mount folder for our postgres service, postgres uses it for initialization but only once after that it uses `/data` for persistence
- So dont delete these folders
--- 
- docker exec -t postgres pg_dump -U user -d finTechDB > seed.sql
- use the above command to get an sql dumb of your database that can be used to replicate your database
- you can then put seed.sql inside /init to intialize your database on docker compose build 
---
# Project Structure so far
<pre>
<code>

├── cmd/
│   └── api/                     # Main application entrypoint
│       └── main.go
├── internal/
│   ├── api/                     # HTTP handlers & routers
│   │   ├── middleware/          # Middlewares: auth, rate limit, session
│   │   │       ├── auth_middleware.go       # checks for access token validity before request is forwarded to endpoints 
│   │   │       └── rate_limit_middlware.go  # limits the rate of requests a user can send both for sign-in and other endpoints
│   │   └── routes.go         # maps different handlers to endpoints, wraps handlers in middlwares and defines middleware chains
│   ├── fee_engine/
│   │   ├── feeengine.go            # contains logic for calculating fee percentage
│   │   └── types.go                # contains structs and types for fee engine
│   ├── handler/                    # contains handler for different endpoints
│   │   ├── auth_handlers.go        # contains auth handlers like sign-in, sign-up, log-out and refresh handlers
│   │   └── calculate_fee.go        # contains handler for fee calculation
│   ├── storage/                    # DB and Redis access
│   │   ├── redis.go                # configuration for the redis go client
│   │   └── postgres.go             # configuration for the postgres  go client
│   └── utils/                      # DB and Redis access
│       └── auth_utils.go           # utility function for auth handlers
├── .env                            # Environment variables
├── go.mod
├── go.sum
├── README.md
├── Dockerfile
├── docker-compose.yaml
├── .gitignore
├── .dockerignore
└── postman_collection.json      # Postman collection

</code>
</pre>
