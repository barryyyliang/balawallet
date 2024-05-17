    
    function addNew(obj){
        window.doWebLog("addnew")
        var table = document.getElementById("paramterTable")
        window.doWebLog("get table")
        var tr = newTr(obj)
        window.doWebLog("get new tr")
        table.append(tr)
    }

    function newTr(obj){
        var tr = document.createElement("tr")
        tr.setAttribute("style","height:50px;")


        var td1 = document.createElement("td")
        td1.innerText = "pamameter name"
        tr.append(td1)

        var ipt = document.createElement("input")
        //ipt.setAttribute("id","address"+number)
        ipt.setAttribute("style","width:300px;font-size: 14px;height:25px;")
        ipt.setAttribute("type","text")
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)

        var td3 = document.createElement("td")
        td3.innerText = "pamameter type"
        tr.append(td3)   
        return tr  

    }

    async function loadABI(){
        var filepath = document.getElementById("filePath").value
        var strjson = await window.doLoadABI(filepath)
        var obj =  JSON.parse(strjson)
        if (obj.Code === 0){
            window.doWebLog(obj.Err)
            var abijson
            try {
                abijson = JSON.parse(obj.Err)
            } catch (e){
                window.doWebLog("js parser error");
                return
            }

            for (var i=0; i<abijson.length;i++){
                if (abijson[i].type === "function"){
                    window.doWebLog(abijson[i].name);
                }
            }
            
        } else {
            MYDIALOG.ShowDiv("MyDiv","fade", obj.Err);
        }
    }

    function initPage(){
        window.doWebLog("init page invokecontract")
    }


