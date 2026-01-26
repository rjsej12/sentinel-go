# 1. Base Image
FROM golang:1.25.6-alpine

# 2. Set Working Directory
WORKDIR /app

# 3. Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# 4. Copy source code
COPY . .

# 5. Build
RUN go build -buildvcs=false -o sentinel-api ./cmd/api 

# 6. Open port 8080
EXPOSE 8080

# 7. Start the application
CMD ["./sentinel-api"]