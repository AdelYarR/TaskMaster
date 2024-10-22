FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make

EXPOSE 8000

CMD ["./taskmaster"]