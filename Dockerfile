FROM golang:1.13.1-alpine3.10 as builder
RUN mkdir /src
ADD . /src/
WORKDIR /src
RUN go build -ldflags "-s -w -X main.version=$(cat VERSION)" -o kubeops
FROM alpine
COPY --from=builder /src/kubeops /app/kubeops
WORKDIR /app
ENTRYPOINT ["/app/kubeops"]
