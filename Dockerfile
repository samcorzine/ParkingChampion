#first stage - builder
FROM golang:1.12.0-stretch as builder
#download modules
COPY go.mod /ParkingChamp/
WORKDIR /ParkingChamp
ENV GO111MODULE=on
RUN go mod download
#copy source
COPY src/ /ParkingChamp
#testing
COPY testfiles /testfiles
RUN go test
#build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version=1.0.0" -o ParkingChamp

#Start from scratch, add binary and zoneinfo
FROM scratch
COPY --from=builder /ParkingChamp/ParkingChamp /ParkingChamp
COPY --from=builder /usr/share/zoneinfo/ /usr/share/zoneinfo/
CMD ["/ParkingChamp"]
