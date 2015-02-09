package main

import (
    dockerapi "github.com/fsouza/go-dockerclient"
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
    domain  string
    entries map[string]HostEntry
}

func NewHosts(docker *dockerapi.Client, path, domain string) *Hosts {
    hosts := &Hosts{
        docker: docker,
        path:   path,
        domain: domain,
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
    hosts.entries["fe00::0"] = HostEntry{"fe00::0", "ip6-localnet", nil}
    hosts.entries["ff00::0"] = HostEntry{"ff00::0", "ip6-mcastprefix", nil}
    hosts.entries["ff02::1"] = HostEntry{"ff02::1", "ip6-allnodes", nil}
    hosts.entries["ff02::2"] = HostEntry{"ff02::2", "ip6-allrouters", nil}

    return hosts
}

func (h *Hosts) WriteFile() {
    file, err := os.Create(h.path)

    if err != nil {
        log.Println("unable to write to", h.path, err)
        return
    }

    defer file.Close()

    for _, entry := range h.entries {
        // <ip>\t<canonical>\t<alias1>\t…\t<aliasN>\n
        file.WriteString(strings.Join(
            append(
                []string{entry.IPAddress, entry.CanonicalHostname},
                entry.Aliases...,
            ),
            "\t",
        ) + "\n")
    }
}

func (h *Hosts) Add(containerId string) {
    h.Lock()
    defer h.Unlock()

    container, err := h.docker.InspectContainer(containerId)
    if err != nil {
        log.Println("unable to inspect container:", containerId, err)
        return
    }

    entry := HostEntry{
        IPAddress:         container.NetworkSettings.IPAddress,
        CanonicalHostname: container.Config.Hostname,
        Aliases:           []string{
            strings.Replace(container.Name[1:], "_", "-"), // could contain "_"
        },
    }

    if h.domain != "" {
        entry.Aliases =
            append(h.entries[containerId].Aliases, container.Config.Hostname+"."+h.domain)
    }

    h.entries[containerId] = entry

    h.WriteFile()
}

func (h *Hosts) Remove(containerId string) {
    h.Lock()
    defer h.Unlock()

    delete(h.entries, containerId)

    h.WriteFile()
}
