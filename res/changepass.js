async function changepass_func_change() {
    var pass = document.getElementById("changepass_input_new").value 
    var confirm = document.getElementById("changepass_input_confirm").value
    
    if (pass === "" || pass === "null" || pass.length < 8){
        MYDIALOG.ShowDiv(await window.doGetMsg("39"))
        return
    }
    if (pass != confirm ){
        MYDIALOG.ShowDiv(await window.doGetMsg("30"))
        return
    }
    var result = await window.doChangeAccountPass(
        document.getElementById("changepass_input_type").value,
        document.getElementById("changepass_input_address").value,
        document.getElementById("changepass_input_old").value,
        pass
    )

    showResult(result)
}

function initPage(){
    gLoadRes()
    document.getElementById("changepass_input_type").value = gParam.type
    document.getElementById("changepass_input_address").value = gParam.address
    document.getElementById("changepass_input_name").value = gParam.name
}
