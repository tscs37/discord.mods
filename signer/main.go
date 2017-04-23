package main

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"github.com/chzyer/readline"
	"go.rls.moe/misc/discord.mods/common"
	_ "go.rls.moe/misc/discord.mods/common/osmode"
	"golang.org/x/crypto/ed25519"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	flagFileToSign = flag.String("file", "mymod.dmod", "File to sign with your key")
	flagKeyFile    = flag.String("key-file", common.SignerKey, "Key with which to sign, if not present it is generated")
	flagNoPasswd   = flag.Bool("no-passwd", false, "Disable Password Entry for keys without password")
)

func main() {
	flag.Parse()
	home, err := common.GetBase()
	if err != nil {
		return
	}
	signerKeyPath := filepath.Join(home, *flagKeyFile)
	if _, err := os.Stat(signerKeyPath); err != nil {
		if err = genKey(signerKeyPath); err != nil {
			log.Fatal(err)
			return
		}
	}

	signerKeyDat, err := ioutil.ReadFile(signerKeyPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	pass, err := readline.Password("Signer Key Password: ")
	if err != nil {
		log.Fatal(err)
		return
	}

	// TODO
	pass = pass

	rawKey, err := base64.RawStdEncoding.DecodeString(string(signerKeyDat))
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(rawKey) != ed25519.PrivateKeySize {
		log.Fatal("Key size mismatch")
		return
	}

	prvKey := ed25519.PrivateKey(rawKey)

	log.Printf("Loaded Key: %x", []byte(prvKey.Public().(ed25519.PublicKey)[:8]))

	msgBytes, err := common.GetFile(*flagFileToSign)
	if err != nil {
		log.Fatal(err)
		return
	}

	signature, err := prvKey.Sign(rand.Reader, msgBytes, crypto.Hash(0))
	if err != nil {
		log.Fatal(err)
		return
	}

	encodedSig := base64.RawStdEncoding.EncodeToString(signature)
	err = common.WriteFile(*flagFileToSign+".sig", []byte(encodedSig), 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Signature: %x", signature[:8])
	return
}

func genKey(path string) error {
	_, prv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(
		base64.RawStdEncoding.EncodeToString(prv)), 0600)
}
