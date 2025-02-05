// Copyright 2017 The go-ethereum Authors
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

// Package keystore implements encrypted storage of secp256k1 private keys.
//
// Keys are stored as encrypted JSON files according to the Web3 Secret Storage specification.
// See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.
package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"

	"balawallet/accounts"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

var (
	ErrLocked  = accounts.NewAuthNeededError("password or unlock")
	ErrNoMatch = errors.New("no key for given address or file")
	ErrDecrypt = errors.New("could not decrypt key with given password")

	// ErrAccountAlreadyExists is returned if an account attempted to import is
	// already present in the keystore.
	ErrAccountAlreadyExists = errors.New("account already exists")
)

//var chainAdapter accounts.IChainAdapter = &eth.EthChainAdapter{}

// KeyStoreType is the reflect type of a keystore backend.
var KeyStoreType = reflect.TypeOf(&KeyStore{})

// Maximum time between wallet refreshes (if filesystem notifications don't work).
const walletRefreshCycle = 3 * time.Second

type IkeyStore interface {
	// Loads and decrypts the key from disk.
	GetKey(addr string, filename string, auth string) (accounts.Key, error)
	// Writes and encrypts the key.
	StoreKey(filename string, k accounts.Key, auth string) error
	// Joins filename with the key directory unless it is already absolute.
	JoinPath(filename string) string
}

// KeyStore manages a key storage directory on disk.
type KeyStore struct {
	storage  IkeyStore            // Storage backend, might be cleartext or encrypted
	cache    *accountCache        // In-memory account cache over the filesystem storage
	changes  chan struct{}        // Channel receiving change notifications from the cache
	unlocked map[string]*unlocked // Currently unlocked account (decrypted private keys)

	wallets     []accounts.Wallet       // Wallet wrappers around the individual key files
	updateFeed  event.Feed              // Event feed to notify wallet additions/removals
	updateScope event.SubscriptionScope // Subscription scope tracking current live listeners
	updating    bool                    // Whether the event notification loop is running

	mu           sync.RWMutex
	importMu     sync.Mutex // Import Mutex locks the import to prevent two insertions from racing
	chainAdapter accounts.IChainAdapter
}

type unlocked struct {
	ukey  accounts.Key
	abort chan struct{}
}

// NewKeyStore creates a keystore for the given directory.
func NewKeyStore(keydir string, scryptN, scryptP int, Adapter accounts.IChainAdapter) *KeyStore {
	keydir, _ = filepath.Abs(keydir)
	ks := &KeyStore{storage: &keyStorePassphrase{keydir, scryptN, scryptP, false, Adapter}, chainAdapter: Adapter}
	ks.init(keydir)
	return ks
}

func (ks *KeyStore) init(keydir string) {
	// Lock the mutex since the account cache might call back with events
	ks.mu.Lock()
	defer ks.mu.Unlock()

	// Initialize the set of unlocked keys and the account cache
	ks.unlocked = make(map[string]*unlocked)
	ks.cache, ks.changes = newAccountCache(keydir, ks.chainAdapter)

	// TODO: In order for this finalizer to work, there must be no references
	// to ks. addressCache doesn't keep a reference but unlocked keys do,
	// so the finalizer will not trigger until all timed unlocks have expired.
	runtime.SetFinalizer(ks, func(m *KeyStore) {
		m.cache.close()
	})
	// Create the initial list of wallets from the cache
	accs := ks.cache.accounts()
	ks.wallets = make([]accounts.Wallet, len(accs))
	for i := 0; i < len(accs); i++ {
		ks.wallets[i] = &keystoreWallet{account: accs[i], keystore: ks}
	}
}

// Wallets implements accounts.Backend, returning all single-key wallets from the
// keystore directory.
func (ks *KeyStore) Wallets() []accounts.Wallet {
	// Make sure the list of wallets is in sync with the account cache
	ks.refreshWallets()

	ks.mu.RLock()
	defer ks.mu.RUnlock()

	cpy := make([]accounts.Wallet, len(ks.wallets))
	copy(cpy, ks.wallets)
	return cpy
}

