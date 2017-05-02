package operations

import (
	"errors"
	"vmkube/model"
)

func DescribeInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Instance Status
	return errors.New("Command Instance Status not implemented!!")
}

func StartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Start Instance
	return errors.New("Command Start Instance not implemented!!")
}

func RestartInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement ReStart Instance
	return errors.New("Command Restart Instance not implemented!!")
}

func StopInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Stop Instance
	return errors.New("Command Stop Instance not implemented!!")
}

func DisableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Disable Instance ??
	return errors.New("Command Disable Instance not implemented!!")
}

func EnableInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Enable Instance ??
	return errors.New("Command Enable Instance not implemented!!")
}

func RecreateInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Re-Create Instance
	return errors.New("Command Re-Create Instance not implemented!!")
}

func DestroyInstance(infrastructure model.Infrastructure, instance model.LocalInstance, cloudInstance model.CloudInstance, isCloud bool) error {
	//TODO: Implement Destroy Instance
	return errors.New("Command Destroy Instance not implemented!!")
}

func AutoFixInfrastructureInstances(infrastructure model.Infrastructure) error {
	return nil
}
