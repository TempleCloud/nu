package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/v1/ping", ping)
	http.HandleFunc("/v1/echo", echo)
	http.HandleFunc("/v1/docker/info", localDocker)
	http.HandleFunc("/v1/docker/images/json", localDocker)
	http.HandleFunc("/v1/docker/containers/json", localDocker)
	http.HandleFunc("/v1/docker/containers/create", localDocker)
	http.HandleFunc("/v1/docker/images/alpine/json", localDocker)
	http.ListenAndServe(":9090", nil)
}

//--------------------------------------------------------------------------------------------------

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong!"))
}

//--------------------------------------------------------------------------------------------------

func echo(w http.ResponseWriter, r *http.Request) {
	toEcho := r.URL.Query().Get("msg")
	w.Write([]byte(toEcho))
}

//--------------------------------------------------------------------------------------------------

func newDockerDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func newDockerClient() *http.Client {
	dockerTransport := &http.Transport{
		Dial: newDockerDial,
	}
	dockerClient := &http.Client{Transport: dockerTransport}
	return dockerClient
}

func buildDockerURL(r *http.Request) string {
	baseDockerURL := "http://localhost" + strings.TrimPrefix(r.URL.Path, "/v1/docker")
	var dockerURL string
	if r.URL.RawQuery == "" {
		dockerURL = baseDockerURL
	} else {
		dockerURL = baseDockerURL + "?" + r.URL.RawQuery
	}
	return dockerURL
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func newDockerRq(r *http.Request) *http.Request {
	dockerURL := buildDockerURL(r)
	dockerRq, err := http.NewRequest(r.Method, dockerURL, r.Body)
	if err != nil {
		log.Fatal("docker proxy request init error: ", err)
	}
	dockerRq.Header.Add("Content-Type", "application/json")
	copyHeader(dockerRq.Header, r.Header)
	return dockerRq
}

func invoke(client *http.Client, request *http.Request) *http.Response {
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("request invocation error: ", err)
	}
	return response
}

func writeResponseBody(responseWriter http.ResponseWriter, response *http.Response) {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("response read error: ", err)
	}
	responseWriter.Write(responseBody)
}

func localDocker(responseWriter http.ResponseWriter, request *http.Request) {
	// build docker proxy request client
	dockerClient := newDockerClient()

	// build docker proxy request
	dockerRq := newDockerRq(request)

	// invoke docker proxy request
	dockerRs := invoke(dockerClient, dockerRq)
	defer dockerRs.Body.Close()

	// process docker proxy response
	writeResponseBody(responseWriter, dockerRs)
}
