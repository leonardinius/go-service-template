ARG arch=amd64
# Builder
FROM --platform=linux/${arch} golang:1.23 as builder
WORKDIR /go/app
COPY . .
RUN make clean build
RUN cp /go/app/bin/service-template-linux-$(go env GOARCH) /go/app/bin/app

# Runtime
FROM --platform=linux/${arch} gcr.io/distroless/static-debian12 as runtime
COPY --from=builder /go/app/bin/app /app
COPY --from=builder /go/app/.env.docker /.env
EXPOSE 8080/tcp
CMD ["/app", "http", "-a", "0.0.0.0:8080"]
