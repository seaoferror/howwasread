FROM golang:1.25.0 AS stage1
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY common/ common/
COPY messagerelay/ messagerelay/
RUN cd messagerelay && CGO_ENABLED=0 GOOS=linux go build -o /messagerelay .

#/app/ (Working directory containing all the source code)
#/app/go.mod
#/app/go.sum
#/app/common/
#/app/chat/ (we execute go build here with .)
#/chat(The compiled Go binary, sitting by itself at the absolute root -o /chat)

FROM scratch
COPY --from=stage1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=stage1 /messagerelay /messagerelay
CMD ["/messagerelay"]

#docker build -f backend/chat/Dockerfile backend