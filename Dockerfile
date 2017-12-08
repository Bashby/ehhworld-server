FROM golang

# Force SSH auth when getting deps for GoLang
RUN echo "[url \"git@bitbucket.org:\"]\n\tinsteadOf = https://bitbucket.org/" >> /root/.gitconfig

# Ignore SSH checking for getting deps via go get from GoLang
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

WORKDIR /go/src/app
COPY ./src .

WORKDIR /go/src/app/bitbucket.org/ehhio/ehhworldserver/cmd/main/

RUN go-wrapper download
RUN go-wrapper install

RUN dep ensure -vendor-only

CMD ["go-wrapper", "run"]
