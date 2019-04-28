#first stage - builder
FROM golang:1.12.0-stretch as builder
COPY src/ /MultiStage
COPY go.mod /MultiStage
WORKDIR /MultiStage
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /parking
COPY testing.json /testing.json
RUN go test

#second stage
FROM scratch
COPY --from=builder /parking /parking
CMD ["/parking"]
