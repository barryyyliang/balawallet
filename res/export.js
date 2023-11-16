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
        var acc = result.split(',')
        window.doWebLog("go to addelement")
        var se = document.getElementById("acc")
        window.doWebLog("get select addelement")
        for (i=0; i<acc.length; i++){
            window.doWebLog(acc[i])
            var par = acc[i].split('-')
            var option = document.createElement("option")
            option.text = par[0]
            option.value = par[0]
            window.doWebLog("add:" + par[0])
            se.add(option)
        }            
    })
}