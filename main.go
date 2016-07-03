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

func localDockerDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func localDocker(w http.ResponseWriter, r *http.Request) {

	log.Printf("Invoking!")

	log.Printf("Creating socket!")

	dockerTransport := &http.Transport{
		Dial: localDockerDial,
	}
	client := &http.Client{Transport: dockerTransport}

	log.Printf("Request: %v", r)
	log.Printf("URL: %v", r.URL)
	log.Printf("RawPath: %v", r.URL.RawPath)
	log.Printf("Path: %v", r.URL.Path)

	dockerURL := "http://localhost" + strings.TrimPrefix(r.URL.Path, "/v1/docker")
	log.Printf("dockerUrl: %v", dockerURL)
	req, err := http.NewRequest(r.Method, dockerURL, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("get request error:", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("body read error:", err)
	}

	res := string(string(body))
	log.Printf("Client got: %v", res)

	w.Write([]byte(body))

}
