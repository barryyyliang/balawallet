package main

import (
	"encoding/json"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	//"github.com/ethereum/go-ethereum/crypto/sha3"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func doSaveDataSpace(name string, content string, auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		value, err2 := db.Get([]byte(name), nil)
		if (err2 == errors.ErrNotFound) || (err2 == nil && len(value) == 0) {
			encryptData, err1 := keystore.EncryptDataV3([]byte(content), []byte(auth), keystore.StandardScryptN, keystore.StandardScryptP)
			if err1 == nil {
				data, _ := json.Marshal(encryptData)
				err3 := db.Put([]byte(name), data, nil)
				if err3 != nil {
					re = getResult(Msg_writedb_fail, err3.Error())
				}
			} else {
				re = getResult(Msg_fail, err1.Error())
			}
		} else if err2 == nil {
			re = getResult(Msg_exist, string(value))
		} else {
			re = getResult(Msg_fail, err2.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}

	return string(re)
}

func doChangeDataSpace(name string, content string, auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		encryptData, err1 := keystore.EncryptDataV3([]byte(content), []byte(auth), keystore.StandardScryptN, keystore.StandardScryptP)
		if err1 == nil {
			data, _ := json.Marshal(encryptData)
			err2 := db.Put([]byte(name), data, nil)
			if err2 != nil {
				re = getResult(Msg_writedb_fail, err2.Error())
			}
		} else {
			re = getResult(Msg_fail, err1.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func doDeleteDataSpace(name string, auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		err2 := db.Delete([]byte(name), nil)
		if err2 != nil {
			re = getResult(Msg_fail, err2.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func doGetDataSpace(auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	data := make(map[string]string)
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		it := db.NewIterator(nil, nil)
		defer it.Release()
		for it.Next() {
			if string(it.Key()) == "data_space_init" {
				continue
			}
			var tmp keystore.CryptoJSON
			if json.Unmarshal(it.Value(), &tmp) == nil {
				dedata, err1 := keystore.DecryptDataV3(tmp, auth)
				if err1 == nil {
					data[string(it.Key())] = string(dedata)
				} else {
					data[string(it.Key())] = err1.Error()
				}
			}
		}
		msg, _ := json.Marshal(data)
		re, _ = json.Marshal(ErrorResult{Msg_success, string(msg)})
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func doCheckDataSpace() string {
	re := getResult(Msg_success, "")
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		//have not init
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			re = getResult(Msg_space_notinited, "")
		} else if err1 == nil { //have inited
			re = getResult(Msg_space_inited, "")
		} else { //other error
			re = getResult(Msg_fail, err1.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}

	return string(re)
}

func doInitDataSpace(auth string) string {
	re := getResult(Msg_success, "")
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		//have not init
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			encryptData, err2 := keystore.EncryptDataV3([]byte("data_space_init"), []byte(auth), keystore.StandardScryptN, keystore.StandardScryptP)
			if err2 == nil {
				data, _ := json.Marshal(encryptData)
				err3 := db.Put([]byte("data_space_init"), data, nil)
				if err3 != nil {
					re = getResult(Msg_writedb_fail, err3.Error())
				}
			} else {
				re = getResult(Msg_fail, err2.Error())
			}
		} else if err1 == nil { //have inited
			re = getResult(Msg_space_inited, "")
		} else { //other error
			re = getResult(Msg_fail, err1.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func checkDataSpaceAuth(auth string) error {
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			return errors.New(doGetMsg(strconv.Itoa(int(Msg_space_notinited))))
		} else if err1 == nil { //have inited
			var tmp keystore.CryptoJSON
			if json.Unmarshal(tag, &tmp) == nil {
				value, err2 := keystore.DecryptDataV3(tmp, auth)
				if err2 == nil && string(value) == "data_space_init" {
					return nil
				} else {
					return errors.New(doGetMsg(strconv.Itoa(int(Msg_password_error))))
				}
			}
		} else { //other error
			return err1
		}
	}
	return nil
}
