FROM golang:latest AS buildstage

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY Makefile go.mod go.sum ./

FROM buildstage AS buildstage2

COPY . .
RUN ls -al

RUN go mod download

FROM buildstage2 AS buildstage3

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main cmd/backend/*.go
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o ./main cmd/backend/*.go

# Deploy Stage
# FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch
FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:latest

WORKDIR /app
COPY --from=buildstage3 /app/main /app/main

RUN echo "Performing other custom copying of resources here..."

CMD ["/app/main"]
