package main

func lib_go() []byte {
	return []byte(`
package main

import (
	"errors"
	"io"
	"io/ioutil"
	"fmt"
	"os"
	"crypto/rand"
	"crypto/sha1"
	"code.google.com/p/go.crypto/pbkdf2"
	"github.com/cryptobox/gocryptobox/strongbox"
	"github.com/cookieo9/resources-go/v2/resources" 
)

var verbosity = 0

func log_level(level int) {
	if level >= 0 {
		verbosity = level
	}
}

func log(level int, format string, args ...interface{}) {
	if level <= verbosity {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func err_exit(msg string, args ...interface{}) {
	log(0, "fatal error: ")
	log(0, msg, args...)
	log(0, "\n")
	os.Exit(1)
}

const PBKDF_ITERS 		= 99999
const KLEN 				= strongbox.KeySize
const PBKDF_SALT_LEN 	= KLEN

func gen_salt() (salt []byte, err error) {
	salt = make([]byte, PBKDF_SALT_LEN)
	n, err := io.ReadFull(rand.Reader, salt)
	if err != nil && n != PBKDF_SALT_LEN {
		err = errors.New("failed to generate salt bytes")
	}
	return
}

func gen_key(password []byte) (salt []byte, key []byte, err error) {
	salt, err = gen_salt()
	if err != nil {
		return
	}

	key = get_key(password, salt)
	return
}

func get_key(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, PBKDF_ITERS, KLEN, sha1.New)
}

/* included here so that go understands that resources-go is
 * a dependency (toad.go will use this function when summoned) */
func get_asset(name string) (value []byte, err error) {
	r, err := resources.DefaultBundle.Find(fmt.Sprintf("assets/%s", name))
	if err != nil {
		return
	}

	rdr, err := r.Open()
	if err != nil {
		return
	}

	value, err = ioutil.ReadAll(rdr)
	return
}
`)
}
