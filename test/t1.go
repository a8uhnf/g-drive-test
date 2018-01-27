package test
import (
	"testing"
)

func TestingTest(t *testing.T) {
	actual := "hello"
	if actual != "ello" {
		t.Fatalf("Expected %s, but get ello", actual)
	}
}

