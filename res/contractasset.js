    function addNew(){
        window.doWebLog("addnew")
        var table = document.getElementById("contractTable")
        window.doWebLog("get table")
        var tr = newTr()
        window.doWebLog("get new tr")
        table.append(tr)
    }

    function newTr(){
        var tr = document.createElement("tr")
        tr.setAttribute("style","height:50px;")

        var ipt0 = document.createElement("input")
        ipt0.setAttribute("style","font-size: 14px;height:25px;")
        ipt0.setAttribute("type","text")
        ipt0.setAttribute("size","10")   
        var td1 = document.createElement("td")
        td1.append(ipt0)
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("id","address")
        ipt.setAttribute("style","font-size: 14px;height:25px;")
        ipt.setAttribute("type","text")
        ipt.setAttribute("size","45") 
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)
        
        var ipt2 = document.createElement("input")
        ipt2.setAttribute("id","name")
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
        bt.setAttribute("value","管理")
        //bt.setAttribute("onclick","changeName("+"'address"+i +"'" +",'name"+i+"')")
        td4.append(bt) 
        tr.append(td4)

        return tr
    }

    function initPage(){
        window.doWebLog("init page contractasset")
    }


