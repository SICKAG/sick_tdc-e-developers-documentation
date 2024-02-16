FROM meinside/alpine-golang:1.16.6-armv7

WORKDIR /app

# Print the current directory contents
COPY . .

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o diogo main.go

ENTRYPOINT ["/app/diogo"]
