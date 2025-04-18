    var number=0
    
    function dataspace_fun_addNew(obj){
        var table = document.getElementById("dataspace_table")
        var tr = addNewTr()
        table.append(tr)
    }


    function addNewTr(){
        var tr = document.createElement("tr")
        var ipt = document.createElement("input")
        ipt.setAttribute("id","dataspace_name"+number)
        ipt.setAttribute("type","text")
        ipt.setAttribute("style","width:200px")
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)
        
        var ipt2 = document.createElement("textarea")
        ipt2.setAttribute("id","dataspace_content"+number)
        ipt2.setAttribute("style","width:500px;height:25px;font-size: 16px;overflow:auto")
        var td3 = document.createElement("td")
        td3.append(ipt2)
        tr.append(td3)

        var td5 = document.createElement("td")               
        var bt2 = document.createElement("input")
        bt2.setAttribute("type","button")
        bt2.setAttribute("class","button1")
        bt2.setAttribute("style","width:100px")
        bt2.setAttribute("value",gGetRes("dataspace_recover"))
        bt2.setAttribute("onclick","dataspace_fun_recover("+"'dataspace_content"+number+"')")
        td5.append(bt2) 
        tr.append(td5)

        var td6 = document.createElement("td")               
        var bt3 = document.createElement("input")
        bt3.setAttribute("type","button")
        bt3.setAttribute("class","button1")
        bt3.setAttribute("style","width:100px")
        bt3.setAttribute("value",gGetRes("dataspace_save"))
        bt3.setAttribute("onclick","dataspace_fun_save("+"'dataspace_name"+number+"'" +",'dataspace_content"+number+"')")
        td6.append(bt3) 
        tr.append(td6)
        
        number++
        return tr
    }


    async function dataspace_fun_delete(id){
        var auth = document.getElementById("dataspace_pass").value
        var dataname = document.getElementById(id).value
        var conf = await MYDIALOG.ShowConfirmDiv(await window.doGetMsg("42"))
        if(conf){
            var result = await window.doDeleteDataSpace(dataname,auth)
            var obj = JSON.parse(result)
            MYDIALOG.ShowDiv(obj.Msg)  
            if (obj.Code === 0){
                loadDataSpace()
            }                     
        } 
        
    }

    async function dataspace_fun_show(name,content){
        var auth = document.getElementById("dataspace_pass").value
        var key = document.getElementById(name).value

        var result = await window.doGetDataItem(auth,key)
        var obj = JSON.parse(result)         
        if (obj.Code === 0){
            document.getElementById(content).value = obj.Msg
            g_items[document.getElementById(name).value] = 1
        }  else {
            MYDIALOG.ShowDiv(obj.Msg) 
        }       
    }

    function dataspace_fun_recover(id){
        var data = document.getElementById(id)
        data.style.width = '500px'
        data.style.height = '25px'
    }

    function dataspace_fun_hide(name,content){
        document.getElementById(content).value = "**********"
        g_items[document.getElementById(name).value] = 0

    }

    async function dataspace_fun_save(name,content){
        var auth = document.getElementById("dataspace_pass").value
        var dataname = document.getElementById(name).value
        var datacontent = document.getElementById(content).value
        var result = await window.doSaveDataSpace(dataname,datacontent,auth)
        var obj = JSON.parse(result)
        MYDIALOG.ShowDiv(obj.Msg)
        if(obj.Code === 0){
            loadDataSpace()
        }
        
    }

    async function dataspace_fun_change(name,content){
        var auth = document.getElementById("dataspace_pass").value
        var dataname = document.getElementById(name).value
        var datacontent = document.getElementById(content).value
        if(g_items[dataname] === 0){
            var change_guide = await window.doGetMsg("45")
            MYDIALOG.ShowDiv(change_guide)
            return
        }
        var tips = await window.doGetMsg("41")
        var conf = await MYDIALOG.ShowConfirmDiv(tips + ":" + "\r\n" + datacontent)
        if(conf){
            var result = await window.doChangeDataSpace(dataname,datacontent,auth)
            showResult(result)
        } 
    }

    async function dataspace_fun_load(){
        var ret = await loadDataSpace()
        if (ret == true){
            MYDIALOG.ShowDiv(await window.doGetMsg("0"))
        }
    }
    
    /*
    在每一行显示对应的功能菜单，数据绑定对应行
     */
    function dataspace_func_showfunc(id){
        var p1 = document.getElementById(id).parentNode.parentNode
        var cells = p1.getElementsByTagName("td")
        var name = cells[0].childNodes[0].id
        var content = cells[1].childNodes[0].id

        cells[0].childNodes[0].classList.add("highlight")
        cells[1].childNodes[0].classList.add("highlight")

        clearFunc()
        addFuncToMenu("dataspace_fun_show('"+name+"','"+content+"')", gGetRes("dataspace_show"))
        addFuncToMenu("dataspace_fun_hide('"+name+"','"+content+"')", gGetRes("dataspace_hide"))
        addFuncToMenu("dataspace_fun_recover('"+content+"')", gGetRes("dataspace_recover"))
        addFuncToMenu("dataspace_fun_change('"+name+"','"+content+"')", gGetRes("dataspace_change"))
        addFuncToMenu("dataspace_fun_delete('"+name+"')", gGetRes("dataspace_delete"))
        onFuncMenu(id)
    }



    function dataspace_func_mouseout(event){
        var p1 = document.getElementById(event.target.id).parentNode.parentNode
        var cells = p1.getElementsByTagName("td")

        cells[0].childNodes[0].classList.remove("highlight")
        cells[1].childNodes[0].classList.remove("highlight")
        onFuncOut(event)
    }

    var g_items = {}
    async function loadDataSpace(){
        var auth = document.getElementById("dataspace_pass").value
        var result =await window.doGetDataSpace(auth)
        if (result != ""){
            var obj = JSON.parse(result)
            if (obj.Code == 0 ){
                g_items = JSON.parse(obj.Msg)
                var table = document.getElementById("dataspace_table")
                for (var n=table.rows.length-1;n>0;n--){
                    table.removeChild(table.rows[n])
                }
                for(var val in g_items)
                {                            
                    var tr = initTr(val,g_items[val])
                    g_items[val] = 0
                    table.append(tr)
                }
                return true
            } else {
                showResult(result)
                return false
            }
        }
        return true
    }

    async function dataspace_fun_init(){
        var auth = document.getElementById("dataspace_pass").value
        if (auth === "" || auth.length < 8){
            var msg = await window.doGetMsg("39")
            MYDIALOG.ShowDiv(msg);
        } else {
            var result = await window.doInitDataSpace(auth)
            showResult(result)
            var obj = JSON.parse(result)
            if (obj.Code === 0){
                bt = document.getElementById("dataspace_load")
                bt.innerText = gGetRes("dataspace_load")
                bt.onclick = dataspace_fun_load
            }
        }

    }

    function initTr(name,content){
        var tr = document.createElement("tr")
        var ipt = document.createElement("input")
        ipt.setAttribute("id","dataspace_name"+number)
        ipt.setAttribute("type","text")
        ipt.setAttribute("style","width:200px")
        ipt.setAttribute("value",name)
        ipt.setAttribute("readonly","readonly")  
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)
        
        var ipt2 = document.createElement("textarea")
        ipt2.setAttribute("id","dataspace_content"+number)
        ipt2.value = content
        ipt2.setAttribute("style","width:500px;height:25px;font-size: 16px;overflow:auto")
        var td3 = document.createElement("td")
        td3.append(ipt2)
        tr.append(td3)

        // var td5 = document.createElement("td")               
        // var bt2 = document.createElement("input")
        // bt2.setAttribute("type","button")
        // bt2.setAttribute("class","button1")
        // bt2.setAttribute("style","width:100px")
        // bt2.setAttribute("value",gGetRes("dataspace_recover"))
        // bt2.setAttribute("onclick","dataspace_fun_recover("+"'dataspace_content"+number+"')")
        // td5.append(bt2) 
        // tr.append(td5)

        // var td6 = document.createElement("td")               
        // var bt3 = document.createElement("input")
        // bt3.setAttribute("type","button")
        // bt3.setAttribute("class","button1")
        // bt3.setAttribute("style","width:100px")
        // bt3.setAttribute("value",gGetRes("dataspace_change"))
        // bt3.setAttribute("onclick","dataspace_fun_change("+"'dataspace_name"+number+"'" +",'dataspace_content"+number+"')")
        // td6.append(bt3) 
        // tr.append(td6)

        // var td4 = document.createElement("td")               
        // var bt = document.createElement("input")
        // bt.setAttribute("type","button")
        // bt.setAttribute("class","button1")
        // bt.setAttribute("style","width:100px")
        // bt.setAttribute("value",gGetRes("dataspace_delete"))
        // bt.setAttribute("onclick","dataspace_fun_delete("+"'dataspace_name"+number+"')")
        // td4.append(bt) 
        // tr.append(td4)

        var td6 = document.createElement("td")               
        var bt4 = document.createElement("input")
        bt4.setAttribute("class","button1")
        bt4.setAttribute("id","dataspace_button"+number)
        bt4.setAttribute("style","width:50px")
        bt4.setAttribute("type","button")
        bt4.setAttribute("value","↓")
        bt4.setAttribute("onmouseover","dataspace_func_showfunc("+"'dataspace_button"+number+"')")
        bt4.setAttribute("onmouseout","dataspace_func_mouseout(event)")
        td6.append(bt4) 
        tr.append(td6)
        
        number++
        return tr
    }

    async function initPage(){
        gLoadRes()
        var result = await window.doCheckDataSpace()
        var obj = JSON.parse(result)
        if (obj.Code === 38) {  //have not init
            bt = document.getElementById("dataspace_load")
            bt.innerText = gGetRes("dataspace_init")
            bt.onclick = dataspace_fun_init
        } else if (obj.Code === 37) { //have inited

        } else {
            showResult(result)
        } 
    }


