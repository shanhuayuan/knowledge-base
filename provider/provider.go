package provider

import (
	"app/global"
	"app/model"
)

// 模型增删改查

func CreateModel(model *model.Model) error {
	return global.OverallRuntime.DataStore.CreateEntry(model)
}

func DeleteModel(model *model.Model) error {
	return global.OverallRuntime.DataStore.DeleteEntry(model)
}

func GetModel(model *model.Model) error {
	return global.OverallRuntime.DataStore.TakeEntry(model)
}

func ListModels(models *[]model.Model) error {
	return global.OverallRuntime.DataStore.FindEntry(model.Model{}, models)
}

// 模型厂商增删改查

func CreateProvider(provider *model.Provider) error {
	return global.OverallRuntime.DataStore.CreateEntry(provider)
}

func DeleteProvider(provider *model.Provider) error {
	return global.OverallRuntime.DataStore.DeleteEntry(provider)
}

func ListProviders(providers *[]model.Provider) error {
	return global.OverallRuntime.DataStore.FindEntry(model.Provider{}, providers)
}

func GetProvider(provider *model.Provider) error {
	return global.OverallRuntime.DataStore.TakeEntry(provider)
}
