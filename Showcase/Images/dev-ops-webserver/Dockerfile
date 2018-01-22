FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

RUN go get -u github.com/golang/lint/golint

ENTRYPOINT ["./buildandrun.sh", "echo", "to do"]
EXPOSE 8080
