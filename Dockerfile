FROM golang:1.26.4 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/bot .


FROM gcr.io/distroless/static-debian12
COPY --from=build /bin/bot /bin/bot
ENTRYPOINT [ "/bin/bot" ]