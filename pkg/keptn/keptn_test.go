package keptn

import (
	"encoding/json"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestNewKeptn(t *testing.T) {
	incomingEvent := cloudevents.New(cloudevents.CloudEventsVersionV02)
	incomingEvent.SetSource("test")
	incomingEvent.SetExtension("shkeptncontext", "test-context")
	incomingEvent.SetDataContentType(cloudevents.ApplicationCloudEventsJSON)

	keptnBase := &KeptnBase{
		Project:            "sockshop",
		Stage:              "dev",
		Service:            "carts",
		TestStrategy:       nil,
		DeploymentStrategy: nil,
		Tag:                nil,
		Image:              nil,
		Labels:             nil,
	}

	marshal, _ := json.Marshal(keptnBase)
	incomingEvent.Data = marshal

	incomingEvent.SetData(marshal)
	incomingEvent.DataEncoded = true
	incomingEvent.DataBinary = true

	type args struct {
		incomingEvent *cloudevents.Event
		opts          KeptnOpts
	}
	tests := []struct {
		name string
		args args
		want *Keptn
	}{
		{
			name: "Get 'in-cluster' Keptn",
			args: args{
				incomingEvent: &incomingEvent,
				opts:          KeptnOpts{},
			},
			want: &Keptn{
				KeptnBase: &KeptnBase{
					Project:            "sockshop",
					Stage:              "dev",
					Service:            "carts",
					TestStrategy:       nil,
					DeploymentStrategy: nil,
					Tag:                nil,
					Image:              nil,
					Labels:             nil,
				},
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
				incomingEvent: &incomingEvent,
				opts: KeptnOpts{
					UseLocalFileSystem:      true,
					ConfigurationServiceURL: "",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
				KeptnBase: &KeptnBase{
					Project:            "sockshop",
					Stage:              "dev",
					Service:            "carts",
					TestStrategy:       nil,
					DeploymentStrategy: nil,
					Tag:                nil,
					Image:              nil,
					Labels:             nil,
				},
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
				incomingEvent: &incomingEvent,
				opts: KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "",
				},
			},
			want: &Keptn{
				KeptnBase: &KeptnBase{
					Project:            "sockshop",
					Stage:              "dev",
					Service:            "carts",
					TestStrategy:       nil,
					DeploymentStrategy: nil,
					Tag:                nil,
					Image:              nil,
					Labels:             nil,
				},
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
				incomingEvent: &incomingEvent,
				opts: KeptnOpts{
					UseLocalFileSystem:      false,
					ConfigurationServiceURL: "custom-config:8080",
					EventBrokerURL:          "custom-eb:8080",
				},
			},
			want: &Keptn{
				KeptnBase: &KeptnBase{
					Project:            "sockshop",
					Stage:              "dev",
					Service:            "carts",
					TestStrategy:       nil,
					DeploymentStrategy: nil,
					Tag:                nil,
					Image:              nil,
					Labels:             nil,
				},
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
			if got, _ := NewKeptn(tt.args.incomingEvent, tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKeptn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeptn_GetKeptnResource(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(200)

			res := &Resource{
				ResourceContent: "dGVzdC1jb250ZW50Cg==",
				ResourceURI:     stringp("test-resource.file"),
			}
			marshal, _ := json.Marshal(res)
			w.Write(marshal)
		}),
	)
	defer ts.Close()

	type fields struct {
		KeptnBase          *KeptnBase
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
		{
			name: "get a resource",
			fields: fields{
				KeptnBase: &KeptnBase{
					Project:            "sockshop",
					Stage:              "dev",
					Service:            "carts",
					TestStrategy:       nil,
					DeploymentStrategy: nil,
					Tag:                nil,
					Image:              nil,
					Labels:             nil,
				},
				eventBrokerURL:     "",
				useLocalFileSystem: false,
				resourceHandler:    NewResourceHandler(ts.URL),
			},
			args: args{
				resource: "test-resource.file",
			},
			want:    "test-content",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Keptn{
				KeptnBase:          tt.fields.KeptnBase,
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
			_ = os.RemoveAll(tt.args.resource)
		})
	}
}

func stringp(s string) *string {
	return &s
}
