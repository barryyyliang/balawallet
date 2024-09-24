
async function onSaveConfig() {
    var address = document.getElementById("nodeAddress").value
    window.doWebLog("onSaveConfig0")
    var cDir = document.getElementById("contractDir").value
    window.doWebLog("onSaveConfig1")
    var json = JSON.stringify({"nodeAddress":address,"contractDir":cDir })
    window.doWebLog(json)

    var result = await window.doConfig(json)
    showResult(result)
}

function initPage() {

    var jsonconfig = JSON.stringify(["nodeAddress","contractDir"])
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
   