// Selfupdate Example - Simple Self-Updating CLI Tool
// This example demonstrates how to build a self-updating Go application
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tekintian/go-selfupdate"
)

var (
	version = "1.0.0"
)

func main() {
	updateFlag := flag.Bool("update", false, "Check for and apply updates")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("go-selfupdate example v%s\n", version)
		return
	}

	if *updateFlag {
		err := doUpdate()
		if err != nil {
			log.Fatalf("Update failed: %v", err)
		}
		fmt.Println("Update successful!")
		return
	}

	fmt.Println("go-selfupdate Example Program")
	fmt.Printf("Version: %s\n", version)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  -update    Check for and apply updates")
	fmt.Println("  -version   Show version information")
	fmt.Println()
	fmt.Println("Example of updating this program:")
	fmt.Println("  ./example -update")
}

func doUpdate() error {
	// In a real application, you would:
	// 1. Check an API endpoint for the latest version
	// 2. Download the update from a secure URL
	// 3. Apply the update with selfupdate.Apply()

	// For demonstration, we'll simulate the update process:
	url := fmt.Sprintf("https://example.com/updates/example-%s", version)
	fmt.Printf("Checking for updates from: %s\n", url)

	// In production, replace this with actual download code:
	// resp, err := http.Get(url)
	// if err != nil {
	//     return fmt.Errorf("failed to download update: %w", err)
	// }
	// defer resp.Body.Close()

	// Apply the update:
	// err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	// if err != nil {
	//     if rerr := selfupdate.RollbackError(err); rerr != nil {
	//         return fmt.Errorf("update failed, rollback also failed: %w", rerr)
	//     }
	//     return fmt.Errorf("update failed: %w", err)
	// }

	fmt.Println("Update simulation complete!")
	fmt.Println("\nIn production, this would:")
	fmt.Println("1. Download the new binary")
	fmt.Println("2. Verify checksum and signature")
	fmt.Println("3. Replace the current executable")
	fmt.Println("4. Restart the program")

	return nil
}

// VersionInfo returns version information
func VersionInfo() string {
	return version
}
