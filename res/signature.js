
function initPage() {
        
} 


async function signature(){
    var param={}
    var info = document.getElementById("info").value
    param.Address = getCurrentAccount()
    if(info != "") {
        param.Info = info
    } else {
        MYDIALOG.ShowDiv("未输入信息")
        return
    }
    window.doWebLog("web: sinature1")
    var result = await window.doSignature(JSON.stringify(param))
    var objJson = JSON.parse(result)
    if (objJson.Code === 0 ){
        document.getElementById("infosign").value = objJson.Msg
    } else {
        MYDIALOG.ShowDiv("签名失败")
    }
}


