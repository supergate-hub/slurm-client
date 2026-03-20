// Package slurmtest provides an in-memory mock slurmrestd server for testing.
//
// Usage:
//
//	mock := slurmtest.NewServer()
//	defer mock.Close()
//
//	mock.AddJob(slurmtest.MockJob{ID: 123, Name: "test"})
//
//	client, _ := slurmclient.NewClient(ctx, slurmclient.AuthOpts{
//	    Endpoint:  mock.URL(),
//	    AuthToken: "test",
//	    Version:   slurmclient.V0040,
//	})
//	jobs, _ := client.Slurm.Jobs().List(ctx)
package slurmtest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"

	v0040 "github.com/supergate-hub/slurm-client/api/v0040"
)

// MockJob represents a mock Slurm job.
type MockJob struct {
	ID   int32
	Name string
}

// MockNode represents a mock Slurm node.
type MockNode struct {
	Name string
}

// MockPartition represents a mock Slurm partition.
type MockPartition struct {
	Name string
}

// MockAccount represents a mock SlurmDB account.
type MockAccount struct {
	Name string
}

// Server is an in-memory mock slurmrestd HTTP server.
type Server struct {
	srv *httptest.Server

	mu         sync.RWMutex
	jobs       []MockJob
	nodes      []MockNode
	partitions []MockPartition
	accounts   []MockAccount
}

// NewServer creates and starts a new mock slurmrestd server.
func NewServer() *Server {
	s := &Server{}
	s.srv = httptest.NewServer(s)
	return s
}

// URL returns the server's base URL.
func (s *Server) URL() string { return s.srv.URL }

// Close shuts down the server.
func (s *Server) Close() { s.srv.Close() }

// AddJob adds a mock job.
func (s *Server) AddJob(j MockJob) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, j)
}

// AddNode adds a mock node.
func (s *Server) AddNode(n MockNode) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes = append(s.nodes, n)
}

// AddPartition adds a mock partition.
func (s *Server) AddPartition(p MockPartition) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.partitions = append(s.partitions, p)
}

// AddAccount adds a mock account.
func (s *Server) AddAccount(a MockAccount) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accounts = append(s.accounts, a)
}

// ServeHTTP handles mock API requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Strip version prefix: /slurm/v0.0.40/jobs → /jobs
	path := r.URL.Path
	parts := strings.SplitN(path, "/", 4) // ["", "slurm", "v0.0.40", "jobs"]
	if len(parts) >= 4 {
		path = "/" + parts[3]
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	switch {
	case r.Method == http.MethodGet && path == "/jobs":
		s.handleListJobs(w)
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/job/"):
		s.handleGetJob(w, strings.TrimPrefix(path, "/job/"))
	case r.Method == http.MethodPost && path == "/job/submit":
		s.handleSubmitJob(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(path, "/job/"):
		w.WriteHeader(200)
	case r.Method == http.MethodGet && path == "/nodes":
		s.handleListNodes(w)
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/node/"):
		s.handleGetNode(w, strings.TrimPrefix(path, "/node/"))
	case r.Method == http.MethodGet && path == "/partitions":
		s.handleListPartitions(w)
	case r.Method == http.MethodGet && path == "/accounts":
		s.handleListAccounts(w)
	case r.Method == http.MethodGet && strings.HasPrefix(path, "/account/"):
		s.handleGetAccount(w, strings.TrimPrefix(path, "/account/"))
	case r.Method == http.MethodGet && path == "/ping":
		s.handlePing(w)
	case path == "/openapi/v3":
		s.handleOpenAPI(w)
	default:
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]any{})
	}
}

func (s *Server) handleListJobs(w http.ResponseWriter) {
	jobs := make([]v0040.V0040JobInfo, len(s.jobs))
	for i, j := range s.jobs {
		id := j.ID
		name := j.Name
		jobs[i] = v0040.V0040JobInfo{JobId: &id, Name: &name}
	}
	json.NewEncoder(w).Encode(v0040.V0040OpenapiJobInfoResp{Jobs: jobs})
}

func (s *Server) handleGetJob(w http.ResponseWriter, idStr string) {
	for _, j := range s.jobs {
		if strings.TrimSpace(idStr) == strings.TrimSpace(json.Number(itoa(j.ID)).String()) {
			id := j.ID
			name := j.Name
			json.NewEncoder(w).Encode(v0040.V0040OpenapiJobInfoResp{
				Jobs: []v0040.V0040JobInfo{{JobId: &id, Name: &name}},
			})
			return
		}
	}
	w.WriteHeader(404)
}

func (s *Server) handleSubmitJob(w http.ResponseWriter, r *http.Request) {
	nextID := int32(len(s.jobs) + 1000)
	json.NewEncoder(w).Encode(v0040.V0040OpenapiJobSubmitResponse{
		JobId: &nextID,
	})
}

func (s *Server) handleListNodes(w http.ResponseWriter) {
	nodes := make([]v0040.V0040Node, len(s.nodes))
	for i, n := range s.nodes {
		name := n.Name
		nodes[i] = v0040.V0040Node{Name: &name}
	}
	json.NewEncoder(w).Encode(v0040.V0040OpenapiNodesResp{Nodes: nodes})
}

func (s *Server) handleGetNode(w http.ResponseWriter, name string) {
	for _, n := range s.nodes {
		if n.Name == name {
			nn := n.Name
			json.NewEncoder(w).Encode(v0040.V0040OpenapiNodesResp{
				Nodes: []v0040.V0040Node{{Name: &nn}},
			})
			return
		}
	}
	w.WriteHeader(404)
}

func (s *Server) handleListPartitions(w http.ResponseWriter) {
	parts := make([]v0040.V0040PartitionInfo, len(s.partitions))
	for i, p := range s.partitions {
		name := p.Name
		parts[i] = v0040.V0040PartitionInfo{Name: &name}
	}
	json.NewEncoder(w).Encode(v0040.V0040OpenapiPartitionResp{Partitions: parts})
}

func (s *Server) handleListAccounts(w http.ResponseWriter) {
	accts := make([]v0040.V0040Account, len(s.accounts))
	for i, a := range s.accounts {
		accts[i] = v0040.V0040Account{Name: a.Name}
	}
	json.NewEncoder(w).Encode(v0040.V0040OpenapiAccountsResp{Accounts: accts})
}

func (s *Server) handleGetAccount(w http.ResponseWriter, name string) {
	for _, a := range s.accounts {
		if a.Name == name {
			json.NewEncoder(w).Encode(v0040.V0040OpenapiAccountsResp{
				Accounts: []v0040.V0040Account{{Name: a.Name}},
			})
			return
		}
	}
	w.WriteHeader(404)
}

func (s *Server) handlePing(w http.ResponseWriter) {
	hostname := "mock-controller"
	mode := "primary"
	json.NewEncoder(w).Encode(v0040.V0040OpenapiPingArrayResp{
		Pings: []v0040.V0040ControllerPing{
			{Hostname: &hostname, Mode: &mode},
		},
	})
}

func (s *Server) handleOpenAPI(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(map[string]any{
		"paths": map[string]any{
			"/slurm/v0.0.40/jobs":    nil,
			"/slurm/v0.0.40/nodes":   nil,
			"/slurmdb/v0.0.40/users": nil,
		},
	})
}

func itoa(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}
