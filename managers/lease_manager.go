package managers


import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"sync"

	"github.com/indkumar8999/ps-tasks/leases"
)


type LeaseManager struct {
	leasesDir string
	leases     map[string]*leases.Lease
	leaseLock *sync.Mutex
}

// NewLeaseManager creates a new LeaseManager
func NewLeaseManager(leasesDir string) (*LeaseManager, error) {
	// Create the leases directory if it doesn't exist
	if _, err := os.Stat(leasesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(leasesDir, 0755); err != nil {
			fmt.Printf("Error creating leases directory: %v\n", err)
			return nil, err
		}
	}

	return &LeaseManager{
		leasesDir:  leasesDir,
		leases:     make(map[string]*leases.Lease),
		leaseLock:  &sync.Mutex{},
	}, nil
}

// AcquireLease acquires a lease for a task
func (lm *LeaseManager) AcquireLease(taskID string, duration time.Duration, username string) (*leases.Lease, error) {
	lm.leaseLock.Lock()
	defer lm.leaseLock.Unlock()
	// Check if the task ID is valid
	if taskID == "" {
		return nil, fmt.Errorf("invalid task ID")
	}
	// Check if the duration is valid
	if duration <= 0 {
		return nil, fmt.Errorf("invalid lease duration")
	}

	// Check if a lease already exists for the task
	if lease, exists := lm.leases[taskID]; exists {
		if !lease.IsExpired() {
			return nil, fmt.Errorf("lease already exists and is not expired")
		}
	}
	// Create a new lease

	lease := leases.NewLease(taskID, duration, username)
	if err := lease.Save(lm.leasesDir); err != nil {
		return nil, err
	}

	lm.leases[lease.ID] = lease
	return lease, nil
}

// ReleaseLease releases a lease for a task
func (lm *LeaseManager) ReleaseLease(leaseID string) error {
	lm.leaseLock.Lock()
	defer lm.leaseLock.Unlock()
	// Check if the lease ID is valid
	if leaseID == "" {
		return fmt.Errorf("invalid lease ID")
	}
	// Check if the lease exists
	lease, exists := lm.leases[leaseID]
	if !exists {
		return fmt.Errorf("lease not found")
	}

	if err := os.Remove(filepath.Join(lm.leasesDir, fmt.Sprintf("%s.json", lease.ID))); err != nil {
		return err
	}

	delete(lm.leases, leaseID)
	return nil
}

// GetLease retrieves a lease by its ID
func (lm *LeaseManager) GetLease(leaseID string) (*leases.Lease, error) {
	lm.leaseLock.Lock()
	defer lm.leaseLock.Unlock()
	// Check if the lease ID is valid
	if leaseID == "" {
		return nil, fmt.Errorf("invalid lease ID")
	}
	// Check if the lease exists
	lease, exists := lm.leases[leaseID]
	if !exists {
		return nil, fmt.Errorf("lease not found")
	}

	return lease, nil
}

// LoadLeases loads all leases from the leases directory
func (lm *LeaseManager) LoadLeases() error {
	lm.leaseLock.Lock()
	defer lm.leaseLock.Unlock()

	files, err := os.ReadDir(lm.leasesDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		leaseID := file.Name()
		lease, err := leases.LoadLease(lm.leasesDir, leaseID)
		if err != nil {
			return err
		}
		lm.leases[lease.ID] = lease
	}

	return nil
}

// CleanupExpiredLeases removes expired leases from the directory
func (lm *LeaseManager) CleanupExpiredLeases() error {
	lm.leaseLock.Lock()
	defer lm.leaseLock.Unlock()

	files, err := os.ReadDir(lm.leasesDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		leaseID := file.Name()
		lease, err := leases.LoadLease(lm.leasesDir, leaseID)
		if err != nil {
			return err
		}
		if lease.IsExpired() {
			if err := os.Remove(filepath.Join(lm.leasesDir, fmt.Sprintf("%s.json", lease.ID))); err != nil {
				return err
			}
			delete(lm.leases, lease.ID)
		}
	}
	return nil
}

// ExtendLease extends the lease duration for a task
func (lm *LeaseManager) ExtendLease(leaseID string, duration time.Duration, username string) error {
	lm.leaseLock.Lock()
	defer lm.leaseLock.Unlock()
	// Check if the lease ID is valid
	if leaseID == "" {
		return fmt.Errorf("invalid lease ID")
	}
	// Check if the duration is valid
	if duration <= 0 {
		return fmt.Errorf("invalid lease duration")
	}
	// Check if the lease exists
	lease, exists := lm.leases[leaseID]
	if !exists {
		return fmt.Errorf("lease not found")
	}
	// Check if the lease is expired
	if lease.IsExpired() {
		return fmt.Errorf("lease is expired")
	}
	// Check if the lease is already extended
	if time.Now().Add(duration).Before(lease.ExpiresAt) {
		return fmt.Errorf("lease is already extended")
	}
	if lease.CreatedBy != username {
		return fmt.Errorf("lease was created by another user")
	}
	// Extend the lease duration
	lease.ExpiresAt = time.Now().Add(duration)
	if err := lease.Save(lm.leasesDir); err != nil {
		return err
	}
	return nil
}

