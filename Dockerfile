FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/cespare/reflex@latest

# copy the source code
COPY . .

# Watch for .go changes and rerun
# CMD ["reflex", "-r", "\\.go$", "--", "go", "run", "./cmd/main.go"]
CMD ["go", "run", "./cmd/main.go"]


# EXPOSE 8080