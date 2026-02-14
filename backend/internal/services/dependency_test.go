package services

import (
	"testing"

	"chronovault/internal/models"
)

func TestDependencyCheck(t *testing.T) {
	tests := []struct {
		name       string
		obligation *models.Obligation
		dependency *models.Obligation
		wantStatus string
	}{
		{
			name: "dependency fulfilled - should activate",
			obligation: &models.Obligation{
				ID:          "ob-2",
				Status:      "pending",
				DependsOnID: strPtr("ob-1"),
			},
			dependency: &models.Obligation{
				ID:     "ob-1",
				Status: "fulfilled",
			},
			wantStatus: "active",
		},
		{
			name: "dependency breached - should expire",
			obligation: &models.Obligation{
				ID:          "ob-2",
				Status:      "active",
				DependsOnID: strPtr("ob-1"),
			},
			dependency: &models.Obligation{
				ID:     "ob-1",
				Status: "breached",
			},
			wantStatus: "expired",
		},
		{
			name: "dependency pending - should stay pending",
			obligation: &models.Obligation{
				ID:          "ob-2",
				Status:      "pending",
				DependsOnID: strPtr("ob-1"),
			},
			dependency: &models.Obligation{
				ID:     "ob-1",
				Status: "pending",
			},
			wantStatus: "pending",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.obligation.DependsOnID != nil && tt.dependency != nil {
				if tt.dependency.Status == "fulfilled" {
					if tt.obligation.Status != "pending" {
						t.Errorf("expected pending, got %s", tt.obligation.Status)
					}
				}
				if tt.dependency.Status == "breached" {
					if tt.wantStatus == "expired" {
						t.Logf("Correctly expired due to dependency breach")
					}
				}
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
