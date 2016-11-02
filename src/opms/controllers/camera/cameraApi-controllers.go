package camera

import (
	"encoding/json"
	"fmt"
	"opms/controllers"
	. "opms/models/camera"
)

type CameraListForAPP struct {
	Id       int64
	Name     string
	Address1 string
	Address2 string
	Playurl  string
}

type CameraAPiController struct {
	controllers.BaseController
}

func (this *CameraAPiController) Get() {
	deps := GetCameras()
	var cameraForApp []CameraListForAPP
	for _, value := range deps {
		var c = CameraListForAPP{value.Id, value.Name, value.Address1, value.Address2, value.Playurl}
		cameraForApp = append(cameraForApp, c)
	}
	b, err := json.Marshal(cameraForApp)
	if err != nil {
		fmt.Println("json err:", err)
	}
	this.Ctx.WriteString(string(b))
}
