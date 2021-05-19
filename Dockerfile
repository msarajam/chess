FROM golang:1.16

WORKDIR /go/src/chess
COPY . .

RUN go get -d -v ./cmd/webServer/
RUN go install -v ./cmd/webServer/main.go

CMD ["main"]

