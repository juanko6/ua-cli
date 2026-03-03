package main

import "testing"

func TestIsAuthExempt(t *testing.T) {
	if !isAuthExempt(loginCmd) {
		t.Fatalf("login should be auth exempt")
	}
	if isAuthExempt(scheduleCmd) {
		t.Fatalf("schedule should not be auth exempt")
	}
}
