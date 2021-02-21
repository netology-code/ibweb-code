FROM golang:1.15-alpine AS build
COPY . /app
WORKDIR /app
ENV CGO_ENABLED 0
RUN go build -o tx ./cmd

FROM scratch
COPY --from=build /app/tx /tx
CMD ["/tx"]