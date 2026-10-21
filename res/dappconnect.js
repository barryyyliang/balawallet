
var dappConnected = false
var dappSessionId = ""
var dappPollTimer = null
var dappAuthorized = false

function initPage() {
    gLoadRes()
    // 页面加载时检查当前连接状态
    checkDappStatus()
}

function reInit() {
    checkDappStatus()
}

function appendLog(msg) {
    var log = document.getElementById("dapp_log")
    var time = new Date().toLocaleTimeString()
    log.value = "[" + time + "] " + msg + "\n" + log.value
}

function parseWcUri(uri) {
    // 解析 WalletConnect v2 URI 格式:
    // wc:<topic>@<version>?<params>
    var result = {}
    try {
        if (!uri.startsWith("wc:")) {
            return null
        }
        var withoutPrefix = uri.substring(3)
        var atIndex = withoutPrefix.indexOf("@")
        if (atIndex === -1) return null

        result.topic = withoutPrefix.substring(0, atIndex)

        var rest = withoutPrefix.substring(atIndex + 1)
        var qIndex = rest.indexOf("?")
        if (qIndex === -1) return null

        result.version = rest.substring(0, qIndex)

        var paramStr = rest.substring(qIndex + 1)
        var params = new URLSearchParams(paramStr)
        result.relay = params.get("relay-protocol") || "irn"
        result.symKey = params.get("symKey") || ""
        result.expiryTimestamp = params.get("expiryTimestamp") || ""

        return result
    } catch (e) {
        return null
    }
}

async function dappConnect() {
    var uri = document.getElementById("dapp_uri_input").value.trim()
    if (uri === "") {
        showMsg("29", "")
        return
    }

    // 解析 URI
    var parsed = parseWcUri(uri)
    if (parsed == null) {
        appendLog("URI 格式错误，请输入正确的 WalletConnect URI")
        MYDIALOG.ShowDiv("URI 格式错误，请输入正确的 WalletConnect URI")
        return
    }

    // 检查过期时间
    if (parsed.expiryTimestamp) {
        var now = Math.floor(Date.now() / 1000)
        if (parseInt(parsed.expiryTimestamp) < now) {
            appendLog("连接 URI 已过期，请获取新的连接 URI")
            MYDIALOG.ShowDiv("连接 URI 已过期，请获取新的连接 URI")
            return
        }
    }

    appendLog("正在连接 WalletConnect Relay Server...")
    appendLog("Topic: " + parsed.topic.substring(0, 16) + "...")
    appendLog("Version: " + parsed.version)
    appendLog("Relay: " + parsed.relay)

    var account = getCurrentAccount()
    if (!account || account === "") {
        appendLog("请先选择一个账户")
        MYDIALOG.ShowDiv("请先选择一个账户")
        return
    }

    // 禁用连接按钮，显示连接中
    document.getElementById("dapp_connect_btn").disabled = true
    document.getElementById("dapp_connect_btn").innerText = "连接中..."

    // 调用后端建立 WebSocket 连接到 Relay
    try {
        var result = await window.doDappConnect(uri, account)
        var obj = JSON.parse(result)
        if (obj.Code === 0) {
            var data = JSON.parse(obj.Msg)
            dappConnected = true
            dappSessionId = data.sessionId || ""
            dappAuthorized = data.authorized || false

            appendLog("已连接到 Relay Server ✓")
            appendLog("等待 Dapp 发送 session proposal...")

            // 更新 UI 为等待状态
            showWaitingStatus(data)
            // 开始轮询请求 (包括 session_propose 的结果)
            startPolling()
        } else {
            appendLog("连接失败: " + obj.Msg)
            MYDIALOG.ShowDiv(obj.Msg)
            resetConnectButton()
        }
    } catch (e) {
        appendLog("连接异常: " + e.message)
        MYDIALOG.ShowDiv("连接失败: " + e.message)
        resetConnectButton()
    }
}

function resetConnectButton() {
    document.getElementById("dapp_connect_btn").disabled = false
    document.getElementById("dapp_connect_btn").innerText = "连接"
}

async function dappDisconnect() {
    if (!dappConnected) return

    try {
        var result = await window.doDappDisconnect(dappSessionId)
        var obj = JSON.parse(result)
        appendLog("已断开连接 ✓")
    } catch (e) {
        appendLog("断开连接异常: " + e.message)
    }

    dappConnected = false
    dappSessionId = ""
    dappAuthorized = false
    stopPolling()
    showDisconnectedStatus()
}

