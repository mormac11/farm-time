# Frontend build stage
FROM node:20-alpine AS frontend

WORKDIR /app/web

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build

# Backend build stage
FROM golang:1.21-alpine AS backend

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
COPY --from=frontend /app/web/dist ./web/dist

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

COPY --from=backend /app/server .

EXPOSE 8080

ENV PORT=8080

CMD ["./server"]
