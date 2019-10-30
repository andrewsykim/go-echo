package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nu7hatch/gouuid"
)

func main() {

	// ENV
	tcpPort := os.Getenv("TCP_PORT")
	if tcpPort == "" {
		tcpPort = "1025"
	}

	nodeName := os.Getenv("NODE_NAME")
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	podIP := os.Getenv("POD_IP")
	serviceAccountName := os.Getenv("SERVICE_ACCOUNT")

	message := ""

	if nodeName != "" {
		message = message + fmt.Sprintf("Welcome, you are connected to node %s.\n", nodeName)
	}

	if podName != "" {
		message = message + fmt.Sprintf("Running on Pod %s.\n", podName)
	}

	if podNamespace != "" {
		message = message + fmt.Sprintf("In namespace %s.\n", podNamespace)
	}

	if podIP != "" {
		message = message + fmt.Sprintf("With IP address %s.\n", podIP)
	}

	if serviceAccountName != "" {
		message = message + fmt.Sprintf("Service %s.\n", serviceAccountName)
	}

	l, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Listening on TCP port", tcpPort)
	defer l.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		select {}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go handleTCPRequest(conn, message)
	}
}

func handleTCPRequest(conn net.Conn, message string) {
	callUuidV4, _ := uuid.NewV4()
	clientId := callUuidV4.String()

	log.Println(clientId + " - TCP connection open.")
	defer conn.Close()
	defer log.Println(clientId + " - TCP connection closed.")

	conn.Write([]byte(message))

	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			return
		}
		data := buf[:size]

		log.Println(clientId+" - Received Raw Data:", data)
		log.Printf(clientId+" - Received Data (converted to string): %s", data)

		data = []byte(string(data) + " " + fmt.Sprintf("%s\n", time.Now()))
		conn.Write(data)
	}
}
