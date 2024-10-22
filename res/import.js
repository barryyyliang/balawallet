async function importAct() {
    var pa = document.getElementById("pass").value
    if (pa == ""){
        var msg = await window.doGetMsg("29")
        MYDIALOG.ShowDiv(msg);
        return
    }
    var repa = document.getElementById("repass").value
    if(repa == ""){
        var msg = await window.doGetMsg("29")
        MYDIALOG.ShowDiv(msg);
        return
    }
    if (pa != repa){
        var msg = await window.doGetMsg("30")
        MYDIALOG.ShowDiv(msg);
        return
    }
    var name = document.getElementById("name").value
    var kk = document.getElementById("kk").value
    var result = await window.doImportAccount(pa,name,kk)
    initMenu()
    showResult(result)
}

function initPage(){
    gLoadRes()
}
