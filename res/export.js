 function exportAct() {
    window.doWebLog("enter exportAct")
    var pa = document.getElementById("pass").value
    window.doWebLog("pa value:"+pa)
    var va = document.getElementById("add").value
    window.doWebLog("add value:" + va)  
    window.doExportAccount(va, pa)
}

function initPage(){
    window.doGetAccounts().then( (result) => {
        window.doWebLog("get result" + result)
        var acc = JSON.parse(result)
        window.doWebLog("go to addelement")
        var se = document.getElementById("acc")
        window.doWebLog("get select addelement")
        for (i=0; i<acc.accounts.length; i++){
            window.doWebLog(acc.accounts[i])
              var option = document.createElement("option")
            option.text = acc.accounts[i].address
            option.value = acc.accounts[i].address
            window.doWebLog("add:" + acc.accounts[i].address)
            se.add(option)
        }            
    })
}