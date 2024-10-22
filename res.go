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
	G_Res["zh-TW"]["/main.html"]["main_about"] = "關於 Web3Wallet"

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

	G_Res["zh-TW"]["/qianbao.html"] = make(map[string]string)
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_type"] = "類型"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_address"] = "地址"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_name"] = "別名"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_balance"] = "餘額"
	G_Res["zh-TW"]["/qianbao.html"]["qianbao_changebutton"] = "修改別名"

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

	G_Res["zh-TW"]["/contractasset.html"] = make(map[string]string)
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_type"] = "合约類型"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_address"] = "合约地址"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_name"] = "資產名稱"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_add"] = "添加"
	G_Res["zh-TW"]["/contractasset.html"]["contractasset_import"] = "導入"

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

	G_Res["en"]["/deploycontract.html"] = make(map[string]string)
	G_Res["en"]["/deploycontract.html"]["deploycontract_file"] = "Contract Bin File"
	G_Res["en"]["/deploycontract.html"]["deploycontract_acc"] = "Password"
	G_Res["en"]["/deploycontract.html"]["deploycontract_deploy"] = "Deploy"

	G_Res["zh-TW"]["/deploycontract.html"] = make(map[string]string)
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_file"] = "合約文件（bin)"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_acc"] = "帳戶密碼"
	G_Res["zh-TW"]["/deploycontract.html"]["deploycontract_deploy"] = "部署"

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

	//deploycontract
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

	G_Res["en"]["/signature.html"] = make(map[string]string)
	G_Res["en"]["/signature.html"]["signature_infomation"] = "Information to be signed"
	G_Res["en"]["/signature.html"]["signature_acc"] = "Password"
	G_Res["en"]["/signature.html"]["signature_sign"] = "Sign"

	G_Res["zh-TW"]["/signature.html"] = make(map[string]string)
	G_Res["zh-TW"]["/signature.html"]["signature_infomation"] = "待簽名訊息"
	G_Res["zh-TW"]["/signature.html"]["signature_acc"] = "帳號密碼"
	G_Res["zh-TW"]["/signature.html"]["signature_sign"] = "簽名"

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
	G_Res["zh-CN"]["/config.html"]["config_address"] = "节点地址"
	G_Res["zh-CN"]["/config.html"]["config_contractdir"] = "合约文件目录"
	G_Res["zh-CN"]["/config.html"]["config_language"] = "系统语言"
	G_Res["zh-CN"]["/config.html"]["config_save"] = "保存"

	G_Res["en"]["/config.html"] = make(map[string]string)
	G_Res["en"]["/config.html"]["config_address"] = "Node Address"
	G_Res["en"]["/config.html"]["config_contractdir"] = "Contract File Dir"
	G_Res["en"]["/config.html"]["config_language"] = "Language"
	G_Res["en"]["/config.html"]["config_save"] = "Save"

	G_Res["zh-TW"]["/config.html"] = make(map[string]string)
	G_Res["zh-TW"]["/config.html"]["config_address"] = "節點位址"
	G_Res["zh-TW"]["/config.html"]["config_contractdir"] = "合約檔案目錄"
	G_Res["zh-TW"]["/config.html"]["config_language"] = "系統語言"
	G_Res["zh-TW"]["/config.html"]["config_save"] = "儲存"

	//about
	G_Res["zh-CN"]["/about.html"] = make(map[string]string)
	G_Res["zh-CN"]["/about.html"]["about_about"] = `web3wallet 提供数字资产管理，智能合约调用。
	
    本软件暂不提供节点服务，要连接到区块链网络，用户可以连接第三方区块链节点或者 RPC 服务

    RPC厂商注册后，可提供一定量的免费服务,如下 RPC 厂商可提供相关服务：
    https://www.infura.io/zh/networks/ethereum
	https://etherscan.io/apis

    使用的的第三方开源项目：    
    https://github.com/ethereum/go-ethereum.git
    https://github.com/webview/webview_go.git

    使用建议和问题反馈：
    barryyyliang@163.com`

	G_Res["en"]["/about.html"] = make(map[string]string)
	G_Res["en"]["/about.html"]["about_about"] = `web3wallet provides digital asset management and smart contract calling functions。
    
	This software does not currently provide node services. To connect to the blockchain network, users can connect to third-party blockchain nodes or RPC services.
	
	After RPC vendors register, they can provide a certain amount of free services, The following RPC vendors can provide related services:
    https://www.infura.io/zh/networks/ethereum
	https://etherscan.io/apis
    
	Third-party open source projects used:    
    https://github.com/ethereum/go-ethereum.git
    https://github.com/webview/webview_go.git
    
	Usage suggestions and problem feedback:
    barryyyliang@163.com`

	G_Res["zh-TW"]["/about.html"] = make(map[string]string)
	G_Res["zh-TW"]["/about.html"]["about_about"] = `web3wallet 提供數位資產管理，智慧合約呼叫功能。
	
    本軟體暫不提供節點服務，要連接到區塊鏈網絡，用戶可以連接第三方區塊鏈節點或 RPC 服務

    RPC廠商註冊後，可提供一定量的免費服務,如下 RPC 廠商可提供相關服務：
    https://www.infura.io/zh/networks/ethereum
	https://etherscan.io/apis

    使用的第三方開源專案：    
    https://github.com/ethereum/go-ethereum.git
    https://github.com/webview/webview_go.git

    使用建議和問題回饋：
    barryyyliang@163.com`

}
