FROM golang:latest AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /rest ./rest

FROM docker:cli

WORKDIR /
COPY --from=build /rest /rest

EXPOSE 80
ENTRYPOINT [ "/rest" ]
