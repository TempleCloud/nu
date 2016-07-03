package main

import (
	"log"
	"net"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/ping", ping)
	http.HandleFunc("/v1/echo", echo)
	http.HandleFunc("/v1/docker/images", dockerImages)
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

func dockerImages(w http.ResponseWriter, r *http.Request) {

	log.Printf("Invoking!")

	log.Printf("Creating socket!")
	dockerSock, err := net.Dial("unix", "/var/run/docker.sock")
	if err != nil {
		panic(err)
	}
	defer dockerSock.Close()

	log.Printf("Writing to socket")
	_, err = dockerSock.Write([]byte("GET /images/json HTTP/1.1\n"))
	_, err = dockerSock.Write([]byte("Host: http\n"))
	_, err = dockerSock.Write([]byte("User-Agent: golang-client\n"))
	_, err = dockerSock.Write([]byte("Accept: */*\n"))
	_, err = dockerSock.Write([]byte("\n"))

	if err != nil {
		log.Fatal("write error:", err)
	}

	buf := make([]byte, 2064)
	n, err := dockerSock.Read(buf[:])
	if err != nil {
		log.Fatal("read error:", err)
	}

	res := string(buf[0:n])
	log.Printf("Docker daemon response: %v", res)

	w.Write([]byte(res))
}
