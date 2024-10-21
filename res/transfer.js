
function ontransfer() {
    var pa = document.getElementById("pass").value
    var add = document.getElementById("add").value
    var toAddress = document.getElementById("toAddress").value
    var amount = document.getElementById("amount").value
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
        if (result.length == 0 ) {
            showMsg("35","")
            return;
        }
        var se = document.getElementById("fee");
        se.value = result;
    })  
}

async function initPage() {    
    await gLoadRes()
    setTimeout(accountChange,100)          
} 
 

