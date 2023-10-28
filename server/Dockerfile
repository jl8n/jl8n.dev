# Start from the latest golang base image
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./src/go.mod ./src/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./src/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  
WORKDIR /root/app

# Since the development environment expects album-art/ to be outside src/ (aunt/uncle)
# the same directory structure needs to be maintained here too
RUN mkdir /root/album-art

# Set the GO_ENV environment variable to "production"
ENV GO_ENV=production

# Copy only the executable from the first stage.
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

EXPOSE 3000

CMD ["./main"]  
