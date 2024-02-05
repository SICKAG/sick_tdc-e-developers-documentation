FROM golang:1.19
WORKDIR /app

COPY . ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o detect-direct main.go

CMD ["/detect-direct"]