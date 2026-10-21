package solchain

import (
	balacommon "balawallet/common"
)

// PersonalSign Solana 个人签名 (暂未实现)
func (s *SolChain) PersonalSign(account, password, data string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")
	return string(re)
}

// SignTypedData Solana 类型化数据签名 (暂未实现)
func (s *SolChain) SignTypedData(account, password, data string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")
	return string(re)
}

// SignTransaction Solana 签名交易 (暂未实现)
func (s *SolChain) SignTransaction(account, password, data string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")
	return string(re)
}

// GetChainId 获取 Solana Chain ID
func (s *SolChain) GetChainId() string {
	return "solana:mainnet"
}
