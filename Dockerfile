#Builder
FROM golang:1.13-alpine as builder

#Upgrade alpine packages... install git.
RUN apk update && apk upgrade && \
    apk add --no-cache git

#RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .  

#RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o deliveryapp.stockservice


#Runner
# Run container
FROM alpine:latest

RUN apk --no-cache add ca-certificates

#RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/deliveryapp.stockservice .

CMD ["./deliveryapp.stockservice"]