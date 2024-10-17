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

	//contractasset
	G_Res["zh-CN"]["/contractasset.html"] = make(map[string]string)
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_type"] = "合约类型" 
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_address"] = "合约地址" 
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_name"] = "资产名称" 
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_add"] = "添加" 
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_import"] = "导入" 

	G_Res["en"]["/contractasset.html"] = make(map[string]string)
	G_Res["en"]["/contractasset.html"]["contractasset_type"] = "Type" 
	G_Res["en"]["/contractasset.html"]["contractasset_address"] = "Address" 
	G_Res["en"]["/contractasset.html"]["contractasset_name"] = "Name" 
	G_Res["en"]["/contractasset.html"]["contractasset_add"] = "Add" 
	G_Res["en"]["/contractasset.html"]["contractasset_import"] = "Import" 

	//contractasset
	G_Res["zh-CN"]["/transfer.html"] = make(map[string]string)
	G_Res["zh-CN"]["/transfer.html"]["transfer_from"] = "转出账户"
	G_Res["zh-CN"]["/transfer.html"]["transfer_name"] = "账户别名"
	G_Res["zh-CN"]["/transfer.html"]["transfer_balance"] = "账户余额"
	G_Res["zh-CN"]["/transfer.html"]["transfer_to"] = "接收账户"
	G_Res["zh-CN"]["/transfer.html"]["transfer_amount"] = "转账金额"
	G_Res["zh-CN"]["/transfer.html"]["transfer_fee"] = "预估费用"
	G_Res["zh-CN"]["/transfer.html"]["transfer_acc"] = "账户密码"
	G_Res["zh-CN"]["/transfer.html"]["transfer_transfer"] = "转账"
	

	G_Res["en"]["/transfer.html"] = make(map[string]string)
	G_Res["en"]["/transfer.html"]["transfer_from"] = "From"
	G_Res["en"]["/transfer.html"]["transfer_name"] = "Name"
	G_Res["en"]["/transfer.html"]["transfer_balance"] = "Balance"
	G_Res["en"]["/transfer.html"]["transfer_to"] = "To"
	G_Res["en"]["/transfer.html"]["transfer_amount"] = "Amount"
	G_Res["en"]["/transfer.html"]["transfer_fee"] = "Fee"
	G_Res["en"]["/transfer.html"]["transfer_acc"] = "Password"
	G_Res["en"]["/transfer.html"]["transfer_transfer"] = "transfer"

}
