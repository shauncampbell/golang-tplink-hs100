// Package hs100 is a public facing package for communicating with hs1xx devices.
package hs100

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const turnOnCommand = `{"system":{"set_relay_state":{"state":1}}}`
const turnOffCommand = `{"system":{"set_relay_state":{"state":0}}}`
const isOnCommand = `{"system":{"get_sysinfo":{}}}`
const currentPowerConsumptionCommand = `{"emeter":{"get_realtime":{},"get_vgain_igain":{}}}`

// Hs100 is a struct representing a particular hs1xx device.
type Hs100 struct {
	Address       string
	commandSender CommandSender
}

// NewHs100 creates a new Hs100 object.
func NewHs100(address string, s CommandSender) *Hs100 {
	return &Hs100{
		Address:       address,
		commandSender: s,
	}
}

// CommandSender is an interface which sends commands to a specific hs1xx device at the specified address.
type CommandSender interface {
	SendCommand(address string, command string) (string, error)
}

// TurnOn turns on the hs1xx device.
func (hs100 *Hs100) TurnOn() error {
	resp, err := hs100.commandSender.SendCommand(hs100.Address, turnOnCommand)
	if err != nil {
		return errors.Wrap(err, "error on sending turn on command for device")
	}

	r, err := parseSetRelayResponse(resp)
	if err != nil {
		return errors.Wrap(err, "Could not parse SystemInformationResponse from device")
	} else if r.errorOccurred() {
		return errors.New("got non zero exit code from device")
	}

	return nil
}

type setRelayResponse struct {
	System struct {
		SetRelayState struct {
			ErrorCode int `json:"err_code"`
		} `json:"set_relay_state"`
	} `json:"system"`
}

func parseSetRelayResponse(response string) (setRelayResponse, error) {
	var result setRelayResponse
	err := json.Unmarshal([]byte(response), &result)
	return result, err
}

func (r *setRelayResponse) errorOccurred() bool {
	return r.System.SetRelayState.ErrorCode != 0
}

// TurnOff turns off an hs1xx device
func (hs100 *Hs100) TurnOff() error {
	resp, err := hs100.commandSender.SendCommand(hs100.Address, turnOffCommand)
	if err != nil {
		return errors.Wrap(err, "error on sending turn on command for device")
	}

	r, err := parseSetRelayResponse(resp)
	if err != nil {
		return errors.Wrap(err, "Could not parse SystemInformationResponse from device")
	} else if r.errorOccurred() {
		return errors.New("got non zero exit code from device")
	}

	return nil
}

// IsOn checks if an hs1xx device is on.
func (hs100 *Hs100) IsOn() (bool, error) {
	resp, err := hs100.commandSender.SendCommand(hs100.Address, isOnCommand)
	if err != nil {
		return false, err
	}

	on, err := isOn(resp)
	if err != nil {
		return false, err
	}

	return on, nil
}

// SystemInformationResponse contains information about the hs1xx device.
type SystemInformationResponse struct {
	System struct {
		SystemInfo struct {
			SoftwareVersion string `json:"sw_ver"`
			HardwareVersion string `json:"hw_ver"`
			Model           string `json:"model"`
			DeviceID        string `json:"deviceId"`
			OemID           string `json:"oemId"`
			HardwareID      string `json:"hwId"`
			MACAddress      string `json:"mac"`
			RelayState      int    `json:"relay_state"`
			Alias           string `json:"alias"`
		} `json:"get_sysinfo"`
	} `json:"system"`
}

// GetInfo gets information about the device.
func (hs100 *Hs100) GetInfo() (*SystemInformationResponse, error) {
	resp, err := hs100.commandSender.SendCommand(hs100.Address, isOnCommand)
	if err != nil {
		return nil, err
	}

	var r SystemInformationResponse
	err = json.Unmarshal([]byte(resp), &r)
	return &r, err
}

func isOn(s string) (bool, error) {
	var r SystemInformationResponse
	err := json.Unmarshal([]byte(s), &r)
	on := r.System.SystemInfo.RelayState == 1
	return on, err
}

// GetName returns the alias of the device.
func (hs100 *Hs100) GetName() (string, error) {
	resp, err := hs100.commandSender.SendCommand(hs100.Address, isOnCommand)

	if err != nil {
		return "", err
	}

	name, err := name(resp)
	if err != nil {
		return "", err
	}

	return name, nil
}

func name(resp string) (string, error) {
	var r SystemInformationResponse
	err := json.Unmarshal([]byte(resp), &r)
	name := r.System.SystemInfo.Alias
	return name, err
}

// GetCurrentPowerConsumption returns the current power consumption available for the device.
func (hs100 *Hs100) GetCurrentPowerConsumption() (PowerConsumption, error) {
	resp, err := hs100.commandSender.SendCommand(hs100.Address, currentPowerConsumptionCommand)
	if err != nil {
		return PowerConsumption{}, errors.Wrap(err, "Could not read from hs100 device")
	}
	return powerConsumption(resp)
}

// PowerConsumption includes the current power consumption for the device.
type PowerConsumption struct {
	Current float32
	Voltage float32
	Power   float32
}

func powerConsumption(resp string) (PowerConsumption, error) {
	var r powerConsumptionResponse
	err := json.Unmarshal([]byte(resp), &r)
	if err != nil {
		return PowerConsumption{}, errors.Wrap(err, "Cannot parse SystemInformationResponse")
	}

	if r.Emeter.ErrorCode != 0 && r.Emeter.ErrorMessage != "" {
		return PowerConsumption{}, fmt.Errorf("error %d: %s", r.Emeter.ErrorCode, r.Emeter.ErrorMessage)
	}
	return r.toPowerConsumption(), nil
}

type powerConsumptionResponse struct {
	Emeter struct {
		ErrorCode    int    `json:"err_code"`
		ErrorMessage string `json:"err_msg"`
		RealTime     struct {
			Current float32 `json:"current"`
			Voltage float32 `json:"voltage"`
			Power   float32 `json:"power"`
		} `json:"get_realtime"`
	} `json:"emeter"`
}

func (r *powerConsumptionResponse) toPowerConsumption() PowerConsumption {
	return PowerConsumption{
		Current: r.Emeter.RealTime.Current,
		Voltage: r.Emeter.RealTime.Voltage,
		Power:   r.Emeter.RealTime.Power,
	}
}
