FROM golang

# Meta
LABEL author="Brendan Ashby <brendan@brendanashby.com>"
LABEL version="0.1.0"

# Build arguments
ARG repo_private_key

# Set ENV VARs
ENV MODE="1"
ENV SERVE="true"
ENV GIN_IMMEDIATE="true"

# Force SSH auth when getting deps for GoLang
RUN echo "[url \"git@bitbucket.org:\"]\n\tinsteadOf = https://bitbucket.org/" >> /root/.gitconfig

# Ignore SSH checking for getting deps via go get from GoLang
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

# Bundle private key to access repos
RUN echo "$repo_private_key" > /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa

# Pull in dependency checker
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Get Dependencies
RUN go get github.com/codegangsta/gin

# Bundle source
COPY . /go/src/app

# Build
WORKDIR /go/src/app/server/
RUN go-wrapper download
RUN go-wrapper install

# Check dependencies
RUN dep ensure -vendor-only

# Start
CMD ["gin", "run"]
