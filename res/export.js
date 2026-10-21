 async function exportAct() {
    var pa = document.getElementById("export_input_pass").value
    var va = document.getElementById("export_select_acc").value
    var result = await window.doExportAccount(va, pa)
    showResult(result)
    
}

function initPage(){
    gLoadRes()
    window.doGetAccounts().then( (result) => {
        if (result != ""){
            var acc = JSON.parse(result)
            var se = document.getElementById("export_select_acc")
            for (i=0; i<acc.accounts.length; i++){
                var option = document.createElement("option")
                option.text = acc.accounts[i].address
                option.value = acc.accounts[i].address
                se.add(option)
            } 
        }
           
    })
}