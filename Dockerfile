FROM golang:alpine
WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -o /banner_app

EXPOSE 9000

CMD ["/banner_app"]