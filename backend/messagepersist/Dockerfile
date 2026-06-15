FROM golang:1.25.0 AS stage1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY common/ common/
COPY messagepersist/ messagepersist/
RUN cd messagepersist && CGO_ENABLED=0 GOOS=linux go build -o /messagepersist .

#/app/ (Working directory containing all the source code)
#/app/go.mod
#/app/go.sum
#/app/common/
#/app/chat/ (we execute go build here with .)
#/chat(The compiled Go binary, sitting by itself at the absolute root -o /chat)

FROM scratch
COPY --from=stage1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=stage1 /messagepersist /messagepersist
CMD ["/messagepersist"]

#docker build -f backend/messagepersist/Dockerfile backend