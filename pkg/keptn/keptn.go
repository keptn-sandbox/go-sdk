package keptn

import (
	"fmt"
	"os"
	"strings"
)

type KeptnOpts struct {
	UseLocalFileSystem      bool
	ConfigurationServiceURL string
	EventBrokerURL          string
}

type Keptn struct {
	Project      string
	Stage        string
	Service      string
	KeptnContext string

	eventBrokerURL     string
	useLocalFileSystem bool
	resourceHandler    *ResourceHandler
}

const configurationServiceURL = "configuration-service:8080"
const defaultEventBrokerURL = "http://event-broker.keptn.svc.cluster.local/keptn"

func NewKeptn(project string, stage string, service string, keptnContext string, opts KeptnOpts) *Keptn {
	k := &Keptn{
		Project:            project,
		Stage:              stage,
		Service:            service,
		KeptnContext:       keptnContext,
		useLocalFileSystem: opts.UseLocalFileSystem,
		resourceHandler:    nil,
	}
	csURL := configurationServiceURL
	if opts.ConfigurationServiceURL != "" {
		csURL = opts.ConfigurationServiceURL
	}

	if opts.EventBrokerURL != "" {
		k.eventBrokerURL = opts.EventBrokerURL
	} else {
		k.eventBrokerURL = defaultEventBrokerURL
	}

	k.resourceHandler = NewResourceHandler(csURL)

	return k
}

func (k *Keptn) GetKeptnResource(resource string) (string, error) {

	// if we run in a runlocal mode we are just getting the file from the local disk
	if k.useLocalFileSystem {
		return _getKeptnResourceFromLocal(resource)
	}

	// get it from Keptn
	requestedResource, err := k.resourceHandler.GetServiceResource(k.Project, k.Stage, k.Service, resource)

	// return Nil in case resource couldnt be retrieved
	if err != nil || requestedResource.ResourceContent == "" {
		fmt.Printf("Keptn Resource not found: %s - %s", resource, err)
		return "", err
	}

	// now store that file on the same directory structure locally
	os.RemoveAll(resource)
	pathArr := strings.Split(resource, "/")
	directory := ""
	for _, pathItem := range pathArr[0 : len(pathArr)-1] {
		directory += pathItem + "/"
	}

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", err
	}
	resourceFile, err := os.Create(resource)
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}
	defer resourceFile.Close()

	_, err = resourceFile.Write([]byte(requestedResource.ResourceContent))

	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}

	return resource, nil
}

/**
 * Retrieves a resource (=file) from the local file system. Basically checks if the file is available and if so returns it
 */
func _getKeptnResourceFromLocal(resource string) (string, error) {
	if _, err := os.Stat(resource); err == nil {
		return resource, nil
	} else {
		return "", err
	}
}
