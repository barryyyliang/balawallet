package main

var G_Msg = make(map[string]map[int]string)

func initMsg() {
	G_Msg["zh"] = make(map[int]string)
	G_Msg["zh"][0] = "操作成功"
	G_Msg["zh"][1] = "创建账户失败"
	G_Msg["zh"][2] = "创建账户成功，写名称失败，可在钱包页面修改名称"
	G_Msg["zh"][3] = "创建账户成功，名称已存在，可在钱包页面修改名称"
	G_Msg["zh"][4] = "写数据库失败"
	G_Msg["zh"][5] = "名称已存在"
	G_Msg["zh"][6] = "私钥格式错误，解析私钥失败"
	G_Msg["zh"][7] = "读取密钥文件失败"
	G_Msg["zh"][8] = "解密密钥文件失败"
	G_Msg["zh"][9] = "地址和私钥不匹配"
	G_Msg["zh"][10] = "打开导出文件失败"
}
