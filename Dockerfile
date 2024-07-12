FROM golang:1.22 as build
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /go/src/zjGPU-device-plugin
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install -ldflags="-s -w"

#replace: gcr.io/distroless/static-debian12
FROM gcr.lank8s.cn/distroless/static-debian12  
COPY --from=build /go/bin/zjGPU-device-plugin /bin/zjGPU-device-plugin

CMD ["/bin/zjGPU-device-plugin"]
