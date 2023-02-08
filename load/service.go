package load

import (
	"errors"
	"fmt"
	"github.com/daijulong/dockser/lib"
	"gopkg.in/yaml.v2"
)

// ServicesContent 服务配置内容（多个）
type ServicesContent map[string]*map[string]interface{}

//------------------------------------------------
//  Service
//------------------------------------------------

// Service 服务
type Service struct {
	Name     string                   //服务名称，对应文件名
	Services map[string]interface{}   //服务，一个文件可能有多个服务
	Add      []*AdditionalInstruction //添加服务时执行的附加指令集
	Remove   []*AdditionalInstruction //移除服务时执行的附加指令集
}

// NewService Service constructor
func NewService() *Service {
	return &Service{}
}

// Load 加载指定名称的 service
func (s *Service) Load(service string) error {
	file := "./compose/services/" + service + ".yml"
	fileContent, err := lib.ReadFile(file, service)
	if err != nil {
		return errors.New("read service[" + service + "] file [" + file + "] fail: " + err.Error())
	}
	serviceMap := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(fileContent), &serviceMap)
	if err != nil {
		return errors.New("parse service[" + service + "] file [" + file + "] fail: " + err.Error())
	}

	s.Name = service
	s.Services = make(map[string]interface{})
	for k, v := range serviceMap {
		// dockser 为指定附加命令的指令，要排除
		if k != "dockser" {
			s.Services[k] = v
		} else {
			//处理附加命令
			//@todo
			fmt.Println("deal service[" + service + "] 处理附加指定")
		}
	}
	return nil
}

// ApplyAddCommands 添加服务时执行附加指令集
func (s *Service) ApplyAddCommands() error {

	return nil
}

// ApplyRemoveCommands 移除服务时执行附加指令集
func (s *Service) ApplyRemoveCommands() error {

	return nil
}

//------------------------------------------------
//  Services
//------------------------------------------------

// Services 服务
type Services struct {
	Contents *ServicesContent
	Services []*Service
}

// NewServices Services constructor
func NewServices() *Services {
	return &Services{}
}

// Load 按名称加载服务设置
func (s *Services) Load(services []string) error {
	if len(services) < 1 {
		return nil
	}
	for _, serviceName := range services {
		service := NewService()
		err := service.Load(serviceName)
		if err != nil {
			return err
		}
		s.Services = append(s.Services, service)
	}
	return nil
}
