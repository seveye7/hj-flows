package main

// https://open.oceanengine.com/labels/7/docs/1696710655781900
type LinkToutiao struct {
	Id            uint64
	Channel       uint32
	AdvertiserId  uint64
	ProjectId     uint64
	PromotionId   uint64
	PromotionName string
	Idfa          string
	Imei          string
	Oaid          string
	Oaid2         string
	Androidid     string
	Os            int64
	Ts            int64
	CallbackUrl   string
	ClickId       string
	Ip            string
}
type Register struct {
	Date     string `validate:"required,len=10"`
	Uid      uint64
	Uuid     string
	CreateAt string `validate:"required,len=19"`
	Channel  uint32 `validate:"gt=0,lte=1000000"`
	Version  string
	Ip       string
	Location string
	Machine  string
	Network  uint32
	DeviceId string `validate:"required,min=6,max=64"`
	IDFA     string
	Tag      string // tag
	Ext1     string
	Ext2     string
	Ext3     string
	Ext4     string
	Ext5     string
	Ext6     string
	Ext7     string
	Ext8     string
	Ext9     string
	Ext10    string
}
