FROM oven/bun:1 AS client-build
WORKDIR /app
COPY apps/client/package.json apps/client/bun.lock* ./
RUN bun install --frozen-lockfile
COPY apps/client/ .
RUN bun run build

FROM golang:1.24-alpine AS api-build
WORKDIR /app
COPY apps/api/ .
RUN go build -o /courrier .

FROM gcr.io/distroless/static-debian12
COPY --from=api-build /courrier /courrier
COPY --from=client-build /app/build /client
EXPOSE 4000
ENTRYPOINT ["/courrier"]
