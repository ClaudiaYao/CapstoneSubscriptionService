FROM golang:latest as builder


RUN mkdir /subscription
WORKDIR /subscription

ADD . ./subscription
COPY . .

COPY Makefile go.mod go.sum .env ./
COPY .env ./cmd/

# Builds your app with optional configuration

RUN go build -o subscription-service ./cmd

# Tells Docker which network port your container listens on
EXPOSE 8081

# Specifies the executable command that runs when the container starts
CMD [ "/subscription/subscription-service"]