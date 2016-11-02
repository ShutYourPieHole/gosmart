// Simple example for the gosmart libraries.
//
// This is a simple demonstration of how to obtain a token from the smartthings
// API using Oauth2 authorization, and how to request the status of some of your
// sensors (in this case, temperature).
//
// This file is part of gosmart, a set of libraries to communicate with
// the Samsumg SmartThings API using Go (golang).
//
// http://github.com/marcopaganini/gosmart
// (C) 2016 by Marco Paganini <paganini@paganini.net>

package main

import (
	"flag"
	"fmt"
	"github.com/marcopaganini/gosmart"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
)

const (
	defaultPort = 4567
)

var (
	flagClient = flag.String("client", "", "OAuth Client ID")
	flagSecret = flag.String("secret", "", "OAuth Secret")

	config oauth2.Config
)

func main() {
	flag.Parse()

	// No date on log messages
	log.SetFlags(0)

	// Attempt to load token from the local storage. If an error occurs
	// of the token is invalid (expired, etc), trigger the OAuth process.
	token, err := gosmart.LoadToken("")
	if err != nil || !token.Valid() {
		// We need client and secret to fetch a new token.
		fmt.Println(*flagClient, *flagSecret)
		if *flagClient == "" || *flagSecret == "" {
			log.Fatalf("Must specify Client ID (--client) and Secret (--secret)")
		}
		// Create new authentication config for our App
		config = oauth2.Config{
			ClientID:     *flagClient,
			ClientSecret: *flagSecret,
			Scopes:       []string{"app"},
			Endpoint:     gosmart.Endpoint,
		}

		gst, err := gosmart.NewAuth(defaultPort, config)
		if err != nil {
			log.Fatalf("Error creating GoSmart struct: %q\n", err)
		}

		fmt.Printf("Please login by visiting http://localhost:%d\n", defaultPort)
		token, err = gst.GetOAuthToken()
		if err != nil {
			log.Fatalf("Error generating token: %q\n", err)
		}

		// Save new token.
		err = gosmart.SaveToken("", token)
		if err != nil {
			log.Fatalf("Error saving token: %q\n", err)
		}
	}

	// Create a client with token
	ctx := context.Background()
	client := config.Client(ctx, token)

	// Retrieve Endpoints URI. All future accesses to the smartthings API
	// for this session should use this URL, followed by the desired URL path.
	endpoint, err := gosmart.GetEndPointsURI(client)
	if err != nil {
		log.Fatalf("Error fetching endpoints: %q\n", err)
		return
	}

	// Fetch /temperature
	resp, err := client.Get(endpoint + "/temperature")
	if err != nil {
		log.Fatalf("Error getting temperature %q\n", err)
		return
	}
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("Temperature content: %s\n", contents)
}
