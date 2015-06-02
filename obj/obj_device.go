package ev3

import (
	. "github.com/zubairhamed/betwixt/api"
	"github.com/zubairhamed/betwixt/core"
	"github.com/zubairhamed/betwixt/core/response"
	"github.com/zubairhamed/betwixt/core/values"
	. "github.com/zubairhamed/betwixt/objects/oma"
	"log"
	"os/exec"
	"strconv"
	"time"
)

type Device struct {
	Model ObjectModel
}

func (o *Device) OnExecute(instanceId int, resourceId int, req Request) Response {
	log.Println("Calling device onExecute", instanceId, resourceId)
	if resourceId == DEVICE_EXEC_REBOOT {
		// Wait 3 seconds before rebooting
		go func() {
			timer := time.NewTimer(time.Second * 3)
			<-timer.C

			o.Reboot()
		}()
	}
	return response.Changed()
}

func (o *Device) OnCreate(instanceId int, resourceId int, req Request) Response {
	return response.Created()
}

func (o *Device) OnDelete(instanceId int, req Request) Response {
	return response.Deleted()
}

func (o *Device) OnRead(instanceId int, resourceId int, req Request) Response {
	if resourceId == -1 {
		// Read Object Instance
	} else {
		// Read Resource Instance
		var val ResponseValue

		resource := o.Model.GetResource(resourceId)
		switch resourceId {
		case 0:
			val = values.String(o.GetManufacturer())
			break

		case 1:
			val = values.String(o.GetModelNumber())
			break

		case 2:
			val = values.String(o.GetSerialNumber())
			break

		case 3:
			val = values.String(o.GetFirmwareVersion())
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
			val = values.Integer(o.GetBatteryLevel())
			break

		case 10:
			val = values.Integer(o.GetMemoryFree())
			break

		case 11:
			val, _ = core.TlvPayloadFromIntResource(resource, o.GetErrorCode())
			break

		case 13:
			val = values.Time(o.GetCurrentTime())
			break

		case 14:
			val = values.String(o.GetTimezone())
			break

		case 15:
			val = values.String(o.GetUtcOffset())
			break

		case 16:
			val = values.String(o.GetSupportedBindingMode())
			break

		default:
			break
		}
		return response.Content(val)
	}
	return response.NotFound()
}

func (o *Device) OnWrite(instanceId int, resourceId int, req Request) Response {
	return response.NotFound()
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

	return values.Empty()
}

func (o *Device) FactoryReset() ResponseValue {
	// unsupported
	return values.Empty()
}

func (o *Device) GetAvailablePowerSources() []int {
	return []int{POWERSOURCE_INTERNAL}
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
	return []int{int(i / 1000)}
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
