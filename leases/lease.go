package leases

import (
	"fmt"
	"encoding/json"
	
	"os"
	"path/filepath"
	"time"
	"github.com/google/uuid"
)

// Lease represents a lease for a task
type Lease struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
}

// NewLease creates a new lease for a task
func NewLease(taskID string, duration time.Duration, username string) *Lease {
	return &Lease{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
		CreatedBy: username,
		UpdatedAt: time.Now(),
		UpdatedBy: username,
	}
}

// IsExpired checks if the lease is expired
func (l *Lease) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

// Save saves the lease to a file
func (l *Lease) Save(leasesDir string) error {
	leaseFile := filepath.Join(leasesDir, fmt.Sprintf("%s.json", l.ID))
	file, err := os.Create(leaseFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serialize the lease to JSON and write to the file
	if err := json.NewEncoder(file).Encode(l); err != nil {
		return err
	}

	return nil
}


// LoadLease loads a lease from a file
func LoadLease(leasesDir, leaseID string) (*Lease, error) {
	leaseFile := filepath.Join(leasesDir, fmt.Sprintf("%s.json", leaseID))
	file, err := os.Open(leaseFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lease Lease
	if err := json.NewDecoder(file).Decode(&lease); err != nil {
		return nil, err
	}

	return &lease, nil
}

