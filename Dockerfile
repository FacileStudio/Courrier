FROM oven/bun:1 AS client-build
WORKDIR /app
COPY apps/client/package.json apps/client/bun.lock* ./
RUN bun install --frozen-lockfile
COPY apps/client/ .
RUN bun run build

FROM golang:1.26-alpine AS api-build
WORKDIR /app
COPY apps/api/ .
RUN go build -mod=vendor -trimpath -ldflags="-s -w" -o /courrier .

FROM gcr.io/distroless/static-debian12
COPY --from=api-build /courrier /courrier
COPY --from=client-build /app/build /client
# The distroless base can carry its own WorkingDir (/home/nonroot on the
# :nonroot variant), which would make a relative ./client resolve there and
# the SPA silently not be served at all. Be explicit.
ENV CLIENT_DIR=/client

EXPOSE 4000
ENTRYPOINT ["/courrier"]
