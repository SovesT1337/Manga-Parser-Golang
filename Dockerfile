FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && apk add --no-cache tzdata ca-certificates
COPY . .
ARG CMD_PATH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app ${CMD_PATH}

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /out/app /app
ENV TZ=Europe/Moscow
USER nonroot:nonroot
ENTRYPOINT ["/app"]

