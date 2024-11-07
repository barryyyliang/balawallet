    var number=0
    
    function addNew(obj){
        var table = document.getElementById("assetTable")
        var tr = addNewTr()
        table.append(tr)
    }


    function addNewTr(){
        var tr = document.createElement("tr")
        var ipt0 = document.createElement("select")
        ipt0.setAttribute("id","type"+number)
        ipt0.setAttribute("style","font-size: 16px;height: 30px;width: 200px;border-radius: 0px;")
        var option = document.createElement("option")
        option.text = "ERC20"
        option.value = "ERC20"
        ipt0.add(option)
        var option2 = document.createElement("option")
        option2.text = "ERC721"
        option2.value = "ERC721"
        ipt0.add(option2)
        
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
        bt.setAttribute("onclick","importContract("+"'address"+number +"'" +",'type"+number+"')")
        td4.append(bt) 
        tr.append(td4)
        
        number++
        return tr
    }


    async function importContract(address,type){
        var add = document.getElementById(address).value
        var ty = document.getElementById(type).value
        var fromadd = document.getElementById("add").value
        var obj = {type:ty,address:add,name:'',from:fromadd}
        var result = await window.doImportContract(JSON.stringify(obj))
        var reobj = JSON.parse(result)
        if(reobj.Code === 0){
            initData()
        } 
        showResult(result)
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


