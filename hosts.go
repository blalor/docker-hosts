package main

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type HostEntry struct {
	IPAddress         string
	CanonicalHostname string
	Aliases           []string
}

type Hosts struct {
	sync.Mutex
	docker  *dockerapi.Client
	path    string
	entries map[string]HostEntry
}

func NewHosts(docker *dockerapi.Client, path string) *Hosts {
	hosts := &Hosts{
		docker: docker,
		path:   path,
	}

	hosts.entries = make(map[string]HostEntry)

	// combination of docker, centos
	hosts.entries["__localhost4"] = HostEntry{
		IPAddress:         "127.0.0.1",
		CanonicalHostname: "localhost",
		Aliases:           []string{"localhost4"},
	}

	hosts.entries["__localhost6"] = HostEntry{
		IPAddress:         "::1",
		CanonicalHostname: "localhost",
		Aliases:           []string{"localhost6", "ip6-localhost", "ip6-loopback"},
	}

	// docker puts these in
	hosts.entries["fe00::0"] = HostEntry{
		IPAddress:         "fe00::0",
		CanonicalHostname: "ip6-localnet",
	}

	hosts.entries["ff00::0"] = HostEntry{
		IPAddress:         "ff00::0",
		CanonicalHostname: "ip6-mcastprefix",
	}

	hosts.entries["ff02::1"] = HostEntry{
		IPAddress:         "ff02::1",
		CanonicalHostname: "ip6-allnodes",
	}

	hosts.entries["ff02::2"] = HostEntry{
		IPAddress:         "ff02::2",
		CanonicalHostname: "ip6-allrouters",
	}

	return hosts
}

func (h *Hosts) WriteFile() error {
	tempFile, err := ioutil.TempFile(os.TempDir(), "hosts")

	if err != nil {
		return err
	}

	for _, entry := range h.entries {
		tempFile.WriteString(strings.Join(
			append(
				[]string{entry.IPAddress, entry.CanonicalHostname},
				entry.Aliases...,
			),
			"\t",
		) + "\n")
	}

	tempFile.Close() // can't close? ignore!

	return os.Rename(tempFile.Name(), h.path)
}

func (h *Hosts) Add(containerId string) {
	h.Lock()
	defer h.Unlock()

	container, err := h.docker.InspectContainer(containerId)
	if err != nil {
		log.Println("unable to inspect container:", containerId, err)
		return
	}

	h.entries[containerId] = HostEntry{
		IPAddress:         container.NetworkSettings.IPAddress,
		CanonicalHostname: container.Config.Hostname,
		// Aliases:           []string{container.Name[1:]}, // could contain "_"
	}

	err = h.WriteFile()
	if err != nil {
		log.Println("unable to write file", err)
	}
}

func (h *Hosts) Remove(containerId string) {
	h.Lock()
	defer h.Unlock()

	delete(h.entries, containerId)

	err := h.WriteFile()
	if err != nil {
		log.Println("unable to write file", err)
	}
}
