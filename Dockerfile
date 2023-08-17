FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build ./main.go

FROM golang:1.20
RUN apt-get update && apt-get install -y poppler-utils wv unrtf tidy
WORKDIR /app
COPY --from=builder /app/main /app/

CMD ["./main"]
