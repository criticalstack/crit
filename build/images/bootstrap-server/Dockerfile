FROM golang:1.13-alpine AS go-builder

RUN apk add --no-cache ca-certificates git

ARG GOPROXY
ARG GOSUMDB

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/bootstrap-server/app/ cmd/bootstrap-server/app/
COPY cmd/bootstrap-server/internal/ cmd/bootstrap-server/internal/
COPY cmd/bootstrap-server/main.go cmd/bootstrap-server/main.go
COPY internal/ internal/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bootstrap-server cmd/bootstrap-server/main.go

############################
# Final stage: Just the executable and bare minimum other files
FROM scratch AS final
WORKDIR /
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /workspace/bootstrap-server /
EXPOSE 8080
ENTRYPOINT ["/bootstrap-server"]
