FROM golang:1.20 AS builder

ADD . /user_reserva_dispo_api

WORKDIR /user_reserva_dispo_api
RUN if [ ! -f go.mod ]; then go mod init user_reserva_dispo_api; fi
RUN go mod tidy
RUN go build -o user_reserva_dispo_api
RUN chmod +x /user_reserva_dispo_api

EXPOSE 8080

ENTRYPOINT ["./user_reserva_dispo_api"]