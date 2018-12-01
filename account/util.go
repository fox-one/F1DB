package account

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"time"

	"github.com/fox-one/foxgo/account"
	"github.com/fox-one/foxgo/request"
	uuid "github.com/satori/go.uuid"
)

type Pin struct {
	code string
	pk   string
}

var EmptyPin = Pin{}

func parsePublicKey(pk string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pk))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub.(*rsa.PublicKey), nil
}

func GetFoxPublicKey(ctx context.Context) (string, error) {
	return account.GetPublicKey(ctx)
}

func NewPin(code, pk string) Pin {
	return Pin{code: code, pk: pk}
}

// UpdatePin: Use EmptyPin as pin when set pin first time
func UpdatePin(ctx context.Context, token string, pin, newPin Pin) error {
	_, err := request.Put(ctx, "account/pin", WithToken(token), WithPin(pin), WithNewPin(newPin))
	return err
}

func WithToken(token string) request.BuildParamFunc {
	return func(p request.Param) error {
		if len(token) == 0 {
			return errors.New("empty token")
		}

		p.SetHeader("Authorization", "Bearer "+token)
		return nil
	}
}

func WithPin(pin Pin) request.BuildParamFunc {
	return func(p request.Param) error {
		if pin == EmptyPin {
			// do nothing
			return nil
		}

		key, err := parsePublicKey(pin.pk)
		if err != nil {
			return err
		}

		h := md5.New()
		h.Write([]byte("fox." + pin.code))
		hashedPin := hex.EncodeToString(h.Sum(nil))

		ts := time.Now().Unix()
		nonce := uuid.Must(uuid.NewV4()).String()

		payload, _ := json.Marshal(map[string]interface{}{
			"hp": hashedPin,
			"t":  ts,
			"n":  nonce,
		})

		hash := sha1.New()
		random := rand.Reader

		data, err := rsa.EncryptOAEP(hash, random, key, payload, nil)
		if err != nil {
			return err
		}

		pinToken := base64.StdEncoding.EncodeToString(data)
		p.SetHeader("fox-client-pin", pinToken)

		return nil
	}
}

func WithNewPin(pin Pin) request.BuildParamFunc {
	return func(p request.Param) error {
		key, err := parsePublicKey(pin.pk)
		if err != nil {
			return err
		}

		h := md5.New()
		h.Write([]byte("fox." + pin.code))
		hashedPin := hex.EncodeToString(h.Sum(nil))

		hash := sha1.New()
		random := rand.Reader

		data, err := rsa.EncryptOAEP(hash, random, key, []byte(hashedPin), nil)
		if err != nil {
			return err
		}

		pinToken := base64.StdEncoding.EncodeToString(data)
		p.SetValue("newPinToken", pinToken)
		p.SetValue("pinType", 1)

		return nil
	}
}
