package protocol

import "testing"

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
	got, ok := VersionFromString("1.21.9")
	if !ok {
		t.Fatalf("expected fallback lookup to succeed")
	}
	if got != V1_21_8 {
		t.Fatalf("expected fallback to V1_21_8, got %v", got)
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

