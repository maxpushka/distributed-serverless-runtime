FROM golang:1.23-alpine AS deps

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download


FROM deps AS build

COPY . .
RUN CGO_ENABLED=0 go build -o ./serverless


FROM scratch

WORKDIR /app
COPY --from=build /serverless /usr/local/bin/serverless
EXPOSE 8080

ENTRYPOINT ["serverless"]