function showWaitingStatus(data) {
    document.getElementById("dapp_status_area").style.display = "block"
    document.getElementById("dapp_status_value").innerText = "已连接 Relay (等待 Dapp 授权...)"
    document.getElementById("dapp_status_value").style.color = "#e8a838"
    document.getElementById("dapp_name_value").innerText = "等待中..."
    document.getElementById("dapp_url_value").innerText = "--"
    document.getElementById("dapp_chain_value").innerText = data.chainId || "--"
    document.getElementById("dapp_account_value").innerText = data.account || getCurrentAccount()

    // 隐藏 icon
    document.getElementById("dapp_icon_area").style.display = "none"

    document.getElementById("dapp_connect_btn").style.display = "none"
    document.getElementById("dapp_disconnect_btn").style.display = "inline-block"
    document.getElementById("dapp_uri_input").disabled = true
}

function showConnectedStatus(data) {
    document.getElementById("dapp_status_area").style.display = "block"
    document.getElementById("dapp_status_value").innerText = "已连接 ✓ (已授权)"
    document.getElementById("dapp_status_value").style.color = "green"
    document.getElementById("dapp_name_value").innerText = data.name || "--"
    document.getElementById("dapp_url_value").innerText = data.url || "--"

    // 显示 Dapp 描述
    if (data.desc) {
        document.getElementById("dapp_desc_value").innerText = data.desc
        document.getElementById("dapp_desc_row").style.display = ""
    }

    // 显示 Dapp 图标
    if (data.icon) {
        document.getElementById("dapp_icon_img").src = data.icon
        document.getElementById("dapp_icon_area").style.display = "block"
    }

    document.getElementById("dapp_chain_value").innerText = data.chainId || "--"
    document.getElementById("dapp_account_value").innerText = data.account || getCurrentAccount()

    document.getElementById("dapp_connect_btn").style.display = "none"
    document.getElementById("dapp_disconnect_btn").style.display = "inline-block"
    document.getElementById("dapp_uri_input").disabled = true
}

function showDisconnectedStatus() {
    document.getElementById("dapp_status_area").style.display = "none"
    document.getElementById("dapp_request_area").style.display = "none"
    document.getElementById("dapp_icon_area").style.display = "none"
    document.getElementById("dapp_desc_row").style.display = "none"

    resetConnectButton()
    document.getElementById("dapp_connect_btn").style.display = "inline-block"
    document.getElementById("dapp_disconnect_btn").style.display = "none"
    document.getElementById("dapp_uri_input").disabled = false
}

async function checkDappStatus() {
    try {
        var result = await window.doDappStatus()
        var obj = JSON.parse(result)
        if (obj.Code === 0 && obj.Msg !== "") {
            var data = JSON.parse(obj.Msg)
            if (data.connected) {
                dappConnected = true
                dappSessionId = data.sessionId || ""
                dappAuthorized = data.authorized || false

                if (dappAuthorized && data.name) {
                    showConnectedStatus(data)
                } else {
                    showWaitingStatus(data)
                }
                startPolling()
            }
        }
    } catch (e) {
        // 忽略错误
    }
}

function startPolling() {
    if (dappPollTimer != null) return
    dappPollTimer = setInterval(async function () {
        if (!dappConnected) {
            stopPolling()
            return
        }
        try {
            var result = await window.doDappPollRequest(dappSessionId)
            var obj = JSON.parse(result)
            if (obj.Code === 0 && obj.Msg !== "" && obj.Msg !== "null") {
                var req = JSON.parse(obj.Msg)
                if (req.method && req.method !== "") {
                    handleDappEvent(req)
                }
            }
        } catch (e) {
            // 忽略轮询错误
        }
    }, 2000)
}

function stopPolling() {
    if (dappPollTimer != null) {
        clearInterval(dappPollTimer)
        dappPollTimer = null
    }
}

