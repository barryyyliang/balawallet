// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethkey

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"balawallet/accounts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

const (
	version = 3
)

type EthChainAdapter struct {
}

func (c *EthChainAdapter) NewAccount(ad string, url accounts.URL) accounts.Account {
	re := EthAccount{}
	re.Init(ad, url)
	return &re
}
func (c *EthChainAdapter) NewKey(id uuid.UUID, pik interface{}) accounts.Key {

	if privatekey, ok := pik.(*ecdsa.PrivateKey); ok {
		return &EthKey{
			Id:         id,
			address:    crypto.PubkeyToAddress(privatekey.PublicKey),
			privateKey: privatekey,
		}
	}
	return nil
}
func (c *EthChainAdapter) BytesToPrivateKey(b []byte) (interface{}, error) {
	return crypto.ToECDSA(b)
}

type EthKey struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	address common.Address
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	privateKey *ecdsa.PrivateKey
}

type plainKeyJSON struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
	Id         string `json:"id"`
	Version    int    `json:"version"`
}

type encryptedKeyJSONV3 struct {
	Address string     `json:"address"`
	Crypto  CryptoJSON `json:"crypto"`
	Id      string     `json:"id"`
	Version int        `json:"version"`
}

type encryptedKeyJSONV1 struct {
	Address string     `json:"address"`
	Crypto  CryptoJSON `json:"crypto"`
	Id      string     `json:"id"`
	Version string     `json:"version"`
}

type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

func (k *EthKey) Is(add string) bool {
	tmp := common.HexToAddress(add)
	return tmp == k.address
}

func (k *EthKey) GetByteAddress() []byte {
	return k.address.Bytes()
}
func (k *EthKey) GetStringAddress() string {
	return k.address.String()
}

func (k *EthKey) SinHash(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, k.privateKey)
}

// func (k *Key) GetPrivateKey() []byte {
// 	return crypto.FromECDSA(k.privateKey)
// }

func (k *EthKey) ZeroKey() {
	if k.privateKey != nil {
		b := k.privateKey.D.Bits()
		clear(b)
	}
}

func (k *EthKey) GetPrivateKey() interface{} {
	return k.privateKey
}
func (k *EthKey) GetBigPrivate() *big.Int {
	if k.privateKey != nil {
		return k.privateKey.D
	}
	return nil
}

func (k *EthKey) MarshalJSON() (j []byte, err error) {
	jStruct := plainKeyJSON{
		hex.EncodeToString(k.address[:]),
		hex.EncodeToString(crypto.FromECDSA(k.privateKey)),
		k.Id.String(),
		version,
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *EthKey) UnmarshalJSON(j []byte) (err error) {
	keyJSON := new(plainKeyJSON)
	err = json.Unmarshal(j, &keyJSON)
	if err != nil {
		return err
	}

	u := new(uuid.UUID)
	*u, err = uuid.Parse(keyJSON.Id)
	if err != nil {
		return err
	}
	k.Id = *u
	addr, err := hex.DecodeString(keyJSON.Address)
	if err != nil {
		return err
	}
	privkey, err := crypto.HexToECDSA(keyJSON.PrivateKey)
	if err != nil {
		return err
	}

	k.address = common.BytesToAddress(addr)
	k.privateKey = privkey

	return nil
}

func newKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) accounts.Key {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}
	key := &EthKey{
		Id:         id,
		address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		privateKey: privateKeyECDSA,
	}
	return key
}

// NewKeyForDirectICAP generates a key whose address fits into < 155 bits so it can fit
// into the Direct ICAP spec. for simplicity and easier compatibility with other libs, we
// retry until the first byte is 0.
func NewKeyForDirectICAP(rand io.Reader) accounts.Key {
	randBytes := make([]byte, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic("key generation: could not read from random source: " + err.Error())
	}
	reader := bytes.NewReader(randBytes)
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), reader)
	if err != nil {
		panic("key generation: ecdsa.GenerateKey failed: " + err.Error())
	}
	key := newKeyFromECDSA(privateKeyECDSA)
	if !strings.HasPrefix(key.GetStringAddress(), "0x00") {
		return NewKeyForDirectICAP(rand)
	}
	return key
}

func newKey(rand io.Reader) (accounts.Key, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return newKeyFromECDSA(privateKeyECDSA), nil
}

func writeTemporaryKeyFile(file string, content []byte) (string, error) {
	// Create the keystore directory with appropriate permissions
	// in case it is not present yet.
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return "", err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := os.CreateTemp(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return "", err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return "", err
	}
	f.Close()
	return f.Name(), nil
}

func writeKeyFile(file string, content []byte) error {
	name, err := writeTemporaryKeyFile(file, content)
	if err != nil {
		return err
	}
	return os.Rename(name, file)
}

// keyFileName implements the naming convention for keyfiles:
// UTC--<created_at UTC ISO8601>-<address hex>
func keyFileName(keyAddr []byte) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), hex.EncodeToString(keyAddr[:]))
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

type EthAccount struct {
	address common.Address `json:"address"` // Ethereum account address derived from the key
	url     accounts.URL   `json:"url"`     // Optional resource locator within a backend
}

func (a *EthAccount) GetAddress() string {
	return a.address.String()
}
func (a *EthAccount) Init(ad string, ul accounts.URL) {
	a.address = common.HexToAddress(ad)
	a.url = ul
}
func (a *EthAccount) GetUrl() accounts.URL {
	return a.url
}
func (a *EthAccount) AddressIsNull() bool {
	return a.address == common.Address{}
}
