package client

import (
	"context"
	"time"

	"k8s.io/client-go/tools/cache"

	"github.com/supergate-hub/slurm-client/pkg/object"
)

// Reader is the interface that wraps the basic Get and List methods.
type Reader interface {
	Get(ctx context.Context, key object.ObjectKey, obj object.Object, opts ...GetOption) error
	List(ctx context.Context, list object.ObjectList, opts ...ListOption) error
}

// Writer is the interface that wraps the basic Create, Update, and Delete methods.
type Writer interface {
	Create(ctx context.Context, obj object.Object, req any, opts ...CreateOption) error
	Update(ctx context.Context, obj object.Object, req any, opts ...UpdateOption) error
	Delete(ctx context.Context, obj object.Object, opts ...DeleteOption) error
}

// Informer is the interface that wraps the Informer methods.
type Informer interface {
	GetInformer(objectType object.ObjectType) InformerCache
	Start(ctx context.Context)
	Stop()
	GetServer() string
	SetServer(server string)
	GetToken() string
	SetToken(token string)
}

// Client is the interface that groups Reader, Writer and Informer methods.
type Client interface {
	Reader
	Writer
	Informer
}

// InformerCache is the interface that wraps the methods for a cached informer.
type InformerCache interface {
	Reader
	Writer
	SetEventHandler(handler cache.ResourceEventHandler)
	UnsetEventHandler()
	Run(stopCh <-chan struct{})
	HasSynced() (bool, error)
	HasStarted() bool
	WaitForSyncList(ctx context.Context, interval time.Duration) error
	WaitForSyncGet(ctx context.Context, interval time.Duration) error
}

// =============================================================================
// Unified Client Interfaces (Newly Added)
// =============================================================================

// SlurmClientInterface is the common interface for Slurm (workload) clients.
type SlurmClientInterface interface {
	Jobs() GenericJobService
	Nodes() GenericNodeService
	// Add other getters as needed
}

// SlurmdbClientInterface is the common interface for Slurmdb (accounting) clients.
type SlurmdbClientInterface interface {
	Accounts() GenericAccountService
	// Add other getters as needed
}

// GenericJobService defines common job operations.
type GenericJobService interface {
	List(ctx context.Context, opts interface{}) (interface{}, error)
	Get(ctx context.Context, id string) (interface{}, error)
	Delete(ctx context.Context, id string) error
}

// GenericNodeService defines common node operations.
type GenericNodeService interface {
	List(ctx context.Context, opts interface{}) (interface{}, error)
	Get(ctx context.Context, name string) (interface{}, error)
}

// GenericAccountService defines common account operations.
type GenericAccountService interface {
	List(ctx context.Context, opts interface{}) (interface{}, error)
	Get(ctx context.Context, name string) (interface{}, error)
}

