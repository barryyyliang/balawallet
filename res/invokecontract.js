    var abijson

    function addNew(table,obj){
        var tr = newTr(obj)
        window.doWebLog("get new tr")
        table.append(tr)
    }

    function newTr(obj){
        var tr = document.createElement("tr")
        tr.setAttribute("style","height:50px;")


        var td1 = document.createElement("td")
        td1.innerText = obj.name
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("style","width:300px;font-size: 14px;height:25px;")
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
        var filepath = document.getElementById("filePath").value
        var strjson = await window.doLoadABI(filepath)
        var obj =  JSON.parse(strjson)
        if (obj.Code === 0){
            window.doWebLog(obj.Err)
            try {
                abijson = JSON.parse(obj.Err)
            } catch (e){
                window.doWebLog("js parser error");
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
                    window.doWebLog(abijson[i].name);
                }                
            }            
        } else {
            MYDIALOG.ShowDiv("MyDiv","fade", obj.Err);
        }
    }

    function changeFunc(){
        window.doWebLog("changeFunc");
        var func = document.getElementById("contractFunc").value
        window.doWebLog("changeFunc2");
        for (var i=0; i<abijson.length;i++){
            window.doWebLog("changeFunc3");
            if (abijson[i].type === "function" && abijson[i].name === func){
                window.doWebLog("changeFunc4");
                document.getElementById("functype").value = abijson[i].stateMutability
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
                break               
            } 
        } 
    }

    function invoke() {
        window.doWebLog("invoke0")
        //var functype = document.getElementById("functype").value
        var funcName = document.getElementById("contractFunc").value
        var i = 0
        for ( ; i<abijson.length;i++){
            if (abijson[i].type === "function" && abijson[i].name === funcName){
                break               
            } 
        } 
        var tmp = new Array()
        tmp.push(abijson[i])
        var func = JSON.stringify(tmp)
        var parems = new Array()
        var inputTable = document.getElementById("paramterTable")
        window.doWebLog("invoke1")
        var caddress = document.getElementById("caddress").value
        window.doWebLog(caddress)
        var from = document.getElementById("add").value
        window.doWebLog(from)
        for (var n=1;n<inputTable.rows.length;n++){
            parems.push(inputTable.rows[n].cells[1].childNodes[0].value)
        }

        var paremsString = JSON.stringify(parems)
        window.doReadContract(from,caddress,func, paremsString)

        // if (functype === "view") {            

        // } else {
            
        // }
    }

    function initPage(){
        window.doWebLog("init page invokecontract")
    }


