FROM golang:1.19-alpine as build
RUN apk --no-cache add wget

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /build/api ./cmd/api
RUN go build -o /build/importcsv ./cmd/importcsv

FROM alpine:3.14
RUN apk --no-cache add ca-certificates curl
RUN mkdir /opt/app
EXPOSE 5000

COPY --from=build /build/api /opt/app/
COPY --from=build /build/importcsv /opt/app/
RUN chmod +x -R /opt/app/

CMD /opt/app/api
# ENTRYPOINT ["/opt/app/entrypoint.sh"]
