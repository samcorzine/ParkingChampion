#first stage - builder
FROM golang:1.12.0-stretch as builder
COPY go.mod /ParkingChamp/
WORKDIR /ParkingChamp
ENV GO111MODULE=on
RUN go mod download
COPY src/ /ParkingChamp
#testing
COPY testfiles /testfiles
RUN go test
RUN CGO_ENABLED=0 GOOS=linux go build -o /parking

#second stage
FROM scratch
COPY --from=builder /parking /parking
CMD ["/parking"]
