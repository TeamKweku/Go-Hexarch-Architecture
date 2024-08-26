FROM golang:1.23-bullseye AS base

WORKDIR /app

# Copy .env file
COPY .env ./

# By copying go.mod and go.sum files first, the dependencies will be
# redownloaded only when these files change.
COPY go.sum go.mod ./
RUN go mod download

COPY . .

FROM base AS build

ARG GO_COMPILER_CACHE

RUN go build -o /out/ ./cmd/server ./cmd/healthcheck

FROM gcr.io/distroless/base-debian11:nonroot AS optimized

ARG PORT

COPY --from=build --chown=nonroot:nonroot /out/ /
COPY --from=base --chown=nonroot:nonroot /app/.env /home/nonroot/

USER nonroot

EXPOSE $PORT

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/healthcheck"]

CMD ["/server"]
