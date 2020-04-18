package models

import (
	orm "go-admin/database"
	"go-admin/tools"
)

type SysConfig struct {
	ConfigId    int    `json:"configId" gorm:"primary_key;auto_increment;"` //编码
	ConfigName  string `json:"configName" gorm:"type:varchar(128);"`        //参数名称
	ConfigKey   string `json:"configKey" gorm:"type:varchar(128);"`         //参数键名
	ConfigValue string `json:"configValue" gorm:"type:varchar(255);"`       //参数键值
	ConfigType  string `json:"configType" gorm:"type:varchar(64);"`         //是否系统内置
	Remark      string `json:"remark" gorm:"type:varchar(128);"`            //备注
	CreateBy    string `json:"createBy" gorm:"type:varchar(128);"`
	UpdateBy    string `json:"updateBy" gorm:"type:varchar(128);"`
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params"  gorm:"-"`
	BaseModel
}

func (SysConfig) TableName() string {
	return "sys_config"
}

// Config 创建
func (e *SysConfig) Create() (SysConfig, error) {
	var doc SysConfig
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取 Config
func (e *SysConfig) Get() (SysConfig, error) {
	var doc SysConfig

	table := orm.Eloquent.Table(e.TableName())
	if e.ConfigId != 0 {
		table = table.Where("config_id = ?", e.ConfigId)
	}

	if e.ConfigKey != "" {
		table = table.Where("config_key = ?", e.ConfigKey)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (e *SysConfig) GetPage(pageSize int, pageIndex int) ([]SysConfig, int, error) {
	var doc []SysConfig

	table := orm.Eloquent.Select("*").Table(e.TableName())

	if e.ConfigName != "" {
		table = table.Where("config_name = ?", e.ConfigName)
	}
	if e.ConfigKey != "" {
		table = table.Where("config_key = ?", e.ConfigKey)
	}
	if e.ConfigType != "" {
		table = table.Where("config_type = ?", e.ConfigType)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table = dataPermission.GetDataScope("sys_config", table)

	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

func (e *SysConfig) Update(id int) (update SysConfig, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("config_id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

func (e *SysConfig) Delete() (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("config_id = ?", e.ConfigId).Delete(&SysConfig{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}
