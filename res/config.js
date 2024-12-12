
async function onSaveConfig() {
    var obj = {}
    obj["ethAddress"] = document.getElementById("config_input_eth_address").value   
    obj["polAddress"] = document.getElementById("config_input_pol_address").value  
    obj["solAddress"] = document.getElementById("config_input_sol_address").value 
    obj["contractDir"] = document.getElementById("config_input_contractDir").value
    obj["language"] = document.getElementById("config_select_language").value
    var json = JSON.stringify(obj)
    var result = await window.doConfig(json)
    showResult(result)
}

function initPage() {
    gLoadRes()
    window.doGetConfig("").then( (result) => {
        var obj = {}
        obj["ethAddress"] = "config_input_eth_address"    
        obj["polAddress"] = "config_input_pol_address"  
        obj["solAddress"] = "config_input_sol_address" 
        obj["contractDir"] = "config_input_contractDir"
        obj["language"] = "config_select_language"
        var acc = JSON.parse(result)
        for (var val in acc){
            try {
                var ele = document.getElementById(obj[val])
                if (ele != null) {
                    ele.value = acc[val]
                } 
            } catch(err){

            }
        }      
    }
    )
} 
   