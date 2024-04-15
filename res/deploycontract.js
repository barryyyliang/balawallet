
function initPage() {
        
} 

function selectfile(){
    document.getElementById("contractFile").click()
}

function changedfile(){
    document.getElementById("filePath").value = document.getElementById("contractFile").value
}
function deploy(){
    window.doWebLog("deploy contract")
    var address = document.getElementById("add").value
    var file = document.getElementById("filePath").value
    window.doWebLog("filename:" + file)
    window.doDeployContract(address,file,"111")
}

