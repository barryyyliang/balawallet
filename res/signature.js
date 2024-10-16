
function initPage() {
    gLoadRes()
} 


async function signature(){
    infosign
    document.getElementById("infosign").value = ""
    var param={}
    var info = document.getElementById("info").value
    var ass = document.getElementById("sigAss").value
    param.Address = getCurrentAccount()
    if(info != "") {
        param.Info = info
    } else {
        MYDIALOG.ShowDiv("未输入需要签名的信息")
        return
    }
    if(ass != ""){
        param.Ass = ass
    } else {
        MYDIALOG.ShowDiv("未输入签名账户的密码")
        return
    }
    window.doWebLog("web: sinature1")
    var result = await window.doSignature(JSON.stringify(param))
    var objJson = JSON.parse(result)
    if (objJson.Code === 0 ){
        document.getElementById("infosign").value = objJson.Msg
    } else {
        MYDIALOG.ShowDiv("签名失败:"+objJson.Msg)
    }
}


