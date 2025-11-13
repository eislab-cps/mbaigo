/*******************************************************************************
 * Copyright (c) 2023 Jan van Deventer
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-2.0/
 *
 * Contributors:
 *   Jan A. van Deventer, Luleå - initial implementation
 *   Thomas Hedeler, Hamburg - initial implementation
 ***************************************************************************SDG*/

package orchestrator

import (
	"context"
	"crypto/x509/pkix"
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/eislab-cps/mbaigo/pkg/components"
	"github.com/eislab-cps/mbaigo/pkg/forms"
	"github.com/eislab-cps/mbaigo/pkg/usecases"
)

// Start starts the Orchestrator system
func Start(ctx context.Context) error {
	// prepare for graceful shutdown
	localCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// instantiate the System
	sys := components.NewSystem("orchestrator", localCtx)

	// Instantiate the husk
	sys.Husk = &components.Husk{
		Description: "provides the URL of a currently available and authorized sought service",
		Certificate: "ABCD",
		Details:     map[string][]string{"Developer": {"Arrowhead"}},
		ProtoPort:   map[string]int{"https": 0, "http": 20103, "coap": 0},
		InfoLink:    "https://github.com/eislab-cps/mbaigo/tree/main/pkg/core/orchestrator",
		DName: pkix.Name{
			CommonName:         sys.Name,
			Organization:       []string{"Synecdoque"},
			OrganizationalUnit: []string{"Systems"},
			Locality:           []string{"Luleå"},
			Province:           []string{"Norrbotten"},
			Country:            []string{"SE"},
		},
	}

	// instantiate a template unit asset
	assetTemplate := initTemplate()
	assetName := assetTemplate.GetName()
	sys.UAssets[assetName] = &assetTemplate

	// Configure the system - use dedicated config directory for core systems
	configPath := "systems/orchestrator/systemconfig.json"
	// Ensure directory exists
	if err := os.MkdirAll("systems/orchestrator", 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v\n", err)
	}
	rawResources, err := usecases.ConfigureWithPath(&sys, configPath)
	if err != nil {
		log.Fatalf("Configuration error: %v\n", err)
	}
	sys.UAssets = make(map[string]*components.UnitAsset) // clear the unit asset map (from the template)
	for _, raw := range rawResources {
		var uac usecases.ConfigurableAsset
		if err := json.Unmarshal(raw, &uac); err != nil {
			log.Fatalf("Resource configuration error: %+v\n", err)
		}
		ua, cleanup := newResource(uac, &sys)
		defer cleanup()
		sys.UAssets[ua.GetName()] = &ua
	}

	// Note: Certificate authentication disabled for basic tutorial
	// Uncomment the following line if you have a CA system running:
	// usecases.RequestCertificate(&sys)

	// Note: Service registration disabled for basic tutorial (orchestrator doesn't need to register itself)
	// Uncomment the following line if you want orchestrator to register its services:
	// usecases.RegisterServices(&sys)

	// start the http handler and server
	go usecases.SetoutServers(&sys)

	// wait for shutdown signal from context or SIGINT
	select {
	case <-localCtx.Done():
		log.Println("shutting down system", sys.Name)
	case <-sys.Sigs:
		log.Println("shutting down system", sys.Name)
		cancel()
	}

	// allow the go routines to be executed, which might take more time than the main routine to end
	time.Sleep(2 * time.Second)
	return nil
}

// Serving handles the resources services. NOTE: it expects those names from the request URL path
func (ua *UnitAsset) Serving(w http.ResponseWriter, r *http.Request, servicePath string) {
	switch servicePath {
	case "squest":
		ua.orchestrate(w, r)
	case "squests":
		ua.orchestrateMultiple(w, r)
	default:
		http.Error(w, "Invalid service request [Do not modify the services subpath in the configuration file]", http.StatusBadRequest)
	}
}

// orchestrate receives a service discovery request and responds with the selected service location if found
func (ua *UnitAsset) orchestrate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		contentType := r.Header.Get("Content-Type")
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			log.Println("Error parsing media type:", err)
			return
		}

		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading discovery request body: %v\n", err)
			return
		}

		questForm, err := usecases.Unpack(bodyBytes, mediaType)
		if err != nil {
			log.Printf("error extracting the discovery request %v\n", err)
		}
		qf, ok := questForm.(*forms.ServiceQuest_v1)
		if !ok {
			log.Println("Problem unpacking the service discovery request form")
			return
		}

		servLocation, err := ua.getServiceURL(*qf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(servLocation) // respond with the selected service location
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method is not supported.", http.StatusNotFound)
	}
}

func (ua *UnitAsset) orchestrateMultiple(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		contentType := r.Header.Get("Content-Type")
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			log.Println("Error parsing media type:", err)
			return
		}

		defer r.Body.Close()
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading discovery request body: %v\n", err)
			return
		}

		questForm, err := usecases.Unpack(bodyBytes, mediaType)
		if err != nil {
			log.Printf("error extracting the discovery request %v\n", err)
		}
		qf, ok := questForm.(*forms.ServiceQuest_v1)
		if !ok {
			log.Println("Problem unpacking the service discovery request form")
			return
		}

		servLocation, err := ua.getServicesURL(*qf)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(servLocation) // respond with the selected service location
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method is not supported.", http.StatusNotFound)
	}
}
