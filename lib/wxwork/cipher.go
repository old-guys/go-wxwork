package lib_wxwork

import (
	"encoding/base64"
	"errors"
	"crypto/aes"
	"io"
	"crypto/cipher"
	"crypto/rand"
	"bytes"
	"sort"
	"crypto/sha1"
	"strings"
	"encoding/hex"
	"regexp"
	"github.com/astaxie/beego"
	"encoding/json"
	"time"

	"github.com/astaxie/beego/utils"
	"strconv"
	"encoding/xml"
)

type Config struct {

}

var (
	Cipher Config
	wxwork_key = beego.AppConfig.String("wxwork_key")
	wxwork_token = beego.AppConfig.String("wxwork_token")
	wxwork_aes_key = beego.AppConfig.String("wxwork_aes_key")
)

type msgXml struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`

	AgentID      string `xml:"AgentID"`

	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`

	Latitude     string `xml:"Latitude"`
	Longitude    string `xml:"Longitude"`
	Precision    string `xml:"Precision"`

	JobId        string `xml:"JobId"`
	JobType      string `xml:"JobType"`
	ErrCode      string `xml:"ErrCode"`
	ErrMsg       string `xml:"ErrMsg"`

	ChangeType   string `xml:"ChangeType"`
	UserID       string `xml:"UserID"`
	NewUserID    string `xml:"NewUserID"`
	Name         string `xml:"Name"`
	Department   string `xml:"Department"`
	Mobile       string `xml:"Mobile"`
	Position     string `xml:"Position"`
	Gender       string `xml:"Gender"`
	Email        string `xml:"Email"`
	Status       string `xml:"Status"`
	Avatar       string `xml:"Avatar"`
	EnglishName  string `xml:"EnglishName"`
	IsLeader     string `xml:"IsLeader"`
	Telephone    string `xml:"Telephone"`
	//ExtAttr

	Id           string `xml:"Id"`
	ParentId     string `xml:"ParentId"`
	Order        string `xml:"Order"`

	TagId         string `xml:"TagId"`
	AddUserItems  string `xml:"AddUserItems"`
	DelUserItems  string `xml:"DelUserItems"`
	AddPartyItems string `xml:"AddPartyItems"`
	DelPartyItems string `xml:"DelPartyItems"`
}

