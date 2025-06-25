package zfs_test

import (
	"testing"

	"github.com/Tualua/gozfsapi/pkgs/zfs"
)

func Test_GetDataset(t *testing.T) {
	// Example dataset name
	datasetName := "data/kvm/desktop"

	// Call the GetDataset function
	ds, err := zfs.GetDataset(datasetName)
	if err != nil {
		t.Fatalf("Failed to get dataset: %v", err)
	}

	// Check if the dataset name matches
	if ds.Name != datasetName {
		t.Errorf("Expected dataset name %s, got %s", datasetName, ds.Name)
	}

	if ds.Type != "filesystem" {
		t.Errorf("Expected dataset type 'filesystem', got '%s'", ds.Type)
	}
	// Additional checks can be added here for other properties of the dataset
}

func Test_ListDatasets(t *testing.T) {
	// Call the ListDatasets function with an empty type to list all datasets
	datasets := zfs.ListDatasets("")
	if len(datasets) == 0 {
		t.Error("Expected to find at least one dataset, but got none")
	}

	// Check if the first dataset is not empty
	if datasets[0] == "" {
		t.Error("Expected first dataset name to be non-empty")
	}
}

func Test_CloneDataset(t *testing.T) {
	// Example source and target dataset names
	testTable := []struct {
		source string
		target string
		failed bool
	}{
		{"data/kvm/desktop/win10-master", "data/kvm/desktop/desktop-vm001", true},    // This should fail because the source dataset is not a snapshot
		{"data/kvm/desktop/win10-master@1", "data/kvm/desktop/desktop-vm001", false}, // This should succeed if the source is a valid snapshot
	}
	for _, tt := range testTable {
		t.Run(
			tt.source+"->"+tt.target,
			func(t *testing.T) {
				err := zfs.CloneDataset(tt.source, tt.target)
				if err != nil && !tt.failed {
					t.Fatalf("Failed to clone dataset: %v", err)
				} else if err == nil && tt.failed {
					t.Fatalf("Expected error when cloning dataset, but got none")
				}
				// Verify that the cloned dataset exists
				if !tt.failed {
					clonedDataset, err := zfs.GetDataset(tt.target)
					if err != nil {
						t.Fatalf("Cloned dataset not found: %v", err)
					}
					if clonedDataset.Name != tt.target {
						t.Errorf("Expected cloned dataset name %s, got %s", tt.target, clonedDataset.Name)
					}
				}
				t.Logf("expected error was: %v\n", err)
			},
		)

	}
}
