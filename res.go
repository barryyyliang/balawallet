package main

var G_Res = make(map[string]map[string]map[string]string)

func initRes() {
	G_Res["zh"] = make(map[string]map[string]string)
	G_Res["zh"]["/main.html"] = make(map[string]string)
	G_Res["zh"]["/main.html"]["qianbao"] = "钱包"
	G_Res["zh"]["/main.html"]["contractasset"] = "合约资产"
	G_Res["zh"]["/main.html"]["transfer"] = "转账"
	G_Res["zh"]["/main.html"]["tools"] = "工具"
	G_Res["zh"]["/main.html"]["system"] = "系统"

	G_Res["en"] = make(map[string]map[string]string)
	G_Res["en"]["/main.html"] = make(map[string]string)
	G_Res["en"]["/main.html"]["qianbao"] = "wallet"
	G_Res["en"]["/main.html"]["contractasset"] = "contract asset"
	G_Res["en"]["/main.html"]["transfer"] = "transfer"
	G_Res["en"]["/main.html"]["tools"] = "tools"
	G_Res["en"]["/main.html"]["system"] = "system"

	G_Res["zh"]["qianbao"] = make(map[string]string)	
}
