async function spacepass_func_change() {
    var pass = document.getElementById("spacepass_input_new").value 
    var confirm = document.getElementById("spacepass_input_confirm").value
    var old = document.getElementById("spacepass_input_old").value
    
    if (pass === "" || pass === "null" || pass.length < 8){
        MYDIALOG.ShowDiv(await window.doGetMsg("39"))
        return
    }

    if (old === "" || old === "null" || old.length < 8){
        MYDIALOG.ShowDiv(await window.doGetMsg("39"))
        return
    }

    if (pass != confirm ){
        MYDIALOG.ShowDiv(await window.doGetMsg("30"))
        return
    }
    var result = await window.doChangeSpacePass(
        document.getElementById("spacepass_input_old").value,
        pass
    )

    showResult(result)
}

function initPage(){
    gLoadRes()
}
