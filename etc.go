package fiberinterstate

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v3"
)

// -------------------------------------------------------------------------------------------------------------------------------------

var (
	green_clr  = color.New(color.FgGreen).SprintFunc()
	white_clr  = color.New(color.Bold, color.FgWhite).SprintFunc()
	cian_clr   = color.New(color.FgCyan).SprintFunc()
	yellow_clr = color.New(color.FgYellow).SprintFunc()

	blue_clr = color.New(color.FgBlue).SprintFunc()
)

// -------------------------------------------------------------------------------------------------------------------------------------

const (
	//socketFile   = "fiberinterstate.sock"
	dialTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	sendInterval = 2 * time.Second
)

func socketPath() string {

	socketFile := fmt.Sprintf("fiberinterstate-%d.sock", os.Getpid())

	if fiber.IsChild() {

		socketFile = fmt.Sprintf("fiberinterstate-%d.sock", os.Getppid())
	}

	//fmt.Println(os.Getpid(), " socketPath()", "socketFile=", socketFile)

	return filepath.Join(os.TempDir(), socketFile)
}

// -------------------------------------------------------------------------------------------------------------------------------------

type InnerMessage struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
	Data []byte `json:"data"`
}

func NewInnerMessage(data []byte) *InnerMessage {
	return &InnerMessage{
		Data: data,
	}
}

func (im *InnerMessage) Encode() []byte {

	bts, err := json.Marshal(im)
	if err != nil {
		panic(err)
	}

	return bts
}

func (im *InnerMessage) Decode() {

	var el InnerMessage

	err := json.Unmarshal(im.Data, &el)
	if err != nil {
		panic(err)
	}

	im.Id = el.Id
	im.Type = el.Type
	im.Data = el.Data

	fmt.Printf("Decode JSON: %+v\n", el)

}

func (im *InnerMessage) SetId(id string) {
	im.Id = id
}

func (im *InnerMessage) SetType(type_ int) {
	im.Type = type_
}

func (im *InnerMessage) SetData(data_ []byte) {
	im.Data = data_
}

// -------------------------------------------------------------------------------------------------------------------------------------

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type EventLogConfig struct {
	isClient bool
	isServer bool
	Tag      string
	Msg      string
}

// | [pid] [datetime] [tag] message

// TODO: | [pid] [datetime] [tag, tag2, tag3] message

func EventLogAdd(tag string, msg ...string) {

	pref := ""

	pref += blue_clr("| ")

	pid := strconv.Itoa(os.Getpid())
	pref += "[" + pid + "] "

	_dt := time.Now().Format("2006-01-02 15:04:05.000")
	pref += "[" + _dt + "] "

	pref += "[" + tag + "] "

	msgConcat := strings.Join(msg, " ")

	pref += msgConcat

	if FiberInterstateConfig.LogEnable {
		fmt.Println(pref)
	}
}

/*
func EventLogAdd(c fiber.Config, status int, tag string, msg ...string) {

	pref := ""

	if status != 200 {
		pref += red_clr("| ")
	} else {
		pref += green_clr("| ")
	}

	pid := strconv.Itoa(os.Getpid())
	pref += "[" + pid + "] "

	_dt := time.Now().Format("2006-01-02 15:04:05.000")
	pref += "[" + _dt + "] "

	_ip := ""
	if c != nil {
		_ip = c.IP()
		pref += "[" + _ip + "] "
	} else {
		pref += "[] "
	}

	arg_silence := false
	arg_nolog := false
	if c != nil {

		arg_silence = GetBoolFromLocals(c, "arg_silence")
		arg_nolog = GetBoolFromLocals(c, "arg_nolog")
	}

	_user := ""
	if c != nil {

		_user = basicauth.UsernameFromContext(c)
	}

	if len(_user) > 0 {
		pref += "[" + green_clr(_user) + "] "
	} else {
		pref += "[] "
	}

	if status != 200 {
		pref += "[" + red_clr(status) + "] "
	} else {
		pref += "[" + green_clr(status) + "] "
	}

	pref += "[" + tag + "] "

	// ----------------------------------------------------------------------------
	// msg join

	msgConcat := strings.Join(msg, " ")

	// ----------------------------------------------------------------------------

	//msgConcat = strings.Replace(msgConcat, conf.ArgFold, white_clr(conf.ArgFold), 1)

	pref += msgConcat

	if tag != "INIT" && !arg_silence {
		fmt.Println(pref)
	}

	// ----------------------------------------------------------------------------


}
*/

/*
func GetStringFromLocals(c fiber.Ctx, key string, defaultValue string) string {
	if val, ok := c.Locals(key).(string); ok {
		return val
	}
	return defaultValue
}


func GetIntFromLocals(c fiber.Ctx, key string, defaultValue int) int {
	if val, ok := c.Locals(key).(int); ok {
		return val
	}
	return defaultValue
}
*/
