package main

import "testing"

func TestHello(t *testing.T) {
	want := "Hello world!"
	got := hello("world")
	if got != want {
		t.Errorf("want: %s got: %s", want, got)
	}
}
