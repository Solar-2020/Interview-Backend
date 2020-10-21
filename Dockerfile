FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . /build

RUN pwd && ls && ls /build

# Build the application
RUN go build -o main /build/cmd/

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM alpine

COPY --from=builder /dist/main /

ENV POSTS_DB_CONNECTION_STRING=postgres://postgres:postgres@185.255.134.117:5432/posts?search_path=posts&sslmode=disable
ENV UPLOAD_DB_CONNECTION_STRING=postgres://postgres:postgres@185.255.134.117:5432/upload?search_path=upload&sslmode=disable

EXPOSE 8099

ADD ./scripts/run.sh /run.sh

ENV GIT_BRANCH="main"

#RUN touch /solar_main.log
#RUN ls /

# Command to run
#CMD ["/main >> /solar_main.log"]
#CMD "/bin/sh -c ls /"
#CMD "/bin/sh"
CMD /run.sh /main /var/log/solar_$GIT_BRANCH.log