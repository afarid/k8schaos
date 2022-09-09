FROM golang:1.19-alpine as builder
RUN mkdir /build
ADD .  /build/
WORKDIR /build
RUN apk add  --no-cache git
RUN CGO_ENABLED=0 GOOS=linux go build  -a -o k8s-chaos cmd/main.go


FROM scratch
COPY --from=builder /build/k8s-chaos .
ENTRYPOINT [ "./k8s-chaos" ]
