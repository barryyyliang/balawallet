# etheremwallet

unc init() {
	// See SEC 2 section 2.7.1
	// curve parameters taken from:
	// http://www.secg.org/sec2-v2.pdf
	theCurve.P, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 0)
	theCurve.N, _ = new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 0)
	theCurve.B, _ = new(big.Int).SetString("0x0000000000000000000000000000000000000000000000000000000000000007", 0)
	theCurve.Gx, _ = new(big.Int).SetString("0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 0)
	theCurve.Gy, _ = new(big.Int).SetString("0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 0)
	theCurve.BitSize = 256
}

// S256 returns a BitCurve which implements secp256k1.
func S256() *BitCurve {
	return theCurve
}



/Users/yongfuliang/go/pkg/mod/github.com/ethereum/go-ethereum@v1.11.5/crypto/secp256k1/curve.go
func (BitCurve *BitCurve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) -----通过私钥得到公钥


// DecryptKey decrypts a key from a json blob, returning the private key itself.
func DecryptKey(keyjson []byte, auth string) (*Key, error) {
	// Parse the json into a simple map to fetch the key version
	m := make(map[string]interface{})
	if err := json.Unmarshal(keyjson, &m); err != nil {
		return nil, err
	}
	// Depending on the version try to parse one way or another
	var (
		keyBytes, keyId []byte
		err             error
	)
	if version, ok := m["version"].(string); ok && version == "1" {
		k := new(encryptedKeyJSONV1)
		if err := json.Unmarshal(keyjson, k); err != nil {
			return nil, err
		}
		keyBytes, keyId, err = decryptKeyV1(k, auth)
	} else {
		k := new(encryptedKeyJSONV3)
		if err := json.Unmarshal(keyjson, k); err != nil {
			return nil, err
		}
		keyBytes, keyId, err = decryptKeyV3(k, auth)
	}
	// Handle any decryption errors and return the key
	if err != nil {
		return nil, err
	}
	key := crypto.ToECDSAUnsafe(keyBytes)
	id, err := uuid.FromBytes(keyId)
	if err != nil {
		return nil, err
	}
	return &Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(key.PublicKey),
		PrivateKey: key,
	}, nil
}

# windows编译
需要下载交叉编译器，建议：
LLVM MinGW
MSYS2
WinLibs

安装 c++编译工具
下载包
https://github.com/niXman/mingw-builds-binaries/releases
解压后，将其下的bin目录，添加到环境变量path


并设置环境变量
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1  

然后编译
go build -ldflags="-H windowsgui"

# windows 运行
安装 WebView2（服务器推荐）：下载Evergreen Standalone Installer（离线安装包，x64 版本）：
https://developer.microsoft.com/zh-cn/microsoft-edge/webview2/?form=MA13LH#download https://developer.microsoft.com/zh-cn/microsoft-edge/webview2/#download-section
以管理员身份运行安装程序，并重启服务器（重要！）

最新的windows系统默认安装了，不需要单独安装，老版本的需要单独安装


