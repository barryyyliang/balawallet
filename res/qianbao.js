async function changeName(address,name){
    window.doWebLog("enter changeName")
    window.doWebLog(address)
    var addressvalue = document.getElementById(address).value
    var namevalue = document.getElementById(name).value
    var result = await window.doChangeAccountName(addressvalue,namevalue)
    showResult(result)
}
function management(address,name){
    window.doWebLog("enter management")
}
function addAccountElement(acc,table){
    window.doWebLog("enter addelement")
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

        var td4 = document.createElement("td")               
        var bt = document.createElement("input")
        bt.setAttribute("class","button1")
        bt.setAttribute("style","width:100px")
        bt.setAttribute("type","button")
        bt.setAttribute("value","修改别名")
        bt.setAttribute("onclick","changeName("+"'address"+i +"'" +",'name"+i+"')")
        td4.append(bt) 
        tr.append(td4)

        table.append(tr)
    }
}


async function initBalance(){
    // 获取每个地址的balance
    window.doWebLog("initBalance")
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
    } else {
        window.doWebLog("get balance conceled, page changed")
   }
}

async function initPage(){
    window.doWebLog("qianbao init page")
    window.doGetAccounts().then( (result) => {
        window.doWebLog("get result" + result)
        var arr = JSON.parse(result)
        window.doWebLog("go to addelement")
        addAccountElement(arr.accounts,document.getElementById("tableWallet"))
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
