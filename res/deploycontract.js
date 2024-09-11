
function initPage() {
        
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
    var file = document.getElementById("filePath").value
    window.doWebLog("filename:" + file)
    var result = await window.doDeployContract(address,file,"111")
    var jsonObj = JSON.parse(result)
    if (jsonObj.Code === 0 ){
        MYDIALOG.ShowDiv("部署交易发送成功，交易 HASH：" + jsonObj.Err)
    } else {
        MYDIALOG.ShowDiv("部署交易发送失败：" + jsonObj.Err)
    }
}

