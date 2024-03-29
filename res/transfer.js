
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
            MYDIALOG.ShowDiv("MyDiv","fade", "转账成功");
        } else {
            MYDIALOG.ShowDiv("MyDiv","fade", obj.Err);
        }
    })
}

function change(){
    window.doWebLog("value changed")
    var address = document.getElementById("add").value 
    document.getElementById("selectedAdd").value =  address         
    window.doGetAccountName(address).then ( (result) => {
        window.doWebLog(result)
        document.getElementById("name").value = result
    }
    )
    window.doGetBalance(address).then ( (result) => {
        document.getElementById("balance").value = result
    }
    )
}

function initPage() {
    change()
    window.doGetFee().then ( (result) => {
        window.doWebLog("get fee:" + result)
        if (result.length == 0 ) {
            window.doWebLog("cannot get fee")
            MYDIALOG.ShowDiv("MyDiv","fade", "无法获得预估转账费用");
            return;
        }
        var se = document.getElementById("fee");
        se.value = result;
    }
    )            
} 
 

