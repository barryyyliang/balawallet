
function addNetworkRow() {
    // 显示模态窗口
    document.getElementById("networkModal").style.display = "block";
    // 清空输入框
    document.getElementById("modal_network_name").value = "";
    document.getElementById("modal_node_address").value = "";
}

async function saveConfig() {
    var obj = {}
    obj["ethAddress"] = document.getElementById("config_input_eth_address").value   
    obj["polAddress"] = document.getElementById("config_input_pol_address").value  
    obj["solAddress"] = document.getElementById("config_input_sol_address").value 
    obj["phrsAddress"] = document.getElementById("config_input_phrs_address").value 
    obj["contractDir"] = document.getElementById("config_input_contractDir").value
    obj["language"] = document.getElementById("config_select_language").value
    
    // 保存新增的网络配置
    var customNetworks = [];
    
    // 获取所有自定义网络行
    var table = document.getElementById("table1");
    for (var i = 0; i < table.rows.length; i++) {
        var row = table.rows[i];
        var networkName = row.getAttribute("data-network-name");
        var nodeAddressInput = row.querySelector(".node-address-input");
        
        if (networkName && nodeAddressInput) {
            var nodeAddr = nodeAddressInput.value;
            if (nodeAddr) {
                // 检查是否已经存在相同名称的网络
                var exists = false;
                for (var j = 0; j < customNetworks.length; j++) {
                    if (customNetworks[j].name === networkName) {
                        exists = true;
                        break;
                    }
                }
                if (!exists) {
                    customNetworks.push({
                        name: networkName,
                        nodeAddress: nodeAddr
                    });
                }
            }
        }
    }
    
    obj["customNetworks"] = customNetworks;
    var json = JSON.stringify(obj)
    window.doWebLog("保存的配置:" + json);
    var result = await window.doConfig(json)
    window.doWebLog("保存结果:" + result);
    showResult(result);
    
    // 尝试刷新main.html中的网络列表
    try {
        if (window.parent && window.parent.refreshMainNet) {
            window.parent.refreshMainNet();
        } else if (window.refreshMainNet) {
            window.refreshMainNet();
        }
    } catch (e) {
        window.doWebLog("刷新网络列表失败: " + e.message);
    }
}

function confirmAddNetwork() {
    // 获取输入值
    var networkName = document.getElementById("modal_network_name").value;
    var nodeAddress = document.getElementById("modal_node_address").value;
    
    // 验证输入
    if (!networkName || !nodeAddress) {
        alert("请填写所有字段");
        return;
    }
    
    // 在表格中添加新行
    var table = document.getElementById("table1");
    var row = table.insertRow(table.rows.length); // 在表格末尾插入
    var cell1 = row.insertCell(0);
    var cell2 = row.insertCell(1);
    var cell3 = row.insertCell(2);
    
    // 设置单元格样式，与内置网络保持一致
    cell1.style.width = "150px";
    
    // 网络名称显示为标签，与内置网络保持一致
    cell1.textContent = networkName;
    
    var nodeAddressInput = document.createElement("input");
    nodeAddressInput.type = "text";
    nodeAddressInput.placeholder = "网络地址";
    nodeAddressInput.className = "node-address-input";
    nodeAddressInput.value = nodeAddress;
    
    var deleteButton = document.createElement("button");
    deleteButton.innerHTML = "&times;";
    deleteButton.className = "button1";
    deleteButton.style.fontSize = "18px";
    deleteButton.style.width = "30px";
    deleteButton.style.height = "30px";
    deleteButton.style.display = "flex";
    deleteButton.style.alignItems = "center";
    deleteButton.style.justifyContent = "center";
    deleteButton.title = "删除";
    deleteButton.onclick = function() {
        row.remove();
        saveConfig(); // 删除网络后保存配置
    };
    
    // 存储网络名称，用于保存配置
    row.setAttribute("data-network-name", networkName);
    
    cell2.appendChild(nodeAddressInput);
    cell3.appendChild(deleteButton);
    
    // 关闭模态窗口
    document.getElementById("networkModal").style.display = "none";
    
    // 新增网络后直接保存配置
    saveConfig();
}

function cancelAddNetwork() {
    // 关闭模态窗口
    document.getElementById("networkModal").style.display = "none";
}

async function onSaveConfig() {
    // 调用保存配置函数
    await saveConfig();
}

function initPage() {
    gLoadRes()
    window.doGetConfig("").then( (result) => {
        window.doWebLog("获取的配置:" + result);
        var obj = {}
        obj["ethAddress"] = "config_input_eth_address"    
        obj["polAddress"] = "config_input_pol_address"  
        obj["solAddress"] = "config_input_sol_address" 
        obj["phrsAddress"] = "config_input_phrs_address" 
        obj["contractDir"] = "config_input_contractDir"
        obj["language"] = "config_select_language"
        var acc = JSON.parse(result)
        window.doWebLog("解析后的配置:" + JSON.stringify(acc));
        window.doWebLog("自定义网络配置:" + JSON.stringify(acc["customNetworks"]));
        for (var val in acc){
            try {
                var ele = document.getElementById(obj[val])
                if (ele != null) {
                    ele.value = acc[val]
                } 
            } catch(err){

            }
        }
        
        // 加载自定义网络配置
        if (acc["customNetworks"] && Array.isArray(acc["customNetworks"])) {
            acc["customNetworks"].forEach(function(network) {
                var table = document.getElementById("table1");
                var row = table.insertRow(table.rows.length); // 在表格末尾插入
                var cell1 = row.insertCell(0);
                var cell2 = row.insertCell(1);
                var cell3 = row.insertCell(2);
                
                // 设置单元格样式，与内置网络保持一致
                cell1.style.width = "150px";
                
                // 网络名称显示为标签，与内置网络保持一致
                cell1.textContent = network.name;
                
                var nodeAddressInput = document.createElement("input");
                nodeAddressInput.type = "text";
                nodeAddressInput.placeholder = "网络地址";
                nodeAddressInput.className = "node-address-input";
                nodeAddressInput.value = network.nodeAddress;
                
                var deleteButton = document.createElement("button");
                deleteButton.innerHTML = "&times;";
                deleteButton.className = "button1";
                deleteButton.style.fontSize = "18px";
                deleteButton.style.width = "30px";
                deleteButton.style.height = "30px";
                deleteButton.style.display = "flex";
                deleteButton.style.alignItems = "center";
                deleteButton.style.justifyContent = "center";
                deleteButton.title = "删除";
                deleteButton.onclick = function() {
                    row.remove();
                    saveConfig(); // 删除网络后保存配置
                };
                
                // 存储网络名称，用于保存配置
                row.setAttribute("data-network-name", network.name);
                
                cell2.appendChild(nodeAddressInput);
                cell3.appendChild(deleteButton);
            });
        }
    }
    )
} 
   