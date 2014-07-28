
## building

I'm using [gpm][gpm] and [gvp][gvp].

    gvp init
    source gvp in
    gpm install
    go build ./...

-- or --

    gvp init
    gvp in gpm install
    gvp in go build ./...

## running

Start the `docker-host` process and give it the path to a file that will be
mounted as `/etc/hosts` in your containers:

    docker-host /path/to/hosts

Optionally specify `DOCKER_HOST` environment variable.

Then start a container:

    docker run -i -t -v /path/to/hosts:/etc/hosts centos /bin/bash

Within the `centos` container, you'll see `/etc/hosts` has an entry for the
container you just started, as well as any other containers already running.
`/etc/hosts` will continue to reflect all of the containers currently running on
this Docker host.

[gpm]: https://github.com/pote/gpm
[gvp]: https://github.com/pote/gvp
