# go-selfupdate

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

A Go library that provides secure, self-updating capabilities for single-file executables.

## Features

- ✅ **Cross-platform support** - Works on Linux, macOS, and Windows
- ✅ **Binary patching** - Reduce download sizes with BSDiff patches
- ✅ **Checksum verification** - Validate updates using SHA256 (or other hash algorithms)
- ✅ **Code signing** - Support for RSA and ECDSA signatures
- ✅ **Atomic updates** - Safe update process with automatic rollback on failure
- ✅ **No external dependencies** - Uses only Go standard library (except internal binarydiff)

## Installation

```bash
go get github.com/tekintian/go-selfupdate
```

## Quick Start

### Basic Update

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/tekintian/go-selfupdate"
)

func doUpdate(url string) error {
    // Download the new binary
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Apply the update
    err = selfupdate.Apply(resp.Body, selfupdate.Options{})
    if err != nil {
        if rerr := selfupdate.RollbackError(err); rerr != nil {
            fmt.Printf("Failed to rollback from bad update: %v\n", rerr)
        }
        return err
    }
    
    fmt.Println("Update successful!")
    return nil
}
```

### Binary Patching

For large binaries, shipping only the difference can significantly reduce download size:

```go
func updateWithPatch(patch io.Reader) error {
    return selfupdate.Apply(patch, selfupdate.Options{
        Patcher: selfupdate.NewBSDiffPatcher(),
    })
}
```

### Checksum Verification

Ensure the update hasn't been corrupted during transmission:

```go
func updateWithChecksum(binary io.Reader, hexChecksum string) error {
    checksum, err := hex.DecodeString(hexChecksum)
    if err != nil {
        return err
    }
    
    return selfupdate.Apply(binary, selfupdate.Options{
        Checksum: checksum,
        Hash:     crypto.SHA256,
    })
}
```

### Cryptographic Signature Verification

Protect against malicious updates by verifying signatures:

```go
var publicKeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEtrVmBxQvheRArXjg2vG1xIprWGuCyESx
MMY8pjmjepSy2kuz+nl9aFLqmr+rDNdYvEBqQaZrYMc6k29gjvoQnQ==
-----END PUBLIC KEY-----`)

func verifiedUpdate(binary io.Reader, hexChecksum, hexSignature string) error {
    checksum, _ := hex.DecodeString(hexChecksum)
    signature, _ := hex.DecodeString(hexSignature)
    
    opts := selfupdate.Options{
        Checksum:  checksum,
        Signature: signature,
        Hash:      crypto.SHA256,
        Verifier:  selfupdate.NewECDSAVerifier(),
    }
    
    if err := opts.SetPublicKeyPEM(publicKeyPEM); err != nil {
        return err
    }
    
    return selfupdate.Apply(binary, opts)
}
```

## API Reference

### Options

```go
type Options struct {
    // TargetPath defines the path to the file to update.
    // Empty string means 'the executable file of the running program'.
    TargetPath string

    // Create TargetPath replacement with this file mode. If zero, defaults to 0755.
    TargetMode os.FileMode

    // Checksum of the new binary to verify against.
    // If nil, no checksum verification is done.
    Checksum []byte

    // Public key to use for signature verification.
    // If nil, no signature verification is done.
    PublicKey crypto.PublicKey

    // Signature to verify the updated file.
    // If nil, no signature verification is done.
    Signature []byte

    // Pluggable signature verification algorithm.
    // If nil, ECDSA is used.
    Verifier Verifier

    // Hash function to generate the checksum.
    // If not set, SHA256 is used.
    Hash crypto.Hash

    // If nil, treat the update as a complete replacement.
    // If non-nil, treat the update as a patch.
    Patcher Patcher

    // Store the old executable file at this path after a successful update.
    // Empty string means the old executable file will be removed.
    OldSavePath string
}
```

### Functions

- `Apply(io.Reader, Options) error` - Apply an update to the target executable
- `RollbackError(error) error` - Extract the rollback error from an update error
- `NewECDSAVerifier() Verifier` - Create an ECDSA signature verifier
- `NewRSAVerifier() Verifier` - Create an RSA signature verifier
- `NewBSDiffPatcher() Patcher` - Create a BSDiff binary patcher

## Security Considerations

**Always verify signatures and checksums** when applying updates. Never apply an update without verifying its authenticity and integrity.

Recommended security practices:

1. Use cryptographic signatures (ECDSA or RSA)
2. Verify checksums before applying
3. Fetch updates over HTTPS
4. Keep private keys secure
5. Use strong hash algorithms (SHA256 or better)

## Building Single-File Binaries

This library only works with single-file executables. If your application has static assets (HTML, CSS, certificates, etc.), embed them into the binary using tools like:

- [embed](https://pkg.go.dev/embed) (Go 1.16+, built-in)
- [go-bindata](https://github.com/jteeuwen/go-bindata)
- [statik](https://github.com/rakyll/statik)

## Update Process

The `Apply` function performs the following steps:

1. **Validate options** - Check signature and key requirements
2. **Apply patch (optional)** - If a patcher is provided, apply it
3. **Verify checksum (optional)** - If provided, verify the checksum
4. **Verify signature (optional)** - If provided, verify the signature
5. **Create new file** - Write updated content to `.target.new`
6. **Rename old file** - Move target to `.target.old`
7. **Rename new file** - Move `.target.new` to target
8. **Cleanup** - Remove `.target.old` (or hide on Windows)

If any step fails after moving the old file, the function attempts to rollback by restoring the old file.

## Limitations

This library intentionally focuses on the core update mechanism. The following are **out of scope**:

- Checking for available updates
- Deciding which update to apply
- Update channels (stable, beta, etc.)
- Hosting update files
- Update metrics and analytics

For complete update management solutions, consider:

- [TUF (The Update Framework)](https://github.com/theupdateframework/go-tuf)
- [Equinox](https://equinox.io)

## Requirements

- Go 1.20 or higher
- Single-file executable binaries

## License

Apache License 2.0

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

