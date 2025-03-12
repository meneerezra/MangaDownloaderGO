FROM golang:1.24.1
LABEL authors="meneerezra"

WORKDIR /usr/local/app
COPY src ./src

WORKDIR /usr/local/app/src

RUN go mod tidy

CMD ["go", "run", "."]