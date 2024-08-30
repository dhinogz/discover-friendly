FROM node:20-slim as node

# 1. Builder step builds tailwindcss styles, templ templates, and Go binary
FROM golang:1.23-bookworm as builder

# Copy node dependencies from base node image
COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/include/node /usr/local/include/node
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
COPY --from=node /opt/yarn-v*/bin/* /usr/local/bin/
COPY --from=node /opt/yarn-v*/lib/* /usr/local/lib/
RUN ln -vs /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm \
    && ln -vs /usr/local/lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx

WORKDIR /app

# Install node deps (tailwindcss)
COPY package-lock.json package.json ./
RUN npm ci 

# Install Go dependencies and templ
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate templ, build CSS and embed assets
COPY . .
RUN templ generate
RUN npx tailwindcss -m -i ./assets/tailwind.css -o ./assets/dist/styles.min.css
RUN go generate ./...

# Build binary
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o web 

# 2. Serve step
FROM debian:12-slim

WORKDIR /app

COPY --from=builder /app/web /app/web

EXPOSE 8090
ENTRYPOINT ["/app/web", "serve", "--http=0.0.0.0:8090"]
