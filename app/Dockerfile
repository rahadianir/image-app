FROM golang:1.24 AS build

# Set destination for COPY
WORKDIR /app

# Copy the source code
COPY . .

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o image-app

# Setup final image
FROM alpine:latest 

WORKDIR /app

# Copy executable from build step
COPY --from=build /app/image-app image-app

EXPOSE 8080

# Run
CMD ["./image-app"]