# Super multi-stage GO image build:

ARG BASE_IMAGE=gcr.io/distroless/base

# -------------------------------------------------#
FROM $BASE_IMAGE:debug as os_config_image

ARG BASE_IMAGE

RUN ["adduser","-h","/","-s","/sbin/nologin","-D","-H","app_user"]

# -------------------------------------------------#
FROM golang:1.12-stretch as build_image

ARG BASE_IMAGE

WORKDIR /go/src/app

COPY *.go ./

RUN go get -d -v ./...

RUN go test -v -cover ./...

RUN go install -v ./...

# -------------------------------------------------#
FROM $BASE_IMAGE as service_image

ARG BASE_IMAGE

LABEL base_image=$BASE_IMAGE

## LABEL owner_team=YOUR_TEAM_NAME

COPY --from=os_config_image --chown=root:root /etc/passwd /etc/group /etc/

COPY --from=build_image --chown=root:root /go/bin/app /app

USER app_user:app_user

ARG LISTEN_PORT="9999"

ENV ASM_BINDADDRESS=":$LISTEN_PORT"

EXPOSE $LISTEN_PORT

CMD ["/app"]
