package commands

// import (
// 	"bufio"
// 	"bytes"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"io/fs"
// 	"math/rand"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"github.com/uhppoted/uhppoted-lib/config"
// )
//
// func (cmd *Daemonize) conf(i info, unpacked bool, grules bool) error {
// 	path := cmd.config
//
// 	fmt.Printf("   ... creating '%s'\n", path)
//
// 	// ... get config from existing uhppoted.conf
// 	cfg := config.NewConfig()
// 	if f, err := os.Open(path); err != nil {
// 		if !os.IsNotExist(err) {
// 			return err
// 		}
// 	} else {
// 		err := cfg.Read(f)
// 		f.Close()
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	// ... write back with added tunnel config
// 	var b strings.Builder
// 	if err := cfg.Write(&b); err != nil {
// 		return err
// 	}
//
// 	return toTextFile(b.String(), path)
// }

// func (cmd *Daemonize) genTLSkeys(i info) (bool, error) {
// 	root := cmd.etc
// 	r := bufio.NewReader(os.Stdin)
//
// 	fmt.Println()
// 	fmt.Printf("     Do you want to create TLS keys and certificates (yes/no)? ")
//
// 	text, err := r.ReadString('\n')
// 	fmt.Println()
// 	if err != nil || strings.TrimSpace(text) != "yes" {
// 		return false, nil
// 	}
//
// 	keys, err := genkeys()
// 	if err != nil {
// 		return false, err
// 	} else if keys == nil {
// 		return false, fmt.Errorf("Invalid TLS key set (%v)", keys)
// 	}
//
// 	list := []struct {
// 		item interface{}
// 		file string
// 	}{
// 		{keys.CA.privateKey, "ca.key"},
// 		{keys.CA.certificate, "ca.cert"},
// 		{keys.server.privateKey, "uhppoted.key"},
// 		{keys.server.certificate, "uhppoted.cert"},
// 		{keys.client.privateKey, "client.key"},
// 		{keys.client.certificate, "client.cert"},
// 	}
//
// 	for _, v := range list {
// 		file := filepath.Join(root, v.file)
//
// 		//	if _, err := os.Stat(file); err != nil {
// 		//		if !os.IsNotExist(err) {
// 		//			return false, err
// 		//		} else if err := toTextFile(string(encode(v.item)), file); err != nil {
// 		//			return false, err
// 		//		} else {
// 		//			fmt.Printf("   ... created %v\n", file)
// 		//		}
// 		//	}
//
// 		if err := toTextFile(string(encode(v.item)), file); err != nil {
// 			return false, err
// 		} else {
// 			fmt.Printf("   ... created %v\n", file)
// 		}
// 	}
//
// 	fmt.Println()
// 	fmt.Println("   ** PLEASE MOVE THE ca.key FILE TO A SECURE LOCATION")
// 	fmt.Println()
// 	fmt.Println("   The supplied client.key file can be installed in a browser to support mutual TLS authentication.")
// 	fmt.Println("   It is provided merely as an example and both the client key and certificate should be removed")
// 	fmt.Println("   and replaced by your own keys and certificates.")
// 	fmt.Println()
// 	fmt.Println("   The client.key file is in PEM format - to convert it to a PKCS12 file for importing into Firefox")
// 	fmt.Println("   execute the following command:")
// 	fmt.Println()
// 	fmt.Println("   openssl pkcs12 -export -in client.cert -inkey client.key -certfile ca.cert -out client.p12")
// 	fmt.Println()
// 	fmt.Println("   ** NB: The generated TLS keys and certificates are for TEST USE ONLY and should be replaced with")
// 	fmt.Println("          your own CA certificate and server and client keys and certificates for production use.")
// 	fmt.Println()
//
// 	return true, nil
// }
