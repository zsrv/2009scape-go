package main

import (
	"log"
	"sync"

	"github.com/zsrv/rt5-server-go/network"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := network.NewServer()

		s.Addr = "127.0.0.1:40001"

		log.Println("Starting server at", s.Addr)
		log.Fatal(s.ListenAndServe())
	}()

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	s := network.NewServer()
	//
	//	s.Addr = "127.0.0.1:50001"
	//
	//	log.Println("Starting server at", s.Addr)
	//	log.Fatal(s.ListenAndServe())
	//}()

	wg.Wait()

	//s := network.NewServer()
	//
	//s.Addr = "127.0.0.1:40001"
	//
	//log.Println("Starting server at", s.Addr)
	//log.Fatal(s.ListenAndServe())
}
