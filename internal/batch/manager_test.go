package batch

import (
	"testing"
	"time"

	"RealityScout/internal/types"
)

func TestSortSuitableDomains(t *testing.T) {
	makeResult := func(domain string, hardRequirements bool, cdn, hot *bool, handshake time.Duration) *types.DetectionResult {
		result := &types.DetectionResult{
			Domain:              domain,
			HardRequirementsMet: hardRequirements,
			TLS:                 &types.TLSResult{HandshakeTime: handshake},
		}
		if cdn != nil || hot != nil {
			result.CDN = &types.CDNResult{}
			if cdn != nil {
				result.CDN.IsCDN = *cdn
			}
			if hot != nil {
				result.CDN.IsHotWebsite = *hot
			}
		}
		return result
	}
	flag := func(value bool) *bool { return &value }

	results := []*types.DetectionResult{
		makeResult("slow-clean", true, flag(false), flag(false), 400*time.Millisecond),
		makeResult("cdn", true, flag(true), flag(false), 10*time.Millisecond),
		makeResult("hot", true, flag(false), flag(true), 10*time.Millisecond),
		makeResult("unknown-signals", true, nil, nil, 20*time.Millisecond),
		makeResult("fast-clean", true, flag(false), flag(false), 30*time.Millisecond),
		makeResult("missing-hard-requirement", false, flag(false), flag(false), time.Millisecond),
	}

	NewManager(nil).sortSuitableDomains(results)

	want := []string{
		"fast-clean",
		"slow-clean",
		"unknown-signals",
		"hot",
		"cdn",
		"missing-hard-requirement",
	}
	for i, result := range results {
		if result.Domain != want[i] {
			t.Fatalf("sort[%d] = %q, want %q", i, result.Domain, want[i])
		}
	}
}

func TestSortSuitableDomainsIsStableForTies(t *testing.T) {
	results := []*types.DetectionResult{
		{Domain: "first", HardRequirementsMet: true},
		{Domain: "second", HardRequirementsMet: true},
	}

	NewManager(nil).sortSuitableDomains(results)

	if results[0].Domain != "first" || results[1].Domain != "second" {
		t.Fatalf("equal-ranked results changed input order: %q, %q", results[0].Domain, results[1].Domain)
	}
}
