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

func localDockerDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func localDocker(w http.ResponseWriter, r *http.Request) {

	// process nu request parameters
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("get request body error:", err)
	}
	payload := string(body)
	queryParams := r.URL.RawQuery
	dockerURL := "http://localhost" + strings.TrimPrefix(r.URL.Path, "/v1/docker") + "?" + queryParams

	log.Printf("dockerURL: %v", dockerURL)

	// build docker proxy request client
	dockerTransport := &http.Transport{
		Dial: localDockerDial,
	}
	client := &http.Client{Transport: dockerTransport}

	// build docker proxy request
	dockerRq, err := http.NewRequest(r.Method, dockerURL, strings.NewReader(payload))
	dockerRq.Header.Add("Content-Type", "application/json")

	// invoke docker proxy request
	dockerRs, err := client.Do(dockerRq)
	if err != nil {
		log.Fatal("get request error:", err)
	}
	defer dockerRs.Body.Close()

	// process docker proxy response
	dockerRsBody, err := ioutil.ReadAll(dockerRs.Body)
	if err != nil {
		log.Fatal("body read error:", err)
	}

	log.Printf("docker proxy response: %v", string(dockerRsBody))
	w.Write(dockerRsBody)

}
