package main

import (
	"io/ioutil"
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

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func dockerImages(w http.ResponseWriter, r *http.Request) {

	log.Printf("Invoking!")

	log.Printf("Creating socket!")
	tr := &http.Transport{
		Dial: fakeDial,
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get("http://localhost/images/json")
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
