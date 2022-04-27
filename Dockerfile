FROM golang:1.17-alpine as Builder

LABEL maintainer="vegarfae@stud.ntnu.no"
LABEL stage=builder

# Set up execution environment in container's GOPATH
WORKDIR /go/src/app/cmd

# Copy relevant folders into container
COPY ./cmd /go/src/app/cmd
COPY ./internal /go/src/app/internal
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server

# To get the time zone data
FROM alpine:latest as alpine-with-tz
RUN apk --no-cache add tzdata zip
WORKDIR /usr/share/zoneinfo

#Compressing the zone data
RUN zip -q -r -0 /zoneinfo.zip .

# Final container
FROM scratch AS final

LABEL maintainer="vegarfae@stud.ntnu.no"

# Root as working directory to copy compiled file to
WORKDIR /

# Retrieve binary from builder container
COPY --from=builder /go/src/app/cmd/server .
COPY ./serviceAccountKey.json .

# Setting time zone data
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine-with-tz /zoneinfo.zip /
ENV TZ=Europe/Berlin

# Fetching the cert hints.
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Instantiate server
CMD ["./server"]

# Time zone and cert inspired from https://github.com/jeremyhuiskamp/golang-docker-scratch/blob/main/Dockerfile.externaltz