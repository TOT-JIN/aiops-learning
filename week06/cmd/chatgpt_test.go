package cmd

import (
	"testing"
)

func TestDeleteResource(t *testing.T) {
	tests := []struct {
		name         string
		namespace    string
		resourceType string
		resourceName string
		expectedMsg  string
		expectError  bool
	}{
		{
			name:         "Valid deployment deletion",
			namespace:    "default",
			resourceType: "deployment",
			resourceName: "nginx-deployment",
			expectedMsg:  "nginx-deployment 删除成功",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := deleteResource(tt.namespace, tt.resourceType, tt.resourceName)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if msg != tt.expectedMsg {
				t.Errorf("expected message: %s, got: %s", tt.expectedMsg, msg)
			}
		})
	}
}
