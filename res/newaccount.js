
function initPage() {
    gLoadRes()
} 

async function newacnt() {
    var pa = document.getElementById("pass").value
    var repa = document.getElementById("repass").value
     window.doWebLog(repa)
    if (pa != repa){
        var msg = await window.doGetMsg("30")
        MYDIALOG.ShowDiv(msg);
        return
    }
    var name = document.getElementById("name").value
    var result = await window.doNewAccount(pa,name)
    var obj = JSON.parse(result)
    MYDIALOG.ShowDiv(obj.Msg);
    initMenu()
}
