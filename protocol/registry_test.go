package protocol

import (
	"strconv"
	"strings"
	"testing"
)

func TestVersionFromStringExact(t *testing.T) {
	got, ok := VersionFromString("1.21.8")
	if !ok {
		t.Fatalf("expected exact version lookup to succeed")
	}
	if got != V1_21_8 {
		t.Fatalf("expected V1_21_8, got %v", got)
	}
}

func TestVersionFromStringFallbackFuturePatch(t *testing.T) {
	got, ok := VersionFromString("1.21.999")
	if !ok {
		t.Fatalf("expected fallback lookup to succeed")
	}
	expected, ok := latestPatchForMinor(1, 21)
	if !ok {
		t.Fatalf("expected to discover at least one 1.21.x version")
	}
	if got != expected {
		t.Fatalf("expected fallback to %v, got %v", expected, got)
	}
}

func TestVersionFromStringFallbackNearestLowerPatch(t *testing.T) {
	got, ok := VersionFromString("1.21.2")
	if !ok {
		t.Fatalf("expected nearest lower patch fallback to succeed")
	}
	if got != V1_21_1 {
		t.Fatalf("expected fallback to V1_21_1, got %v", got)
	}
}

func TestVersionFromStringRejectNonRelease(t *testing.T) {
	if _, ok := VersionFromString("25w37a"); ok {
		t.Fatalf("expected snapshot string to be rejected without explicit generated support")
	}
}

func latestPatchForMinor(major, minor int) (Version, bool) {
	bestPatch := -1
	var best Version
	for _, v := range Versions {
		mj, mn, p, ok := parseVersion(v.String())
		if !ok || mj != major || mn != minor {
			continue
		}
		if p > bestPatch {
			bestPatch = p
			best = v
		}
	}
	return best, bestPatch >= 0
}

func parseVersion(s string) (int, int, int, bool) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return 0, 0, 0, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, 0, false
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, 0, false
	}
	return major, minor, patch, true
}
