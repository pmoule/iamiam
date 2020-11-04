FROM golang:1.15.3 AS builder
ENV GO111MODULE=on
RUN mkdir iamiam
WORKDIR /iamiam
COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . /iamiam
RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o iam

FROM prom/busybox:latest
RUN mkdir iamiam
WORKDIR /iamiam
COPY --from=builder iamiam/iam .
COPY README.md .
COPY LICENSE .
CMD [ "/iamiam/iam" ]