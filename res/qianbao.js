function changeName(address,name){
    window.doWebLog("enter changeName")
    window.doWebLog(address)
    var addressvalue = document.getElementById(address).value
    var namevalue = document.getElementById(name).value
    window.doChangeAccountName(addressvalue,namevalue)
}
function management(address,name){
    window.doWebLog("enter management")

}
function addAccountElement(acc,table){
    window.doWebLog("enter addelement")
    for (i=0; i<acc.length; i++){
        var par = acc[i].split('-')

        var tr = document.createElement("tr")
        tr.setAttribute("style","height:50px;")
        var td1 = document.createElement("td")
        td1.append("Eth: ")
        tr.append(td1)

        var ipt = document.createElement("input")
        ipt.setAttribute("id","address"+i)
        ipt.setAttribute("style","font-size: 14px;height:25px;")
        ipt.setAttribute("type","text")
        ipt.setAttribute("value",par[0])
        ipt.setAttribute("size","45")   
        ipt.setAttribute("disabled","")         
        var td2 = document.createElement("td")
        td2.append(ipt)
        tr.append(td2)


        var td3 = document.createElement("td")
        var ipt2 = document.createElement("input")
        ipt2.setAttribute("id","name"+i)
        ipt2.setAttribute("style","font-size: 14px;height:25px;")
        ipt2.setAttribute("type","text")
        ipt2.setAttribute("value",par[1])
        ipt2.setAttribute("size","20")  
        td3.append(ipt2)
        tr.append(td3)

        var td4 = document.createElement("td")
        var ipt3 = document.createElement("input")
        ipt3.setAttribute("id","balance"+i)
        ipt3.setAttribute("style","font-size: 14px;height:25px;")
        ipt3.setAttribute("type","text")
        ipt3.setAttribute("size","20")  
        td4.append(ipt3)
        tr.append(td4)

        var td4 = document.createElement("td")               
        var bt = document.createElement("input")
        bt.setAttribute("class","button2")
        bt.setAttribute("type","button")
        bt.setAttribute("value","修改别名")
        bt.setAttribute("onclick","changeName("+"'address"+i +"'" +",'name"+i+"')")
        td4.append(bt) 
        tr.append(td4)

        table.append(tr)
    }
    //ShowDiv("MyDiv","fade")

}

async function initBalance(){
    // 获取每个地址的balance
    window.doWebLog("initBalance")
    var table = document.getElementById("table1")
    for(var i=1; i< table.rows.length;i++){
        //window.doGetBalance(table.rows[i].cells[1].childNodes[0].value).then ( (result) => {
            //table.rows[i].cells[3].childNodes[0].value = result
        //})    
        table.rows[i].cells[3].childNodes[0].value = await window.doGetBalance(table.rows[i].cells[1].childNodes[0].value)
    }
}

function initPage(){
    window.doWebLog("qianbao init page")
    window.doGetAccounts().then( (result) => {
        window.doWebLog("get result" + result)
        var arr = result.split(',')
        window.doWebLog("go to addelement")
        addAccountElement(arr,document.getElementById("table1"))
    })

    setTimeout(initBalance,1000)
}
