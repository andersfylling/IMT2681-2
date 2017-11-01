# Not configured for heroku
# Docker only
FROM golang:1.8
MAINTAINER https://github.com/andersfylling

WORKDIR /go/src/github.com/andersfylling/IMT2681-2
COPY . .

# Get Glide for package management
RUN curl https://glide.sh/get | sh
RUN glide install

# The URL needs to be replaced. Please see extra details given with this submission.
ENV IMT_DATABASE_MONGODB_URL mongodb://<dbuser>:<dbpassword>@ds045099.mlab.com:45099/core
ENV IMT_DATABASE_MONGODB_DATABASE core
