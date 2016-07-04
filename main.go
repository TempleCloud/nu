package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/v1/ping", ping)
	router.GET("/v1/echo", echo)
	router.HEAD("/v1/docker/*command", dockerHandler)
	router.OPTIONS("/v1/docker/*command", dockerHandler)
	router.GET("/v1/docker/*command", dockerHandler)
	router.PUT("/v1/docker/*command", dockerHandler)
	router.POST("/v1/docker/*command", dockerHandler)
	router.DELETE("/v1/docker/*command", dockerHandler)
	router.PATCH("/v1/docker/*command", dockerHandler)
	log.Fatal(http.ListenAndServe(":9090", router))
}

//--------------------------------------------------------------------------------------------------

func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("pong!"))
}

//--------------------------------------------------------------------------------------------------

func echo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	toEcho := r.URL.Query().Get("msg")
	w.Write([]byte(toEcho))
}

//--------------------------------------------------------------------------------------------------

func dockerHandler(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {

	// build docker proxy request client
	dockerClient := newDockerClient()

	// build docker proxy request
	dockerRq := newDockerRq(request, params)

	// invoke docker proxy request
	dockerRs := invoke(dockerClient, dockerRq)
	defer dockerRs.Body.Close()

	// process docker proxy response
	writeResponseBody(responseWriter, dockerRs)
}

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

func buildDockerURL(request *http.Request, params httprouter.Params) string {
	baseDockerURL := "http://localhost" + params.ByName("command")
	var dockerURL string
	if request.URL.RawQuery == "" {
		dockerURL = baseDockerURL
	} else {
		dockerURL = baseDockerURL + "?" + request.URL.RawQuery
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

func newDockerRq(request *http.Request, params httprouter.Params) *http.Request {
	dockerURL := buildDockerURL(request, params)
	dockerRq, err := http.NewRequest(request.Method, dockerURL, request.Body)
	if err != nil {
		log.Fatal("docker proxy request init error: ", err)
	}
	dockerRq.Header.Add("Content-Type", "application/json")
	copyHeader(dockerRq.Header, request.Header)
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
