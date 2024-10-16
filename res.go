package main

var G_Res = make(map[string]map[string]map[string]string)

func initRes() {
	G_Res["zh-CN"] = make(map[string]map[string]string)
	G_Res["en"] = make(map[string]map[string]string)
	G_Res["zh-TW"] = make(map[string]map[string]string)

	//main
	G_Res["zh-CN"]["/main.html"] = make(map[string]string)
	G_Res["zh-CN"]["/main.html"]["qianbao"] = "钱包"
	G_Res["zh-CN"]["/main.html"]["contractasset"] = "合约资产"
	G_Res["zh-CN"]["/main.html"]["transfer"] = "转账"
	G_Res["zh-CN"]["/main.html"]["tools"] = "工具"
	G_Res["zh-CN"]["/main.html"]["system"] = "系统"

	G_Res["en"]["/main.html"] = make(map[string]string)
	G_Res["en"]["/main.html"]["qianbao"] = "Wallet"
	G_Res["en"]["/main.html"]["contractasset"] = "Contract Asset"
	G_Res["en"]["/main.html"]["transfer"] = "Transfer"
	G_Res["en"]["/main.html"]["tools"] = "Tools"
	G_Res["en"]["/main.html"]["system"] = "System"

	//qianbao
	G_Res["zh-CN"]["/qianbao.html"] = make(map[string]string)
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_type"] = "类型"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_address"] = "地址"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_name"] = "别名"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_balance"] = "余额"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_changebutton"] = "修改别名"

	G_Res["en"]["/qianbao.html"] = make(map[string]string)
	G_Res["en"]["/qianbao.html"]["qianbao_type"] = "Type"
	G_Res["en"]["/qianbao.html"]["qianbao_address"] = "Address"
	G_Res["en"]["/qianbao.html"]["qianbao_name"] = "Name"
	G_Res["en"]["/qianbao.html"]["qianbao_balance"] = "Balance"
	G_Res["en"]["/qianbao.html"]["qianbao_changebutton"] = "Change Name"
}
