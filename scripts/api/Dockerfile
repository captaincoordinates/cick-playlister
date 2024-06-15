FROM golang:1.22.3-alpine

WORKDIR /src
COPY cmd ./cmd
COPY internal ./internal
COPY go.mod .
COPY go.sum .

ENV GOOS=windows
ENV GOARCH=amd64
WORKDIR /src/cmd/cick-playlister

ARG date_tag
ENV output_path="/src/dist/${date_tag}/cick-playlister.exe"

CMD ["sh", "-c", "go build -o ${output_path} cick-playlister.go"]
