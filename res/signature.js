
function initPage() {
    gLoadRes()
} 


async function signature(){
    document.getElementById("infosign").value = ""
    var param={}
    var info = document.getElementById("info").value
    var ass = document.getElementById("sigAss").value
    param.Address = getCurrentAccount()
    if(info != "") {
        param.Info = info
    } else {
        showMsg("29","")
        return
    }
    if(ass != ""){
        param.Ass = ass
    } else {
        showMsg("29","")
        return
    }
    var result = await window.doSignature(JSON.stringify(param))
    var objJson = JSON.parse(result)
    if (objJson.Code === 0 ){
        document.getElementById("infosign").value = objJson.Msg
    } else {
        MYDIALOG.ShowDiv(objJson.Msg)
    }
}

async function verify(){
    var param={}
    var info = document.getElementById("info").value
    var sig = document.getElementById("infosign").value
    param.Address = getCurrentAccount()
    if(info == "" || sig == "") {        
        showMsg("29","")
        return
    } else {
        param.Info = info
        param.Sign = sig
    }

    var result = await window.doSignatureVerify(JSON.stringify(param))
    window.doWebLog(result)
    showResult(result)    
}


