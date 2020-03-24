package keptn

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewKeptn(t *testing.T) {
	type args struct {
		project      string
		stage        string
		service      string
		keptnContext string
		opts         KeptnOpts
	}
	tests := []struct {
		name string
		args args
		want *Keptn
	}{
		{
			name: "Get 'in-cluster' Keptn",
			args: args{
				project:      "sockshop",
				stage:        "dev",
				service:      "carts",
				keptnContext: "test-context",
				opts:         KeptnOpts{},
			},
			want: &Keptn{
				Project:            "sockshop",
				Stage:              "dev",
				Service:            "carts",
				KeptnContext:       "test-context",
				eventBrokerURL:     defaultEventBrokerURL,
				useLocalFileSystem: false,
				resourceHandler: &ResourceHandler{
					BaseURL:    configurationServiceURL,
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
		{
			name: "Get local Keptn",
			args: args{
				project:      "sockshop",
				stage:        "dev",
				service:      "carts",
				keptnContext: "test-context",
				opts: KeptnOpts{
					UseLocalFileSystem:      true,
					ConfigurationServiceURL: "",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
				Project:            "sockshop",
				Stage:              "dev",
				Service:            "carts",
				KeptnContext:       "test-context",
				eventBrokerURL:     defaultEventBrokerURL,
				useLocalFileSystem: true,
				resourceHandler: &ResourceHandler{
					BaseURL:    configurationServiceURL,
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
		{
			name: "Get Keptn with custom configuration service URL",
			args: args{
				project:      "sockshop",
				stage:        "dev",
				service:      "carts",
				keptnContext: "test-context",
				opts: KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
				Project:            "sockshop",
				Stage:              "dev",
				Service:            "carts",
				KeptnContext:       "test-context",
				eventBrokerURL:     defaultEventBrokerURL,
				useLocalFileSystem: false,
				resourceHandler: &ResourceHandler{
					BaseURL:    "custom-config:8080",
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
		{
			name: "Get Keptn with custom event brokerURL",
			args: args{
				project:      "sockshop",
				stage:        "dev",
				service:      "carts",
				keptnContext: "test-context",
				opts: KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "custom-eb:8080",
				},
			},
			want: &Keptn{
				Project:            "sockshop",
				Stage:              "dev",
				Service:            "carts",
				KeptnContext:       "test-context",
				eventBrokerURL:     "custom-eb:8080",
				useLocalFileSystem: false,
				resourceHandler: &ResourceHandler{
					BaseURL:    "custom-config:8080",
					AuthHeader: "",
					AuthToken:  "",
					HTTPClient: &http.Client{},
					Scheme:     "http",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKeptn(tt.args.project, tt.args.stage, tt.args.service, tt.args.keptnContext, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeptn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeptn_GetKeptnResource(t *testing.T) {
	type fields struct {
		Project            string
		Stage              string
		Service            string
		KeptnContext       string
		eventBrokerURL     string
		useLocalFileSystem bool
		resourceHandler    *ResourceHandler
	}
	type args struct {
		resource string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Keptn{
				Project:            tt.fields.Project,
				Stage:              tt.fields.Stage,
				Service:            tt.fields.Service,
				KeptnContext:       tt.fields.KeptnContext,
				eventBrokerURL:     tt.fields.eventBrokerURL,
				useLocalFileSystem: tt.fields.useLocalFileSystem,
				resourceHandler:    tt.fields.resourceHandler,
			}
			got, err := k.GetKeptnResource(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKeptnResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKeptnResource() got = %v, want %v", got, tt.want)
			}
		})
	}
}
