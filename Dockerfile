FROM golang:1.19.0-alpine3.16 as builder
WORKDIR /opt

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY api ./api
COPY cmd ./cmd
COPY internal ./internal
COPY config ./config

RUN CGO_ENABLED=0 go build -o auther ./cmd/service/main.go

FROM busybox
WORKDIR /bin
COPY --from=builder /opt/auther /bin/auther
ENV GIN_MODE=release
CMD /bin/auther