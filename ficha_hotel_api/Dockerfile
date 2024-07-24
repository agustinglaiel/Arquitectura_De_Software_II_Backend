FROM golang:1.20 AS builder

ADD . /ficha_hotel_api

WORKDIR /ficha_hotel_api
RUN if [ ! -f go.mod ]; then go mod init ficha_hotel_api; fi
RUN go mod tidy
RUN go build -o ficha_hotel_api
RUN chmod +x /ficha_hotel_api

EXPOSE 8080

ENTRYPOINT ["./ficha_hotel_api"]