package main

import (
	"flag"
	dockerapi "github.com/fsouza/go-dockerclient"
	"log"
	"os"
)

func getopt(name, def string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return def
}

func assert(err error) {
	if err != nil {
		log.Fatal("docker-hosts: ", err)
	}
}

func main() {
	flag.Parse()

	docker, err := dockerapi.NewClient(getopt("DOCKER_HOST", "unix:///var/run/docker.sock"))
	assert(err)

	hosts := NewHosts(docker, flag.Arg(0))

	// set up to handle events early, so we don't miss anything while doing the
	// initial population
	events := make(chan *dockerapi.APIEvents)
	assert(docker.AddEventListener(events))

	containers, err := docker.ListContainers(dockerapi.ListContainersOptions{})
	assert(err)

	for _, listing := range containers {
		go hosts.Add(listing.ID)
	}

	log.Println("docker-hosts: Listening for Docker events...")
	for msg := range events {
		switch msg.Status {
		case "start":
			go hosts.Add(msg.ID)

		case "die":
			go hosts.Remove(msg.ID)
		}
	}

	log.Fatal("docker-hosts: docker event loop closed") // todo: reconnect?

}
