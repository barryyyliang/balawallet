    var number=0
    
    function addNew(obj){
        var table = document.getElementById("assetTable")
        var tr = addNewTr()
        table.append(tr)
    }

    function initOption(ipt){
        var net = getMainNet()
        if (net === "Ethereum" || net === "Polygon"){
            var option = document.createElement("option")
            option.text = "ERC20"
            option.value = "ERC20"
            ipt.add(option)
            var option2 = document.createElement("option")
            option2.text = "ERC721"
            option2.value = "ERC721"
            ipt.add(option2)
            var option3 = document.createElement("option")
            option3.text = "Other"
            option3.value = "Other"
            ipt.add(option3)
        } else if(net === "Solana") {
            var option = document.createElement("option")
            option.text = "SPL"
            option.value = "SPL"
            ipt.add(option)
        }
    }

    function addNewTr(){
        var tr = document.createElement("tr")
        var ipt0 = document.createElement("select")
        ipt0.setAttribute("id","type"+number)
        ipt0.setAttribute("style","font-size: 16px;height: 30px;width: 200px;border-radius: 0px;")
        initOption(ipt0)
        
        var td1 = document.createElement("td")
        td1.append(ipt0)
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("id","address"+number)
        ipt.setAttribute("type","text")
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)
        
        var ipt2 = document.createElement("input")
        ipt2.setAttribute("id","name"+number)
        ipt2.setAttribute("type","text")
        ipt2.setAttribute("style","width:200px")
   
        var td3 = document.createElement("td")
        td3.append(ipt2)
        tr.append(td3)

        var td4 = document.createElement("td")               
        var bt = document.createElement("input")
        //bt.setAttribute("style","width:50px")
        bt.setAttribute("type","button")
        bt.setAttribute("class","button1")
        bt.setAttribute("value",gGetRes("contractasset_import"))
        bt.setAttribute("onclick","importContract("+"'address"+number + "'" + "," +"'type"+number+"'" + "," + "'name"+number +"'" + ")")
        td4.append(bt) 
        tr.append(td4)

        number++
        return tr
    }


    async function importContract(address,type,name){
        var add = document.getElementById(address).value
        var ty = document.getElementById(type).value
        var fromadd = document.getElementById("add").value
        var na = document.getElementById(name).value
        var obj = {type:ty,address:add,name:na,from:fromadd}
        var result = await window.doImportContract(JSON.stringify(obj))
        var reobj = JSON.parse(result)
        if(reobj.Code === 0){
            initData()
        } 
        showResult(result)
    }

    async function contractasset_func_delete(address){        
        var conf = await MYDIALOG.ShowConfirmDiv(await window.doGetMsg("42"))
        if(conf){
            var add = document.getElementById(address).value
            var result = await window.doDeleteContract(add)
            var reobj = JSON.parse(result)
            if(reobj.Code === 0){
                initData()
            } 
            showResult(result)
        } 
    }

    function contractasset_func_invoke(address){
        gParam = {}
        gParam.invokeContract = {}
        gParam.invokeContract.address = document.getElementById(address).value
        changePage("tools","invokecontract")
    }

    function contractasset_func_showfunc(id){
        var p1 = document.getElementById(id).parentNode.parentNode
        var cells = p1.getElementsByTagName("td")
        var address = cells[1].childNodes[0].id

        cells[0].childNodes[0].classList.add("highlight")
        cells[1].childNodes[0].classList.add("highlight")
        cells[2].childNodes[0].classList.add("highlight")
    
        clearFunc()
        addFuncToMenu("contractasset_func_invoke('"+address+"')", gGetRes("contractasset_invoke"))
        addFuncToMenu("contractasset_func_delete('"+address+"')", gGetRes("contractasset_delete"))
        onFuncMenu(id)
    }

    function contractasset_func_mouseout(id){
        var p1 = document.getElementById(id).parentNode.parentNode
        var cells = p1.getElementsByTagName("td")

        cells[0].childNodes[0].classList.remove("highlight")
        cells[1].childNodes[0].classList.remove("highlight")
        cells[2].childNodes[0].classList.remove("highlight")
    }


    function initTr(obj){
        var tr = document.createElement("tr")

        var ipt0 = document.createElement("input")
        ipt0.setAttribute("id","type"+number)
        ipt0.setAttribute("style","width: 200px;")
        ipt0.setAttribute("readonly","readonly")
        var td1 = document.createElement("td")
        td1.append(ipt0)
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("id","address"+number)
        ipt.setAttribute("type","text")
        ipt.setAttribute("readonly","readonly")
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)
        
        var ipt2 = document.createElement("input")
        ipt2.setAttribute("id","name"+number)
        ipt2.setAttribute("style","width:200px")
        ipt2.setAttribute("type","text")
        ipt2.setAttribute("size","20")  
        ipt2.setAttribute("readonly","readonly")
        var td3 = document.createElement("td")
        td3.append(ipt2)
        tr.append(td3)

        ipt0.setAttribute("value",obj.type) 
        ipt.setAttribute("value",obj.address) 
        ipt2.setAttribute("value",obj.name) 

        // var td4 = document.createElement("td")               
        // var bt = document.createElement("input")
        // //bt.setAttribute("style","width:50px")
        // bt.setAttribute("type","button")
        // bt.setAttribute("class","button1")
        // bt.setAttribute("value",gGetRes("contractasset_delete"))
        // bt.setAttribute("onclick","contractasset_func_delete("+"'address"+number+"')")
        // td4.append(bt) 
        // tr.append(td4)

        // var td5 = document.createElement("td")               
        // var bt2 = document.createElement("input")
        // //bt.setAttribute("style","width:50px")
        // bt2.setAttribute("type","button")
        // bt2.setAttribute("class","button1")
        // bt2.setAttribute("value",gGetRes("contractasset_invoke"))
        // bt2.setAttribute("onclick","contractasset_func_invoke("+"'address"+number+"')")
        // td5.append(bt2) 
        // tr.append(td5)

        var td6 = document.createElement("td")               
        var bt3 = document.createElement("input")
        bt3.setAttribute("class","button1")
        bt3.setAttribute("id","contractasset_button"+number)
        bt3.setAttribute("style","width:50px")
        bt3.setAttribute("type","button")
        bt3.setAttribute("value","↓")
        bt3.setAttribute("onmouseover","contractasset_func_showfunc("+"'contractasset_button"+number+"')")
        bt3.setAttribute("onmouseout","contractasset_func_mouseout("+"'contractasset_button"+number+"')")
        td6.append(bt3) 
        tr.append(td6)




        number++
        return tr
    }

    function initData(){
        window.doGetContract().then( (result) => {
            if (result != ""){
                var obj = JSON.parse(result)
                var table = document.getElementById("assetTable")
                for (var n=table.rows.length-1;n>0;n--){
                    table.removeChild(table.rows[n])
                }
                for(var i = 0; i < obj.contracts.length;i++)
                {                    
                    var tr = initTr(obj.contracts[i])
                    table.append(tr)
                }
            }
        })
    }

    function initPage(){
        gLoadRes()
        initData()
    }


