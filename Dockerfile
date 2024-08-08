FROM golang:latest AS builder

# create base working directory inside container
WORKDIR /app

#copying the dependencies first
#they are most unlikely to change
#so layer caching will make docker build fast
COPY go.mod go.sum ./

# Install all the dependencies
RUN go mod tidy

# Now copy all
COPY . .

# build the go application
RUN CGO_ENABLED=0 GOOS=linux go build -o BookServer .

#install bash


# pull alpine latest image, it contains necessary files for running go binaries
FROM alpine:latest
RUN apk  --no-cache add bash
WORKDIR /app

# COPY everything from last image working directory
COPY --from=builder /app/BookServer .
ENTRYPOINT ["./BookServer"]