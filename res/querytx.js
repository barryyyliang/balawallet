
function initPage() {
    gLoadRes()
} 


async function queryByHash(){
    window.doWebLog("web:queryByHash")
    var hs = document.getElementById("txhash").value
    if (hs != ""){
        var result = await window.doQueryTx(hs)
        var objJson = JSON.parse(result)
        window.doWebLog(JSON.stringify(objJson))
        document.getElementById("txinfo").value = JSON.stringify(objJson, null, 2)
    } else {
        var msg = await doGetMsg("29")
        MYDIALOG.ShowDiv(msg)
    }
}


