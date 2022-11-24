FROM golang:1.18

EXPOSE 9001/udp
EXPOSE 9000/udp
EXPOSE 53/udp
EXPOSE 53/tcp
EXPOSE 80/tcp
EXPOSE 443/tcp

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app

CMD ["app"]