// refreshWallets retrieves the current account list and based on that does any
// necessary wallet refreshes.
func (ks *KeyStore) refreshWallets() {
	// Retrieve the current list of accounts
	ks.mu.Lock()
	accs := ks.cache.accounts()

	// Transform the current list of wallets into the new one
	var (
		wallets = make([]accounts.Wallet, 0, len(accs))
		events  []accounts.WalletEvent
	)

	for _, account := range accs {
		// Drop wallets while they were in front of the next account
		for len(ks.wallets) > 0 && ks.wallets[0].URL().Cmp(account.GetUrl()) < 0 {
			events = append(events, accounts.WalletEvent{Wallet: ks.wallets[0], Kind: accounts.WalletDropped})
			ks.wallets = ks.wallets[1:]
		}
		// If there are no more wallets or the account is before the next, wrap new wallet
		if len(ks.wallets) == 0 || ks.wallets[0].URL().Cmp(account.GetUrl()) > 0 {
			wallet := &keystoreWallet{account: account, keystore: ks}

			events = append(events, accounts.WalletEvent{Wallet: wallet, Kind: accounts.WalletArrived})
			wallets = append(wallets, wallet)
			continue
		}
		// If the account is the same as the first wallet, keep it
		if ks.wallets[0].Accounts()[0] == account {
			wallets = append(wallets, ks.wallets[0])
			ks.wallets = ks.wallets[1:]
			continue
		}
	}
	// Drop any leftover wallets and set the new batch
	for _, wallet := range ks.wallets {
		events = append(events, accounts.WalletEvent{Wallet: wallet, Kind: accounts.WalletDropped})
	}
	ks.wallets = wallets
	ks.mu.Unlock()

	// Fire all wallet events and return
	for _, event := range events {
		ks.updateFeed.Send(event)
	}
}

// Subscribe implements accounts.Backend, creating an async subscription to
// receive notifications on the addition or removal of keystore wallets.
func (ks *KeyStore) Subscribe(sink chan<- accounts.WalletEvent) event.Subscription {
	// We need the mutex to reliably start/stop the update loop
	ks.mu.Lock()
	defer ks.mu.Unlock()

	// Subscribe the caller and track the subscriber count
	sub := ks.updateScope.Track(ks.updateFeed.Subscribe(sink))

	// Subscribers require an active notification loop, start it
	if !ks.updating {
		ks.updating = true
		go ks.updater()
	}
	return sub
}

// updater is responsible for maintaining an up-to-date list of wallets stored in
// the keystore, and for firing wallet addition/removal events. It listens for
// account change events from the underlying account cache, and also periodically
// forces a manual refresh (only triggers for systems where the filesystem notifier
// is not running).
func (ks *KeyStore) updater() {
	for {
		// Wait for an account update or a refresh timeout
		select {
		case <-ks.changes:
		case <-time.After(walletRefreshCycle):
		}
		// Run the wallet refresher
		ks.refreshWallets()

		// If all our subscribers left, stop the updater
		ks.mu.Lock()
		if ks.updateScope.Count() == 0 {
			ks.updating = false
			ks.mu.Unlock()
			return
		}
		ks.mu.Unlock()
	}
}

// HasAddress reports whether a key with the given address is present.
func (ks *KeyStore) HasAddress(addr string) bool {
	return ks.cache.hasAddress(addr)
}

// Accounts returns all key files present in the directory.
func (ks *KeyStore) Accounts() []accounts.Account {
	return ks.cache.accounts()
}

