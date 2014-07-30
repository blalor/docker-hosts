FROM blalor/centos:latest
MAINTAINER Brian Lalor <blalor@bravo5.org>

ADD release/docker-hosts /usr/local/bin/

USER nobody
ENTRYPOINT ["/usr/local/bin/docker-hosts"]
