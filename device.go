package lpr

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"net/http"
	"time"
)

type Device struct {
	Url      string
	Username string
	Password string

	client  http.Client
	hasher  hash.Hash
	nonce   string
	nonce64 string
}

func NewDevice(url, username, password string, timeout time.Duration) *Device {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return &Device{
		Url:      url,
		Username: username,
		Password: password,

		client:  http.Client{Timeout: timeout},
		hasher:  sha1.New(),
		nonce:   string(buf),
		nonce64: base64.StdEncoding.EncodeToString(buf),
	}
}

func (device *Device) StartPullingRecognitions(callback func(rec *Recognition, err error)) {
	var pullAddress string
	var err error
	until := time.Now()

	for {
		if until.Before(time.Now()) {
			pullAddress, err = device.createPullPointSubscription()
			if err != nil {
				callback(nil, err)
				time.Sleep(5 * time.Second)
				continue
			}
			until = time.Now().Add(2 * time.Minute)
		}

		rec, err := device.pullMessages(pullAddress)
		if err != nil {
			callback(nil, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if rec != nil {
			callback(rec, nil)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (device *Device) createPullPointSubscriptionXml() string {
	token, nonce64, timestamp := device.generateCredentials()
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?><s:Envelope xmlns:s=\"http://www.w3.org/2003/05/soap-envelope\"><s:Header><Action mustUnderstand=\"1\" xmlns=\"http://www.w3.org/2005/08/addressing\">http://www.onvif.org/ver10/events/wsdl/EventPortType/CreatePullPointSubscriptionRequest</Action><Security s:mustUnderstand=\"1\" xmlns=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd\"><UsernameToken><Username>" + device.Username + "</Username><Password Type=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordDigest\">" + token + "</Password><Nonce EncodingType=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary\">" + nonce64 + "</Nonce><Created xmlns=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd\">" + timestamp + "</Created></UsernameToken></Security></s:Header><s:Body><CreatePullPointSubscription xmlns=\"http://www.onvif.org/ver10/events/wsdl\"><InitialTerminationTime>PT300S</InitialTerminationTime></CreatePullPointSubscription></s:Body></s:Envelope>"
}

func (device *Device) createPullPointSubscription() (string, error) {
	xml := device.createPullPointSubscriptionXml()

	body, err := device.doRequest(xml)
	if err != nil {
		return "", fmt.Errorf("failed to create pull point subscription: %s", err)
	}

	address, err := parsePullAddress(body)
	if err != nil {
		return "", fmt.Errorf("failed to create pull point subscription: %s", err)
	}

	return address, nil
}

func (device *Device) pullMessagesXml(pullAddress string) string {
	token, nonce64, timestamp := device.generateCredentials()
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?><s:Envelope xmlns:s=\"http://www.w3.org/2003/05/soap-envelope\"><s:Header><Action mustUnderstand=\"1\" xmlns=\"http://www.w3.org/2005/08/addressing\">http://www.onvif.org/ver10/events/wsdl/PullPointSubscription/PullMessagesRequest</Action><Security s:mustUnderstand=\"1\" xmlns=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd\"><UsernameToken><Username>" + device.Username + "</Username><Password Type=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordDigest\">" + token + "</Password><Nonce EncodingType=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary\">" + nonce64 + "</Nonce><Created xmlns=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd\">" + timestamp + "</Created></UsernameToken></Security><wsa:To>" + pullAddress + "</wsa:To></s:Header><s:Body><PullMessages xmlns=\"http://www.onvif.org/ver10/events/wsdl\"><Timeout>PT3S</Timeout><MessageLimit>10</MessageLimit></PullMessages></s:Body></s:Envelope>"
}

func (device *Device) pullMessages(pullAddress string) (*Recognition, error) {
	xml := device.pullMessagesXml(pullAddress)

	body, err := device.doRequest(xml)
	if err != nil {
		return nil, fmt.Errorf("failed to pull messages: %s", err)
	}

	return parseRecognition(body), nil
}

func (device *Device) doRequest(body string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, device.Url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/soap+xml")
	req.Header.Set("Charset", "utf-8")

	res, err := device.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (device *Device) generateCredentials() (string, string, string) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	device.hasher.Reset()
	device.hasher.Write([]byte(device.nonce + timestamp + device.Password))
	return base64.StdEncoding.EncodeToString(device.hasher.Sum(nil)), device.nonce64, timestamp
}
