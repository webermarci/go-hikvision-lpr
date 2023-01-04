package lpr

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	deviceUrl = "http://localhost:1234"
	username  = "username"
	password  = "password"
	timeout   = time.Second
)

func TestStartPullingRecognitions(t *testing.T) {
	device := NewDevice(deviceUrl, username, password, timeout)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
				return
			}

			if strings.Contains(string(body), "CreatePullPointSubscriptionRequest") {
				resBody := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><env:Envelope xmlns:env=\"http://www.w3.org/2003/05/soap-envelope\" xmlns:soapenc=\"http://www.w3.org/2003/05/soap-encoding\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\" xmlns:tt=\"http://www.onvif.org/ver10/schema\" xmlns:tds=\"http://www.onvif.org/ver10/device/wsdl\" xmlns:trt=\"http://www.onvif.org/ver10/media/wsdl\" xmlns:timg=\"http://www.onvif.org/ver20/imaging/wsdl\" xmlns:tev=\"http://www.onvif.org/ver10/events/wsdl\" xmlns:tptz=\"http://www.onvif.org/ver20/ptz/wsdl\" xmlns:tan=\"http://www.onvif.org/ver20/analytics/wsdl\" xmlns:tst=\"http://www.onvif.org/ver10/storage/wsdl\" xmlns:ter=\"http://www.onvif.org/ver10/error\" xmlns:dn=\"http://www.onvif.org/ver10/network/wsdl\" xmlns:tns1=\"http://www.onvif.org/ver10/topics\" xmlns:tmd=\"http://www.onvif.org/ver10/deviceIO/wsdl\" xmlns:wsdl=\"http://schemas.xmlsoap.org/wsdl\" xmlns:wsoap12=\"http://schemas.xmlsoap.org/wsdl/soap12\" xmlns:http=\"http://schemas.xmlsoap.org/wsdl/http\" xmlns:d=\"http://schemas.xmlsoap.org/ws/2005/04/discovery\" xmlns:wsadis=\"http://schemas.xmlsoap.org/ws/2004/08/addressing\" xmlns:wsnt=\"http://docs.oasis-open.org/wsn/b-2\" xmlns:wsa=\"http://www.w3.org/2005/08/addressing\" xmlns:wstop=\"http://docs.oasis-open.org/wsn/t-1\" xmlns:wsrf-bf=\"http://docs.oasis-open.org/wsrf/bf-2\" xmlns:wsntw=\"http://docs.oasis-open.org/wsn/bw-2\" xmlns:wsrf-rw=\"http://docs.oasis-open.org/wsrf/rw-2\" xmlns:wsaw=\"http://www.w3.org/2006/05/addressing/wsdl\" xmlns:wsrf-r=\"http://docs.oasis-open.org/wsrf/r-2\" xmlns:trc=\"http://www.onvif.org/ver10/recording/wsdl\" xmlns:tse=\"http://www.onvif.org/ver10/search/wsdl\" xmlns:trp=\"http://www.onvif.org/ver10/replay/wsdl\" xmlns:tnshik=\"http://www.hikvision.com/2011/event/topics\" xmlns:hikwsd=\"http://www.onvifext.com/onvif/ext/ver10/wsdl\" xmlns:hikxsd=\"http://www.onvifext.com/onvif/ext/ver10/schema\" xmlns:tas=\"http://www.onvif.org/ver10/advancedsecurity/wsdl\"><env:Header><wsa:Action>http://www.onvif.org/ver10/events/wsdl/EventPortType/CreatePullPointSubscriptionResponse</wsa:Action></env:Header><env:Body><tev:CreatePullPointSubscriptionResponse><tev:SubscriptionReference><wsa:Address>http://192.168.84.71/onvif/Events/PullSubManager_2022-12-30T18:43:58Z_6</wsa:Address></tev:SubscriptionReference><wsnt:CurrentTime>2022-12-30T18:43:58Z</wsnt:CurrentTime><wsnt:TerminationTime>2022-12-30T18:46:58Z</wsnt:TerminationTime></tev:CreatePullPointSubscriptionResponse></env:Body></env:Envelope>"
				fmt.Fprint(w, resBody)
			} else if strings.Contains(string(body), "PullMessagesRequest") {
				resBody := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><env:Envelope xmlns:env=\"http://www.w3.org/2003/05/soap-envelope\" xmlns:soapenc=\"http://www.w3.org/2003/05/soap-encoding\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\" xmlns:tt=\"http://www.onvif.org/ver10/schema\" xmlns:tds=\"http://www.onvif.org/ver10/device/wsdl\" xmlns:trt=\"http://www.onvif.org/ver10/media/wsdl\" xmlns:timg=\"http://www.onvif.org/ver20/imaging/wsdl\" xmlns:tev=\"http://www.onvif.org/ver10/events/wsdl\" xmlns:tptz=\"http://www.onvif.org/ver20/ptz/wsdl\" xmlns:tan=\"http://www.onvif.org/ver20/analytics/wsdl\" xmlns:tst=\"http://www.onvif.org/ver10/storage/wsdl\" xmlns:ter=\"http://www.onvif.org/ver10/error\" xmlns:dn=\"http://www.onvif.org/ver10/network/wsdl\" xmlns:tns1=\"http://www.onvif.org/ver10/topics\" xmlns:tmd=\"http://www.onvif.org/ver10/deviceIO/wsdl\" xmlns:wsdl=\"http://schemas.xmlsoap.org/wsdl\" xmlns:wsoap12=\"http://schemas.xmlsoap.org/wsdl/soap12\" xmlns:http=\"http://schemas.xmlsoap.org/wsdl/http\" xmlns:d=\"http://schemas.xmlsoap.org/ws/2005/04/discovery\" xmlns:wsadis=\"http://schemas.xmlsoap.org/ws/2004/08/addressing\" xmlns:wsnt=\"http://docs.oasis-open.org/wsn/b-2\" xmlns:wsa=\"http://www.w3.org/2005/08/addressing\" xmlns:wstop=\"http://docs.oasis-open.org/wsn/t-1\" xmlns:wsrf-bf=\"http://docs.oasis-open.org/wsrf/bf-2\" xmlns:wsntw=\"http://docs.oasis-open.org/wsn/bw-2\" xmlns:wsrf-rw=\"http://docs.oasis-open.org/wsrf/rw-2\" xmlns:wsaw=\"http://www.w3.org/2006/05/addressing/wsdl\" xmlns:wsrf-r=\"http://docs.oasis-open.org/wsrf/r-2\" xmlns:trc=\"http://www.onvif.org/ver10/recording/wsdl\" xmlns:tse=\"http://www.onvif.org/ver10/search/wsdl\" xmlns:trp=\"http://www.onvif.org/ver10/replay/wsdl\" xmlns:tnshik=\"http://www.hikvision.com/2011/event/topics\" xmlns:hikwsd=\"http://www.onvifext.com/onvif/ext/ver10/wsdl\" xmlns:hikxsd=\"http://www.onvifext.com/onvif/ext/ver10/schema\" xmlns:tas=\"http://www.onvif.org/ver10/advancedsecurity/wsdl\"><env:Header><wsa:Action>http://www.onvif.org/ver10/events/wsdl/PullPointSubscription/PullMessagesResponse</wsa:Action></env:Header><env:Body><tev:PullMessagesResponse><tev:CurrentTime>2022-12-30T15:18:27Z</tev:CurrentTime><tev:TerminationTime>2022-12-30T15:26:31Z</tev:TerminationTime><wsnt:NotificationMessage><wsnt:Topic Dialect=\"http://www.onvif.org/ver10/tev/topicExpression/ConcreteSet\">tns1:RuleEngine/VehicleDetector/tnshik:Vehicle</wsnt:Topic><wsnt:Message><tt:Message UtcTime=\"2022-12-30T15:18:27Z\" PropertyOperation=\"Changed\"><tt:Source><tt:SimpleItem Name=\"VideoSourceConfigurationToken\" Value=\"VideoSourceToken\"/><tt:SimpleItem Name=\"VideoAnalyticsConfigurationToken\" Value=\"VideoAnalyticsToken\"/><tt:SimpleItem Name=\"Rule\" Value=\"MyVehicleDetector\"/></tt:Source><tt:Data><tt:SimpleItem Name=\"PlateNumber\" Value=\"PP7069\"/><tt:SimpleItem Name=\"Likelihood\" Value=\"71\"/><tt:SimpleItem Name=\"Nation\" Value=\"EU\"/><tt:SimpleItem Name=\"Country\" Value=\"Netherlands\"/><tt:SimpleItem Name=\"VehicleLaneNumber\" Value=\"1\"/><tt:SimpleItem Name=\"VehicleDirection\" Value=\"forward\"/><tt:SimpleItem Name=\"PictureUri\" Value=\"http://192.168.84.71/doc/ui/images/plate/202212301618275900.jpg\"/></tt:Data></tt:Message></wsnt:Message></wsnt:NotificationMessage></tev:PullMessagesResponse></env:Body></env:Envelope>"
				fmt.Fprint(w, resBody)
			} else {
				t.Error("invalid request")
				return
			}
		})
		http.ListenAndServe(":1234", nil)
	}()

	channel := make(chan struct{})

	go device.StartPullingRecognitions(func(rec *Recognition, err error) {
		if err != nil {
			t.Error(err)
			return
		}

		if rec == nil {
			t.Error("recognition expected")
			return
		}

		channel <- struct{}{}
	})

	for range channel {
		return
	}
}

func TestGenerateCredentials(t *testing.T) {
	device := NewDevice(deviceUrl, username, password, timeout)
	token, nonce64, timestamp := device.generateCredentials()

	if token == "" {
		t.Error("token value expected")
	}

	if nonce64 == "" {
		t.Error("nonce value expected")
	}

	if timestamp == "" {
		t.Error("timestamp value expected")
	}
}

func BenchmarkCreatePullPointSubscriptionXml(b *testing.B) {
	device := NewDevice(deviceUrl, username, password, timeout)
	for i := 0; i < b.N; i++ {
		device.createPullPointSubscriptionXml()
	}
}

func BenchmarkPullMessagesXml(b *testing.B) {
	device := NewDevice(deviceUrl, username, password, timeout)
	for i := 0; i < b.N; i++ {
		device.pullMessagesXml("http://localhost:1234")
	}
}

func BenchmarkGenerateCredentials(b *testing.B) {
	device := NewDevice(deviceUrl, username, password, timeout)
	for i := 0; i < b.N; i++ {
		device.generateCredentials()
	}
}
