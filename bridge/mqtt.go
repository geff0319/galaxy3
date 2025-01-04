package bridge

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ge-fei-fan/gefflog"
	"time"
)

type Subscribe struct {
	topic    string
	qos      byte
	callback mqtt.MessageHandler
}
type MqttInfo struct {
	Broker               string      `json:"broker"`
	Port                 int         `json:"port"`
	ClientID             string      `json:"clientID"`
	UserName             string      `json:"userName"`
	Password             string      `json:"password"`
	SuccessSubscribeList []string    `json:"successSubscribeList"`
	SubscribeList        []Subscribe `json:"-"`
}
type MqttClient struct {
	Mqtt   MqttInfo
	opt    *mqtt.ClientOptions
	Client mqtt.Client
}

var MqttC = &MqttClient{
	Mqtt: MqttInfo{},
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	MqttC.Mqtt.subscribe(client)
	gefflog.Info(fmt.Sprintf("Connected"))
}
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	gefflog.Info(fmt.Sprintf("mqtt Connect lost: %v", err))
}

func (mi *MqttInfo) subscribe(client mqtt.Client) {
	if client == nil {
		return
	}
	for _, s := range mi.SubscribeList {
		token := client.Subscribe(s.topic, s.qos, s.callback)
		if token.WaitTimeout(10) && token.Error() != nil {
			gefflog.Err(fmt.Sprintf("mqtt 订阅%s失败: %v", s.topic, token.Error()))
		} else {
			mi.SuccessSubscribeList = append(mi.SuccessSubscribeList, s.topic)
		}
	}
	gefflog.Info(fmt.Sprintf("订阅成功：%v", mi.SuccessSubscribeList))
}

func (mc *MqttClient) init() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mc.Mqtt.Broker, mc.Mqtt.Port))
	opts.SetClientID(mc.Mqtt.ClientID)
	opts.SetUsername(mc.Mqtt.UserName)
	opts.SetPassword(mc.Mqtt.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.CleanSession = false
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	MqttC.opt = opts
	MqttC.Client = mqtt.NewClient(opts)
	mc.initSubscribe()
}

func (mc *MqttClient) initSubscribe() {
	mc.Mqtt.SuccessSubscribeList = make([]string, 0, 50)
	mc.Mqtt.SubscribeList = make([]Subscribe, 0, 50)
	mc.Mqtt.SubscribeList = append(mc.Mqtt.SubscribeList, Subscribe{
		topic:    fmt.Sprintf(BILIBIL_DOWNLOAD_TOPIC, mc.Mqtt.ClientID),
		qos:      0,
		callback: BilibiliPubHandler,
	}, Subscribe{
		topic:    fmt.Sprintf(YOUTUBE_DOWNLOAD_TOPIC, mc.Mqtt.ClientID),
		qos:      0,
		callback: YoutubePubHandler,
	}, Subscribe{
		topic:    fmt.Sprintf(TWITTER_DOWNLOAD_TOPIC, mc.Mqtt.ClientID),
		qos:      0,
		callback: TwitterPubHandler,
	})
}

func (mc *MqttClient) connect() error {
	if token := mc.Client.Connect(); token.WaitTimeout(10*time.Second) && token.Error() != nil {
		gefflog.Err("mqtt连接失败：" + token.Error().Error())
		return token.Error()
	}
	return nil
}

func (mc *MqttClient) disconnect() {
	if mc.Client != nil && mc.Client.IsConnected() {
		mc.Client.Disconnect(3)
	}
}

func (mc *MqttClient) status() bool {
	return mc.Client.IsConnected()
}

func (mc *MqttClient) getSubscribe() {
	mc.Client.IsConnectionOpen()
}
func (a *App) ConnectMqtt(mqttInfo string) FlagResult {

	if err := json.Unmarshal([]byte(mqttInfo), &MqttC.Mqtt); err != nil {
		return FlagResult{false, "初始化mqtt信息失败!"}
	}
	fmt.Println(MqttC.Mqtt)
	MqttC.init()
	err := MqttC.connect()
	if err != nil {
		return FlagResult{Flag: false, Data: "mqtt连接失败:" + err.Error()}
	}
	return FlagResult{Flag: true, Data: "mqtt连接成功"}
}
func (a *App) DisconnectMqtt() FlagResult {
	MqttC.disconnect()
	return FlagResult{Flag: true, Data: "mqtt断开连接成功"}
}
func (a *App) StatusMqtt() FlagResultWithData {
	if MqttC.Client == nil {
		return FlagResultWithData{Flag: true, Msg: "查询成功", Data: false}
	}
	return FlagResultWithData{Flag: true, Msg: "查询成功", Data: MqttC.Client.IsConnected()}
}
