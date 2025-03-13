FROM golang:latest AS build_base
WORKDIR /src
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN export GIT_COMMIT=$(git rev-list -1 HEAD) && \
CGO_ENABLED=0 go build -ldflags "-X main.GitCommit=$GIT_COMMIT" -o hello-world main.go
RUN chmod +x hello-world

FROM alpine:3.19
COPY --from=build_base /src/hello-world /opt/hello-world/hello-world
COPY img/* /opt/hello-world/
WORKDIR /opt/hello-world
ENTRYPOINT ["/opt/hello-world/hello-world"]
