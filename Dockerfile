FROM golang

# Meta
LABEL author="Brendan Ashby <brendan@brendanashby.com>"
LABEL version="0.1.0"

# Force SSH auth when getting deps for GoLang
RUN echo "[url \"git@bitbucket.org:\"]\n\tinsteadOf = https://bitbucket.org/" >> /root/.gitconfig

# Ignore SSH checking for getting deps via go get from GoLang
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/server/

RUN ls -al
RUN cd /root/.ssh && ls -al
RUN go-wrapper download
RUN go-wrapper install

RUN dep ensure -vendor-only

CMD ["go-wrapper", "run"]
