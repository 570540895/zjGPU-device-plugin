FROM golang:1.22 as builder
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /go/src/zjGPU-device-plugin
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install -ldflags="-s -w"

#replace: gcr.lank8s.cn/distroless/static-debian12
FROM gcr.io/distroless/static-debian12
COPY --from=builder /go/bin/zjGPU-device-plugin /bin/zjGPU-device-plugin

CMD ["/bin/zjGPU-device-plugin"]
