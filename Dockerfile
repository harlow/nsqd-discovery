FROM centurylink/ca-certs
MAINTAINER Harlow Ward "harlow@hward.com"
WORKDIR /app
COPY nsqd-discovery /app/
ENTRYPOINT ["./nsqd-discovery"]
