package api

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, code, status, result interface{}) {
	if status != result {
		t.Errorf("Expected %v (type %v)---> %v (type %v) - Got %v (type %v)",
			code, reflect.TypeOf(code),
			status, reflect.TypeOf(status),
			result, reflect.TypeOf(result))
	} else {
		t.Logf("Test ok! %v ---> %v", code, result)
	}
}

type code2status struct {
	code, status int
}

var c2sArray []code2status = []code2status{
	{200, 200},
	{12000, 200},
	{1200, 200},
	{2200, 200},
	{3200, 200},
	{220022, 200},
	{140022, 400},
	{140422, 404},
	{1200, 200},
	{1400, 400},
	{14003, 400},
	{14004, 400},
	{1403, 403},
	{14030, 403},
	{1404, 404},
	{14040, 404},
	{1405, 405},
	{1503, 503},
	{140010, 400},
}

func TestCode2Status(t *testing.T) {

	for _, c2s := range c2sArray {
		expect(t, c2s.code, c2s.status, trickCode2Status(c2s.code))
	}

}
