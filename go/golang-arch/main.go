package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofrs/uuid"
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
	letterRunes  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	currentKeyId string
	keys         = map[string]key{}
	aesSecretKey []byte
)

type key struct {
	key       []byte
	createdAt time.Time
}

type UserClaims struct {
	jwt.StandardClaims       // This struct contains the basic required fields like Issuer, ExpiresAt, Subject etc
	SessionID          int64 // pretty common to include a session id in the custom claims section
}

// It is recommended taht you use a custom Valid method for the UserClaims struct
func (uc *UserClaims) Valid() error {
	if !uc.VerifyExpiresAt(time.Now().UnixNano(), true) {
		return fmt.Errorf("ERROR - JWT token has expired")
	}
	if uc.SessionID == 0 {
		return fmt.Errorf("ERROR - JWT token has invalid session id")
	}
	return nil
}

type person struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	generateHMACKey()
	generateAESKey()
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

	fmt.Println("Starting web server")
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/encode", handleEncode)
	http.HandleFunc("/decode", handleDecode)
	http.HandleFunc("/show_cookie_example", handleShowCookieExample)
	http.HandleFunc("/submit_cookie_example", handleSubmitCookieExample)
	http.ListenAndServe(":8000", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
  <title>Golang Arch for simple webapp depicting usage of various encryption examples</title>
</head>
<body>
  <h1>Home page</h1>
  <p>This is the home page for this application. It is a simple webapp depicting the usage of various ways of encrytion/decryption, hashing, generating digests etc required in authentication mechanism</p>
</body>
</html>
	`
	if _, err := io.WriteString(w, html); err != nil {
		log.Error("ERROR - Could not write html to the response writer")
		return
	}
}

func handleEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
		return
	}
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	var ps []person
	if r.Method != http.MethodPost {
		log.Error("ERROR - This path only handles a POST request")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		log.Error(fmt.Sprintf("ERROR - Could not decode json - %v", err))
	}
	log.Info(ps)
}

func handleShowCookieExample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Error("ERROR - This path only handles a GET request")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c, err := r.Cookie("session")
	if err != nil {
		c = &http.Cookie{}
	}
	html := `
<!DOCTYPE html>
<html>
<head>
  <title>HMAC Example</title>
</head>
<body>
  <h1>Form for HMAC</h1>
<p>Cookie value : ` + c.Value + `</p>
  <form action="/submit_cookie_example" method="post">
    <input type="email" name="email" />
    <input type="submit" />
  </form>
</body>
</html>
	`
	if _, err := io.WriteString(w, html); err != nil {
		log.Error("ERROR - Could not write html to the response writer")
		return
	}
}

func handleSubmitCookieExample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Error("ERROR - This path only handles a POST request")
		http.Redirect(w, r, "/show_cookie_example", http.StatusSeeOther) // http.StatusSeeOther - go to this other location and dont forward form values
		return
	}
	email := r.FormValue("email")
	if len(email) == 0 {
		log.Error("ERROR - The email sent via the form is an empty string")
		http.Redirect(w, r, "/show_cookie_example", http.StatusSeeOther)
		return
	}
	signedCookieValue, err := signMessage([]byte(email))
	if err != nil {
		log.Error("ERROR - Could not create a signature for the value to put in the cookie")
		http.Redirect(w, r, "/show_cookie_example", http.StatusSeeOther)
		return
	}
	cookieVal := fmt.Sprintf("%s|%s", string(signedCookieValue), email)
	c := &http.Cookie{
		Name:  "session",
		Value: cookieVal,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/show_cookie_example", http.StatusSeeOther)
}

// ----------------------------------------------------------------------------------------------------
// ----------------------------------- Base64 encoding and decoding -----------------------------------
// ------------------------- This part is about base64 encoding and decoding --------------------------
// ----------------------------------------------------------------------------------------------------

func b64encode(msg string) string {
	return base64.StdEncoding.EncodeToString([]byte(msg))
}

func b64decode(msg string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return "", fmt.Errorf("Could not base64 decode string because of error - %v", err)
	}
	return string(bytes), nil
}

// ----------------------------------------------------------------------------------------------------
// ---------------------------------- AES encryption and decryption -----------------------------------
// ------------------------ This part is about aes encryption and decryption --------------------------
// --------------------- AES uses a symetric key for encryption and decryption ------------------------
// ----------------------------------------------------------------------------------------------------

// The same password and same initialization vector (or called salt here) needs to be used for encryption and decryption mechanism
// Also, important to remember that the salt size needs to be at least aes.BlockSize in length
// Usually, the encrypted message is sent to the client along with the initialization vector.
// And the password/secret is shared between the server and the client.
// That way the client can use the shared secret and the initialization vector to decode the encrypted message
func aesEncryptDecrypt(key []byte, msg string) (string, error) {
	// AES-128 requires a 16 bit key and AES-256 requires a 32 bit key
	// In this step we are generating the cipher that we will use to encrypt our input message
	b, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Could not generate cipher to encrypt message - %v", err)
	}

	// This is an initialization vector of all zeros. We could have used some random chars as well.
	// However a salt generated from random chars need to be sent back to the client to decrypt.
	// To keep things simple we are using an initialization vector of all zeros.
	iv := make([]byte, aes.BlockSize)

	// Pass the cipher and a salt/initialization vector to cipher.NewCTR
	// to get back a stream where the encrypted message will get written to
	s := cipher.NewCTR(b, iv)

	// Use the stream above and an empty bytes buffer to create a stream writer.
	// The stream writer will take your input message, encrypt it and write it out to the bytes buffer.
	buff := &bytes.Buffer{}
	sw := cipher.StreamWriter{
		S: s,
		W: buff,
	}

	// The stream writer is gonna take your input message and write it out to the bytes buffer mentioned above
	if _, err := sw.Write([]byte(msg)); err != nil {
		return "", fmt.Errorf("Could not encrypt message - %v", err)
	}

	// Finally, get the bytes from the bytes buffer
	out := buff.Bytes()
	return string(out), nil
}

// This is a similar function as the one above except it is different in the sense that it wraps around a writer.
// This can be any writer. A bytes buffer or a http response stream writer.
// So, whatever message you pass to this encrypted stream writer, it takes it, encrypts it and writes it to the destination writer.
// Example use case :
// wtr := &bytes.Buffer{}
// encWriter, err := encryptedWriter(wtr, aesSecretKey)
// _, err := io.WriteString(encWriter, "This is the message i am trying to encrypt")
// if err != nil {
// panic(err)
// }
// encrypted := wtr.String()
// fmt.Println(encrypted)
func encryptedWriter(w io.Writer, key []byte) (io.Writer, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Could not generate cipher to encrypt input - %v", err)
	}

	// Make an initialization vector of all 0 bytes. That way the initialization vector/salt is of 16 bytes but it is the same value always.
	iv := make([]byte, aes.BlockSize)

	s := cipher.NewCTR(b, iv)
	return cipher.StreamWriter{
		S: s,
		W: w,
	}, nil
}

// ----------------------------------------------------------------------------------------------------
// --------------------- Basic auth plaintext password bcrypt hashing and verifying -------------------
// ---------------------------------- This part is about hashing --------------------------------------
// ----------------------------------------------------------------------------------------------------

// How to generate a hashed string from a password using a plain text password for basic auth
func hashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ERROR - Could not generate hashed string from password - %v", err)
	}
	return bytes, nil
}

// How to compare a hashed password and a user provided password to validate a login
func comparePassword(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return fmt.Errorf("ERROR - Passwords did not match - %v", err)
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------
// ----------------------------- HMAC Message signing and verifying signature -------------------------
// ---------------------------------- This part is about signing --------------------------------------
// ----------------------------------------------------------------------------------------------------

func signMessage(msg []byte) ([]byte, error) {
	// hmac.New takes any hashing function that returns a hash.Hash interface
	// and a private key to create the signature.
	// Later this same key is used to verify the signature as well.
	// The length of the private key should match whatever hashing algo you have chosen.
	// For sha512, that size is 64
	h := hmac.New(sha512.New, keys[currentKeyId].key)
	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("ERROR - Encountered error while hashing message - %v", err)
	}
	signature := h.Sum(nil)
	return signature, nil
}

// Server gets the message and the signature. Server sends it to the user as the Bearer token.
// The user sends it back to the server. So you get the original message and the signature again.
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

// ----------------------------------------------------------------------------------------------------
// ----------------------------- JWT token creation and parsing and verifying -------------------------
// ----------------------------------------------------------------------------------------------------

// {JWT Standard fields}.{Your fields}.Signature
// JWT tokens can be signed with HMAC or with RSA/ECDSA. The difference being that the HMAC signing
// requires the private key to be shared because it's the same key that signs a message and validates a signature.
// Whereas with RSA/ECDSA, you can sign a message with a private key but can validate with a public key.
func createToken(c *UserClaims) (string, error) {
	// To create a jwt token from a claim, you need an object that satisfies the SigningMethod interface.
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	// Use a private key for signing and public key for validation if using something like RSA or ECDSA.
	// OR if you are using something like HMAC then use the same shared key for signing and validation.
	return t.SignedString(keys[currentKeyId].key)
}

// This func is used when a user passes back a token with claims and we need to parse it, verify signature, validate it
// and extract the claims information so that we can use information from the claim to deduce other things.
func parseToken(signedToken string) (*UserClaims, error) {
	// ParseWithClaims checks the signature of the token and also checks if the token is valid.
	// The keyFunc passed at the end of jwt.ParseWithClaims takes an unverified token, inspects it's headers
	// for things like the key id (kid) or something similar and returns the key that needs to be used to verify the signature.
	t, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Verify that the algo with which the token is signed is the same as what you are expecting.
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("ERROR - The signing algo of the token and what we expected do not match, so cant use the shared key to verify signature")
		}
		// There can be multiple private keys with which a token could have been encrypted as well.
		// It is possible to pull down the key id information from the header.
		// Using multiple keys allows us to be able to easily rotate keys.
		keyId, ok := t.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("ERROR - Invalid or missing header for key id")
		}
		k, ok := keys[keyId]
		if !ok {
			return nil, fmt.Errorf("ERROR - key id is not valid")
		}
		return k.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("ERROR - Encountered an error when parsing the token - %v", err)
	}
	// The Valid field is populated when ParseWithClaims is called
	if !t.Valid {
		return nil, fmt.Errorf("ERROR - Token is not valid")
	}
	// The tokens' Claims field is a jwt.Claims interface. So, we need to type cast it into our concrete struct type UserClaims.
	claims := t.Claims.(*UserClaims)
	return claims, nil
}

// ----------------------------------------------------------------------------------------------------
// ------------------------------------------- Helper funcs -------------------------------------------
// ----------------------------------------------------------------------------------------------------

func generateHMACKey() error {
	privateKey := []byte(randStringRunes(64))
	id, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("ERROR - could not generate uuid - %v", err)
	}
	keys[id.String()] = key{key: privateKey, createdAt: time.Now()}
	currentKeyId = id.String()
	return nil
}

func generateAESKey() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randStringRunes(64)), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Could not generate aes secret key - %v", err)
	}
	aesSecretKey = hashedPassword[:32]
	return nil
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
