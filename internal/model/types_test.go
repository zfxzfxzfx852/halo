package model

import "testing"

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
