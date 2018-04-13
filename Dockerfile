FROM golang:alpine as builder
WORKDIR /go/src/github.com/Adictes/pets-health
COPY . .
RUN go build main.go

FROM alpine:latest
COPY --from=builder /go/src/github.com/Adictes/pets-health .
CMD [ "./main" ]