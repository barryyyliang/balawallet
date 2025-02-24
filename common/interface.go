package common

const (
	Nettype_ethereum int = 0
	Nettype_polygon  int = 1
	Nettype_solana   int = 2
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
