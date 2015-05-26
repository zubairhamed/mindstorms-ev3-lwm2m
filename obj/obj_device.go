package ev3

import (
	. "github.com/zubairhamed/go-lwm2m/api"
	"github.com/zubairhamed/go-lwm2m/core"
	. "github.com/zubairhamed/go-lwm2m/objects/oma"
	"github.com/zubairhamed/goap"
	"time"
	"os/exec"
	"log"
	"strconv"
)

type Device struct {
	Model ObjectModel
}

func (o *Device) OnExecute(instanceId int, resourceId int) goap.CoapCode {
	log.Println("Calling device onExecute", instanceId, resourceId)
	if resourceId == DEVICE_EXEC_REBOOT {
		// Wait 3 seconds before rebooting
		go func() {
			timer := time.NewTimer(time.Second * 3)
			<- timer.C

			o.Reboot()
		}()
	}
	return goap.COAPCODE_204_CHANGED
}

func (o *Device) OnCreate(instanceId int, resourceId int) goap.CoapCode {
	return goap.COAPCODE_201_CREATED
}

func (o *Device) OnDelete(instanceId int) goap.CoapCode {
	return goap.COAPCODE_202_DELETED
}

func (o *Device) OnRead(instanceId int, resourceId int) (ResponseValue, goap.CoapCode) {
	if resourceId == -1 {
		// Read Object Instance
	} else {
		// Read Resource Instance
		var val ResponseValue

		resource := o.Model.GetResource(resourceId)
		switch resourceId {
		case 0:
			val = core.NewStringValue(o.GetManufacturer())
			break

		case 1:
			val = core.NewStringValue(o.GetModelNumber())
			break

		case 2:
			val = core.NewStringValue(o.GetSerialNumber())
			break

		case 3:
			val = core.NewStringValue(o.GetFirmwareVersion())
			break

		case 6:
			val, _ = core.TlvPayloadFromIntResource(resource, o.GetAvailablePowerSources())
			break

		case 7:
			val, _ = core.TlvPayloadFromIntResource(resource, o.GetPowerSourceVoltage())
			break

		case 8:
			val, _ = core.TlvPayloadFromIntResource(resource, o.GetPowerSourceCurrent())
			break

		case 9:
			val = core.NewIntegerValue(o.GetBatteryLevel())
			break

		case 10:
			val = core.NewIntegerValue(o.GetMemoryFree())
			break

		case 11:
			val, _ = core.TlvPayloadFromIntResource(resource, o.GetErrorCode())
			break

		case 13:
			val = core.NewTimeValue(o.GetCurrentTime())
			break

		case 14:
			val = core.NewStringValue(o.GetTimezone())
			break

		case 15:
			val = core.NewStringValue(o.GetUtcOffset())
			break

		case 16:
			val = core.NewStringValue(o.GetSupportedBindingMode())
			break

		default:
			break
		}
		return val, goap.COAPCODE_205_CONTENT
	}
	return core.NewEmptyValue(), goap.COAPCODE_404_NOT_FOUND
}

func (o *Device) OnWrite(instanceId int, resourceId int) goap.CoapCode {
	return goap.COAPCODE_404_NOT_FOUND
}

func (o *Device) GetManufacturer() string {
	return "LEGO"
}

func (o *Device) GetModelNumber() string {
	return "Lego Mindstorms EV3"
}

func (o *Device) GetSerialNumber() string {
	return "12345"
}

func (o *Device) GetFirmwareVersion() string {
	return "1.0"
}

func (o *Device) Reboot() ResponseValue {
	// shutdown -r now
	err := exec.Command("shutdown", "-r", "now").Run()
	if err != nil {
		log.Println(err)
	}

	return core.NewEmptyValue()
}

func (o *Device) FactoryReset() ResponseValue {
	// unsupported
	return core.NewEmptyValue()
}

func (o *Device) GetAvailablePowerSources() []int {
	return []int{ POWERSOURCE_INTERNAL }
}

func (o *Device) GetPowerSourceVoltage() []int {
	out, err := exec.Command("cat", "/sys/class/power_supply/legoev3-battery/voltage_now").Output()

	if err != nil {
		log.Println(err)
	}
	s := string(out[:len(out)-1])
	i, e := strconv.Atoi(s)
	if e != nil {
		log.Println(e)
	}
	return []int{ int(i/1000) }
}

func (o *Device) GetPowerSourceCurrent() []int {
	// /sys/class/power_supply/legoev3-battery/current_now
	return []int{125, 900}
}

func (o *Device) GetBatteryLevel() int {
	// Unknown
	return 0
}

func (o *Device) GetMemoryFree() int {
	return 0
}

func (o *Device) GetErrorCode() []int {
	return []int{0}
}

func (o *Device) ResetErrorCode() string {
	return ""
}

func (o *Device) GetCurrentTime() time.Time {
	return time.Now()
}

func (o *Device) GetTimezone() string {
	return "+8:00"
}

func (o *Device) GetUtcOffset() string {
	return "+8:00"
}

func (o *Device) GetSupportedBindingMode() string {
	return "U"
}

func NewDeviceObject(reg Registry) *Device {
	return &Device{
		Model: reg.GetModel(OBJECT_LWM2M_DEVICE),
	}
}

