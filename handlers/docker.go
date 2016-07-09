package handlers

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

//--------------------------------------------------------------------------------------------------

// DockerProxy forwards the request onto the configured 'docker daemon' and write back the result
func DockerProxy(c *gin.Context) {

	// build docker daemon proxy request client
	dockerClient := newDockerClient()

	// build docker proxy request
	dockerRq := newDockerRequest(c)

	// invoke docker proxy request
	dockerRs := invoke(dockerClient, dockerRq)
	dockerResponseBody, _ := ioutil.ReadAll(dockerRs.Body)
	defer dockerRs.Body.Close()

	c.String(http.StatusOK, string(dockerResponseBody))
}

//--------------------------------------------------------------------------------------------------

func newDockerClient() *http.Client {
	dockerTransport := &http.Transport{
		Dial: newDockerDial,
	}
	dockerClient := &http.Client{Transport: dockerTransport}
	return dockerClient
}

func newDockerDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func newDockerRequest(c *gin.Context) *http.Request {

	dockerURL := buildDockerURL(c.Param("command"), c.Request.URL.RawQuery)
	dockerRq, err := http.NewRequest(c.Request.Method, dockerURL, c.Request.Body)
	if err != nil {
		log.Fatal("docker proxy request init error: ", err)
	}
	dockerRq.Header.Add("Content-Type", "application/json")
	copyHeader(dockerRq.Header, c.Request.Header)

	return dockerRq
}

func buildDockerURL(dockerCmdURL string, dockerCmdQueryParams string) string {

	baseDockerURL := "http://localhost" + dockerCmdURL

	var dockerURL string
	if dockerCmdQueryParams == "" {
		dockerURL = baseDockerURL
	} else {
		dockerURL = baseDockerURL + "?" + dockerCmdQueryParams
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

func invoke(client *http.Client, request *http.Request) *http.Response {

	response, err := client.Do(request)

	if err != nil {
		log.Fatal("request invocation error: ", err)
	}
	return response
}
