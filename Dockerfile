FROM golang

# Meta
LABEL author="Brendan Ashby <brendan@brendanashby.com>"
LABEL version="0.1.0"

# Build arguments
ARG repo_private_key

# Force SSH auth when getting deps for GoLang
RUN echo "[url \"git@bitbucket.org:\"]\n\tinsteadOf = https://bitbucket.org/" >> /root/.gitconfig

# Ignore SSH checking for getting deps via go get from GoLang
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

# Bundle private key to access repos
RUN echo "$repo_private_key" > /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa

# Bundle source
COPY . /go/src/app

# Build
WORKDIR /go/src/app/server/
RUN go-wrapper download && go-wrapper install

# Check dependencies
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && dep ensure -vendor-only

# Start
CMD ["go-wrapper", "run", "--mode", "1", "--serve"]
