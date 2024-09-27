    var abijson

    function addNew(table,obj){
        var tr = newTr(obj)
        window.doWebLog("get new tr")
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
                MYDIALOG.ShowDiv("js parser error");
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
                for (var n=inputTable.rows.length-1;n>0;n--){
                    inputTable.removeChild(inputTable.rows[n])
                }
                for (var n=0; n<abijson[i].inputs.length;n++){
                    addNew(inputTable,abijson[i].inputs[n])
                } 
                var outputTable = document.getElementById("outTable")
                for (var n=outputTable.rows.length-1;n>0;n--){
                    outputTable.removeChild(outputTable.rows[n])
                }
                for (var n=0; n<abijson[i].outputs.length;n++){
                    addNew(outputTable,abijson[i].outputs[n])
                }
                
                window.doWebLog(abijson[i].stateMutability )
                if (abijson[i].stateMutability === "view"){                    
                    document.getElementById("btInvoke").innerHTML = "查询"
                    document.getElementById("outTable").style.visibility =  "visible"
                    document.getElementById("assTable").style.visibility = "hidden"
                } else {
                    document.getElementById("btInvoke").innerHTML = "调用"
                    document.getElementById("outTable").style.visibility =  "hidden"
                    document.getElementById("assTable").style.visibility = "visible"
                }
                break               
            } 
        } 
    }

    async function invoke() {
        window.doWebLog("invoke0")
        var functype
        var funcName = document.getElementById("contractFunc").value
        var caddress = document.getElementById("caddress").value
        var from = document.getElementById("add").value
        if(caddress === ""){
            MYDIALOG.ShowDiv("合约地址为空")
            return
        }
        if(funcName === ""){
            MYDIALOG.ShowDiv("合约参数为空，先加载 abi 文件")
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

        for (var n=1;n<inputTable.rows.length;n++){
            parems.push(inputTable.rows[n].cells[1].childNodes[0].value)
        }

        var paremsString = JSON.stringify(parems)       

        if (functype === "view") {  
            var result = await window.doReadContract(from,caddress,func, paremsString)
            var obj = JSON.parse(result)
            if (obj.Code === 0) {
                window.doWebLog("invoke 2")
                var outsTable = document.getElementById("outTable")
                var outs = JSON.parse(obj.Msg)
                window.doWebLog("invoke 3:"+outs)
                for (var n=1;n<outsTable.rows.length;n++){
                    if (n-1 < outs.length) {
                        outsTable.rows[n].cells[1].childNodes[0].value = outs[n-1]
                    } else {
                        MYDIALOG.ShowDiv("输出参数少于需要的数量")
                    }                    
                }
            } else {
                MYDIALOG.ShowDiv(obj.Msg)
            }
        } else {
            var ass = document.getElementById("ass").value
            if(ass === ""){
                MYDIALOG.ShowDiv("发送交易需要账户密码，不能为空")
                return
            }
            var result = await window.doInvokeContract(from,caddress,func, paremsString,ass)
            var obj = JSON.parse(result)
            if(obj.Code === 0){
                MYDIALOG.ShowDiv("交易发送成功，\r\n 交易 ID：\r\n"+obj.Msg)
            } else {
                MYDIALOG.ShowDiv(obj.Msg)
            }

        }
    }

    async function initPage(){
        window.doWebLog("init page invokecontract")
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


