package pass

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/TechMaster/core/logger"
	"github.com/TechMaster/eris"
	"golang.org/x/crypto/argon2"
)

type Argon2id struct {
	params Params
}

type Params struct {
	// The amount of memory used by the algorithm (in kibibytes).
	Memory uint32

	// The number of iterations over the memory.
	Iterations uint32

	// The number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// Length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32

	// Length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

var (
	// ErrInvalidHash in returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrInvalidHash = eris.New("argon2id: hash is not in the correct format")

	// ErrIncompatibleVersion in returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = eris.New("argon2id: incompatible version of argon2")
)

/*
Băm password sử dụng DefaultParams
*/
func (argon Argon2id) Hash(password string) (hash string) {
	salt, err := generateRandomBytes(argon.params.SaltLength)
	if err != nil {
		logger.Log2(err)
		return ""
	}

	key := argon2.IDKey([]byte(password), salt, argon.params.Iterations, argon.params.Memory, argon.params.Parallelism, argon.params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon.params.Memory, argon.params.Iterations, argon.params.Parallelism, b64Salt, b64Key)
	return hash
}

/*
password: raw input password
hash: password được mã hoá lưu trong CSDL
*/
func (argon Argon2id) Compare(password string, hash string) bool {
	params, salt, key, err := DecodeHash(hash)
	if err != nil {
		logger.Log2(err)
		return false
	}

	otherKey := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))

	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return false
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 { //Password đúng
		return true
	}
	return false
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, eris.NewFrom(err)
	}

	return b, nil
}

// DecodeHash expects a hash created from this package, and parses it to return the params used to
// create it, as well as the salt and key (password hash).
func DecodeHash(hash string) (params *Params, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params = &Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, eris.NewFrom(err)
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, eris.NewFrom(err)
	}
	params.SaltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, eris.NewFrom(err)
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
