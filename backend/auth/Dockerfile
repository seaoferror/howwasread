FROM golang:1.25.0 AS stage1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY common/ common/
COPY auth/ auth/
RUN cd auth && CGO_ENABLED=0 GOOS=linux go build -o /auth .

FROM scratch
COPY --from=stage1 /auth /auth
CMD ["/auth"]