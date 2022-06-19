FROM golang:1.18

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /chat-app

EXPOSE 8000

CMD [ "/chat-app" ]
