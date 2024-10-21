    var abijson

    function addNew(table,obj){
        var tr = newTr(obj)
        table.append(tr)
    }

    function newTr(obj){
        var tr = document.createElement("tr")

        var td1 = document.createElement("td")
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

    async function loadABI(){
        var filepath = document.getElementById("abiFile").value
        var strjson = await window.doLoadABI(filepath)
        var obj =  JSON.parse(strjson)
        if (obj.Code === 0){
            try {
                abijson = JSON.parse(obj.Msg)
            } catch (e){
                MYDIALOG.ShowDiv(e.message);
                return
            }
            var se = document.getElementById("contractFunc")
            se.options.length = 0
            for (var i=0; i<abijson.length;i++){
                if (abijson[i].type === "function"){
                    var option = document.createElement("option")
                    option.text = abijson[i].name
                    option.value = abijson[i].name
                    se.add(option)
                }                
            }
            changeFunc()            
        } else {
            MYDIALOG.ShowDiv(obj.Msg);
        }
    }

    function changeFunc(){
        var func = document.getElementById("contractFunc").value
        for (var i=0; i<abijson.length;i++){
            if (abijson[i].type === "function" && abijson[i].name === func){
                var inputTable = document.getElementById("paramterTable")
                for (var n=inputTable.rows.length-1;n>1;n--){
                    inputTable.removeChild(inputTable.rows[n])
                }
                for (var n=0; n<abijson[i].inputs.length;n++){
                    addNew(inputTable,abijson[i].inputs[n])
                } 
                var outputTable = document.getElementById("outTable")
                for (var n=outputTable.rows.length-1;n>1;n--){
                    outputTable.removeChild(outputTable.rows[n])
                }
                for (var n=0; n<abijson[i].outputs.length;n++){
                    addNew(outputTable,abijson[i].outputs[n])
                }

                if (abijson[i].stateMutability === "view"){                    
                    document.getElementById("invokecontract_invoke").innerHTML = gGetRes("invokecontract_query")
                    document.getElementById("outTableDiv").style.visibility =  "visible"
                    document.getElementById("assDiv").style.visibility = "hidden"
                    document.getElementById("assDiv").style.height = "0px"
                } else {
                    document.getElementById("invokecontract_invoke").innerHTML = gGetRes("invokecontract_invoke")
                    document.getElementById("outTableDiv").style.visibility =  "hidden"
                    document.getElementById("assDiv").style.visibility = "visible"
                    document.getElementById("assDiv").style.height = ""
                }
                break               
            } 
        } 
    }

    async function invoke() {
        var functype
        var funcName = document.getElementById("contractFunc").value
        var caddress = document.getElementById("caddress").value
        var from = document.getElementById("add").value
        if(caddress === ""){
            var msg = await window.doGetMsg("29")
            MYDIALOG.ShowDiv(msg)
            return
        }
        if(funcName === ""){
            var msg = await window.doGetMsg("31")
            MYDIALOG.ShowDiv(msg)
            return
        }
        var i = 0
        for ( ; i<abijson.length;i++){
            if (abijson[i].type === "function" && abijson[i].name === funcName){
                functype = abijson[i].stateMutability
                break               
            } 
        } 
        var tmp = new Array()
        tmp.push(abijson[i])
        var func = JSON.stringify(tmp)
        var parems = new Array()
        var inputTable = document.getElementById("paramterTable")

        for (var n=2;n<inputTable.rows.length;n++){
            parems.push(inputTable.rows[n].cells[1].childNodes[0].value)
        }

        var paremsString = JSON.stringify(parems)       

        if (functype === "view") {  
            var result = await window.doReadContract(from,caddress,func, paremsString)
            var obj = JSON.parse(result)
            if (obj.Code === 0) {
                var outsTable = document.getElementById("outTable")
                var outs = JSON.parse(obj.Msg)
                for (var n=2;n<outsTable.rows.length;n++){
                    if (n-2 < outs.length) {
                        outsTable.rows[n].cells[1].childNodes[0].value = outs[n-2]
                    } else {
                        showMsg("32","")
                    }                    
                }
            } else {
                MYDIALOG.ShowDiv(obj.Msg)
            }
        } else {
            var ass = document.getElementById("invokecontract_ass").value
            if(ass === ""){
                showMsg("33","")
                return
            }
            var result = await window.doInvokeContract(from,caddress,func, paremsString,ass)
            var obj = JSON.parse(result)
            if(obj.Code === 0){
                showMsg("34",obj.Msg)
            } else {
                MYDIALOG.ShowDiv(obj.Msg)
            }

        }
    }

    async function initPage(){
        gLoadRes()
        var result = await window.doGetAbiFile()
        var jsonObj = JSON.parse(result)
        if (jsonObj.Code === 0 ){
            var names = JSON.parse(jsonObj.Msg)
            var seFile = document.getElementById("abiFile")
            for(var i=0;i<names.length;i++){
                var option = document.createElement("option")
                option.text = names[i]
                option.value = names[i]
                seFile.add(option)
            }
        } else {
            MYDIALOG.ShowDiv(jsonObj.Msg)
        }
    }


