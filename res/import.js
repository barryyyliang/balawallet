async function importAct() {
    window.doWebLog("进入函数importAct")
    var pa = document.getElementById("pass").value
    if (pa == ""){
        MYDIALOG.ShowDiv("请输入密码");
        return
    }
    window.doWebLog("密码：")
    window.doWebLog(pa)
    var repa = document.getElementById("repass").value
    if(repa == ""){
        MYDIALOG.ShowDiv("请输入确认密码");
        return
    }
    window.doWebLog("重输密码：")
    window.doWebLog(repa)
    if (pa != repa){
        window.doWebLog("dont match")
        MYDIALOG.ShowDiv("两次输入的密码不一致");
        return
    }
    var name = document.getElementById("name").value
    var kk = document.getElementById("kk").value
    var result = await window.doImportAccount(pa,name,kk)
    showResult(result)
}

function initPage(){
    gLoadRes()
}
