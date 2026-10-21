package pagefunc

import (
	"balawallet/common"
	"balawallet/serverapi"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

// ================== 数据结构 ==================

// WCSession 管理 WalletConnect v2 会话
type WCSession struct {
	PairingTopic string `json:"pairingTopic"` // pairing topic (from URI)
	SessionTopic string `json:"sessionTopic"` // session topic (negotiated)
	Version      string `json:"version"`
	Relay        string `json:"relay"`
	SymKey       string `json:"symKey"`     // pairing 对称密钥 (hex)
	SessionKey   []byte `json:"-"`          // session 对称密钥 (bytes)
	SelfPrivKey  []byte `json:"-"`          // X25519 私钥 (32 bytes)
	SelfPubKey   []byte `json:"-"`          // X25519 公钥 (32 bytes)
	Account      string `json:"account"`    // 授权的钱包地址
	ChainId      string `json:"chainId"`    // chain ID
	Connected    bool   `json:"connected"`  // 是否已连接
	Authorized   bool   `json:"authorized"` // 是否已授权 session
	DappName     string `json:"dappName"`   // Dapp 名称
	DappUrl      string `json:"dappUrl"`    // Dapp URL
	DappIcon     string `json:"dappIcon"`   // Dapp 图标
	DappDesc     string `json:"dappDesc"`   // Dapp 描述
	ExpireAt     int64  `json:"expireAt"`   // 过期时间
	ProposalId   int64  `json:"proposalId"` // session proposal ID
}

// WCRequest Dapp 发来的签名请求
type WCRequest struct {
	Id     int64       `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
	Topic  string      `json:"topic"`
}

// DappConnectResult 返回给前端的连接结果
type DappConnectResult struct {
	SessionId  string `json:"sessionId"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	Icon       string `json:"icon"`
	Desc       string `json:"desc"`
	ChainId    string `json:"chainId"`
	Account    string `json:"account"`
	Connected  bool   `json:"connected"`
	Authorized bool   `json:"authorized"`
}

// DappApproveParam 签名批准参数
type DappApproveParam struct {
	SessionId string `json:"sessionId"`
	RequestId int64  `json:"requestId"`
	Method    string `json:"method"`
	Data      string `json:"data"`
	Account   string `json:"account"`
	Password  string `json:"password"`
}

// JSON-RPC 通用消息
type jsonRPCRequest struct {
	Id      int64       `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRPCResponse struct {
	Id      int64       `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Relay 发布/订阅消息
type relayPublishParams struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Tag     int    `json:"tag"`
	Prompt  bool   `json:"prompt"`
}

type relaySubscribeParams struct {
	Topic string `json:"topic"`
}

// Relay 订阅推送消息
type relaySubscriptionParams struct {
	Id   string                `json:"id"`
	Data relaySubscriptionData `json:"data"`
}

type relaySubscriptionData struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
	Tag     int    `json:"tag"`
}

// WalletConnect v2 envelope (type1 = AES-256-GCM)
type wcEnvelope struct {
	Type    int    `json:"type"`
	IV      string `json:"iv"`
	Data    string `json:"data"`
	Sealbox string `json:"sealbox"` // optional
}

// session_propose params 结构
type sessionProposeParams struct {
	Proposer       proposerInfo      `json:"proposer"`
	RequiredNS     map[string]nsInfo `json:"requiredNamespaces"`
	OptionalNS     map[string]nsInfo `json:"optionalNamespaces"`
	Relays         []relayInfo       `json:"relays"`
	ProposerPubKey string            `json:"proposerPublicKey"`
}

type proposerInfo struct {
	PublicKey string       `json:"publicKey"`
	Metadata  dappMetadata `json:"metadata"`
}

type dappMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Icons       []string `json:"icons"`
}

type nsInfo struct {
	Chains  []string `json:"chains"`
	Methods []string `json:"methods"`
	Events  []string `json:"events"`
}

type relayInfo struct {
	Protocol string `json:"protocol"`
}

// session_request params 结构
type sessionRequestParams struct {
	Request sessionRequestDetail `json:"request"`
	ChainId string               `json:"chainId"`
}

type sessionRequestDetail struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

// ================== 全局状态 ==================

var (
	wcSession     *WCSession
	wcSessionLock sync.RWMutex
	wcRequests    []WCRequest
	wcRequestLock sync.RWMutex
	wcConn        *websocket.Conn
	wcConnLock    sync.Mutex
	wcDone        chan struct{}
	wcMsgIdSeq    int64 // 已弃用，使用 generateRequestId()
	wcSubId       string // subscription id

	// Ed25519 密钥对 (用于 Relay 认证 JWT)
	relayAuthKey     ed25519.PrivateKey
	relayAuthPubKey  ed25519.PublicKey
	relayAuthKeyOnce sync.Once
)

const (
	relayURL              = "wss://relay.reown.com"            // Reown 是 WalletConnect 新品牌，此域名未被封锁
	relayAudience         = "wss://relay.reown.com"            // JWT aud 字段必须用官方标准值
	defaultRelayProjectId = "71a57a73ef1ce118b70fa5d4ba061c91" // 默认 Project ID (需到 cloud.reown.com 注册)
)

// ================== 公开的前端绑定函数 ==================

