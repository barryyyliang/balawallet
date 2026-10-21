
function ontransfer() {
    var pa = document.getElementById("transfer_pass").value
    var add = document.getElementById("transfer_input_selectedAdd").value
    var toAddress = document.getElementById("transfer_input_toAddress").value
    var amount = document.getElementById("transfer_input_amount").value
    window.doTransfer(add, toAddress,amount,pa).then((result) => {
        var obj = JSON.parse(result)
        MYDIALOG.ShowDiv(obj.Msg);
    })
}

async function accountChange(){
    var address = document.getElementById("add").value 
    document.getElementById("transfer_input_selectedAdd").value =  address         
    ret = await window.doGetAccountName(address)
    document.getElementById("transfer_input_name").value = ret
    window.doGetBalance(address).then ((result) =>{
        var jsonObj = JSON.parse(result)
        if (jsonObj.Code === 0){
            var blObj = JSON.parse(jsonObj.Msg)
            document.getElementById("transfer_input_balance").value = blObj.Balance
            getFee()
        } else {
            MYDIALOG.ShowDiv(jsonObj.Msg);
        }
    })        
}

function getFee(){
    window.doGetFee().then ( (result) => {
        if (result.length == 0 ) {
            showMsg("22","")
            return;
        }
        var se = document.getElementById("transfer_input_fee");
        se.value = result;
    })  
}

async function initPage() {    
    await gLoadRes()
    setTimeout(accountChange,100)          
} 
 