func (c *Config) Decrypt(args map[string]interface{}) (decrypted_data interface{}, err error) {
	key := args["key"].(string)
	if len(key) == 0 {
		key = wxwork_key
	}

	aes_key := args["aes_key"].(string)
	if len(aes_key) == 0 {
		aes_key = wxwork_aes_key
	}

	data := args["data"].(string)

	byte_aes_key, err := base64.StdEncoding.DecodeString(aes_key + "=")
	Logger.Info("wxwork === decrypt: byte_aes_key =", string(byte_aes_key), "err = ", err)
	if err != nil { return nil, err }

	encrypted_data, err := base64.StdEncoding.DecodeString(data)
	Logger.Info("wxwork === decrypt: decode encrypted_data =", string(encrypted_data), "err =", err)
	if err != nil { return nil, err }

	bak_len := len(byte_aes_key)
	if len(encrypted_data) % bak_len != 0 {
		Logger.Info("wxwork === decrypt: ciphertext size is not multiple of aes key length")
		return nil, errors.New("crypto/cipher: ciphertext size is not multiple of aes key length")
	}

	block, err := aes.NewCipher(byte_aes_key)
	if err != nil {
		Logger.Info("wxwork === decrypt: NewCipher err =", err)
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		Logger.Info("wxwork === decrypt: cipher get iv err =", err)
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	plainData := make([]byte, len(encrypted_data))
	blockMode.CryptBlocks(plainData, encrypted_data)

	// remove PKCS7 padding
	plainData = Utils.PKCS7Padding.RemovePadding(plainData)

	// remove 16 Bytes random string
	plainData = plainData[16 : len(plainData)]

	// 4 Bytes length
	length := Utils.NetworkOrder.BytesToInt(plainData[0:4])

	plainData = plainData[4 : len(plainData)]
	msg_data := string(plainData[0 : length]) // Msg
	received_key := string(plainData[length : len(plainData)]) // key

	Logger.Info("wxwork === decrypt: msg_data =", msg_data)
	Logger.Info("wxwork === decrypt: received_key =", received_key)

	if key != received_key {
		Logger.Info("wxwork === decrypt: key error, received_key =", received_key, "key =", key)
		return nil, errors.New("key error")
	}

	reg := regexp.MustCompile(`^[\d]*$`)
	if reg.MatchString(msg_data) {
		return msg_data, nil
	}

	reg = regexp.MustCompile(`^<xml+`)
	if !reg.MatchString(msg_data) {
		err = json.Unmarshal([]byte(msg_data), &decrypted_data)
		if err != nil {
			Logger.Info("wxwork === decrypt: parse json failure, msg_data format is string, err =", err)

			return msg_data, err
		} else {
			decrypted_data = c.InterfaceToString(decrypted_data)
			Logger.Info("wxwork === decrypt: parse json success, decrypted_data =", decrypted_data)
			return decrypted_data, nil
		}
	}

	mx := msgXml{}
	err = xml.Unmarshal([]byte(msg_data), &mx)
	if err != nil {
		Logger.Info("wxwork === decrypt: parse xml failure, err =", err)
		return nil, err
	}

	by, err := json.Marshal(mx)
	if err != nil {
		Logger.Info("wxwork === decrypt: json(struct, xml) to byte failure, mx =", mx, "err =", err)
		return nil, err
	}

	err = json.Unmarshal(by, &decrypted_data)
	if err != nil {
		Logger.Info("wxwork === decrypt: byte to map failure, by =", string(by), "err =", err)
		return nil, err
	}

	Logger.Info("wxwork === decrypt: decrypted_data =", decrypted_data)

	return decrypted_data, err
}

func (c *Config) Encrypt(args map[string]interface{}) (string, error)  {
	params := args["params"]
	json_data := ""

	// reflect.TypeOf(params).Name()
	bts, err := json.Marshal(params)
	if err != nil {
		json_data = params.(string)
	} else {
		json_data = string(bts)
	}
	Logger.Info("wxwork === encrypt: json_data =", json_data)

	key := args["key"].(string)
	if len(key) == 0 {
		key = wxwork_key
	}

	aes_key := args["aes_key"].(string)
	if len(aes_key) == 0 {
		aes_key = wxwork_aes_key
	}

	byte_aes_key, err := base64.StdEncoding.DecodeString(aes_key + "=")
	if err != nil { return "", nil }

	// add 16 Bytes random string
	// random_str := []byte("abcdefghijklmnop")
	random_str := utils.RandomCreateBytes(16)
	Logger.Info("wxwork === encrypt: random_str =", string(random_str))

	// 4 Bytes length
	len_hex_str := Utils.NetworkOrder.IntToBytes(len(json_data))
	Logger.Info("wxwork === encrypt: len_hex_str =", string(len_hex_str))
	Logger.Info("wxwork === encrypt: len_hex_str byte length =", len(len_hex_str))

	data := bytes.Join([][]byte{random_str, len_hex_str, []byte(json_data), []byte(key)}, nil)
	Logger.Info("wxwork === encrypt: data =", string(data))
	Logger.Info("wxwork === encrypt: data byte length =", len(data))

	// add PKCS7 padding
	data = Utils.PKCS7Padding.AddPadding(data)
	Logger.Info("wxwork === encrypt:  encrypt_data =", string(data))
	Logger.Info("wxwork === encrypt:  encrypt_data byte length =", len(data))

	bak_len := len(byte_aes_key)
	if len(data) % bak_len != 0 {
		Logger.Info("wxwork === encrypt: ciphertext size is not multiple of aes key length")
		return "", errors.New("crypto/cipher: ciphertext size is not multiple of aes key length")
	}

	block, err := aes.NewCipher(byte_aes_key)
	if err != nil {
		Logger.Info("wxwork === encrypt: NewCipher err =", err)
		return "", err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		Logger.Info("wxwork === encrypt: cipher get iv err =", err)
		return "", err
	}

	encrypted_data := make([]byte, len(data))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(encrypted_data, data)

	encypted_base64_data := base64.StdEncoding.EncodeToString(encrypted_data)
	Logger.Info("wxwork === encrypt: encypted_base64_data =", encypted_base64_data)

	return encypted_base64_data, nil
}

func (c *Config) Sign(args map[string]interface{}) (string, error) {
	nonce := args["nonce"].(string)
	if len(nonce) == 0 {
		nonce = string(utils.RandomCreateBytes(8))
	}

	timestamp := args["timestamp"].(string)
	if len(timestamp) == 0 {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	token := args["token"].(string)
	if len(token) == 0 {
		token = beego.AppConfig.String("wxwork_token")
	}

	encrypt := args["encrypt"].(string)
	Logger.Info("wxwork === sign: token =", token, ", nonce =", nonce, ", timestamp =", timestamp, ", encrypt =", encrypt)

	data := []string{encrypt, nonce, timestamp, token}
	sort.Strings(data)

	str := strings.Join(data, "")
	h := sha1.New()
	if _, err := h.Write([]byte(str)); err != nil {
		Logger.Info("wxwork === sign: sha1 write err =", err)
		return "", err
	}

	bs := h.Sum(nil)
	signature := hex.EncodeToString(bs)

	Logger.Info("wxwork === sign: signature =", signature)

	return signature, nil
}

func (c *Config) JsapiSign(args map[string]interface{}) (string, error) {
	jsapi_ticket := args["jsapi_ticket"].(string)
	url := args["url"].(string)

	noncestr := args["noncestr"].(string)
	if len(noncestr) == 0 {
		noncestr = string(utils.RandomCreateBytes(8))
	}

	timestamp := args["timestamp"].(string)
	if len(timestamp) == 0 {
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	}

	arr := []string{"jsapi_ticket=", jsapi_ticket, "&noncestr=", noncestr, "&timestamp=", timestamp, "&url=", url}
	data := strings.Join(arr, "")

	Logger.Info("wxwork === jsapi_sign: data =", data)

	h := sha1.New()
	if _, err := h.Write([]byte(data)); err != nil {
		Logger.Info("wxwork === jsapi_sign: sha1 write err =", err)
		return "", err
	}

	bs := h.Sum(nil)
	signature := hex.EncodeToString(bs)

	Logger.Info("wxwork === jsapi_sign: signature =", signature)

	return signature, nil
}

func (c *Config) InterfaceToString(obj interface{}) (data interface{}) {
	data = obj

	switch obj := obj.(type) {
		case float64:
			data = strconv.FormatFloat(obj, 'f', -1, 64)
		case float32:
			data = strconv.FormatFloat(float64(obj), 'f', -1, 32)

	}

	return data
}