FROM golang:1.25.0 AS stage1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY common/ common/
COPY messagepreprocess/ messagepreprocess/
RUN cd messagepreprocess && CGO_ENABLED=0 GOOS=linux go build -o /messagepreprocess .

#/app/ (Working directory containing all the source code)
#/app/go.mod
#/app/go.sum
#/app/common/
#/app/chat/ (we execute go build here with .)
#/chat(The compiled Go binary, sitting by itself at the absolute root -o /chat)

FROM scratch
COPY --from=stage1 /messagepreprocess /messagepreprocess
CMD ["/messagepreprocess"]