// Delete deletes the key matched by account if the passphrase is correct.
// If the account contains no filename, the address must match a unique key.
func (ks *KeyStore) Delete(a accounts.Account, passphrase string) error {
	// Decrypting the key isn't really necessary, but we do
	// it anyway to check the password and zero out the key
	// immediately afterwards.
	a, key, err := ks.getDecryptedKey(a, passphrase)
	if key != nil {
		//zeroKey(key.PrivateKey)
		key.ZeroKey()
	}
	if err != nil {
		return err
	}
	// The order is crucial here. The key is dropped from the
	// cache after the file is gone so that a reload happening in
	// between won't insert it into the cache again.
	err = os.Remove(a.GetUrl().Path)
	if err == nil {
		ks.cache.delete(a)
		ks.refreshWallets()
	}
	return err
}

// SignHash calculates a ECDSA signature for the given hash. The produced
// signature is in the [R || S || V] format where V is 0 or 1.
func (ks *KeyStore) SignHash(a accounts.Account, hash []byte) ([]byte, error) {
	// Look up the key to sign with and abort if it cannot be found
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	unlockedKey, found := ks.unlocked[a.GetAddress()]
	if !found {
		return nil, ErrLocked
	}
	// Sign the hash using plain ECDSA operations
	//return crypto.Sign(hash, unlockedKey.PrivateKey)
	return unlockedKey.ukey.SinHash(hash)
}

// SignTx signs the given transaction with the requested account.
func (ks *KeyStore) SignEthTx(a accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Look up the key to sign with and abort if it cannot be found
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	unlockedKey, found := ks.unlocked[a.GetAddress()]
	if !found {
		return nil, ErrLocked
	}
	// Depending on the presence of the chain ID, sign with 2718 or homestead
	signer := types.LatestSignerForChainID(chainID)
	pi := unlockedKey.ukey.GetPrivateKey()
	if prikey, ok := pi.(*ecdsa.PrivateKey); ok {
		return types.SignTx(tx, signer, prikey)
	}
	return nil, errors.New("can not get private key")
	//return types.SignTx(tx, signer, *ecdsa.PrivateKey(unlockedKey.GetPrivateKey()))
}

// SignHashWithPassphrase signs hash if the private key matching the given address
// can be decrypted with the given passphrase. The produced signature is in the
// [R || S || V] format where V is 0 or 1.
func (ks *KeyStore) SignHashWithPassphrase(a accounts.Account, passphrase string, hash []byte) (signature []byte, err error) {
	_, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return nil, err
	}
	defer key.ZeroKey()
	return key.SinHash(hash)
}

// SignTxWithPassphrase signs the transaction if the private key matching the
// given address can be decrypted with the given passphrase.
func (ks *KeyStore) SignTxWithPassphrase(a accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	_, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return nil, err
	}
	//defer zeroKey(key.PrivateKey)
	defer key.ZeroKey()
	// Depending on the presence of the chain ID, sign with or without replay protection.
	signer := types.LatestSignerForChainID(chainID)
	pi := key.GetPrivateKey()
	if prikey, ok := pi.(*ecdsa.PrivateKey); ok {
		return types.SignTx(tx, signer, prikey)
	}
	return nil, errors.New("can not get private key")
	//return types.SignTx(tx, signer, key.PrivateKey)
}

// Unlock unlocks the given account indefinitely.
func (ks *KeyStore) Unlock(a accounts.Account, passphrase string) error {
	return ks.TimedUnlock(a, passphrase, 0)
}

// Lock removes the private key with the given address from memory.
func (ks *KeyStore) Lock(addr string) error {
	ks.mu.Lock()
	unl, found := ks.unlocked[addr]
	ks.mu.Unlock()
	if found {
		ks.expire(addr, unl, time.Duration(0)*time.Nanosecond)
	}
	return nil
}

