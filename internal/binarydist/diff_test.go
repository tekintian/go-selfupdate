package binarydist

import (
	"bytes"
	"os"
	"testing"
)

var diffT = []struct {
	old *os.File
	new *os.File
}{
	{
		old: mustWriteRandFile("test.old", 1e3),
		new: mustWriteRandFile("test.new", 1e3),
	},
	{
		old: mustOpen("testdata/sample.old"),
		new: mustOpen("testdata/sample.new"),
	},
}

func TestDiff(t *testing.T) {
	for _, s := range diffT {
		// Test that Diff produces a patch
		patch := new(bytes.Buffer)
		err := Diff(s.old, s.new, patch)
		if err != nil {
			t.Fatalf("Diff failed: %v", err)
		}

		// Verify the patch is not empty
		if patch.Len() == 0 {
			t.Fatalf("Diff produced empty patch")
		}

		t.Logf("diff %s %s: patch size = %d bytes", s.old.Name(), s.new.Name(), patch.Len())

		// Test round-trip: old + patch should = new
		result := new(bytes.Buffer)
		err = Patch(s.old, result, patch)
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		// Read original new file
		_, err = s.new.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		newBytes := mustReadAll(s.new)

		// Compare result with expected
		if !bytes.Equal(result.Bytes(), newBytes) {
			t.Fatalf("Round-trip failed: result does not match expected new file")
		}
	}
}
