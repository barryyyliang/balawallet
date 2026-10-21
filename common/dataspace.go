package common

import (
	"encoding/json"
	"strconv"

	"balawallet/accounts"
	"balawallet/accounts/keystore"

	//"github.com/ethereum/go-ethereum/crypto/sha3"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

func DoSaveDataSpace(name string, content string, auth string) string {
	re := GetResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = GetResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
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
					re = GetResult(Msg_writedb_fail, err3.Error())
				}
			} else {
				re = GetResult(Msg_fail, err1.Error())
			}
		} else if err2 == nil {
			re = GetResult(Msg_exist, string(value))
		} else {
			re = GetResult(Msg_fail, err2.Error())
		}
	} else {
		re = GetResult(Msg_fail, err.Error())
	}

	return string(re)
}

func DoChangeDataSpace(name string, content string, auth string) string {
	re := GetResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = GetResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		encryptData, err1 := keystore.EncryptDataV3([]byte(content), []byte(auth), keystore.StandardScryptN, keystore.StandardScryptP)
		if err1 == nil {
			data, _ := json.Marshal(encryptData)
			err2 := db.Put([]byte(name), data, nil)
			if err2 != nil {
				re = GetResult(Msg_writedb_fail, err2.Error())
			}
		} else {
			re = GetResult(Msg_fail, err1.Error())
		}
	} else {
		re = GetResult(Msg_fail, err.Error())
	}
	return string(re)
}

func DoDeleteDataSpace(name string, auth string) string {
	re := GetResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = GetResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		err2 := db.Delete([]byte(name), nil)
		if err2 != nil {
			re = GetResult(Msg_fail, err2.Error())
		}
	} else {
		re = GetResult(Msg_fail, err.Error())
	}
	return string(re)
}

func DoGetDataSpace(auth string) string {
	re := GetResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = GetResult(Msg_fail, check.Error())
		return string(re)
	}
	data := make(map[string]string)
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
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

			/*
				var tmp accounts.CryptoJSON
				if json.Unmarshal(it.Value(), &tmp) == nil {
					dedata, err1 := keystore.DecryptDataV3(tmp, auth)
					if err1 == nil {
						data[string(it.Key())] = string(dedata)
					} else {
						data[string(it.Key())] = err1.Error()
					}
				}
			*/
			data[string(it.Key())] = "**********"
		}
		msg, _ := json.Marshal(data)
		re, _ = json.Marshal(ErrorResult{Msg_success, string(msg)})
	} else {
		re = GetResult(Msg_fail, err.Error())
	}
	return string(re)
}

func DoGetDataItem(auth, key string) string {
	var re []byte
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = GetResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		value, err2 := db.Get([]byte(key), nil)
		if err2 == nil {
			var tmp accounts.CryptoJSON
			if json.Unmarshal(value, &tmp) == nil {
				dedata, err3 := keystore.DecryptDataV3(tmp, auth)
				if err3 == nil {
					re, err3 = json.Marshal(ErrorResult{Msg_success, string(dedata)})
					if err3 != nil {
						re = GetResult(Msg_fail, err3.Error())
					}
				} else {
					re = GetResult(Msg_fail, err3.Error())
				}
			}
		} else {
			re = GetResult(Msg_fail, err2.Error())
		}
	} else {
		re = GetResult(Msg_fail, err.Error())
	}
	return string(re)
}

func DoCheckDataSpace() string {
	re := GetResult(Msg_success, "")
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		//have not init
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			re = GetResult(Msg_space_notinited, "")
		} else if err1 == nil { //have inited
			re = GetResult(Msg_space_inited, "")
		} else { //other error
			re = GetResult(Msg_fail, err1.Error())
		}
	} else {
		re = GetResult(Msg_fail, err.Error())
	}

	return string(re)
}

func DoInitDataSpace(auth string) string {
	re := GetResult(Msg_success, "")
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
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
					re = GetResult(Msg_writedb_fail, err3.Error())
				}
			} else {
				re = GetResult(Msg_fail, err2.Error())
			}
		} else if err1 == nil { //have inited
			re = GetResult(Msg_space_inited, "")
		} else { //other error
			re = GetResult(Msg_fail, err1.Error())
		}
	} else {
		re = GetResult(Msg_fail, err.Error())
	}
	return string(re)
}

func DoChangeSpacePass(oldpass, newpass string) string {
	re := GetResult(Msg_success, "")
	err := checkDataSpaceAuth(oldpass)
	if err != nil {
		re = GetResult(Msg_fail, err.Error())
		return string(re)
	}

	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err != nil {
		re = GetResult(Msg_fail, err.Error())
		return string(re)
	}

	it := db.NewIterator(nil, nil)
	defer it.Release()
	var batch leveldb.Batch
	var errortag bool = true
	for it.Next() {
		var tmp accounts.CryptoJSON
		err = json.Unmarshal(it.Value(), &tmp)
		if err == nil {
			dedata, err1 := keystore.DecryptDataV3(tmp, oldpass)
			if err1 == nil {
				encryptData, err1 := keystore.EncryptDataV3(dedata, []byte(newpass), keystore.StandardScryptN, keystore.StandardScryptP)
				if err1 == nil {
					data, _ := json.Marshal(encryptData)
					batch.Put(it.Key(), data)
				} else {
					re = GetResult(Msg_fail, err1.Error())
					errortag = false
					break
				}
			} else {
				re = GetResult(Msg_fail, err1.Error())
				errortag = false
				break
			}
		} else {
			re = GetResult(Msg_fail, err.Error())
			errortag = false
			break
		}
	}
	if errortag {
		db.Write(&batch, nil)
	}

	return string(re)
}

func checkDataSpaceAuth(auth string) error {
	db, err := leveldb.OpenFile(G_DataDir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			return errors.New(GetMsg(strconv.Itoa(int(Msg_space_notinited))))
		} else if err1 == nil { //have inited
			var tmp accounts.CryptoJSON
			if json.Unmarshal(tag, &tmp) == nil {
				value, err2 := keystore.DecryptDataV3(tmp, auth)
				if err2 == nil && string(value) == "data_space_init" {
					return nil
				} else {
					return errors.New(GetMsg(strconv.Itoa(int(Msg_password_error))))
				}
			}
		} else { //other error
			return err1
		}
	}
	return nil
}
