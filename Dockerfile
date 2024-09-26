


FROM golang:1.23-bookworm

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN apt update && apt install -y ca-certificates
RUN go build -o main .
#RUN go build -o /out/bot .



#RUN apk add ca-certificates


EXPOSE 8086

CMD [ "./main" ]
