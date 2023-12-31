package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// declaring a base url constant value for API client
const apiBaseUrl = "http://localhost:8080/kv"

func list() error {
	// creating a http GET request client call
	resp, err := http.Get(apiBaseUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d %s", resp.StatusCode, resp.Status)
	}
	var keys []string
	// Decoding the response body and updating key list
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		log.Printf("error receiving: %s", err)
	}

	// loop over all keys in the key list and print
	for _, key := range keys {
		fmt.Println(key)
	}
	return nil
}

func set(key string) error {
	// creating a http POST request client call
	resp, err := http.Post(fmt.Sprintf("%s/%s", apiBaseUrl, key), "application/octet-stream", os.Stdin)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Key  string
		Size int
	}
	// Decoding the response body and updating data struct
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("error receiving: %s", err)
	}
	// Returning the values from data struct
	fmt.Printf("%s: %d bytes\n", data.Size, data.Key)

	return nil
}

func get(key string) error {
	// creating a http GET request client call
	resp, err := http.Get(fmt.Sprintf("%s/%s", apiBaseUrl, key))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Returning the values from response directly to standard output i.e. terminal window
	_, err = io.Copy(os.Stdout, resp.Body)
	return err
}

func delete(key string) error {
	// create a new HTTP client
	client := &http.Client{}

	// create a new DELETE request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", apiBaseUrl, key), nil)
	if err != nil {
		return err
	}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var reply struct {
		Data []byte
		Size int
	}
	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		log.Printf("error receiving: %s", err)
	}
	fmt.Printf("deleted %s\n value: %s\nfreed memory %d bytes\n", key, reply.Data, reply.Size)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: kv get|set|delete|list [key]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("error: wrong number of arguments")
	}

	switch flag.Arg(0) {
	case "get":
		key := flag.Arg(1)
		if key == "" {
			log.Fatalf("error: missing key")
		}
		if err := get(key); err != nil {
			log.Fatal(err)
		}
	case "set":
		key := flag.Arg(1)
		if key == "" {
			log.Fatalf("error: missing key")
		}
		if err := set(key); err != nil {
			log.Fatal(err)
		}
	case "delete":
		key := flag.Arg(1)
		if key == "" {
			log.Fatalf("error: missing key")
		}
		if err := delete(key); err != nil {
			log.Fatal(err)
		}
	case "list":
		if err := list(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("error: unknown command: %q", flag.Arg(0))
	}
}
