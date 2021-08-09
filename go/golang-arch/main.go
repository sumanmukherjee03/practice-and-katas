package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/form3tech-oss/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	slurpOutputType  = "slurp"
	streamOutputType = "stream"
)

var (
	allowedOutTypes = map[string]bool{
		slurpOutputType:  true,
		streamOutputType: true,
	}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	privateKey  []byte
)

type person struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFormatter(&log.JSONFormatter{})
	privateKey = []byte(randStringRunes(64))
	signals := make(chan os.Signal)

	// The Notify function will pass the incoming signals that you provided, in this case os.Interrupt
	// to the signals channel, which you can then read from to customize how you handle OS signals
	// This comes from CTRL+c or kill -2 <pid>
	signal.Notify(signals, os.Interrupt)

	// Process the OS interrupt signal in a goroutine
	go func() {
		s := <-signals
		errorf("Received OS signal - %v", s)
	}()

	http.HandleFunc("/encode", handleEncode)
	http.HandleFunc("/decode", handleDecode)
	http.ListenAndServe(":8080", nil)
}

func handleEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Error("ERROR - This path only handles a GET request")
		return
	}
	p1 := &person{
		FirstName: "John",
		LastName:  "Doe",
	}
	p2 := &person{
		FirstName: "Jane",
		LastName:  "Doe",
	}
	ps := []*person{p1, p2}
	err := json.NewEncoder(w).Encode(ps)
	if err != nil {
		log.Error(fmt.Sprintf("ERROR - Could not encode into json - %v", err))
	}
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	var ps []person
	if r.Method != "POST" {
		log.Error("ERROR - This path only handles a POST request")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		log.Error(fmt.Sprintf("ERROR - Could not decode json - %v", err))
	}
	log.Info(ps)
}

// How to generate a hashed string from a password using a plain text password
func hashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ERROR - Could not generate hashed string from password - %v", err)
	}
	return bytes, nil
}

func comparePassword(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return fmt.Errorf("ERROR - Passwords did not match - %v", err)
	}
	return nil
}

func signMessage(msg []byte) ([]byte, error) {
	// hmac.New takes any hashing function that returns a hash.Hash interface
	// and a private key to create the signature.
	// Later this same key is used to verify the signature as well.
	// The length of the private key should match whatever hashing algo you have chosen.
	// For sha512, that size is 64
	h := hmac.New(sha512.New, privateKey)
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("ERROR - Encountered error while hashing message - %v", err)
	}
	signature := h.Sum(nil)
	return signature, nil
}

// You get the message and the signature. Send it across to the user as the Bearer token.
// The user sends it back to you. So you get the original message and the signature again.
// Then you compare the message and the signature to verify if the message has been tampered with.
func checkSig(msg, signature []byte) error {
	newSignature, err := signMessage(msg)
	if err != nil {
		return fmt.Errorf("ERROR - Encountered error while comparing signed message with signature - %v", err)
	}
	if !hmac.Equal(newSignature, signature) {
		return fmt.Errorf("ERROR - Signed message and signature do not match - %v", err)
	}
	return nil
}

// {JWT Standard fields}.{Your fields}.Signature
// JWT tokens can be signed with HMAC or with RSA/ECDSA. The difference being that the HMAC signing
// requires the private key to be shared because it's the same key that signs a message and validates a signature.
// Whereas with RSA/ECDSA, you can sign a message with a private key but can validate with a public key.
type UserClaims struct {
	jwt.StandardClaims
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(2)
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
