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