FROM golang:latest

WORKDIR /output

COPY . /output

RUN go mod tidy

RUN go build -o main .

EXPOSE 3000

CMD ["/output/main"]
