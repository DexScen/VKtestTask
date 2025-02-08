package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"time"
)

//Я не нашел как правильно с заданной архитектурой переиспользовать пакеты из backend
//Поэтому тут повторяется Container и ListContainer

type Container struct {
	IP          string    `json:"ip"`
	PingTime    time.Time `json:"pingtime"`
	SuccessDate time.Time `json:"successdate"`
}

type ListContainer []Container

// Wait for backend to start
func waitForBackend() {
	url := "http://backend:8080/containers"
	timeout := 30 * time.Second
	start := time.Now()

	for time.Since(start) < timeout {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		}
		time.Sleep(2 * time.Second)
	}
}

// Get ips of running containers
func getContainerIPs() ([]string, error) {
	cmd := exec.Command("docker", "ps", "-q")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var ips []string
	containers := bytes.Split(out.Bytes(), []byte("\n"))

	for _, container := range containers {
		if len(container) == 0 {
			continue
		}

		cmd := exec.Command("docker", "inspect", "-f", "{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}", string(container))
		var ipOut bytes.Buffer
		cmd.Stdout = &ipOut
		err := cmd.Run()
		if err != nil {
			continue
		}
		ip := ipOut.String()
		ips = append(ips, ip)
	}

	return ips, nil
}

// Initial containers info collection
func init() {
	ips, err := getContainerIPs()
	if err != nil {
		log.Println("Couldn't get container ips")
		return
	}
	var list ListContainer
	for _, ip := range ips {
		c := Container{
			IP:          ip[:len(ip)-1],
			PingTime:    time.Now(),
			SuccessDate: time.Now(),
		}
		list = append(list, c)
	}

	waitForBackend()

	url := "http://backend:8080/containers"
	jsonData, err := json.Marshal(list)
	if err != nil {
		log.Printf("json marshalling err in init containers: %v\n", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("init containers db request err: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("backend вернул статус: %s\n", resp.Status)
		return
	}

	log.Println("initial containers sent to DB")
}

// Ping containers
func isReachable(ip string, port string) bool {
	address := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// Get containers via GET
func getContainers() ([]Container, error) {
	resp, err := http.Get("http://backend:8080/containers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var containers []Container
	err = json.NewDecoder(resp.Body).Decode(&containers)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// Send containers via post
func updateContainers(containers []Container) {
	jsonData, _ := json.Marshal(containers)

	resp, err := http.Post("http://backend:8080/containers", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("error updating containers:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("updated containers sent back success")
}

func main() {
	log.Println("Pinger started")

	//костыль чтобы не переписывать все приложение =) он никогда не обрабатывает запросы, просто так я могу его пингануть
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println(w, "pinger pinged")
		})
		http.ListenAndServe(":8081", nil)
	}()

	for {
		log.Println("pinger: getting containers")
		containers, err := getContainers()
		if err != nil {
			log.Println("pinger: err while getting containers:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for i := range containers {
			timeSync := time.Now()
			containers[i].PingTime = timeSync
			if isReachable(containers[i].IP, "8080") || isReachable(containers[i].IP, "5432") || isReachable(containers[i].IP, "8081") {
				log.Println("pinger: container available with ip:", containers[i].IP)
				containers[i].SuccessDate = timeSync
			} else {
				log.Println("pinger: container unavailable with ip:", containers[i].IP)
			}
		}

		updateContainers(containers)

		time.Sleep(5 * time.Second)
	}

}