// TimedUnlock unlocks the given account with the passphrase. The account
// stays unlocked for the duration of timeout. A timeout of 0 unlocks the account
// until the program exits. The account must match a unique key file.
//
// If the account address is already unlocked for a duration, TimedUnlock extends or
// shortens the active unlock timeout. If the address was previously unlocked
// indefinitely the timeout is not altered.
func (ks *KeyStore) TimedUnlock(a accounts.Account, passphrase string, timeout time.Duration) error {
	a, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return err
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()
	u, found := ks.unlocked[a.GetAddress()]
	if found {
		if u.abort == nil {
			// The address was unlocked indefinitely, so unlocking
			// it with a timeout would be confusing.
			//zeroKey(key.PrivateKey)
			key.ZeroKey()
			return nil
		}
		// Terminate the expire goroutine and replace it below.
		close(u.abort)
	}
	if timeout > 0 {
		u = &unlocked{ukey: key, abort: make(chan struct{})}
		go ks.expire(a.GetAddress(), u, timeout)
	} else {
		u = &unlocked{ukey: key}
	}
	ks.unlocked[a.GetAddress()] = u
	return nil
}

// Find resolves the given account into a unique entry in the keystore.
func (ks *KeyStore) Find(a accounts.Account) (accounts.Account, error) {
	ks.cache.maybeReload()
	ks.cache.mu.Lock()
	a, err := ks.cache.find(a)
	ks.cache.mu.Unlock()
	return a, err
}

func (ks *KeyStore) getDecryptedKey(a accounts.Account, auth string) (accounts.Account, accounts.Key, error) {
	a, err := ks.Find(a)
	if err != nil {
		return a, nil, err
	}
	key, err := ks.storage.GetKey(a.GetAddress(), a.GetUrl().Path, auth)
	return a, key, err
}

func (ks *KeyStore) expire(addr string, u *unlocked, timeout time.Duration) {
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case <-u.abort:
		// just quit
	case <-t.C:
		ks.mu.Lock()
		// only drop if it's still the same key instance that dropLater
		// was launched with. we can check that using pointer equality
		// because the map stores a new pointer every time the key is
		// unlocked.
		if ks.unlocked[addr] == u {
			//zeroKey(u.PrivateKey)
			u.ukey.ZeroKey()
			delete(ks.unlocked, addr)
		}
		ks.mu.Unlock()
	}
}

// NewAccount generates a new key and stores it into the key directory,
// encrypting it with the passphrase.
func (ks *KeyStore) NewAccount(passphrase string) (accounts.Account, error) {
	_, account, err := storeNewKey(ks.storage, crand.Reader, passphrase, ks.chainAdapter)
	if err != nil {
		//return accounts.Account{}, err
		return ks.chainAdapter.NewNullAccount(), err
	}
	// Add the account to the cache immediately rather
	// than waiting for file system notifications to pick it up.
	ks.cache.add(account)
	ks.refreshWallets()
	return account, nil
}

// Export exports as a JSON key, encrypted with newPassphrase.
func (ks *KeyStore) Export(a accounts.Account, passphrase, newPassphrase string) (keyJSON []byte, err error) {
	_, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return nil, err
	}
	var N, P int
	if store, ok := ks.storage.(*keyStorePassphrase); ok {
		N, P = store.scryptN, store.scryptP
	} else {
		N, P = StandardScryptN, StandardScryptP
	}
	return EncryptKey(key, newPassphrase, N, P)
}

// Import stores the given encrypted JSON key into the key directory.
func (ks *KeyStore) ImportEth(keyJSON []byte, passphrase, newPassphrase string) (accounts.Account, error) {
	key, err := DecryptKey(keyJSON, passphrase, ks.chainAdapter)
	if key != nil {
		defer key.ZeroKey()
	}
	if err != nil {
		return ks.chainAdapter.NewNullAccount(), err
	}
	ks.importMu.Lock()
	defer ks.importMu.Unlock()

	if ks.cache.hasAddress(key.GetStringAddress()) {
		a := ks.cache.byAddr[key.GetStringAddress()]
		var re accounts.Account
		if len(a) == 1 {
			re = a[0]
		} else {
			re = ks.chainAdapter.NewNullAccount()
		}
		return re, ErrAccountAlreadyExists
	}
	return ks.importKey(key, newPassphrase)
}

