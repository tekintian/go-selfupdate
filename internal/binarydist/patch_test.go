package binarydist

import (
	"bytes"
	"testing"
)

func TestPatch(t *testing.T) {
	oldFile := mustWriteRandFile("test.old", 1e3)
	newFile := mustWriteRandFile("test.new", 1e3)

	// Create a patch using our Diff function
	patchBuf := new(bytes.Buffer)
	err := Diff(oldFile, newFile, patchBuf)
	if err != nil {
		panic(err)
	}

	// Reset old file for reading
	_, err = oldFile.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	// Apply the patch
	resultBuf := new(bytes.Buffer)
	err = Patch(oldFile, resultBuf, patchBuf)
	if err != nil {
		t.Fatal("Patch failed:", err)
	}

	// Verify result matches new
	_, err = newFile.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	newBytes := mustReadAll(newFile)

	if !bytes.Equal(resultBuf.Bytes(), newBytes) {
		t.Fatalf("produced different output")
	}

	t.Logf("Patch successful: %d bytes", resultBuf.Len())
}

func TestPatchHk(t *testing.T) {
	resultBuf := new(bytes.Buffer)

	err := Patch(mustOpen("testdata/sample.old"), resultBuf, mustOpen("testdata/sample.patch"))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("got %d bytes", resultBuf.Len())
	expected := mustReadAll(mustOpen("testdata/sample.new"))
	if !bytes.Equal(resultBuf.Bytes(), expected) {
		t.Fatalf("produced different output")
	}
}
