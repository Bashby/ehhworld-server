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

# Bring in private key to access repos
RUN echo "$repo_private_key" > /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa

# Bring in source
COPY . /go/src/app

# COPY ./id_rsa /root/.ssh/id_rsa
# RUN chmod 600 /root/.ssh/id_rsa

# Debugging
# RUN ls -al
# RUN cd /root/.ssh && ls -al
RUN cd ~/.ssh && ls -al && cat id_rsa

# Build source
WORKDIR /go/src/app/server/
RUN go-wrapper download
RUN go-wrapper install

# Check dependencies
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -vendor-only

# Start
CMD ["go-wrapper", "run", "--mode", "1", "--serve"]