func (ks *KeyStore) ImportKey(key accounts.Key, passphrase string) (accounts.Account, error) {
	ks.importMu.Lock()
	defer ks.importMu.Unlock()
	defer key.ZeroKey()
	if ks.cache.hasAddress(key.GetStringAddress()) {
		a := ks.cache.byAddr[key.GetStringAddress()]
		var re accounts.Account
		if len(a) == 1 {
			re = a[0]
		} else {
			re = ks.chainAdapter.NewNullAccount()
		}
		return re, ErrAccountAlreadyExists
	}
	return ks.importKey(key, passphrase)
}

func (ks *KeyStore) importKey(key accounts.Key, passphrase string) (accounts.Account, error) {
	//a := accounts.Account{Address: key.Address, URL: accounts.URL{Scheme: KeyStoreScheme, Path: ks.storage.JoinPath(keyFileName(key.Address))}}
	account := ks.chainAdapter.NewAccount(key.GetStringAddress(), accounts.URL{Scheme: accounts.KeyStoreScheme, Path: ks.storage.JoinPath(keyFileName(key.GetStringAddress()))})
	if err := ks.storage.StoreKey(account.GetUrl().Path, key, passphrase); err != nil {
		return ks.chainAdapter.NewNullAccount(), err
	}
	ks.cache.add(account)
	ks.refreshWallets()
	return account, nil
}

// Update changes the passphrase of an existing account.
func (ks *KeyStore) Update(a accounts.Account, passphrase, newPassphrase string) error {
	a, key, err := ks.getDecryptedKey(a, passphrase)
	if err != nil {
		return err
	}
	return ks.storage.StoreKey(a.GetUrl().Path, key, newPassphrase)
}

// ImportPreSaleKey decrypts the given Ethereum presale wallet and stores
// a key file in the key directory. The key file is encrypted with the same passphrase.
func (ks *KeyStore) ImportPreSaleKey(keyJSON []byte, passphrase string) (accounts.Account, error) {
	a, _, err := importPreSaleKey(ks.storage, keyJSON, passphrase)
	if err != nil {
		return a, err
	}
	ks.cache.add(a)
	ks.refreshWallets()
	return a, nil
}

// isUpdating returns whether the event notification loop is running.
// This method is mainly meant for tests.
func (ks *KeyStore) isUpdating() bool {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	return ks.updating
}

// zeroKey zeroes a private key in memory.
func zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	clear(b)
}

func storeNewKey(ks IkeyStore, rand io.Reader, auth string, chAdapter accounts.IChainAdapter) (accounts.Key, accounts.Account, error) {
	key, err := chAdapter.NewKey(rand)
	if err != nil {
		return nil, nil, err
	}
	// a := accounts.Account{
	// 	Address: key.Address.Bytes(),
	// 	URL:     accounts.URL{Scheme: KeyStoreScheme, Path: ks.JoinPath(keyFileName(key.Address))},
	// }
	a := chAdapter.NewAccount(key.GetStringAddress(), accounts.URL{Scheme: accounts.KeyStoreScheme, Path: ks.JoinPath(keyFileName(key.GetStringAddress()))})

	if err := ks.StoreKey(a.GetUrl().Path, key, auth); err != nil {
		//zeroKey(key.PrivateKey)
		key.ZeroKey()
		return nil, a, err
	}
	return key, a, err
}

// keyFileName implements the naming convention for keyfiles:
// UTC--<created_at UTC ISO8601>-<address hex>
func keyFileName(keyAddr string) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), keyAddr)
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

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

func aesCBCDecrypt(key, cipherText, iv []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCBCDecrypter(aesBlock, iv)
	paddedPlaintext := make([]byte, len(cipherText))
	decrypter.CryptBlocks(paddedPlaintext, cipherText)
	plaintext := pkcs7Unpad(paddedPlaintext)
	if plaintext == nil {
		return nil, ErrDecrypt
	}
	return plaintext, err
}
