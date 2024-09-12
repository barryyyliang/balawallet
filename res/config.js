
async function onSaveConfig() {
    var address = document.getElementById("nodeAddress").value

    var json = JSON.stringify({'nodeAddress':address})
    window.doWebLog(json)

    var result = await window.doConfig(json)
    showResult(result)
}

function initPage() {

    var jsonconfig = JSON.stringify(["nodeAddress","usernames"])
    window.doWebLog(jsonconfig)
    
    window.doGetConfig(jsonconfig).then( (result) => {
        window.doWebLog("get result" + result)
        var acc = JSON.parse(result)
        for (var val in acc){
            document.getElementById(val).value = acc[val] 
        }      
    }
    )
} 
   