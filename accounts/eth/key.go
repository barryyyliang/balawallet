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

package eth

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"

	"balawallet/accounts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

type EthChainAdapter struct {
}

func (c *EthChainAdapter) NewAccount(ad string, url accounts.URL) accounts.Account {
	re := EthAccount{}
	re.Init(ad, url)
	return &re
}
func (c *EthChainAdapter) NewNullAccount() accounts.Account {
	re := EthAccount{}
	return &re
}
func (c *EthChainAdapter) InitKey(id uuid.UUID, pik interface{}) accounts.Key {

	if privatekey, ok := pik.(*ecdsa.PrivateKey); ok {
		return &EthKey{
			Id:         id,
			address:    crypto.PubkeyToAddress(privatekey.PublicKey),
			privateKey: privatekey,
		}
	}
	return nil
}
func (c *EthChainAdapter) NewKey(rand io.Reader) (accounts.Key, error) {
	return newKey(rand)
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

func (k *EthKey) SetID(id uuid.UUID) {
	k.Id = id
}

func (k *EthKey) GetID() string {
	return k.Id.String()
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
	jStruct := accounts.PlainKeyJSON{
		hex.EncodeToString(k.address[:]),
		hex.EncodeToString(crypto.FromECDSA(k.privateKey)),
		k.Id.String(),
		accounts.Version,
	}
	j, err = json.Marshal(jStruct)
	return j, err
}

func (k *EthKey) UnmarshalJSON(j []byte) (err error) {
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

func NewKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) accounts.Key {
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
	key := NewKeyFromECDSA(privateKeyECDSA)
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
	return NewKeyFromECDSA(privateKeyECDSA), nil
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
func (a *EthAccount) Cmp(ac accounts.Account) bool {
	return a.address.String() == ac.GetAddress()
}
