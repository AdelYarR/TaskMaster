FROM golang:latest

RUN mkdir /app
ADD . /app/
WORKDIR /app/cmd/main
RUN make
WORKDIR /app
CMD [ "/app/cmd/main/main" ]    