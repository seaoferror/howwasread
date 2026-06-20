FROM golang:1.26.0 AS stage1
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY common/ common/
COPY notification/ notification/
RUN cd notification && CGO_ENABLED=0 GOOS=linux go build -o /notification .

#/app/ (Working directory containing all the source code)
#/app/go.mod
#/app/go.sum
#/app/common/
#/app/chat/ (we execute go build here with .)
#/chat(The compiled Go binary, sitting by itself at the absolute root -o /chat)

FROM scratch
COPY --from=stage1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=stage1 /notification /notification
CMD ["/notification"]

#docker build -f backend/chat/Dockerfile backend