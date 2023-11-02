
function initPage() {
    
} 

function newacnt() {
    window.doWebLog("进入函数")
    var pa = document.getElementById("pass").value
    window.doWebLog("密码：")
    window.doWebLog(pa)
    var repa = document.getElementById("repass").value
    window.doWebLog("重输密码：")
    window.doWebLog(repa)
    if (pa != repa){
        window.doWebLog("dont match")
        return
    }
    var name = document.getElementById("name").value
    window.doNewAccount(pa,name)
}
