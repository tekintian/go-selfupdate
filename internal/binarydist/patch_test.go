package binarydist

import (
	"bytes"
	"os"
	"testing"
)

func TestPatch(t *testing.T) {
	old := mustWriteRandFile("test.old", 1e3)
	new := mustWriteRandFile("test.new", 1e3)

	// Create a patch using our Diff function
	patch := new(bytes.Buffer)
	err := Diff(old, new, patch)
	if err != nil {
		panic(err)
	}

	// Reset old file for reading
	_, err = old.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	// Apply the patch
	result := new(bytes.Buffer)
	err = Patch(old, result, patch)
	if err != nil {
		t.Fatal("Patch failed:", err)
	}

	// Verify result matches new
	_, err = new.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	newBytes := mustReadAll(new)

	if !bytes.Equal(result.Bytes(), newBytes) {
		if n := matchlen(result.Bytes(), newBytes); n > -1 {
			t.Fatalf("produced different output at pos %d", n)
		} else {
			t.Fatalf("produced different output")
		}
	}

	t.Logf("Patch successful: %d bytes", result.Len())
}

func TestPatchHk(t *testing.T) {
	result := new(bytes.Buffer)

	err := Patch(mustOpen("testdata/sample.old"), result, mustOpen("testdata/sample.patch"))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("got %d bytes", result.Len())
	if n := fileCmpBytes(result.Bytes(), mustReadAll(mustOpen("testdata/sample.new"))); n > -1 {
		t.Fatalf("produced different output at pos %d", n)
	}
}

// fileCmpBytes compares a byte slice with expected bytes
func fileCmpBytes(a, b []byte) int64 {
	if len(a) != len(b) {
		return int64(len(a))
	}

	for i := range a {
		if a[i] != b[i] {
			return int64(i)
		}
	}
	return -1
}
