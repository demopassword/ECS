# Build stage
FROM arm64v8/golang:1.18.2-alpine AS build

WORKDIR /source

COPY main.go .

RUN go mod init example.com/m/v2 \
 && go get github.com/gin-gonic/gin \
 && go build -o ./main

# Runtime stage
FROM arm64v8/alpine

WORKDIR /app

COPY --from=build /source/main ./

RUN adduser --disabled-password app

RUN apk add --no-cache gcompat && apk add --no-cache libc6-compat && apk add --no-cache curl && apk add --no-cache stress-ng && mv /usr/bin/stress-ng /usr/bin/stress &&  chmod +x ./main && chown app /app
USER app

#ENTRYPOINT ["./main"]
