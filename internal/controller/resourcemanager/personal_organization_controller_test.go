// SPDX-License-Identifier: AGPL-3.0-only

package resourcemanager

import "testing"

func TestHashPersonalOrgName(t *testing.T) {
	first := hashPersonalOrgName("uid-123")
	second := hashPersonalOrgName("uid-123")
	if first != second {
		t.Fatalf("hashPersonalOrgName() not stable: %q vs %q", first, second)
	}
	if first == "" {
		t.Fatal("hashPersonalOrgName() returned empty string")
	}
	if hashPersonalOrgName("uid-456") == first {
		t.Fatal("hashPersonalOrgName() returned same value for different inputs")
	}
}
