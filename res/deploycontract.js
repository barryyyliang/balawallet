
async function initPage() {
    gLoadRes()
    var result = await window.doGetBinFile()
    var jsonObj = JSON.parse(result)
    if (jsonObj.Code === 0 ){
        var names = JSON.parse(jsonObj.Msg)
        var seFile = document.getElementById("binFile")
        for(var i=0;i<names.length;i++){
            var option = document.createElement("option")
            option.text = names[i]
            option.value = names[i]
            seFile.add(option)
        }
    } else {
        MYDIALOG.ShowDiv(jsonObj.Msg)
    }
} 

function selectfile(){
    document.getElementById("contractFile").click()
}

function changedfile(){
    document.getElementById("filePath").value = document.getElementById("contractFile").value
}
async function deploy(){
    window.doWebLog("deploy contract")
    var address = document.getElementById("add").value
    var file = document.getElementById("binFile").value
    var ass = document.getElementById("ass").value
    window.doWebLog("filename:" + file)
    var result = await window.doDeployContract(address,file,ass)
    var jsonObj = JSON.parse(result)
    if (jsonObj.Code === 0 ){
        MYDIALOG.ShowDiv("部署交易发送成功，交易 HASH：" + jsonObj.Msg)
    } else {
        MYDIALOG.ShowDiv("部署交易发送失败：" + jsonObj.Msg)
    }
}