// 处理来自后端的各种事件
function handleDappEvent(req) {
    switch (req.method) {
        case "_session_authorized":
            // Dapp 已授权 session，更新 UI
            dappAuthorized = true
            appendLog("✅ Dapp 已授权 session!")
            if (req.params) {
                var info = req.params
                if (typeof info === "string") {
                    info = JSON.parse(info)
                }
                appendLog("Dapp: " + (info.name || "Unknown"))
                appendLog("URL: " + (info.url || ""))
                showConnectedStatus({
                    name: info.name,
                    url: info.url,
                    icon: info.icon,
                    desc: info.desc,
                    chainId: info.chainId,
                    account: info.account
                })
            }
            break

        case "_session_disconnected":
            // Dapp 断开了 session
            appendLog("⚠️ Dapp 已断开连接")
            dappConnected = false
            dappAuthorized = false
            stopPolling()
            showDisconnectedStatus()
            break

        case "personal_sign":
        case "eth_sign":
        case "eth_signTypedData":
        case "eth_signTypedData_v4":
        case "eth_sendTransaction":
            // 签名/交易请求
            showSignRequest(req)
            break

        default:
            appendLog("收到请求: " + req.method)
            showSignRequest(req)
            break
    }
}

function showSignRequest(req) {
    document.getElementById("dapp_request_area").style.display = "block"
    document.getElementById("dapp_request_method").innerText = req.method || "--"
    document.getElementById("dapp_request_id").value = req.id || "0"

    var dataStr = ""
    if (req.params) {
        if (typeof req.params === "string") {
            // 如果是 hex 格式的 personal_sign 数据，尝试解码显示
            if (req.method === "personal_sign" && req.params.startsWith("0x")) {
                try {
                    var hexStr = req.params.substring(2)
                    var decoded = ""
                    for (var i = 0; i < hexStr.length; i += 2) {
                        decoded += String.fromCharCode(parseInt(hexStr.substr(i, 2), 16))
                    }
                    dataStr = "原始 (hex): " + req.params + "\n\n解码: " + decoded
                } catch (e) {
                    dataStr = req.params
                }
            } else {
                dataStr = req.params
            }
        } else {
            dataStr = JSON.stringify(req.params, null, 2)
        }
    }
    document.getElementById("dapp_request_data").value = dataStr
    document.getElementById("dapp_request_pass").value = ""

    appendLog("📝 收到签名请求: " + req.method + " (ID: " + req.id + ")")
}

async function dappApprove() {
    var pass = document.getElementById("dapp_request_pass").value
    if (pass === "") {
        showMsg("33", "")
        return
    }

    var method = document.getElementById("dapp_request_method").innerText
    var requestId = parseInt(document.getElementById("dapp_request_id").value) || 0

    // 从显示的数据中提取原始数据用于签名
    var rawData = document.getElementById("dapp_request_data").value
    // 如果数据包含"原始 (hex)"格式，提取 hex 部分
    if (rawData.indexOf("原始 (hex): ") === 0) {
        var lineEnd = rawData.indexOf("\n")
        if (lineEnd > 0) {
            rawData = rawData.substring("原始 (hex): ".length, lineEnd)
        }
    }

    var account = getCurrentAccount()

    appendLog("正在处理签名请求...")

    try {
        var param = JSON.stringify({
            sessionId: dappSessionId,
            requestId: requestId,
            method: method,
            data: rawData,
            account: account,
            password: pass
        })
        var result = await window.doDappApprove(param)
        var obj = JSON.parse(result)
        if (obj.Code === 0) {
            appendLog("✅ 签名成功!")
            appendLog("签名: " + obj.Msg.substring(0, 20) + "...")
            appendLog("结果已通过 Relay 发送给 Dapp")
            MYDIALOG.ShowDiv("签名成功")
        } else {
            appendLog("❌ 签名失败: " + obj.Msg)
            MYDIALOG.ShowDiv(obj.Msg)
        }
    } catch (e) {
        appendLog("签名异常: " + e.message)
        MYDIALOG.ShowDiv("签名失败: " + e.message)
    }

    document.getElementById("dapp_request_area").style.display = "none"
    document.getElementById("dapp_request_pass").value = ""
}

async function dappReject() {
    var requestId = parseInt(document.getElementById("dapp_request_id").value) || 0
    appendLog("已拒绝签名请求")

    try {
        var param = JSON.stringify({
            sessionId: dappSessionId,
            requestId: requestId,
            reject: true
        })
        await window.doDappReject(param)
        appendLog("拒绝信息已通过 Relay 发送给 Dapp")
    } catch (e) {
        // 忽略错误
    }

    document.getElementById("dapp_request_area").style.display = "none"
    document.getElementById("dapp_request_pass").value = ""
}
