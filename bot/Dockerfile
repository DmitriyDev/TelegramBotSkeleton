FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go get -u gopkg.in/tucnak/telebot.v2

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o bot .

CMD ["/go/src/app/bot"]