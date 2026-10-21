package common

const (
	Nettype_ethereum string = "0"
	Nettype_polygon  string = "1"
	Nettype_solana   string = "2"
	Nettype_pharos   string = "3"
)

type ErrorResult struct {
	Code uint
	Msg  string
}

type ResBalance struct {
	Address string
	Balance string
}

type InputImportContract struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Name    string `json:"name"`
	From    string `json:"from"`
}

type ComContract struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Name    string `json:"name"`
}
