FROM golang:1.20.6-alpine3.17 as builder
RUN mkdir /src
ADD . /src/
WORKDIR /src
RUN go build -ldflags "-s -w -X main.version=$(cat VERSION)" -o kubeops
FROM alpine
COPY --from=builder /src/kubeops /app/kubeops
WORKDIR /app
ENTRYPOINT ["/app/kubeops"]
