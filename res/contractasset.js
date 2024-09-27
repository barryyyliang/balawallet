    var number=0
    
    function addNew(obj){
        window.doWebLog("addnew")
        var table = document.getElementById("contractTable")
        window.doWebLog("get table")
        var tr = addNewTr()
        window.doWebLog("get new tr")
        table.append(tr)
    }


    function addNewTr(){
        var tr = document.createElement("tr")
        tr.setAttribute("style","height:50px;")

        var ipt0 = document.createElement("select")
        ipt0.setAttribute("id","type"+number)
        ipt0.setAttribute("style","font-size: 14px;height: 30px;width: 200px;border-radius: 0px;")
        var option = document.createElement("option")
        option.text = "ERC20"
        option.value = "ERC20"
        ipt0.add(option)
        var option2 = document.createElement("option")
        option2.text = "NFT"
        option2.value = "NFT"
        ipt0.add(option2)
        
        var td1 = document.createElement("td")
        td1.append(ipt0)
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("id","address"+number)
        ipt.setAttribute("style","font-size: 14px;height:25px;")
        ipt.setAttribute("type","text")
        ipt.setAttribute("size","45") 
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)
        
        var ipt2 = document.createElement("input")
        ipt2.setAttribute("id","name"+number)
        ipt2.setAttribute("style","font-size: 14px;height:25px;")
        ipt2.setAttribute("type","text")
        ipt2.setAttribute("size","20")  
 
        var td3 = document.createElement("td")
        td3.append(ipt2)
        tr.append(td3)

        var td4 = document.createElement("td")               
        var bt = document.createElement("input")
        bt.setAttribute("style","font-size: 14px;height:25px;")
        bt.setAttribute("type","button")
        bt.setAttribute("value","导入")
        bt.setAttribute("onclick","importContract("+"'address"+number +"'" +",'type"+number+"')")
        td4.append(bt) 
        tr.append(td4)


        number++
        return tr
    }


    async function importContract(address,type){
        var add = document.getElementById(address).value
        var ty = document.getElementById(type).value
        var obj = {type:ty,address:add,name:''}

        window.doWebLog("import contract:" + add)
        var result = await window.doImportContract(JSON.stringify(obj))
        showResult(result)
    }

    function initTr(obj){
        var tr = document.createElement("tr")
        tr.setAttribute("style","height:50px;")

        var ipt0 = document.createElement("input")
        ipt0.setAttribute("id","type"+number)
        ipt0.setAttribute("style","font-size: 14px;height:25px;width: 200px;")
        ipt0.setAttribute("readonly","readonly")
        var td1 = document.createElement("td")
        td1.append(ipt0)
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("id","address"+number)
        ipt.setAttribute("style","font-size: 14px;height:25px;")
        ipt.setAttribute("type","text")
        ipt.setAttribute("size","45") 
        ipt.setAttribute("readonly","readonly")
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)
        
        var ipt2 = document.createElement("input")
        ipt2.setAttribute("id","name"+number)
        ipt2.setAttribute("style","font-size: 14px;height:25px;")
        ipt2.setAttribute("type","text")
        ipt2.setAttribute("size","20")  
        ipt2.setAttribute("readonly","readonly")
        var td3 = document.createElement("td")
        td3.append(ipt2)
        tr.append(td3)

        if(JSON.stringify(obj)!="{}"){
            ipt0.setAttribute("value",obj.type) 
            ipt.setAttribute("value",obj.address) 
            ipt2.setAttribute("value",obj.name) 
        } else {    
            var td4 = document.createElement("td")               
            var bt = document.createElement("input")
            bt.setAttribute("style","font-size: 14px;height:25px;")
            bt.setAttribute("type","button")
            bt.setAttribute("value","导入")
            bt.setAttribute("onclick","importContract("+"'address"+number +"'" +",'type"+number+"')")
            td4.append(bt) 
            tr.append(td4)
        }

        number++
        return tr
    }

    function initPage(){
        window.doWebLog("init page contractasset")
        window.doGetContract().then( (result) => {
            if (result != ""){
                var obj = JSON.parse(result)
                window.doWebLog("get result json" + obj)
                for(var i = 0; i < obj.contracts.length;i++)
                {
                    var table = document.getElementById("contractTable")
                    var tr = initTr(obj.contracts[i])
                    table.append(tr)
                }
            }
        })
    }


