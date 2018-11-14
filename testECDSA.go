package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/stone-ch/FabricDevTutorial/ecdsa"
)

func main() {
	// GenerateKey
	c := elliptic.P256()
	r := rand.Reader
	priv, err := ecdsa.GenerateKey(c, r)
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	if !c.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		fmt.Printf("public key invalid: %s", err)
	}

    hashed := []byte("testing")
	r2, s, err := ecdsa.Sign(rand.Reader, priv, hashed)
	if err != nil {
		return
	}
    fmt.Printf("r2: %s, s: %s\n", r2, s)

    verify := ecdsa.Verify(&priv.PublicKey, hashed, r2, s)
    fmt.Printf("verify result is: %s \n", verify)

}
