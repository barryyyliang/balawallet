
async function onSaveConfig() {
    var obj = {}
    obj["nodeAddress"] = document.getElementById("nodeAddress").value    
    obj["contractDir"] = document.getElementById("contractDir").value
    obj["language"] = document.getElementById("language").value
    var json = JSON.stringify(obj)
    var result = await window.doConfig(json)
    showResult(result)
}

function initPage() {
    gLoadRes()
    window.doGetConfig("").then( (result) => {
        var acc = JSON.parse(result)
        for (var val in acc){
            document.getElementById(val).value = acc[val] 
        }      
    }
    )
} 
   