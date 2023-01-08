FROM golang:1.18.9-alpine3.16

#LABEL maintainer="liemkg1234@gmail.com"

RUN mkdir /server
WORKDIR /server

ENV GO111MODULE=on CGO_ENABLED=0

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN #go build -o /server/main /server/main.go

RUN go build -o /build

EXPOSE 6868

CMD [ "/build" ]

#docker build --tag golang .
#docker run --name backend_golang -p 6868:6868 -d golang:latest