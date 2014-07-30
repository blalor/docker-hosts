
## building

This project uses [gpm][gpm] and [gvp][gvp].  Both must be available on your
path.

    make

-- or --

    gvp init
    source gvp in
    gpm install
    go build -v -o stage/docker-host ./...

## running

Start the `docker-host` process and give it the path to a file that will be
mounted as `/etc/hosts` in your containers:

    docker-host /path/to/hosts

Optionally specify `DOCKER_HOST` environment variable.

Then start a container:

    docker run -i -t -v /path/to/hosts:/etc/hosts:ro centos /bin/bash

Within the `centos` container, you'll see `/etc/hosts` has an entry for the
container you just started, as well as any other containers already running.
`/etc/hosts` will continue to reflect all of the containers currently running on
this Docker host.

The **only** container that should have write access to the generated hosts file
is the container running this application.

## running in Docker

Create an empty file at `/var/lib/docker/hosts` and make it mode `0644` and
owned by `nobody:nobody`.

    docker run \
        -d \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v /var/lib/docker/hosts:/srv/hosts \
        blalor/docker-hosts --domain-name=dev.docker /srv/hosts

[gpm]: https://github.com/pote/gpm
[gvp]: https://github.com/pote/gvp
