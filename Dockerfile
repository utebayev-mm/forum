FROM golang:1.16.7
LABEL project="Forum" \
      authors="utebayev.mm and darwin939" \
      link="https://git.01.alem.school/utebayev.mm/forum"
WORKDIR /forum
COPY . .
COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN go build -o main .

EXPOSE 8080
CMD ["./main"]