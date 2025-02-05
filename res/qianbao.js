async function changeName(address,name){
    var addressvalue = document.getElementById(address).value
    var namevalue = document.getElementById(name).value
    var result = await window.doChangeAccountName(addressvalue,namevalue)
    showResult(result)
}

async function changePass(address,name){
    window.doWebLog("enter changePass")
    gParam = {}
    gParam.type="eth"
    gParam.address=document.getElementById(address).value
    gParam.name=document.getElementById(name).value
    changePage("qianbao","changepass")
}

function management(address,name){
    window.doWebLog("enter management")
}

function qianbao_func_showfunc(id){
    var p1 = document.getElementById(id).parentNode.parentNode
    var cells = p1.getElementsByTagName("td")
    var address = cells[1].childNodes[0].id
    var name = cells[2].childNodes[0].id

    cells[1].childNodes[0].classList.add("highlight")
    cells[2].childNodes[0].classList.add("highlight")
    cells[3].childNodes[0].classList.add("highlight")

    clearFunc()
    addFuncToMenu("changeName('"+address+"','"+name+"')", gGetRes("qianbao_changebutton"))
    addFuncToMenu("changePass('"+address+"','"+name+"')", gGetRes("qianbao_changepass"))
    onFuncMenu(id)
}

function qianbao_func_mouseout(id){
    var p1 = document.getElementById(id).parentNode.parentNode
    var cells = p1.getElementsByTagName("td")

    cells[1].childNodes[0].classList.remove("highlight")
    cells[2].childNodes[0].classList.remove("highlight")
    cells[3].childNodes[0].classList.remove("highlight")
}


function addAccountElement(acc,table){
    for (i=0; i<acc.length; i++){
        var par = acc[i]

        var tr = document.createElement("tr")
        tr.setAttribute("style","height:50px;")
        var td1 = document.createElement("td")
        td1.append("Eth: ")
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("id","address"+i)
        ipt.setAttribute("type","text")
        ipt.setAttribute("value",par.address)
        ipt.setAttribute("readonly","readonly")       //此处不用disable属性，disable属性会导致自定义copy菜单无法拷贝内容  
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)

        var td3 = document.createElement("td")
        var ipt2 = document.createElement("input")
        ipt2.setAttribute("id","name"+i)
        ipt2.setAttribute("style","width:200px;")
        ipt2.setAttribute("type","text")
        ipt2.setAttribute("value",par.name)
        td3.append(ipt2)
        tr.append(td3)

        var td4 = document.createElement("td")
        var ipt3 = document.createElement("input")
        ipt3.setAttribute("id","balance"+i)
        ipt3.setAttribute("style","width:200px;")
        ipt3.setAttribute("type","text")
        ipt3.setAttribute("readonly","readonly") 
        td4.append(ipt3)
        tr.append(td4)

        var td6 = document.createElement("td")               
        var bt3 = document.createElement("input")
        bt3.setAttribute("class","button1")
        bt3.setAttribute("id","qianbao_button"+i)
        bt3.setAttribute("style","width:50px")
        bt3.setAttribute("type","button")
        bt3.setAttribute("value","↓")
        bt3.setAttribute("onmouseover","qianbao_func_showfunc("+"'qianbao_button"+i+"')")
        bt3.setAttribute("onmouseout","qianbao_func_mouseout("+"'qianbao_button"+i+"')")

        
        td6.append(bt3) 
        tr.append(td6)

        // var td4 = document.createElement("td")               
        // var bt = document.createElement("input")
        // bt.setAttribute("class","button1")
        // //bt.setAttribute("style","width:100px")
        // bt.setAttribute("type","button")
        // bt.setAttribute("value",gGetRes("qianbao_changebutton"))
        // bt.setAttribute("onclick","changeName("+"'address"+i +"'" +",'name"+i+"')")
        // td4.append(bt) 
        // tr.append(td4)

        // var td5 = document.createElement("td")               
        // var bt2 = document.createElement("input")
        // bt2.setAttribute("class","button1")
        // //bt.setAttribute("style","width:100px")
        // bt2.setAttribute("type","button")
        // bt2.setAttribute("value",gGetRes("qianbao_changepass"))
        // bt2.setAttribute("onclick","changePass("+"'address"+i +"'" +",'name"+i+"')")
        // td5.append(bt2) 
        // tr.append(td5)




        table.append(tr)
    }
}


async function initBalance(){
    // 获取每个地址的balance
    if ( getCurrentPage() === "qianbao"){
        var table = document.getElementById("tableWallet")
        for(var i =1; i< table.rows.length;i++){
            result = await window.doGetBalance(table.rows[i].cells[1].childNodes[0].value)
            var jsonObj = JSON.parse(result)
            if (jsonObj.Code === 0){
                var blObj = JSON.parse(jsonObj.Msg)
                table.rows[i].cells[3].childNodes[0].value = blObj.Balance
            } else {
                MYDIALOG.ShowDiv(jsonObj.Msg);
                return
            }
        }   
    } 
}

async function initPage(){
    gLoadRes()    
    var se = document.getElementById("add")
    se.options.length = 0

    window.doGetAccounts().then( (result) => {
        if (result != ""){
            var acc = JSON.parse(result)
            addAccountElement(acc.accounts,document.getElementById("tableWallet"))
            
            if (acc != ""){    
                var se = document.getElementById("add")            
                for (i=0; i<acc.accounts.length; i++){
                    var option = document.createElement("option")
                    option.text = acc.accounts[i].address
                    option.value = acc.accounts[i].address
                    se.add(option)
                }
            }
        }
    })
    var re = await window.doCheckNode()
    var obj = JSON.parse(re)
    if(obj.Code === 0){
        setTimeout(initBalance,500)
    }
    // var result = await window.doGetAccounts()
    // var arr = JSON.parse(result)
    // addAccountElement(arr.accounts,document.getElementById("table1"))
    // initBalance()    
}
