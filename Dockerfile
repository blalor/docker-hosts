FROM blalor/centos:latest
MAINTAINER Brian Lalor <blalor@bravo5.org>

ADD release/docker-hosts /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/docker-hosts"]
