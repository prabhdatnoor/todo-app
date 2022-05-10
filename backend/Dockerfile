##Dockerfile I would like to use later, because it has latest version of node.
##This will be the dockerfile for the backend, won't be used yet though

FROM golang:latest


WORKDIR /code

#COPY go.mod /code/
COPY go.mod /code/
# RUN \
#     apt-get update && \


# default to port 80 for node, and 9229 and 9230 (tests) for debug
ARG PORT=5001
ENV PORT $PORT
EXPOSE $PORT 9229 9230


# check every 30s to ensure this service returns HTTP 200
# HEALTHCHECK --interval=30s \
#     CMD node healthcheck.js

# copy in our source code last, as it changes the most
#COPY . /code

# if you want to use npm start instead, then use `docker run --init in production`
# so that signals are passed properly. Note the code in index.js is needed to catch Docker signals
# using node here is still more graceful stopping then npm with --init afaik
# I still can't come up with a good production way to run with npm and graceful shutdown
CMD ["go", "build", "main.go", "models/models.go"]
