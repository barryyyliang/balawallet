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
	G_Res["zh-CN"]["/main.html"]["main_close"] = "关闭"
	G_Res["zh-CN"]["/main.html"]["main_copy"] = "拷贝"
	G_Res["zh-CN"]["/main.html"]["main_paste"] = "粘贴"
	G_Res["zh-CN"]["/main.html"]["main_cut"] = "剪切"
	G_Res["zh-CN"]["/main.html"]["main_clear"] = "清除"
	G_Res["zh-CN"]["/main.html"]["main_deploy"] = "部署合约"
	G_Res["zh-CN"]["/main.html"]["main_invoke"] = "调用合约"
	G_Res["zh-CN"]["/main.html"]["main_query"] = "查询交易"
	G_Res["zh-CN"]["/main.html"]["main_sign"] = "数字签名"
	G_Res["zh-CN"]["/main.html"]["main_newaccount"] = "新建账户"
	G_Res["zh-CN"]["/main.html"]["main_importaccount"] = "导入账户"
	G_Res["zh-CN"]["/main.html"]["main_exportaccount"] = "导出账户"
	G_Res["zh-CN"]["/main.html"]["main_config"] = "系统设置"
	G_Res["zh-CN"]["/main.html"]["main_about"] = "关于 Web3Wallet"
	

	G_Res["en"]["/main.html"] = make(map[string]string)
	G_Res["en"]["/main.html"]["qianbao"] = "Wallet"
	G_Res["en"]["/main.html"]["contractasset"] = "Contract Asset"
	G_Res["en"]["/main.html"]["transfer"] = "Transfer"
	G_Res["en"]["/main.html"]["tools"] = "Tools"
	G_Res["en"]["/main.html"]["system"] = "System"
	G_Res["en"]["/main.html"]["main_close"] = "Close"
	G_Res["en"]["/main.html"]["main_copy"] = "Copy"
	G_Res["en"]["/main.html"]["main_paste"] = "Paste"
	G_Res["en"]["/main.html"]["main_cut"] = "Cut"
	G_Res["en"]["/main.html"]["main_clear"] = "Clear"
	G_Res["en"]["/main.html"]["main_deploy"] = "Deploy Contract"
	G_Res["en"]["/main.html"]["main_invoke"] = "Invoke Contract"
	G_Res["en"]["/main.html"]["main_query"] = "Query Transaction"
	G_Res["en"]["/main.html"]["main_sign"] = "Sign Message"
	G_Res["en"]["/main.html"]["main_newaccount"] = "Create Account"
	G_Res["en"]["/main.html"]["main_importaccount"] = "Import Account"
	G_Res["en"]["/main.html"]["main_exportaccount"] = "Export Account"
	G_Res["en"]["/main.html"]["main_config"] = "System Config"
	G_Res["en"]["/main.html"]["main_about"] = "About Web3Wallet"
	

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

	//deploycontract
	G_Res["zh-CN"]["/deploycontract.html"] = make(map[string]string)
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_file"] = "合约文件（bin)"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_acc"] = "账户密码"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_deploy"] = "部署"
	

	G_Res["en"]["/deploycontract.html"] = make(map[string]string)
	G_Res["en"]["/deploycontract.html"]["deploycontract_file"] = "Contract Bin File"
	G_Res["en"]["/deploycontract.html"]["deploycontract_acc"] = "Password"
	G_Res["en"]["/deploycontract.html"]["deploycontract_deploy"] = "Deploy"

	//deploycontract
	G_Res["zh-CN"]["/invokecontract.html"] = make(map[string]string)
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_contractaddress"] = "合约地址"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_contractabi"] = "合约ABI"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_load"] = "加载"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_function"] = "合约函数"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_invoke"] = "调用"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_query"] = "查询"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_acc"] = "账户密码"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_inputparameter"] = "输入参数"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_parameter"] = "参数"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_inputvalue"] = "值"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_inputtype"] = "类型"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_return"] = "返回结果"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_returnname"] = "返回值"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_returnvalue"] = "值"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_returntype"] = "类型"

	G_Res["en"]["/invokecontract.html"] = make(map[string]string)
	G_Res["en"]["/invokecontract.html"]["invokecontract_contractaddress"] = "Contract Address"
	G_Res["en"]["/invokecontract.html"]["invokecontract_contractabi"] = "Contract ABI"
	G_Res["en"]["/invokecontract.html"]["invokecontract_load"] = "Load"
	G_Res["en"]["/invokecontract.html"]["invokecontract_function"] = "Function"
	G_Res["en"]["/invokecontract.html"]["invokecontract_invoke"] = "Invoke"
	G_Res["en"]["/invokecontract.html"]["invokecontract_query"] = "Query"
	G_Res["en"]["/invokecontract.html"]["invokecontract_acc"] = "Password"
	G_Res["en"]["/invokecontract.html"]["invokecontract_inputparameter"] = "Input Parameters"
	G_Res["en"]["/invokecontract.html"]["invokecontract_parameter"] = "Parameters"
	G_Res["en"]["/invokecontract.html"]["invokecontract_inputvalue"] = "Value"
	G_Res["en"]["/invokecontract.html"]["invokecontract_inputtype"] = "Type"
	G_Res["en"]["/invokecontract.html"]["invokecontract_return"] = "The Results"
	G_Res["en"]["/invokecontract.html"]["invokecontract_returnname"] = "Return Name"
	G_Res["en"]["/invokecontract.html"]["invokecontract_returnvalue"] = "Value"
	G_Res["en"]["/invokecontract.html"]["invokecontract_returntype"] = "Type"

	//deploycontract
	G_Res["zh-CN"]["/querytx.html"] = make(map[string]string)
	G_Res["zh-CN"]["/querytx.html"]["querytx_txhash"] = "交易 hash"
	G_Res["zh-CN"]["/querytx.html"]["querytx_query"] = "查询"	

	G_Res["en"]["/querytx.html"] = make(map[string]string)
	G_Res["en"]["/querytx.html"]["querytx_txhash"] = "TX Hash"
	G_Res["en"]["/querytx.html"]["querytx_query"] = "Query"

}
