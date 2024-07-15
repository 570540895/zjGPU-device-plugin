FROM golang:1.22 as builder
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /go/src/zjgpu-device-plugin
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install -ldflags="-s -w"

FROM gcr.io/distroless/static-debian12
COPY --from=builder /go/bin/zjgpu-device-plugin /bin/zjgpu-device-plugin

CMD ["/bin/zjgpu-device-plugin"]
