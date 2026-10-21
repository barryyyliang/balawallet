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

package sol

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"

	"balawallet/accounts"

	"filippo.io/edwards25519"
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
)

const (
	// Number of bytes in a pubkey.
	PublicKeyLength = 32
	// Maximum length of derived pubkey seed.
	MaxSeedLength = 32
	// Maximum number of seeds.
	MaxSeeds = 16
	// Number of bytes in a signature.
	SignatureLength = 64

	// Number of bytes in a private key.
	PrivateKeyLength = ed25519.PrivateKeySize

	// // Maximum string length of a base58 encoded pubkey.
	// MaxBase58Length = 44
)

var _ accounts.IChainAdapter = (*SolChainAdapter)(nil)
var _ accounts.Key = (*SolKey)(nil)
var _ accounts.Account = (*SolAccount)(nil)

type SolChainAdapter struct {
}

func (c *SolChainAdapter) NewAccount(ad string, url accounts.URL) accounts.Account {
	re := SolAccount{}
	re.Init(ad, url)
	return &re
}
func (c *SolChainAdapter) NewNullAccount() accounts.Account {
	re := SolAccount{}
	return &re
}
func (c *SolChainAdapter) InitKey(id uuid.UUID, pik interface{}) accounts.Key {
	prik, ok := pik.(*ed25519.PrivateKey)
	if ok {
		return &SolKey{
			Id:         id,
			address:    prik.Public().(ed25519.PublicKey),
			privateKey: prik,
		}
	} else {
		fmt.Println("key type:" + reflect.TypeOf(pik).String())
		fmt.Println("InitKey error")
	}

	return nil
}
func (c *SolChainAdapter) NewKey(rand io.Reader) (accounts.Key, error) {
	pub, priv, err := ed25519.GenerateKey(rand)
	if err != nil {
		return nil, err
	}
	var publicKey ed25519.PublicKey
	copy(publicKey[:], pub)

	_, err = validatePrivateKey(priv)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}
	key := &SolKey{
		Id:         id,
		address:    priv.Public().(ed25519.PublicKey),
		privateKey: &priv,
	}
	return key, nil
}

func (c *SolChainAdapter) BytesToPrivateKey(b []byte) (interface{}, error) {
	var tmp = ed25519.PrivateKey(b)
	return &tmp, nil
}

type SolKey struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	address ed25519.PublicKey
	// we only store privkey as pubkey/address can be derived from it
	// privkey in this struct is always in plaintext
	privateKey *ed25519.PrivateKey
}

func (k *SolKey) SetID(id uuid.UUID) {
	k.Id = id
}

func (k *SolKey) GetID() string {
	return k.Id.String()
}

func (k *SolKey) Is(add string) bool {
	val, err := base58.Decode(add)
	if err != nil {
		return false
	}

	if len(val) != PublicKeyLength {
		fmt.Errorf("invalid length, expected %v, got %d", PublicKeyLength, len(val))
		return false
	}

	var out ed25519.PublicKey = val

	return out.Equal(k.address)
}

func (k *SolKey) GetByteAddress() []byte {
	return k.address
}
func (k *SolKey) GetStringAddress() string {
	tmp := base58.Encode(k.address[:])
	fmt.Println("sol GetStringAddress:" + tmp)
	return tmp
}

func (k *SolKey) SignHash(hash []byte) ([]byte, error) {

	return nil, fmt.Errorf("not support")
}
func (k *SolKey) SignMessage(message []byte) ([]byte, error) {

	signData, err := k.privateKey.Sign(rand.Reader, message, crypto.Hash(0))
	if err != nil {
		fmt.Println("sol sign hash error:" + err.Error())
	}

	return signData, err
}

// func (k *Key) GetPrivateKey() []byte {
// 	return crypto.FromECDSA(k.privateKey)
// }

func (k *SolKey) ZeroKey() {
	clear(*k.privateKey)
}

func (k *SolKey) GetPrivateKey() interface{} {
	return k.privateKey
}
func (k *SolKey) GetByte() []byte {

	return (*k.privateKey)[:]
}

func (k *SolKey) MarshalJSON() (j []byte, err error) {
	jStruct := accounts.PlainKeyJSON{
		Address:    base58.Encode(k.address[:]),
		PrivateKey: base58.Encode((*k.privateKey)[:]),
		Id:         k.Id.String(),
		Version:    accounts.Version,
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *SolKey) UnmarshalJSON(j []byte) (err error) {
	keyJSON := new(accounts.PlainKeyJSON)
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
	addr, err := base58.Decode(keyJSON.Address)
	if err != nil {
		return err
	}
	privkey, err := base58.Decode(keyJSON.PrivateKey)
	if err != nil {
		return err
	}

	_, err = validatePrivateKey(privkey)
	if err != nil {
		return err
	}

	if ed25519.PublicKey(addr).Equal(ed25519.PrivateKey(privkey).Public().(ed25519.PublicKey)) {
		fmt.Println("address :" + keyJSON.Address + " != " + base58.Encode(ed25519.PrivateKey(privkey).Public().(ed25519.PublicKey)))
		return nil
	}
	tmp := ed25519.PrivateKey(privkey)

	k.privateKey = &tmp
	k.address = k.privateKey.Public().(ed25519.PublicKey)

	return nil
}

type SolAccount struct {
	address ed25519.PublicKey `json:"address"` // Ethereum account address derived from the key
	url     accounts.URL      `json:"url"`     // Optional resource locator within a backend
}

func (a *SolAccount) GetAddress() string {
	return base58.Encode(a.address)
}
func (a *SolAccount) Init(ad string, ul accounts.URL) {
	tmp, err := base58.Decode(ad)
	if err != nil {
		fmt.Println("init solaccount error:" + err.Error())
		return
	}
	a.address = ed25519.PublicKey(tmp)
	a.url = ul
}
func (a *SolAccount) GetUrl() accounts.URL {
	return a.url
}
func (a *SolAccount) AddressIsNull() bool {
	return a.address.Equal(ed25519.PublicKey{})
}
func (a *SolAccount) Cmp(ac accounts.Account) bool {
	tmp := base58.Encode(a.address)

	return tmp == ac.GetAddress()
}

func validatePrivateKey(b []byte) (bool, error) {
	if len(b) != ed25519.PrivateKeySize {
		return false, fmt.Errorf("invalid private key size, expected %v, got %d", ed25519.PrivateKeySize, len(b))
	}
	// check if the public key is on the ed25519 curve
	pub := ed25519.PrivateKey(b).Public().(ed25519.PublicKey)
	if !isOnCurve(pub) {
		return false, errors.New("the corresponding public key is NOT on the ed25519 curve")
	}
	return true, nil
}

// Check if the provided `b` is on the ed25519 curve.
func isOnCurve(b []byte) bool {
	if len(b) != ed25519.PublicKeySize {
		return false
	}
	_, err := new(edwards25519.Point).SetBytes(b)
	isOnCurve := err == nil
	return isOnCurve
}

func NewKeyFromBase58(p58 string) (accounts.Key, error) {

	res, err := base58.Decode(p58)
	if err != nil {
		return nil, err
	}

	var priv = ed25519.PrivateKey(res)

	_, err = validatePrivateKey(priv)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}
	key := &SolKey{
		Id:         id,
		address:    priv.Public().(ed25519.PublicKey),
		privateKey: &priv,
	}
	return key, nil
}
