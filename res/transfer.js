
function ontransfer() {
    window.doWebLog("enter transfer")
    var pa = document.getElementById("pass").value
    window.doWebLog("pa value:"+pa)
    var add = document.getElementById("add").value
    window.doWebLog("add value:" + add) 
    var toAddress = document.getElementById("toAddress").value
    window.doWebLog("toAddress value:" + toAddress)  
    var amount = document.getElementById("amount").value
    window.doWebLog("amount value:" + amount) 
    window.doTransfer(add, toAddress,amount,pa).then((result) => {
        var obj = JSON.parse(result)
        if (obj.Code === 0) {
            MYDIALOG.ShowDiv("转账成功");
        } else {
            MYDIALOG.ShowDiv(obj.Msg);
        }
    })
}

async function accountChange(){
    var address = document.getElementById("add").value 
    document.getElementById("selectedAdd").value =  address         
    ret = await window.doGetAccountName(address)
    document.getElementById("name").value = ret
    window.doGetBalance(address).then ((result) =>{
        var jsonObj = JSON.parse(result)
        if (jsonObj.Code === 0){
            var blObj = JSON.parse(jsonObj.Msg)
            document.getElementById("balance").value = blObj.Balance
            getFee()
        } else {
            MYDIALOG.ShowDiv(jsonObj.Msg);
        }
    })        
}

function getFee(){
    window.doGetFee().then ( (result) => {
        window.doWebLog("get fee:" + result)
        if (result.length == 0 ) {
            window.doWebLog("cannot get fee")
            MYDIALOG.ShowDiv("无法获得预估转账费用");
            return;
        }
        var se = document.getElementById("fee");
        se.value = result;
    })  
}

function initPage() {    
    setTimeout(accountChange,100)          
} 
 

