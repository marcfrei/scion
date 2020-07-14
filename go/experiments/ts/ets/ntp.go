package ets;

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/beevik/ntp"
)

var ntpLog = log.New(ioutil.Discard, "[ets/ntp] ", log.LstdFlags) 

func GetNTPClockOffset(host string) (time.Duration, error) {

	ntpLog.Printf("[%s] ----------------------", host)
	ntpLog.Printf("[%s] NTP protocol version %d", host, 4)

	r, err := ntp.Query(host)
	if err != nil {
		return 0, err
	}
	err = r.Validate();
	if err != nil {
		return 0, err
	}

	ntpLog.Printf("[%s]  LocalTime: %v\n", host, time.Now())
	ntpLog.Printf("[%s]   XmitTime: %v\n", host, r.Time)
	ntpLog.Printf("[%s]    RefTime: %v\n", host, r.ReferenceTime)
	ntpLog.Printf("[%s]        RTT: %v\n", host, r.RTT)
	ntpLog.Printf("[%s]     Offset: %v\n", host, r.ClockOffset)
	ntpLog.Printf("[%s]       Poll: %v\n", host, r.Poll)
	ntpLog.Printf("[%s]  Precision: %v\n", host, r.Precision)
	ntpLog.Printf("[%s]    Stratum: %v\n", host, r.Stratum)
	ntpLog.Printf("[%s]      RefID: 0x%08x\n", host, r.ReferenceID)
	ntpLog.Printf("[%s]  RootDelay: %v\n", host, r.RootDelay)
	ntpLog.Printf("[%s]   RootDisp: %v\n", host, r.RootDispersion)
	ntpLog.Printf("[%s]   RootDist: %v\n", host, r.RootDistance)
	ntpLog.Printf("[%s]   MinError: %v\n", host, r.MinError)
	ntpLog.Printf("[%s]       Leap: %v\n", host, r.Leap)
	ntpLog.Printf("[%s]   KissCode: \"%v\"\n", host, r.KissCode)

	return r.ClockOffset, nil
}