FROM golang:1.15.3-alpine3.12
WORKDIR /build

# Take the current space cloud version as a argument
ARG SC_VERSION=0.21.5

# Copy all the source files
COPY . .
# Install the required packages
RUN apk --no-cache add ca-certificates wget unzip

# Build SC
RUN GOOS=linux CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o app .

FROM alpine:3.12
ARG SC_VERSION=0.21.5

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=0 /build/mission-control/build /root/.space-cloud/mission-control-v$SC_VERSION/build
COPY --from=0 /build/app .

CMD ["./app", "run"]
