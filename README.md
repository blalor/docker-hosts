# Simplified Docker container hostname resolution

`docker-hosts` (yes, a terrible name) maintains a file in the format of
`/etc/hosts` that contains IP addresses and hostnames of Docker containers. When
the generated file is mounted at `/etc/hosts` within your Docker container it
provides simple hostname resolution.  This allows you to set up `redis` and
`web` containers where the `web` container is able to connect to `redis` via its
hostname.  You can optionally provide a domain like `dev.docker`, so
`redis.dev.docker` is a usable alias, as well.

This utility was inspired by Michael Crosby's
[skydock](https://github.com/crosbymichael/skydock) project.  `docker-hosts` and
skydock (paired with skydns) work in much the same way: the container lifecycle
is monitored via the Docker daemon's events, and resolvable hostnames are made
available to appropriately-configured containers.  The end result is that you
have a simple way of connecting containers together on the same Docker host,
without having to resort to links or manual configuration.  This does *not*
provide a solution to container connectivity across Docker hosts.  For that you
should look at something like Jeff Lindsay's
[registrator](https://github.com/progrium/registrator).

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

Create an empty file at `/var/lib/docker/hosts`, make it mode `0644` and owned
by `nobody:nobody`.

    docker run \
        -d \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v /var/lib/docker/hosts:/srv/hosts \
        blalor/docker-hosts --domain-name=dev.docker /srv/hosts

[gpm]: https://github.com/pote/gpm
[gvp]: https://github.com/pote/gvp
