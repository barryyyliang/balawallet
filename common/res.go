package common

var G_Res = make(map[string]map[string]map[string]string)

func InitRes() {
	G_Res["zh-CN"] = make(map[string]map[string]string)
	G_Res["en"] = make(map[string]map[string]string)
	G_Res["zh-TW"] = make(map[string]map[string]string)
	G_Res["common"] = make(map[string]map[string]string)

	//common
	G_Res["common"]["zh-CN"] = make(map[string]string)
	G_Res["common"]["zh-CN"]["old_pass"] = "旧密码"
	G_Res["common"]["zh-CN"]["new_pass"] = "新密码"
	G_Res["common"]["zh-CN"]["confirm_pass"] = "确认密码"
	G_Res["common"]["zh-CN"]["change"] = "修改"

	G_Res["common"]["zh-TW"] = make(map[string]string)
	G_Res["common"]["zh-TW"]["old_pass"] = "舊密碼"
	G_Res["common"]["zh-TW"]["new_pass"] = "新密碼"
	G_Res["common"]["zh-TW"]["confirm_pass"] = "確認密碼"
	G_Res["common"]["zh-TW"]["change"] = "修改"

	G_Res["common"]["en"] = make(map[string]string)
	G_Res["common"]["en"]["old_pass"] = "Old Password"
	G_Res["common"]["en"]["new_pass"] = "New Password"
	G_Res["common"]["en"]["confirm_pass"] = "Confirm Password"
	G_Res["common"]["en"]["change"] = "Change"

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
	G_Res["zh-CN"]["/main.html"]["main_about"] = "关于 BalaWallet"
	G_Res["zh-CN"]["/main.html"]["dataspace"] = "数据空间"
	G_Res["zh-CN"]["/main.html"]["main_ok"] = "确定"
	G_Res["zh-CN"]["/main.html"]["main_cancel"] = "取消"
	G_Res["zh-CN"]["/main.html"]["main_changespacepass"] = "数据空间密码"

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
	G_Res["en"]["/main.html"]["main_about"] = "About BalaWallet"
	G_Res["en"]["/main.html"]["dataspace"] = "Data Space"
	G_Res["en"]["/main.html"]["main_ok"] = "Ok"
	G_Res["en"]["/main.html"]["main_cancel"] = "Cancel"
	G_Res["en"]["/main.html"]["main_changespacepass"] = "Data Space Password"

	G_Res["zh-TW"]["/main.html"] = make(map[string]string)
	G_Res["zh-TW"]["/main.html"]["qianbao"] = "錢包"
	G_Res["zh-TW"]["/main.html"]["contractasset"] = "合約資產"
	G_Res["zh-TW"]["/main.html"]["transfer"] = "轉帳"
	G_Res["zh-TW"]["/main.html"]["tools"] = "工具"
	G_Res["zh-TW"]["/main.html"]["system"] = "系統"
	G_Res["zh-TW"]["/main.html"]["main_close"] = "關閉"
	G_Res["zh-TW"]["/main.html"]["main_copy"] = "複製"
	G_Res["zh-TW"]["/main.html"]["main_paste"] = "粘貼"
	G_Res["zh-TW"]["/main.html"]["main_cut"] = "剪切"
	G_Res["zh-TW"]["/main.html"]["main_clear"] = "清除"
	G_Res["zh-TW"]["/main.html"]["main_deploy"] = "部署合約"
	G_Res["zh-TW"]["/main.html"]["main_invoke"] = "調用合約"
	G_Res["zh-TW"]["/main.html"]["main_query"] = "查詢交易"
	G_Res["zh-TW"]["/main.html"]["main_sign"] = "數位簽名"
	G_Res["zh-TW"]["/main.html"]["main_newaccount"] = "新建帳戶"
	G_Res["zh-TW"]["/main.html"]["main_importaccount"] = "導入帳戶"
	G_Res["zh-TW"]["/main.html"]["main_exportaccount"] = "導出帳戶"
	G_Res["zh-TW"]["/main.html"]["main_config"] = "系統設定"
	G_Res["zh-TW"]["/main.html"]["main_about"] = "關於 BalaWallet"
	G_Res["zh-TW"]["/main.html"]["dataspace"] = "數據空間"
	G_Res["zh-TW"]["/main.html"]["main_ok"] = "確定"
	G_Res["zh-TW"]["/main.html"]["main_cancel"] = "取消"
	G_Res["zh-TW"]["/main.html"]["main_changespacepass"] = "數據空間密碼"

	//qianbao
	G_Res["zh-CN"]["/qianbao.html"] = make(map[string]string)
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_type"] = "类型"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_address"] = "地址"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_name"] = "别名"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_balance"] = "余额"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_changebutton"] = "修改别名"
	G_Res["zh-CN"]["/qianbao.html"]["qianbao_changepass"] = "修改密码"

	G_Res["en"]["/qianbao.html"] = make(map[string]string)
	G_Res["en"]["/qianbao.html"]["qianbao_type"] = "Type"
	G_Res["en"]["/qianbao.html"]["qianbao_address"] = "Address"
	G_Res["en"]["/qianbao.html"]["qianbao_name"] = "Name"
	G_Res["en"]["/qianbao.html"]["qianbao_balance"] = "Balance"
	G_Res["en"]["/qianbao.html"]["qianbao_changebutton"] = "Change Name"
	G_Res["en"]["/qianbao.html"]["qianbao_changepass"] = "Change Password"

	G_Res["zh-TW"]["/qianbao.html"] = make(map[string]string)
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_type"] = "類型"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_address"] = "地址"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_name"] = "別名"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_balance"] = "餘額"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_changebutton"] = "修改別名"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_changepass"] = "修改密码"

	//contractasset
	G_Res["zh-CN"]["/contractasset.html"] = make(map[string]string)
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_type"] = "合约类型"
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_address"] = "合约地址"
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_name"] = "资产名称"
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_add"] = "添加"
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_import"] = "导入"
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_delete"] = "删除"
	G_Res["zh-CN"]["/contractasset.html"]["contractasset_invoke"] = "调用"

	G_Res["en"]["/contractasset.html"] = make(map[string]string)
	G_Res["en"]["/contractasset.html"]["contractasset_type"] = "Type"
	G_Res["en"]["/contractasset.html"]["contractasset_address"] = "Address"
	G_Res["en"]["/contractasset.html"]["contractasset_name"] = "Name"
	G_Res["en"]["/contractasset.html"]["contractasset_add"] = "Add"
	G_Res["en"]["/contractasset.html"]["contractasset_import"] = "Import"
	G_Res["en"]["/contractasset.html"]["contractasset_delete"] = "Delete"
	G_Res["en"]["/contractasset.html"]["contractasset_invoke"] = "Inovke"

	G_Res["zh-TW"]["/contractasset.html"] = make(map[string]string)
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_type"] = "合约類型"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_address"] = "合约地址"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_name"] = "資產名稱"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_add"] = "添加"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_import"] = "導入"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_delete"] = "删除"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_invoke"] = "調用"

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

	G_Res["zh-TW"]["/transfer.html"] = make(map[string]string)
	G_Res["zh-TW"]["/transfer.html"]["transfer_from"] = "轉出帳戶"
	G_Res["zh-TW"]["/transfer.html"]["transfer_name"] = "帳戶别名"
	G_Res["zh-TW"]["/transfer.html"]["transfer_balance"] = "帳戶餘額"
	G_Res["zh-TW"]["/transfer.html"]["transfer_to"] = "接收帳戶"
	G_Res["zh-TW"]["/transfer.html"]["transfer_amount"] = "轉帳金額"
	G_Res["zh-TW"]["/transfer.html"]["transfer_fee"] = "費用"
	G_Res["zh-TW"]["/transfer.html"]["transfer_acc"] = "帳戶密碼"
	G_Res["zh-TW"]["/transfer.html"]["transfer_transfer"] = "轉出"

	//deploycontract
	G_Res["zh-CN"]["/deploycontract.html"] = make(map[string]string)
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_file"] = "合约文件（bin)"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_acc"] = "账户密码"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_deploy"] = "部署"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontractabi_file"] = "ABI文件(abi)"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_fee"] = "预估费用"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_estimate"] = "预估"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_inputparameter"] = "构造参数"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_parameter"] = "参数"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_inputvalue"] = "值"
	G_Res["zh-CN"]["/deploycontract.html"]["deploycontract_inputtype"] = "类型"

	G_Res["en"]["/deploycontract.html"] = make(map[string]string)
	G_Res["en"]["/deploycontract.html"]["deploycontract_file"] = "Contract Bin File"
	G_Res["en"]["/deploycontract.html"]["deploycontract_acc"] = "Password"
	G_Res["en"]["/deploycontract.html"]["deploycontract_deploy"] = "Deploy"
	G_Res["en"]["/deploycontract.html"]["deploycontractabi_file"] = "ABI File(abi)"
	G_Res["en"]["/deploycontract.html"]["deploycontract_fee"] = "Estimate Fee"
	G_Res["en"]["/deploycontract.html"]["deploycontract_estimate"] = "Estimate"
	G_Res["en"]["/deploycontract.html"]["deploycontract_inputparameter"] = "Construct Parameter"
	G_Res["en"]["/deploycontract.html"]["deploycontract_parameter"] = "Parameter"
	G_Res["en"]["/deploycontract.html"]["deploycontract_inputvalue"] = "Value"
	G_Res["en"]["/deploycontract.html"]["deploycontract_inputtype"] = "Type"

	G_Res["zh-TW"]["/deploycontract.html"] = make(map[string]string)
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_file"] = "合約文件（bin)"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_acc"] = "帳戶密碼"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_deploy"] = "部署"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontractabi_file"] = "ABI文件(abi)"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_fee"] = "預估費用"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_estimate"] = "預估"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_inputparameter"] = "構建參數"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_parameter"] = "參數"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_inputvalue"] = "值"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_inputtype"] = "類型"

	//invokecontract
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
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_fee"] = "预估费用"
	G_Res["zh-CN"]["/invokecontract.html"]["invokecontract_estimate"] = "预估"

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
	G_Res["en"]["/invokecontract.html"]["invokecontract_fee"] = "Estimate Fee"
	G_Res["en"]["/invokecontract.html"]["invokecontract_estimate"] = "Estimate"

	G_Res["zh-TW"]["/invokecontract.html"] = make(map[string]string)
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_contractaddress"] = "合約地址"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_contractabi"] = "合約ABI"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_load"] = "載入"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_function"] = "合約函數"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_invoke"] = "調用"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_query"] = "查詢"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_acc"] = "帳戶密碼"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_inputparameter"] = "輸入參數"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_parameter"] = "參數"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_inputvalue"] = "值"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_inputtype"] = "類型"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_return"] = "回傳結果"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_returnname"] = "傳回值"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_returnvalue"] = "值"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_returntype"] = "類型"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_fee"] = "預估費用"
	G_Res["zh-TW"]["/invokecontract.html"]["invokecontract_estimate"] = "預估"

	//querytx
	G_Res["zh-CN"]["/querytx.html"] = make(map[string]string)
	G_Res["zh-CN"]["/querytx.html"]["querytx_txhash"] = "交易 hash"
	G_Res["zh-CN"]["/querytx.html"]["querytx_query"] = "查询"

	G_Res["en"]["/querytx.html"] = make(map[string]string)
	G_Res["en"]["/querytx.html"]["querytx_txhash"] = "TX Hash"
	G_Res["en"]["/querytx.html"]["querytx_query"] = "Query"

	G_Res["zh-TW"]["/querytx.html"] = make(map[string]string)
	G_Res["zh-TW"]["/querytx.html"]["querytx_txhash"] = "交易 hash"
	G_Res["zh-TW"]["/querytx.html"]["querytx_query"] = "查詢"

	//signature
	G_Res["zh-CN"]["/signature.html"] = make(map[string]string)
	G_Res["zh-CN"]["/signature.html"]["signature_infomation"] = "待签名信息"
	G_Res["zh-CN"]["/signature.html"]["signature_acc"] = "账户密码"
	G_Res["zh-CN"]["/signature.html"]["signature_sign"] = "签名"
	G_Res["zh-CN"]["/signature.html"]["signature_verify"] = "验证签名"

	G_Res["en"]["/signature.html"] = make(map[string]string)
	G_Res["en"]["/signature.html"]["signature_infomation"] = "Information to be signed"
	G_Res["en"]["/signature.html"]["signature_acc"] = "Password"
	G_Res["en"]["/signature.html"]["signature_sign"] = "Sign"
	G_Res["en"]["/signature.html"]["signature_verify"] = "Verify"

	G_Res["zh-TW"]["/signature.html"] = make(map[string]string)
	G_Res["zh-TW"]["/signature.html"]["signature_infomation"] = "待簽名訊息"
	G_Res["zh-TW"]["/signature.html"]["signature_acc"] = "帳號密碼"
	G_Res["zh-TW"]["/signature.html"]["signature_sign"] = "簽名"
	G_Res["zh-TW"]["/signature.html"]["signature_verify"] = "驗證簽名"

	//newaccount
	G_Res["zh-CN"]["/newaccount.html"] = make(map[string]string)
	G_Res["zh-CN"]["/newaccount.html"]["newaccount_name"] = "账户别名"
	G_Res["zh-CN"]["/newaccount.html"]["newaccount_acc"] = "输入密码"
	G_Res["zh-CN"]["/newaccount.html"]["newaccount_reacc"] = "重输密码"
	G_Res["zh-CN"]["/newaccount.html"]["newaccount_create"] = "生成账户"

	G_Res["en"]["/newaccount.html"] = make(map[string]string)
	G_Res["en"]["/newaccount.html"]["newaccount_name"] = "Name"
	G_Res["en"]["/newaccount.html"]["newaccount_acc"] = "Password"
	G_Res["en"]["/newaccount.html"]["newaccount_reacc"] = "Re-enter password"
	G_Res["en"]["/newaccount.html"]["newaccount_create"] = "Create"

	G_Res["zh-TW"]["/newaccount.html"] = make(map[string]string)
	G_Res["zh-TW"]["/newaccount.html"]["newaccount_name"] = "帳戶別名"
	G_Res["zh-TW"]["/newaccount.html"]["newaccount_acc"] = "輸入密碼"
	G_Res["zh-TW"]["/newaccount.html"]["newaccount_reacc"] = "重輸密碼"
	G_Res["zh-TW"]["/newaccount.html"]["newaccount_create"] = "生成帳戶"

	//export
	G_Res["zh-CN"]["/export.html"] = make(map[string]string)
	G_Res["zh-CN"]["/export.html"]["export_hint"] = "提示：将导出明文私钥，请妥善保管，秘钥丢失，会导致秘钥管理的数字资产丢失，如果想要备份加密的钱包文件，请在钱包的account目录下直接获取"
	G_Res["zh-CN"]["/export.html"]["export_account"] = "账户"
	G_Res["zh-CN"]["/export.html"]["export_acc"] = "账户密码"
	G_Res["zh-CN"]["/export.html"]["export_export"] = "导出"

	G_Res["en"]["/export.html"] = make(map[string]string)
	G_Res["en"]["/export.html"]["export_hint"] = "Tip: The plain text private key will be exported, please keep it properly. Loss of the secret key will result in the loss of the digital assets managed by the secret key. If you want to back up the encrypted wallet file, please obtain it directly from the account directory of the wallet."
	G_Res["en"]["/export.html"]["export_account"] = "Account"
	G_Res["en"]["/export.html"]["export_acc"] = "Password"
	G_Res["en"]["/export.html"]["export_export"] = "Export"

	G_Res["zh-TW"]["/export.html"] = make(map[string]string)
	G_Res["zh-TW"]["/export.html"]["export_hint"] = "提示：將導出明文私鑰，請妥善保管，秘鑰遺失，會導致秘鑰管理的數位資產遺失，如果想要備份加密的錢包文件，請在錢包的account目錄下直接獲取"
	G_Res["zh-TW"]["/export.html"]["export_account"] = "帳戶"
	G_Res["zh-TW"]["/export.html"]["export_acc"] = "帳戶密碼"
	G_Res["zh-TW"]["/export.html"]["export_export"] = "導出"

	//import
	G_Res["zh-CN"]["/import.html"] = make(map[string]string)
	G_Res["zh-CN"]["/import.html"]["import_private"] = "私钥"
	G_Res["zh-CN"]["/import.html"]["import_name"] = "账户别名"
	G_Res["zh-CN"]["/import.html"]["import_acc"] = "输入密码"
	G_Res["zh-CN"]["/import.html"]["import_reacc"] = "重输密码"
	G_Res["zh-CN"]["/import.html"]["import_import"] = "导入账户"

	G_Res["en"]["/import.html"] = make(map[string]string)
	G_Res["en"]["/import.html"]["import_private"] = "Private Key"
	G_Res["en"]["/import.html"]["import_name"] = "Name"
	G_Res["en"]["/import.html"]["import_acc"] = "Password"
	G_Res["en"]["/import.html"]["import_reacc"] = "Re-enter password"
	G_Res["en"]["/import.html"]["import_import"] = "Import"

	G_Res["zh-TW"]["/import.html"] = make(map[string]string)
	G_Res["zh-TW"]["/import.html"]["import_private"] = "私鑰"
	G_Res["zh-TW"]["/import.html"]["import_name"] = "帳戶別名"
	G_Res["zh-TW"]["/import.html"]["import_acc"] = "輸入密碼"
	G_Res["zh-TW"]["/import.html"]["import_reacc"] = "重輸密碼"
	G_Res["zh-TW"]["/import.html"]["import_import"] = "導入帳戶"

	//config
	G_Res["zh-CN"]["/config.html"] = make(map[string]string)
	G_Res["zh-CN"]["/config.html"]["config_eth_address"] = "以太坊地址"
	G_Res["zh-CN"]["/config.html"]["config_contractdir"] = "合约文件目录"
	G_Res["zh-CN"]["/config.html"]["config_language"] = "系统语言"
	G_Res["zh-CN"]["/config.html"]["config_save"] = "保存"
	G_Res["zh-CN"]["/config.html"]["config_pol_address"] = "Polygon 地址"
	G_Res["zh-CN"]["/config.html"]["config_sol_address"] = "Solana 地址"
	G_Res["zh-CN"]["/config.html"]["config_phrs_address"] = "Pharos 地址"

	G_Res["en"]["/config.html"] = make(map[string]string)
	G_Res["en"]["/config.html"]["config_eth_address"] = "Eth Address"
	G_Res["en"]["/config.html"]["config_contractdir"] = "Contract File Dir"
	G_Res["en"]["/config.html"]["config_language"] = "Language"
	G_Res["en"]["/config.html"]["config_save"] = "Save"
	G_Res["en"]["/config.html"]["config_pol_address"] = "Polygon Address"
	G_Res["en"]["/config.html"]["config_sol_address"] = "Solana Address"
	G_Res["en"]["/config.html"]["config_phrs_address"] = "Pharos Address"

	G_Res["zh-TW"]["/config.html"] = make(map[string]string)
	G_Res["zh-TW"]["/config.html"]["config_eth_address"] = "以太坊位址"
	G_Res["zh-TW"]["/config.html"]["config_contractdir"] = "合約檔案目錄"
	G_Res["zh-TW"]["/config.html"]["config_language"] = "系統語言"
	G_Res["zh-TW"]["/config.html"]["config_save"] = "儲存"
	G_Res["zh-TW"]["/config.html"]["config_pol_address"] = "Polygon 位址"
	G_Res["zh-TW"]["/config.html"]["config_sol_address"] = "Solana 位址"
	G_Res["zh-TW"]["/config.html"]["config_phrs_address"] = "Pharos 位址"

	//about
	G_Res["zh-CN"]["/about.html"] = make(map[string]string)
	G_Res["zh-CN"]["/about.html"]["about_about"] = `BalaWallet 提供数字资产管理，智能合约调用。
	
    本软件暂不提供节点服务，要连接到区块链网络，用户可以连接第三方区块链节点或者 RPC 服务

    RPC厂商注册后，可提供一定量的免费服务,如下 RPC 厂商可提供相关服务：
    https://www.infura.io/zh/networks/ethereum
	https://etherscan.io/apis
	https://zan.top/home/node-service

    使用的的第三方开源项目：    
    https://github.com/ethereum/go-ethereum.git
    https://github.com/webview/webview_go.git

    使用建议和问题反馈：
    balawallet@163.com`

	G_Res["en"]["/about.html"] = make(map[string]string)
	G_Res["en"]["/about.html"]["about_about"] = `BalaWallet provides digital asset management and smart contract calling functions。
    
	This software does not currently provide node services. To connect to the blockchain network, users can connect to third-party blockchain nodes or RPC services.
	
	After RPC vendors register, they can provide a certain amount of free services, The following RPC vendors can provide related services:
    https://www.infura.io/zh/networks/ethereum
	https://etherscan.io/apis
	https://zan.top/home/node-service
    
	Third-party open source projects used:    
    https://github.com/ethereum/go-ethereum.git
    https://github.com/webview/webview_go.git
    
	Usage suggestions and problem feedback:
    balawallet@163.com`

	G_Res["zh-TW"]["/about.html"] = make(map[string]string)
	G_Res["zh-TW"]["/about.html"]["about_about"] = `BalaWallet 提供數位資產管理，智慧合約呼叫功能。
	
    本軟體暫不提供節點服務，要連接到區塊鏈網絡，用戶可以連接第三方區塊鏈節點或 RPC 服務

    RPC廠商註冊後，可提供一定量的免費服務,如下 RPC 廠商可提供相關服務：
    https://www.infura.io/zh/networks/ethereum
	https://etherscan.io/apis
	https://zan.top/home/node-service

    使用的第三方開源專案：    
    https://github.com/ethereum/go-ethereum.git
    https://github.com/webview/webview_go.git

    使用建議和問題回饋：
    balawallet@163.com`

	//dataspace
	G_Res["zh-CN"]["/dataspace.html"] = make(map[string]string)
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_name"] = "数据名称"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_content"] = "数据内容"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_add"] = "添加"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_delete"] = "删除"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_recover"] = "还原"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_save"] = "保存"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_change"] = "修改"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_auth"] = "个人空间密码"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_load"] = "加载"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_init"] = "初始化密码"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_show"] = "显示"
	G_Res["zh-CN"]["/dataspace.html"]["dataspace_hide"] = "隐藏"

	G_Res["zh-TW"]["/dataspace.html"] = make(map[string]string)
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_name"] = "數據名稱"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_content"] = "數據內容"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_add"] = "添加"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_delete"] = "删除"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_recover"] = "還原"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_save"] = "保存"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_change"] = "修改"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_auth"] = "個人空間密碼"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_load"] = "加載"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_init"] = "初始化密碼"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_show"] = "顯示"
	G_Res["zh-TW"]["/dataspace.html"]["dataspace_hide"] = "隱藏"

	G_Res["en"]["/dataspace.html"] = make(map[string]string)
	G_Res["en"]["/dataspace.html"]["dataspace_name"] = "Data Name"
	G_Res["en"]["/dataspace.html"]["dataspace_content"] = "Content"
	G_Res["en"]["/dataspace.html"]["dataspace_add"] = "Add"
	G_Res["en"]["/dataspace.html"]["dataspace_delete"] = "Delete"
	G_Res["en"]["/dataspace.html"]["dataspace_recover"] = "Recover"
	G_Res["en"]["/dataspace.html"]["dataspace_save"] = "Save"
	G_Res["en"]["/dataspace.html"]["dataspace_change"] = "Change"
	G_Res["en"]["/dataspace.html"]["dataspace_auth"] = "Personal Data Space Password"
	G_Res["en"]["/dataspace.html"]["dataspace_load"] = "Load"
	G_Res["en"]["/dataspace.html"]["dataspace_init"] = "Initialize Password"
	G_Res["en"]["/dataspace.html"]["dataspace_show"] = "Show"
	G_Res["en"]["/dataspace.html"]["dataspace_hide"] = "Hide"

	//changepass
	G_Res["zh-CN"]["/changepass.html"] = make(map[string]string)
	G_Res["zh-CN"]["/changepass.html"]["changepass_type"] = "账户类型"
	G_Res["zh-CN"]["/changepass.html"]["changepass_address"] = "账户地址"
	G_Res["zh-CN"]["/changepass.html"]["changepass_name"] = "账户名称"
	G_Res["zh-CN"]["/changepass.html"]["changepass_old"] = G_Res["common"]["zh-CN"]["old_pass"]
	G_Res["zh-CN"]["/changepass.html"]["changepass_new"] = G_Res["common"]["zh-CN"]["new_pass"]
	G_Res["zh-CN"]["/changepass.html"]["changepass_confirm"] = G_Res["common"]["zh-CN"]["confirm_pass"]
	G_Res["zh-CN"]["/changepass.html"]["changepass_change"] = G_Res["common"]["zh-CN"]["change"]
	G_Res["zh-CN"]["/changepass.html"]["changepass_hint"] = "提示：请妥善保管密码，忘记密码，无法恢复，会导致秘钥管理的数字资产丢失，请备份账户私钥"

	G_Res["zh-TW"]["/changepass.html"] = make(map[string]string)
	G_Res["zh-TW"]["/changepass.html"]["changepass_type"] = "賬戶類型"
	G_Res["zh-TW"]["/changepass.html"]["changepass_address"] = "賬戶地址"
	G_Res["zh-TW"]["/changepass.html"]["changepass_name"] = "名稱"
	G_Res["zh-TW"]["/changepass.html"]["changepass_old"] = G_Res["common"]["zh-TW"]["old_pass"]
	G_Res["zh-TW"]["/changepass.html"]["changepass_new"] = G_Res["common"]["zh-TW"]["new_pass"]
	G_Res["zh-TW"]["/changepass.html"]["changepass_confirm"] = G_Res["common"]["zh-TW"]["confirm_pass"]
	G_Res["zh-TW"]["/changepass.html"]["changepass_change"] = G_Res["common"]["zh-TW"]["change"]
	G_Res["zh-TW"]["/changepass.html"]["changepass_hint"] = "提示：請妥善保管密碼，忘記密碼將無法恢復，這將導致密鑰管理的數字資產丟失，請備份賬戶私鑰。"

	G_Res["en"]["/changepass.html"] = make(map[string]string)
	G_Res["en"]["/changepass.html"]["changepass_type"] = "Type"
	G_Res["en"]["/changepass.html"]["changepass_address"] = "Address"
	G_Res["en"]["/changepass.html"]["changepass_name"] = "Name"
	G_Res["en"]["/changepass.html"]["changepass_old"] = G_Res["common"]["en"]["old_pass"]
	G_Res["en"]["/changepass.html"]["changepass_new"] = G_Res["common"]["en"]["new_pass"]
	G_Res["en"]["/changepass.html"]["changepass_confirm"] = G_Res["common"]["en"]["confirm_pass"]
	G_Res["en"]["/changepass.html"]["changepass_change"] = G_Res["common"]["en"]["change"]
	G_Res["en"]["/changepass.html"]["changepass_hint"] = "Note: Please keep your password safe. Forgetting your password cannot be recovered, which will lead to the loss of digital assets in key management. Please back up your account private key."

	//spacepass
	G_Res["zh-CN"]["/spacepass.html"] = make(map[string]string)
	G_Res["zh-CN"]["/spacepass.html"]["spacepass_old"] = G_Res["common"]["zh-CN"]["old_pass"]
	G_Res["zh-CN"]["/spacepass.html"]["spacepass_new"] = G_Res["common"]["zh-CN"]["new_pass"]
	G_Res["zh-CN"]["/spacepass.html"]["spacepass_confirm"] = G_Res["common"]["zh-CN"]["confirm_pass"]
	G_Res["zh-CN"]["/spacepass.html"]["spacepass_change"] = G_Res["common"]["zh-CN"]["change"]
	G_Res["zh-CN"]["/spacepass.html"]["spacepass_hint"] = "提示：请妥善保管密码，忘记密码，无法恢复。"

	G_Res["zh-TW"]["/spacepass.html"] = make(map[string]string)
	G_Res["zh-TW"]["/spacepass.html"]["spacepass_old"] = G_Res["common"]["zh-TW"]["old_pass"]
	G_Res["zh-TW"]["/spacepass.html"]["spacepass_new"] = G_Res["common"]["zh-TW"]["new_pass"]
	G_Res["zh-TW"]["/spacepass.html"]["spacepass_confirm"] = G_Res["common"]["zh-TW"]["confirm_pass"]
	G_Res["zh-TW"]["/spacepass.html"]["spacepass_change"] = G_Res["common"]["zh-TW"]["change"]
	G_Res["zh-TW"]["/spacepass.html"]["spacepass_hint"] = "提示：請妥善保管密碼，忘記密碼將無法恢復。"

	G_Res["en"]["/spacepass.html"] = make(map[string]string)
	G_Res["en"]["/spacepass.html"]["spacepass_old"] = G_Res["common"]["en"]["old_pass"]
	G_Res["en"]["/spacepass.html"]["spacepass_new"] = G_Res["common"]["en"]["new_pass"]
	G_Res["en"]["/spacepass.html"]["spacepass_confirm"] = G_Res["common"]["en"]["confirm_pass"]
	G_Res["en"]["/spacepass.html"]["spacepass_change"] = G_Res["common"]["en"]["change"]
	G_Res["en"]["/spacepass.html"]["spacepass_hint"] = "Note: Please keep your password safe. Forgetting your password cannot be recovered."

}
