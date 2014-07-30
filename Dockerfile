FROM blalor/centos:latest
MAINTAINER Brian Lalor <blalor@bravo5.org>

ADD release/docker-hosts /usr/local/bin/

## should *not* run as root, but needs access to /var/run/docker.sock, which
## should *not* be accessible by nobody. *sigh*
#USER nobody
ENTRYPOINT ["/usr/local/bin/docker-hosts"]
