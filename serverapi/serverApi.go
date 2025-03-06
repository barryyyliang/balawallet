package serverapi

type Bchain interface {
	NewAccount(pass string, walletName string) string
	GetAccounts() string
	ChangeAccountName(address string, name string) string
	ImportAccount(pa string, name string, kk string) string
	ExportAccount(add string, pa string) string
	GetAccountName(address string) string
	Transfer(from string, to string, value string, pass string) string
	GetBalance(add string) string
	GetFee() string
	GetContract() string
	ImportContract(cons string) string
	DeployContract(from string, file string, pass string) string
	CheckNode() string
	LoadABI(file string) string
	ReadContract(from, to, funcAbi, paremters string) string
	QueryTx(hs string) string
	Signature(info string) string
	GetAbiFile() string
	GetBinFile() string
	InvokeContract(from, conAddr, funcAbi, param, pass string) string
	DeleteContract(address string) string
	SignatureVerify(info string) string
	ChangeAccountPass(atype, address, old, new string) string
	InvokeEstimate(from, conAddr, funcAbi, param string) string
	GetContractBalance(assetType, conAddr, ownerAddr string) string
}

var server = new(ServerAPI)

func GetServer(in Bchain) *ServerAPI {
	server.Chain = in
	return server
}
func SetServer(in Bchain) {
	server.setServer(in)
}

type ServerAPI struct {
	Chain Bchain
}

func (in *ServerAPI) setServer(bl Bchain) {
	in.Chain = bl
}

func (in *ServerAPI) DoNewAccount(pass string, walletName string) string {
	return in.Chain.NewAccount(pass, walletName)
}
func (in *ServerAPI) DoGetAccounts() string {
	return in.Chain.GetAccounts()
}
func (in *ServerAPI) DoChangeAccountName(address string, name string) string {
	return in.Chain.ChangeAccountName(address, name)

}
func (in *ServerAPI) DoImportAccount(pa string, name string, kk string) string {
	return in.Chain.ImportAccount(pa, name, kk)
}
func (in *ServerAPI) DoExportAccount(add string, pa string) string {
	return in.Chain.ExportAccount(add, pa)
}
func (in *ServerAPI) DoGetAccountName(address string) string {
	return in.Chain.GetAccountName(address)
}
func (in *ServerAPI) DoTransfer(from string, to string, value string, pass string) string {
	return in.Chain.Transfer(from, to, value, pass)
}
func (in *ServerAPI) DoGetBalance(add string) string {
	return in.Chain.GetBalance(add)
}
func (in *ServerAPI) DoGetFee() string {
	return in.Chain.GetFee()
}
func (in *ServerAPI) DoGetContract() string {
	return in.Chain.GetContract()
}
func (in *ServerAPI) DoImportContract(cons string) string {
	return in.Chain.ImportContract(cons)
}
func (in *ServerAPI) DoDeployContract(from string, file string, pass string) string {
	return in.Chain.DeployContract(from, file, pass)
}
func (in *ServerAPI) DoCheckNode() string {
	return in.Chain.CheckNode()
}
func (in *ServerAPI) DoLoadABI(file string) string {
	return in.Chain.LoadABI(file)
}
func (in *ServerAPI) DoReadContract(from, to, funcAbi, paremters string) string {
	return in.Chain.ReadContract(from, to, funcAbi, paremters)
}
func (in *ServerAPI) DoQueryTx(hs string) string {
	return in.Chain.QueryTx(hs)
}
func (in *ServerAPI) DoSignature(info string) string {
	return in.Chain.Signature(info)
}
func (in *ServerAPI) DoGetAbiFile() string {
	return in.Chain.GetAbiFile()
}
func (in *ServerAPI) DoGetBinFile() string {
	return in.Chain.GetBinFile()
}
func (in *ServerAPI) DoInvokeContract(from, conAddr, funcAbi, param, pass string) string {
	return in.Chain.InvokeContract(from, conAddr, funcAbi, param, pass)
}
func (in *ServerAPI) DoDeleteContract(address string) string {
	return in.Chain.DeleteContract(address)
}
func (in *ServerAPI) DoSignatureVerify(info string) string {
	return in.Chain.SignatureVerify(info)
}
func (in *ServerAPI) DoChangeAccountPass(atype, address, old, new string) string {
	return in.Chain.ChangeAccountPass(atype, address, old, new)
}
func (in *ServerAPI) DoInvokeEstimate(from, conAddr, funcAbi, param string) string {
	return in.Chain.InvokeEstimate(from, conAddr, funcAbi, param)
}
func (in *ServerAPI) DoGetContractBalance(assetType, conAddr, ownerAddr string) string {
	return in.Chain.GetContractBalance(assetType, conAddr, ownerAddr)
}
