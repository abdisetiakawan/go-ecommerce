# Pin specific version for stability
FROM golang:1.23-bullseye AS build-base

WORKDIR /app

# Copy only files required to install dependencies (better layer caching)
COPY go.mod go.sum ./

# Use cache mount to speed up install of existing dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

FROM build-base AS dev

# Install air for hot reload & delve for debugging
RUN go install github.com/air-verse/air@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .
CMD ["air", "-c", ".air.toml"]

FROM build-base AS build-production

# Add non-root user
RUN useradd -u 1001 nonroot

COPY . .

# Compile application with consistent binary name
RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o go-ecommerce ./cmd/web

FROM scratch

WORKDIR /

# Copy the passwd file for non-root user
COPY --from=build-production /etc/passwd /etc/passwd

# Copy the app binary from the build stage
COPY --from=build-production /app/go-ecommerce /go-ecommerce


# Use non-root user
USER nonroot

# Indicate expected port
EXPOSE 3000

CMD ["/go-ecommerce"]