# --- 1. Build frontend ---
FROM node:20 AS frontend
WORKDIR /app
COPY web/frontend ./web/frontend
WORKDIR /app/web/frontend
RUN npm install && npm run build

# --- 2. Build Go app ---
FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy built frontend into correct embed path
COPY --from=frontend /app/web/frontend/dist ./cmd/server/dist

# Now web/frontend/dist/* exists — Go embed will succeed
RUN CGO_ENABLED=0 GOOS=linux go build -o shortlink ./cmd/server

# --- 3. Final image ---
FROM alpine:3.20
WORKDIR /srv
COPY --from=builder /app/shortlink /usr/local/bin/
COPY --from=builder /app/web/frontend/dist ./web/frontend/dist

EXPOSE 8080
ENTRYPOINT ["shortlink"]
