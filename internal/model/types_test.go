package model

import (
	"strings"
	"testing"
)

func TestNewID(t *testing.T) {
	id := NewID("evt")
	if !strings.HasPrefix(id, "evt_") {
		t.Fatalf("id %q missing prefix", id)
	}
	// 8 字节 -> 16 个十六进制字符，加前缀与下划线
	if len(id) != len("evt_")+16 {
		t.Fatalf("id length = %d", len(id))
	}
	if NewID("evt") == id {
		t.Fatal("two ids should differ")
	}
}

func TestRiskSeverityOrdering(t *testing.T) {
	if RiskNone.Severity() != 0 {
		t.Fatal("none should be 0")
	}
	if !(RiskLow.Severity() < RiskMedium.Severity() &&
		RiskMedium.Severity() < RiskHigh.Severity() &&
		RiskHigh.Severity() < RiskCritical.Severity()) {
		t.Fatal("risk severity must be strictly ordered")
	}
}