// DoDappConnect 连接到 Dapp (通过 WalletConnect v2 Relay)
func DoDappConnect(uri string, account string) string {
	// 解析 WalletConnect URI
	parsed, err := parseWCUri(uri)
	if err != nil {
		re := common.GetResult(common.Msg_param_error, err.Error())
		return string(re)
	}

	// 检查过期时间
	if parsed.ExpireAt > 0 {
		now := time.Now().Unix()
		if parsed.ExpireAt < now {
			re := common.GetResult(common.Msg_fail, "URI has expired")
			return string(re)
		}
	}

	// 断开已有连接
	disconnectInternal()

	chainId := serverapi.GetCurrentServer().DoGetChainId()

	// 生成 X25519 密钥对 (用于后续 Session 密钥交换)
	selfPrivKey := make([]byte, 32)
	rand.Read(selfPrivKey)
	selfPubKey, err := curve25519.X25519(selfPrivKey, curve25519.Basepoint)
	if err != nil {
		re := common.GetResult(common.Msg_fail, "Generate X25519 key: "+err.Error())
		return string(re)
	}

	// 创建会话
	wcSessionLock.Lock()
	wcSession = &WCSession{
		PairingTopic: parsed.Topic,
		Version:      parsed.Version,
		Relay:        parsed.Relay,
		SymKey:       parsed.SymKey,
		SelfPrivKey:  selfPrivKey,
		SelfPubKey:   selfPubKey,
		Account:      account,
		ChainId:      chainId,
		Connected:    false,
		Authorized:   false,
		ExpireAt:     parsed.ExpireAt,
	}
	wcSessionLock.Unlock()

	// 清空请求队列
	wcRequestLock.Lock()
	wcRequests = make([]WCRequest, 0)
	wcRequestLock.Unlock()

	// 连接 Relay Server
	err = connectRelay(parsed.Topic, parsed.SymKey)
	if err != nil {
		wcSessionLock.Lock()
		wcSession = nil
		wcSessionLock.Unlock()
		re := common.GetResult(common.Msg_fail, "Connect relay failed: "+err.Error())
		return string(re)
	}

	wcSessionLock.Lock()
	wcSession.Connected = true
	wcSessionLock.Unlock()

	fmt.Println("WalletConnect: connected to relay, subscribed to topic:", parsed.Topic[:16]+"...")

	// 返回初始连接结果（等待 session_propose）
	result := DappConnectResult{
		SessionId:  parsed.Topic[:16],
		Name:       "Waiting for Dapp...",
		Url:        "",
		ChainId:    chainId,
		Account:    account,
		Connected:  true,
		Authorized: false,
	}

	resultJson, _ := json.Marshal(result)
	re, _ := json.Marshal(common.ErrorResult{common.Msg_success, string(resultJson)})
	return string(re)
}

// DoDappDisconnect 断开 Dapp 连接
func DoDappDisconnect(sessionId string) string {
	disconnectInternal()
	re := common.GetResult(common.Msg_success, "")
	return string(re)
}

// DoDappStatus 获取 Dapp 连接状态
func DoDappStatus() string {
	wcSessionLock.RLock()
	defer wcSessionLock.RUnlock()

	if wcSession != nil && wcSession.Connected {
		result := DappConnectResult{
			SessionId:  wcSession.PairingTopic[:16],
			Name:       wcSession.DappName,
			Url:        wcSession.DappUrl,
			Icon:       wcSession.DappIcon,
			Desc:       wcSession.DappDesc,
			ChainId:    wcSession.ChainId,
			Account:    wcSession.Account,
			Connected:  wcSession.Connected,
			Authorized: wcSession.Authorized,
		}
		resultJson, _ := json.Marshal(result)
		re, _ := json.Marshal(common.ErrorResult{common.Msg_success, string(resultJson)})
		return string(re)
	}

	re, _ := json.Marshal(common.ErrorResult{common.Msg_success, ""})
	return string(re)
}

// DoDappPollRequest 轮询来自 Dapp 的签名请求
func DoDappPollRequest(sessionId string) string {
	wcRequestLock.Lock()
	defer wcRequestLock.Unlock()

	if len(wcRequests) > 0 {
		req := wcRequests[0]
		wcRequests = wcRequests[1:]
		reqJson, _ := json.Marshal(req)
		re, _ := json.Marshal(common.ErrorResult{common.Msg_success, string(reqJson)})
		return string(re)
	}

	re, _ := json.Marshal(common.ErrorResult{common.Msg_success, ""})
	return string(re)
}

// DoDappApprove 批准签名请求
func DoDappApprove(param string) string {
	var approveParam DappApproveParam
	err := json.Unmarshal([]byte(param), &approveParam)
	if err != nil {
		re := common.GetResult(common.Msg_param_error, err.Error())
		return string(re)
	}

	if approveParam.Password == "" {
		re := common.GetResult(common.Msg_password_null, "")
		return string(re)
	}

	method := approveParam.Method
	data := approveParam.Data
	account := approveParam.Account

	fmt.Println("WalletConnect: processing request, method=" + method)

	var signResult string
	sv := serverapi.GetCurrentServer()
	switch method {
	case "personal_sign", "eth_sign":
		signResult = sv.DoPersonalSign(account, approveParam.Password, data)
	case "eth_signTypedData", "eth_signTypedData_v4":
		signResult = sv.DoSignTypedData(account, approveParam.Password, data)
	case "eth_sendTransaction":
		signResult = sv.DoSignTransaction(account, approveParam.Password, data)
	default:
		signResult = sv.DoPersonalSign(account, approveParam.Password, data)
	}

	// 解析签名结果
	var signObj common.ErrorResult
	json.Unmarshal([]byte(signResult), &signObj)

	if signObj.Code == common.Msg_success {
		// 签名成功，通过 Relay 发送结果给 Dapp
		sendSignResult(approveParam.RequestId, signObj.Msg, nil)
	} else {
		// 签名失败，发送错误给 Dapp
		sendSignResult(approveParam.RequestId, "", fmt.Errorf("%s", signObj.Msg))
	}

	return signResult
}

