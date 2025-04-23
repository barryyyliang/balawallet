
async function initPage() {
    gLoadRes()
    var result = await window.doGetBinFile()
    var jsonObj = JSON.parse(result)
    if (jsonObj.Code === 0 ){
        if (jsonObj.Msg.length != 0){
            var names = JSON.parse(jsonObj.Msg)
            var seFile = document.getElementById("binFile")
            for(var i=0;i<names.length;i++){
                var option = document.createElement("option")
                option.text = names[i]
                option.value = names[i]
                seFile.add(option)
            }
        }

    } else {
        MYDIALOG.ShowDiv(jsonObj.Msg)
    }
    deploy_func_loadABI()
} 

function selectfile(){
    document.getElementById("contractFile").click()
}

function changedfile(){
    document.getElementById("filePath").value = document.getElementById("contractFile").value
}


async function deploy_func_loadABI(){
    var seFile = document.getElementById("deploycontractabiFile")
    seFile.options.length = 0
    var option = document.createElement("option")
    option.text = ""
    option.value = ""
    seFile.add(option)
    var result = await window.doGetAbiFile()
    var jsonObj = JSON.parse(result)
    if (jsonObj.Code === 0 ){
        var names = JSON.parse(jsonObj.Msg)

        for(var i=0;i<names.length;i++){
            var option = document.createElement("option")
            option.text = names[i]
            option.value = names[i]
            seFile.add(option)
        }
    } else {
        MYDIALOG.ShowDiv(jsonObj.Msg)
    }
    deploycontract_func_abichange()
}

async function deploycontract_func_abichange(){
    window.doWebLog("deploycontract_func_abichange1")
    var inputTable = document.getElementById("deploycontract_paramterTable")
    for (var n=inputTable.rows.length-1;n>1;n--){
        inputTable.removeChild(inputTable.rows[n])
    }


    document.getElementById("deploycontract_parameter_div").style.visibility = "hidden"
    document.getElementById("deploycontract_parameter_div").style.height = "0px"





    var filepath = document.getElementById("deploycontractabiFile").value
    if (filepath === ""){
        window.doWebLog("deploycontract_func_abichange2")
        return
    }
    var strjson = await window.doLoadABI(filepath)
    var obj =  JSON.parse(strjson)
    if (obj.Code === 0){
        var abijson
        try {
            abijson = JSON.parse(obj.Msg)
        } catch (e){
            MYDIALOG.ShowDiv(e.message);
            return
        }
        window.doWebLog("deploycontract_func_abichange3")
        for (var i=0; i<abijson.length;i++){
            if (abijson[i].type === "constructor"){
                window.doWebLog("deploycontract_func_abichange4")
                document.getElementById("deploycontract_parameter_div").style.visibility = "visible"
                document.getElementById("deploycontract_parameter_div").style.height = ""
                for (var n=0; n<abijson[i].inputs.length;n++){
                    addNew(inputTable,abijson[i].inputs[n])
                } 
            }                
        }                     
    } else {
        MYDIALOG.ShowDiv(obj.Msg);
    }
}

function addNew(table,obj){
    var tr = newTr(obj)
    table.append(tr)
}
function newTr(obj){
    var tr = document.createElement("tr")

    var td1 = document.createElement("td")
    td1.setAttribute("style","width:100px;")
    td1.innerText = obj.name
    tr.append(td1)

    var ipt = document.createElement("input")
    ipt.setAttribute("type","text")
    var td2 = document.createElement("td")
    td2.append(ipt)
    tr.append(td2)

    var td3 = document.createElement("td")
    td3.innerText = obj.type
    tr.append(td3)   
    return tr  
}

async function deploy(){
    var address = document.getElementById("add").value
    var file = document.getElementById("binFile").value
    var ass = document.getElementById("deploycontract_ass").value
    var result = await window.doDeployContract(address,file,ass)
    var jsonObj = JSON.parse(result)
    if (jsonObj.Code === 0 ){
        showMsg("34",jsonObj.Msg)
    } else {
        MYDIALOG.ShowDiv(jsonObj.Msg)
    }
}

