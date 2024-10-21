
function initPage() {
    gLoadRes()
} 


async function queryByHash(){
    var hs = document.getElementById("txhash").value
    if (hs != ""){
        var result = await window.doQueryTx(hs)
        var objJson = JSON.parse(result)
        document.getElementById("txinfo").value = JSON.stringify(objJson, null, 2)
    } else {
        var msg = await window.doGetMsg("29")
        MYDIALOG.ShowDiv(msg)
    }
}


