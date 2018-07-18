FROM node:10 as builder
WORKDIR /app
COPY ./frontend/ /app

ENV VUE_APP_FAUCET_URL http://116.62.62.39:30200

RUN npm i yarn -g && VUE_APP_FAUCET_URL=$VUE_APP_FAUCET_URL yarn install && yarn build

FROM golang:1.10.3-alpine3.7 as go-builder
ENV GOPATH       /root/go
ENV REPO_PATH    $GOPATH/src/github.com/irisnet/explorer/server
ENV PATH         $GOPATH/bin:$PATH

RUN mkdir -p GOPATH REPO_PATH

COPY ./server/ $REPO_PATH
WORKDIR $REPO_PATH

RUN apk add --no-cache make git && go get github.com/golang/dep/cmd/dep && dep ensure && make build


FROM alpine:3.7
WORKDIR /app/server
COPY --from=builder /app/dist/ /app/frontend/dist
COPY --from=go-builder /root/go/src/github.com/irisnet/explorer/server/build/ /app/server/
CMD ['./irisplorer']