// DoDappReject 拒绝签名请求
func DoDappReject(param string) string {
	var rejectParam struct {
		RequestId int64 `json:"requestId"`
	}
	json.Unmarshal([]byte(param), &rejectParam)

	// 通过 Relay 发送拒绝消息给 Dapp
	if rejectParam.RequestId > 0 {
		sendSignResult(rejectParam.RequestId, "", fmt.Errorf("User rejected the request"))
	}

	fmt.Println("WalletConnect: request rejected")
	re := common.GetResult(common.Msg_success, "")
	return string(re)
}

// DoDappSendRequest 手动发送签名请求 (用于本地测试)
func DoDappSendRequest(reqJson string) string {
	var req WCRequest
	err := json.Unmarshal([]byte(reqJson), &req)
	if err != nil {
		re := common.GetResult(common.Msg_param_error, err.Error())
		return string(re)
	}

	wcRequestLock.Lock()
	wcRequests = append(wcRequests, req)
	wcRequestLock.Unlock()

	re := common.GetResult(common.Msg_success, "")
	return string(re)
}

// ================== Relay WebSocket 通信 ==================

// resolveViaDoH 通过 Cloudflare DNS-over-HTTPS 解析域名的真实 IP，绕过本地 DNS 污染
func resolveViaDoH(host string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	reqURL := fmt.Sprintf("https://1.1.1.1/dns-query?name=%s&type=A", url.QueryEscape(host))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/dns-json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("WalletConnect: DoH query failed:", err)
		return ""
	}
	defer resp.Body.Close()

	var result struct {
		Answer []struct {
			Type int    `json:"type"`
			Data string `json:"data"`
		} `json:"Answer"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	for _, ans := range result.Answer {
		if ans.Type == 1 { // A record
			return ans.Data
		}
	}
	return ""
}

// connectRelay 连接到 WalletConnect Relay 并订阅 pairing topic
func connectRelay(topic string, symKeyHex string) error {
	// 生成 JWT auth token
	authToken, err := generateRelayAuthJWT()
	if err != nil {
		return fmt.Errorf("generate auth jwt: %w", err)
	}

	// 构造 relay URL (使用 auth 和 projectId 查询参数)
	projectId := common.G_Config["wcProjectId"]
	if projectId == "" {
		projectId = defaultRelayProjectId
	}
	relayAddr := fmt.Sprintf("%s/?auth=%s&projectId=%s", relayURL, url.QueryEscape(authToken), projectId)

	// 通过 Cloudflare DoH 解析真实 IP，绕过 DNS 污染
	relayHost := "relay.reown.com"
	realIP := resolveViaDoH(relayHost)

	dialer := websocket.Dialer{
		HandshakeTimeout: 15 * time.Second,
		TLSClientConfig: &tls.Config{
			ServerName: relayHost,
		},
	}

	// 如果 DoH 解析到了真实 IP，则使用自定义 NetDialContext 绕过本地 DNS
	if realIP != "" {
		fmt.Println("WalletConnect: DoH resolved", relayHost, "->", realIP)
		dialer.NetDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 将 relay.walletconnect.com:443 替换为 realIP:443
			_, port, _ := net.SplitHostPort(addr)
			return (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, network, net.JoinHostPort(realIP, port))
		}
	}

	// 设置 Origin header 避免被 relay 拒绝
	header := http.Header{}
	header.Set("Origin", "https://balawallet.com")

	fmt.Println("WalletConnect: dialing relay:", relayURL)
	conn, resp, err := dialer.Dial(relayAddr, header)
	if err != nil {
		if resp != nil {
			fmt.Printf("WalletConnect: relay HTTP status: %d %s\n", resp.StatusCode, resp.Status)
		}
		return fmt.Errorf("websocket dial: %w", err)
	}
	if resp != nil {
		resp.Body.Close()
	}

	fmt.Println("WalletConnect: connected to relay successfully")

	wcConnLock.Lock()
	wcConn = conn
	wcDone = make(chan struct{})
	wcConnLock.Unlock()

	// 订阅 pairing topic
	err = subscribeTopic(topic)
	if err != nil {
		conn.Close()
		return fmt.Errorf("subscribe topic: %w", err)
	}

	// 启动消息接收循环
	go relayReadLoop()

	return nil
}

// generateRequestId 生成 WC v2 relay 要求的大随机数 Request ID
// 格式: 时间戳(毫秒) * 1000 + 随机数，确保足够大且唯一
func generateRequestId() int64 {
	timestamp := time.Now().UnixMilli()
	randomPart, _ := rand.Int(rand.Reader, big.NewInt(1000))
	return timestamp*1000 + randomPart.Int64()
}

// subscribeTopic 向 Relay 订阅一个 topic
func subscribeTopic(topic string) error {
	msg := jsonRPCRequest{
		Id:      generateRequestId(),
		JsonRPC: "2.0",
		Method:  "irn_subscribe",
		Params: relaySubscribeParams{
			Topic: topic,
		},
	}
	return writeJSON(msg)
}

// relayReadLoop 持续从 WebSocket 读取消息
func relayReadLoop() {
	defer func() {
		wcConnLock.Lock()
		if wcDone != nil {
			close(wcDone)
			wcDone = nil
		}
		wcConnLock.Unlock()
	}()

	for {
		wcConnLock.Lock()
		conn := wcConn
		wcConnLock.Unlock()

		if conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WalletConnect: relay read error:", err)
			return
		}

		// 调试：打印原始消息前200字符
		msgStr := string(message)
		if len(msgStr) > 200 {
			msgStr = msgStr[:200] + "..."
		}
		fmt.Println("WalletConnect: raw relay message:", msgStr)

		go handleRelayMessage(message)
	}
}

// handleRelayMessage 处理从 Relay 收到的消息
func handleRelayMessage(message []byte) {
	// 首先尝试解析为 JSON-RPC request (来自 relay 的推送)
	var rpcReq jsonRPCRequest
	err := json.Unmarshal(message, &rpcReq)
	if err != nil {
		fmt.Println("WalletConnect: parse relay message error:", err)
		return
	}

	// 处理订阅推送 (irn_subscription)
	if rpcReq.Method == "irn_subscription" {
		handleSubscription(rpcReq)
		return
	}

	// 处理订阅响应
	if rpcReq.Id > 0 && rpcReq.Method == "" {
		// 这是一个 JSON-RPC 响应 (subscribe 或 publish 的回复)
		var rpcResp jsonRPCResponse
		json.Unmarshal(message, &rpcResp)
		if rpcResp.Result != nil {
			if subId, ok := rpcResp.Result.(string); ok {
				wcSubId = subId
				fmt.Println("WalletConnect: subscribed, id:", subId)
			}
		}
		return
	}
}

// handleSubscription 处理订阅推送的消息
func handleSubscription(rpcReq jsonRPCRequest) {
	// 解析 params
	paramsBytes, err := json.Marshal(rpcReq.Params)
	if err != nil {
		return
	}

	var subParams relaySubscriptionParams
	err = json.Unmarshal(paramsBytes, &subParams)
	if err != nil {
		return
	}

	topic := subParams.Data.Topic
	encMessage := subParams.Data.Message
	tag := subParams.Data.Tag

	// 确认收到消息 (ACK)
	ackSubscription(rpcReq.Id)

	// 解密消息
	wcSessionLock.RLock()
	session := wcSession
	wcSessionLock.RUnlock()
	if session == nil {
		return
	}

	var decrypted []byte

	// 判断用哪个密钥解密 (pairing 密钥 或 session 密钥)
	if topic == session.PairingTopic {
		// 用 pairing symKey 解密
		decrypted, err = decryptMessage(encMessage, session.SymKey)
	} else if topic == session.SessionTopic && session.SessionKey != nil {
		// 用 session key 解密
		decrypted, err = decryptMessageWithKey(encMessage, session.SessionKey)
	} else {
		fmt.Println("WalletConnect: unknown topic:", topic[:16]+"...")
		return
	}

	if err != nil {
		fmt.Println("WalletConnect: decrypt error:", err, "tag:", tag)
		return
	}

	fmt.Println("WalletConnect: received decrypted message on topic:", topic[:16]+"..., tag:", tag)

	// 解析内部 JSON-RPC 消息
	var innerRPC jsonRPCRequest
	err = json.Unmarshal(decrypted, &innerRPC)
	if err != nil {
		fmt.Println("WalletConnect: parse inner rpc error:", err)
		fmt.Println("WalletConnect: raw message:", string(decrypted))
		return
	}

	fmt.Println("WalletConnect: inner method:", innerRPC.Method, "id:", innerRPC.Id)

	// 根据方法类型处理
	switch innerRPC.Method {
	case "wc_sessionPropose":
		handleSessionPropose(innerRPC)
	case "wc_sessionRequest":
		handleSessionRequest(innerRPC)
	case "wc_sessionDelete":
		handleSessionDelete(innerRPC)
	case "wc_sessionPing":
		handleSessionPing(innerRPC)
	case "wc_sessionEvent":
		fmt.Println("WalletConnect: session event received")
	default:
		fmt.Println("WalletConnect: unhandled method:", innerRPC.Method)
	}
}

// ackSubscription 确认收到订阅消息
func ackSubscription(id int64) {
	resp := jsonRPCResponse{
		Id:      id,
		JsonRPC: "2.0",
		Result:  true,
	}
	writeJSON(resp)
}

// ================== Session Proposal 处理 ==================

// handleSessionPropose 处理 Dapp 的 session proposal
func handleSessionPropose(rpc jsonRPCRequest) {
	paramsBytes, _ := json.Marshal(rpc.Params)

	// WC v2 的 sessionPropose params 是一个结构体
	var proposal sessionProposeParams
	err := json.Unmarshal(paramsBytes, &proposal)
	if err != nil {
		fmt.Println("WalletConnect: parse session propose error:", err)
		return
	}

	fmt.Println("WalletConnect: session proposal from:", proposal.Proposer.Metadata.Name)
	fmt.Println("WalletConnect: url:", proposal.Proposer.Metadata.Url)

	// 更新会话信息
	wcSessionLock.Lock()
	if wcSession != nil {
		wcSession.DappName = proposal.Proposer.Metadata.Name
		wcSession.DappUrl = proposal.Proposer.Metadata.Url
		wcSession.DappDesc = proposal.Proposer.Metadata.Description
		if len(proposal.Proposer.Metadata.Icons) > 0 {
			wcSession.DappIcon = proposal.Proposer.Metadata.Icons[0]
		}
		wcSession.ProposalId = rpc.Id
	}
	wcSessionLock.Unlock()

	// 自动批准 session proposal
	approveSessionProposal(rpc.Id, proposal)
}

// approveSessionProposal 批准 session proposal
func approveSessionProposal(proposalId int64, proposal sessionProposeParams) {
	wcSessionLock.RLock()
	session := wcSession
	wcSessionLock.RUnlock()
	if session == nil {
		return
	}

	// 使用会话创建时生成的 X25519 公钥作为 responder 公钥
	responderPubKey := hex.EncodeToString(session.SelfPubKey)

	// 获取 proposer 的 X25519 公钥
	// WC v2 标准: proposer.publicKey 是主字段，proposerPublicKey 是兜底
	proposerPubKey := proposal.Proposer.PublicKey
	if proposerPubKey == "" {
		proposerPubKey = proposal.ProposerPubKey
	}
	fmt.Println("WalletConnect: proposer.publicKey:", proposal.Proposer.PublicKey)
	fmt.Println("WalletConnect: proposerPublicKey:", proposal.ProposerPubKey)
	fmt.Println("WalletConnect: using pubkey (len):", len(proposerPubKey), "hex:", proposerPubKey)

	// 使用 X25519 ECDH + HKDF 派生 session 对称密钥
	sessionKey, err := deriveSessionKey(session.SelfPrivKey, proposerPubKey)
	if err != nil {
		fmt.Println("WalletConnect: derive session key error:", err)
		return
	}

	// 生成 session topic = SHA256(sharedKey)
	topicHash := sha256.Sum256(sessionKey)
	sessionTopic := hex.EncodeToString(topicHash[:])

	// 保存 session 信息
	wcSessionLock.Lock()
	wcSession.SessionTopic = sessionTopic
	wcSession.SessionKey = sessionKey
	wcSession.Authorized = true
	wcSessionLock.Unlock()

	// 订阅 session topic
	subscribeTopic(sessionTopic)

	// 构造 approval 消息
	chainIdFull := "eip155:" + session.ChainId
	approval := map[string]interface{}{
		"relay": map[string]string{
			"protocol": "irn",
		},
		"responderPublicKey": responderPubKey,
		"controller": map[string]interface{}{
			"publicKey": responderPubKey,
			"metadata": map[string]interface{}{
				"name":        "Bala Wallet",
				"description": "A multi-chain crypto wallet",
				"url":         "https://balawallet.com",
				"icons":       []string{},
			},
		},
	}

	// 发送 session settle 到 session topic (批准)
	settleMsg := map[string]interface{}{
		"namespaces": map[string]interface{}{
			"eip155": map[string]interface{}{
				"accounts": []string{chainIdFull + ":" + session.Account},
				"methods": []string{
					"eth_sendTransaction",
					"personal_sign",
					"eth_sign",
					"eth_signTypedData",
					"eth_signTypedData_v4",
				},
				"events": []string{"chainChanged", "accountsChanged"},
				"chains": []string{chainIdFull},
			},
		},
		"controller": map[string]interface{}{
			"publicKey": responderPubKey,
			"metadata": map[string]interface{}{
				"name":        "Bala Wallet",
				"description": "A multi-chain crypto wallet",
				"url":         "https://balawallet.com",
				"icons":       []string{},
			},
		},
		"expiry": time.Now().Unix() + 7*24*3600, // 7 天过期
		"relay": map[string]string{
			"protocol": "irn",
		},
	}

	// 1. 回复 pairing topic 上的 session proposal (approval response)
	approvalResp := jsonRPCResponse{
		Id:      proposalId,
		JsonRPC: "2.0",
		Result:  approval,
	}
	approvalBytes, _ := json.Marshal(approvalResp)

	// 用 pairing symKey 加密并发送到 pairing topic
	encrypted, err := encryptMessage(approvalBytes, session.SymKey)
	if err != nil {
		fmt.Println("WalletConnect: encrypt approval error:", err)
		return
	}
	publishToTopic(session.PairingTopic, encrypted, 1101) // tag 1101 = session_propose response

	// 2. 在 session topic 上发送 wc_sessionSettle
	settleRPC := jsonRPCRequest{
		Id:      generateRequestId(),
		JsonRPC: "2.0",
		Method:  "wc_sessionSettle",
		Params:  settleMsg,
	}
	settleBytes, _ := json.Marshal(settleRPC)

	// 用 session key 加密并发送
	encSettled, err := encryptMessageWithKey(settleBytes, sessionKey)
	if err != nil {
		fmt.Println("WalletConnect: encrypt settle error:", err)
		return
	}
	publishToTopic(sessionTopic, encSettled, 1108) // tag 1108 = session_settle request

	fmt.Println("WalletConnect: session approved, session topic:", sessionTopic[:16]+"...")

	// 将授权信息加入请求队列，前端会显示
	wcRequestLock.Lock()
	wcRequests = append(wcRequests, WCRequest{
		Id:     0,
		Method: "_session_authorized",
		Params: map[string]interface{}{
			"name":    session.DappName,
			"url":     session.DappUrl,
			"icon":    session.DappIcon,
			"desc":    session.DappDesc,
			"account": session.Account,
			"chainId": session.ChainId,
		},
	})
	wcRequestLock.Unlock()
}

// ================== Session Request 处理 ==================

// handleSessionRequest 处理 Dapp 发来的签名/交易请求
func handleSessionRequest(rpc jsonRPCRequest) {
	paramsBytes, _ := json.Marshal(rpc.Params)

	var reqParams sessionRequestParams
	err := json.Unmarshal(paramsBytes, &reqParams)
	if err != nil {
		fmt.Println("WalletConnect: parse session request error:", err)
		return
	}

	method := reqParams.Request.Method
	params := reqParams.Request.Params

	fmt.Println("WalletConnect: session request, method:", method)

	// 将请求数据序列化
	var dataStr string
	if len(params) > 0 {
		switch method {
		case "personal_sign":
			// personal_sign params: [message, address]
			if msg, ok := params[0].(string); ok {
				dataStr = msg
			}
		case "eth_sign":
			// eth_sign params: [address, message]
			if len(params) >= 2 {
				if msg, ok := params[1].(string); ok {
					dataStr = msg
				}
			}
		case "eth_signTypedData", "eth_signTypedData_v4":
			// params: [address, typedDataJson]
			if len(params) >= 2 {
				if td, ok := params[1].(string); ok {
					dataStr = td
				} else {
					tdBytes, _ := json.Marshal(params[1])
					dataStr = string(tdBytes)
				}
			}
		case "eth_sendTransaction":
			// params: [txObject]
			txBytes, _ := json.Marshal(params[0])
			dataStr = string(txBytes)
		default:
			dataBytes, _ := json.Marshal(params)
			dataStr = string(dataBytes)
		}
	}

	// 添加到请求队列，前端轮询获取
	wcRequestLock.Lock()
	wcRequests = append(wcRequests, WCRequest{
		Id:     rpc.Id,
		Method: method,
		Params: dataStr,
		Topic:  "", // topic 由后端管理
	})
	wcRequestLock.Unlock()
}

// handleSessionDelete 处理会话断开
func handleSessionDelete(rpc jsonRPCRequest) {
	fmt.Println("WalletConnect: session deleted by Dapp")

	// 确认
	ackResponse(rpc.Id)

	// 标记为断开
	wcSessionLock.Lock()
	if wcSession != nil {
		wcSession.Connected = false
		wcSession.Authorized = false
	}
	wcSessionLock.Unlock()

	// 通知前端
	wcRequestLock.Lock()
	wcRequests = append(wcRequests, WCRequest{
		Id:     0,
		Method: "_session_disconnected",
		Params: "Dapp disconnected",
	})
	wcRequestLock.Unlock()
}

// handleSessionPing 回复 ping
func handleSessionPing(rpc jsonRPCRequest) {
	ackResponse(rpc.Id)
}

// ackResponse 发送通用确认响应
func ackResponse(id int64) {
	wcSessionLock.RLock()
	session := wcSession
	wcSessionLock.RUnlock()
	if session == nil || session.SessionKey == nil {
		return
	}

	resp := jsonRPCResponse{
		Id:      id,
		JsonRPC: "2.0",
		Result:  true,
	}
	respBytes, _ := json.Marshal(resp)
	encrypted, err := encryptMessageWithKey(respBytes, session.SessionKey)
	if err != nil {
		return
	}
	publishToTopic(session.SessionTopic, encrypted, 1109)
}

// sendSignResult 将签名结果通过 Relay 发送给 Dapp
func sendSignResult(requestId int64, result string, signErr error) {
	wcSessionLock.RLock()
	session := wcSession
	wcSessionLock.RUnlock()
	if session == nil || session.SessionKey == nil {
		return
	}

	var resp jsonRPCResponse
	if signErr != nil {
		resp = jsonRPCResponse{
			Id:      requestId,
			JsonRPC: "2.0",
			Error: &rpcError{
				Code:    5000,
				Message: signErr.Error(),
			},
		}
	} else {
		resp = jsonRPCResponse{
			Id:      requestId,
			JsonRPC: "2.0",
			Result:  result,
		}
	}

	respBytes, _ := json.Marshal(resp)
	encrypted, err := encryptMessageWithKey(respBytes, session.SessionKey)
	if err != nil {
		fmt.Println("WalletConnect: encrypt sign result error:", err)
		return
	}

	publishToTopic(session.SessionTopic, encrypted, 1109) // tag 1109 = session_request response
	fmt.Println("WalletConnect: sign result sent to Dapp")
}

// ================== Relay 消息发布 ==================

// publishToTopic 向指定 topic 发布加密消息
func publishToTopic(topic string, encryptedMsg string, tag int) {
	msg := jsonRPCRequest{
		Id:      generateRequestId(),
		JsonRPC: "2.0",
		Method:  "irn_publish",
		Params: relayPublishParams{
			Topic:   topic,
			Message: encryptedMsg,
			TTL:     300,
			Tag:     tag,
			Prompt:  false,
		},
	}
	err := writeJSON(msg)
	if err != nil {
		fmt.Println("WalletConnect: publish error:", err)
	}
}

// writeJSON 安全地写入 WebSocket 消息
func writeJSON(v interface{}) error {
	wcConnLock.Lock()
	defer wcConnLock.Unlock()
	if wcConn == nil {
		return fmt.Errorf("connection closed")
	}
	return wcConn.WriteJSON(v)
}

// ================== 加密/解密 (ChaCha20-Poly1305) ==================

// encryptMessage 使用 hex 编码的 symKey 加密消息，返回 base64 编码的 type0 envelope
func encryptMessage(plaintext []byte, symKeyHex string) (string, error) {
	keyBytes, err := hex.DecodeString(symKeyHex)
	if err != nil {
		return "", fmt.Errorf("decode symKey: %w", err)
	}
	return encryptMessageWithKey(plaintext, keyBytes)
}

// encryptMessageWithKey 使用 ChaCha20-Poly1305 加密，返回 base64 编码的 type0 envelope
func encryptMessageWithKey(plaintext []byte, key []byte) (string, error) {
	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return "", fmt.Errorf("chacha20poly1305: %w", err)
	}

	// 生成随机 nonce (12 bytes)
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	// 加密: ciphertext + tag
	sealed := aead.Seal(nil, nonce, plaintext, nil)

	// WC v2 type0 envelope: type(1 byte, 0x00) + nonce(12) + sealed(ciphertext + tag)
	envelope := make([]byte, 0, 1+len(nonce)+len(sealed))
	envelope = append(envelope, 0x00) // type0
	envelope = append(envelope, nonce...)
	envelope = append(envelope, sealed...)

	return base64.StdEncoding.EncodeToString(envelope), nil
}

// decryptMessage 使用 hex 编码的 symKey 解密 base64 编码的消息
func decryptMessage(encMessageB64 string, symKeyHex string) ([]byte, error) {
	keyBytes, err := hex.DecodeString(symKeyHex)
	if err != nil {
		return nil, fmt.Errorf("decode symKey: %w", err)
	}
	return decryptMessageWithKey(encMessageB64, keyBytes)
}

// decryptMessageWithKey 解密 base64 编码的 WC v2 envelope
func decryptMessageWithKey(encMessageB64 string, key []byte) ([]byte, error) {
	encBytes, err := base64.StdEncoding.DecodeString(encMessageB64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode message: %w", err)
	}

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, fmt.Errorf("chacha20poly1305: %w", err)
	}

	nonceSize := aead.NonceSize()

	// WC v2 envelope 格式: type(1 byte) + nonce(12) + sealed(ciphertext + tag)
	if len(encBytes) < 1 {
		return nil, fmt.Errorf("decrypt failed: message too short")
	}

	envType := encBytes[0]
	payload := encBytes[1:] // 跳过 type 字节

	if envType == 0 {
		// type0: 对称密钥解密
		if len(payload) < nonceSize+aead.Overhead() {
			return nil, fmt.Errorf("decrypt failed: type0 payload too short")
		}
		nonce := payload[:nonceSize]
		ciphertext := payload[nonceSize:]
		return aead.Open(nil, nonce, ciphertext, nil)
	}

	// 兜底：尝试无 type 前缀解密 (兼容某些实现)
	if len(encBytes) >= nonceSize+aead.Overhead() {
		nonce := encBytes[:nonceSize]
		ciphertext := encBytes[nonceSize:]
		plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
		if err == nil {
			return plaintext, nil
		}
	}

	return nil, fmt.Errorf("decrypt failed: unsupported envelope type %d", envType)
}

// deriveSessionKey 使用 X25519 ECDH + HKDF-SHA256 派生 session 对称密钥
func deriveSessionKey(selfPrivKey []byte, peerPubKeyHex string) ([]byte, error) {
	if peerPubKeyHex == "" {
		return nil, fmt.Errorf("peer public key is empty")
	}

	peerPubKey, err := hex.DecodeString(peerPubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("decode peer pubkey: %w", err)
	}

	if len(peerPubKey) != 32 {
		return nil, fmt.Errorf("invalid peer pubkey length: got %d bytes, want 32", len(peerPubKey))
	}

	// X25519 ECDH 共享密钥
	sharedSecret, err := curve25519.X25519(selfPrivKey, peerPubKey)
	if err != nil {
		return nil, fmt.Errorf("x25519 ecdh: %w", err)
	}

	// HKDF-SHA256 派生 32 字节密钥
	hkdfReader := hkdf.New(sha256.New, sharedSecret, nil, nil)
	derivedKey := make([]byte, 32)
	if _, err := hkdfReader.Read(derivedKey); err != nil {
		return nil, fmt.Errorf("hkdf expand: %w", err)
	}

	return derivedKey, nil
}

// ================== 工具函数 ==================

// parseWCUri 解析 WalletConnect URI
type wcUriParsed struct {
	Topic    string
	Version  string
	Relay    string
	SymKey   string
	ExpireAt int64
}

func parseWCUri(uri string) (*wcUriParsed, error) {
	if !strings.HasPrefix(uri, "wc:") {
		return nil, fmt.Errorf("invalid WalletConnect URI: must start with 'wc:'")
	}

	withoutPrefix := uri[3:]
	atIndex := strings.Index(withoutPrefix, "@")
	if atIndex == -1 {
		return nil, fmt.Errorf("invalid WalletConnect URI: missing '@'")
	}

	topic := withoutPrefix[:atIndex]
	rest := withoutPrefix[atIndex+1:]

	qIndex := strings.Index(rest, "?")
	if qIndex == -1 {
		return nil, fmt.Errorf("invalid WalletConnect URI: missing parameters")
	}

	version := rest[:qIndex]
	paramStr := rest[qIndex+1:]

	params, err := url.ParseQuery(paramStr)
	if err != nil {
		return nil, fmt.Errorf("invalid WalletConnect URI: bad parameters")
	}

	result := &wcUriParsed{
		Topic:   topic,
		Version: version,
		Relay:   params.Get("relay-protocol"),
		SymKey:  params.Get("symKey"),
	}

	if expStr := params.Get("expiryTimestamp"); expStr != "" {
		exp := new(big.Int)
		exp.SetString(expStr, 10)
		result.ExpireAt = exp.Int64()
	}

	return result, nil
}

// disconnectInternal 断开 WebSocket 连接并清理状态
func disconnectInternal() {
	// 关闭 WebSocket
	wcConnLock.Lock()
	if wcConn != nil {
		wcConn.Close()
		wcConn = nil
	}
	wcConnLock.Unlock()

	// 等待 readLoop 结束
	wcConnLock.Lock()
	done := wcDone
	wcConnLock.Unlock()
	if done != nil {
		select {
		case <-done:
		case <-time.After(2 * time.Second):
		}
	}

	// 清除会话
	wcSessionLock.Lock()
	wcSession = nil
	wcSessionLock.Unlock()

	wcRequestLock.Lock()
	wcRequests = make([]WCRequest, 0)
	wcRequestLock.Unlock()

	fmt.Println("WalletConnect: disconnected")
}

// ================== Ed25519 JWT 认证 (WalletConnect Relay Auth) ==================

// getRelayAuthKeyPair 获取或生成 Ed25519 密钥对
func getRelayAuthKeyPair() (ed25519.PublicKey, ed25519.PrivateKey) {
	relayAuthKeyOnce.Do(func() {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Println("WalletConnect: generate ed25519 key error:", err)
			return
		}
		relayAuthPubKey = pub
		relayAuthKey = priv
	})
	return relayAuthPubKey, relayAuthKey
}

// generateRelayAuthJWT 生成用于 WalletConnect Relay 认证的 JWT
// 规范: https://specs.walletconnect.com/2.0/specs/clients/core/relay/relay-client-auth
func generateRelayAuthJWT() (string, error) {
	pubKey, privKey := getRelayAuthKeyPair()
	if privKey == nil {
		return "", fmt.Errorf("ed25519 key pair not initialized")
	}

	// 构造 did:key 格式的 issuer
	iss := encodeDidKey(pubKey)

	// 生成随机 sub (nonce)
	subBytes := make([]byte, 32)
	rand.Read(subBytes)
	sub := hex.EncodeToString(subBytes)

	now := time.Now().Unix()
	exp := now + 86400 // 24 小时有效期

	// JWT Header
	header := map[string]string{
		"alg": "EdDSA",
		"typ": "JWT",
	}

	// JWT Payload
	payload := map[string]interface{}{
		"iss": iss,
		"sub": sub,
		"aud": relayAudience,
		"iat": now,
		"exp": exp,
		"act": "client_auth",
	}

	// base64url 编码 header 和 payload
	headerBytes, _ := json.Marshal(header)
	payloadBytes, _ := json.Marshal(payload)

	headerB64 := base64URLEncode(headerBytes)
	payloadB64 := base64URLEncode(payloadBytes)

	// 签名数据 = header.payload
	signingInput := headerB64 + "." + payloadB64

	// 使用 Ed25519 签名
	signature := ed25519.Sign(privKey, []byte(signingInput))

	// 组装 JWT = header.payload.signature
	jwt := signingInput + "." + base64URLEncode(signature)

	return jwt, nil
}

// encodeDidKey 将 Ed25519 公钥编码为 did:key 格式
// 格式: did:key:z<base58btc(multicodec_header + publicKey)>
func encodeDidKey(pubKey ed25519.PublicKey) string {
	// Ed25519 multicodec header: 0xed, 0x01
	multicodecHeader := []byte{0xed, 0x01}

	// 拼接 header + 公钥
	data := make([]byte, 0, len(multicodecHeader)+len(pubKey))
	data = append(data, multicodecHeader...)
	data = append(data, pubKey...)

	// Base58btc 编码
	encoded := base58.Encode(data)

	// 前缀 z (multibase 标识) + base58 编码
	return "did:key:z" + encoded
}

// base64URLEncode 标准 base64url 编码 (无 padding)
func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}
