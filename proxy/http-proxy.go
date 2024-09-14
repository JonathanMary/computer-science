package main

import (
	"log"
	"syscall"
)

func main() {
	proxyAddr := &syscall.SockaddrInet4{
		Addr: [4]byte{0, 0, 0, 0},
		Port: 8000,
	}
	serverAddr := &syscall.SockaddrInet4{
		Addr: [4]byte{127, 0, 0, 1},
		Port: 9000,
	}

	/* Connect to proxy */
	proxyfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		log.Fatal("Error creating proxy socket: ", err)
	}
	defer syscall.Close(proxyfd)
	syscall.SetsockoptInt(proxyfd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 0x1) // Allow for proxy to restart without blocking port

	if err := syscall.Bind(proxyfd, proxyAddr); err != nil {
		log.Fatal("Errror binding socket: ", err)
	}
	if syscall.Listen(proxyfd, syscall.SOMAXCONN); err != nil {
		log.Fatal("Error listening to port: ", err)
	}
	log.Print("Listening to port: ", proxyAddr.Port)

	for {
		/* Connect to client */
		clientfd, clientAddr, err := syscall.Accept(proxyfd)
		if err != nil {
			log.Fatal("Error accepting client connection: ", err)
		}
		defer syscall.Close(clientfd)
		log.Print("New connection from: ", clientAddr)

		buffer := make([]byte, 4096)
		nclient, _, err := syscall.Recvfrom(clientfd, buffer, 0x0)
		if err != nil {
			log.Fatal("Error receiving from client: ", err)
		}
		if nclient == 0 {
			continue
		}

		/* Connect to server */
		serverfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
		if err != nil {
			log.Fatal("Error creating server socket: ", err)
		}
		defer syscall.Close(serverfd)

		if err := syscall.Connect(serverfd, serverAddr); err != nil {
			/* Handle disconnected server without breaking connection */
			syscall.Sendto(clientfd, []byte("HTTP/1.0 502 Bad Gateway\r\n\r\n"), 0x0, clientAddr)
			syscall.Close(clientfd)
			log.Print("Bad Gateway: ", err)
			continue
		}
		log.Print("Connected to: ", serverAddr.Port)

		if err := syscall.Sendto(serverfd, buffer[:nclient], 0x0, serverAddr); err != nil {
			log.Fatal("Error sending to server: ", err)
		}

		for {
			nserver, _, err := syscall.Recvfrom(serverfd, buffer, 0x0)
			if err != nil {
				log.Fatal("Error receiving from server: ", err)
			}
			if nserver == 0 {
				break
			}

			syscall.Sendto(clientfd, buffer[:nserver], 0x0, clientAddr)
		}
	}
}
