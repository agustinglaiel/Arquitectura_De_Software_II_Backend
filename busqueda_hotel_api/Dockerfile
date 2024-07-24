FROM golang:1.20 AS builder

ADD . /busqueda_hotel_api

WORKDIR /busqueda_hotel_api
RUN if [ ! -f go.mod ]; then go mod init busqueda_hotel_api; fi
RUN go mod tidy
RUN go build -o busqueda_hotel_api
RUN chmod +x /busqueda_hotel_api

EXPOSE 8070

ENTRYPOINT ["./busqueda_hotel_api"]