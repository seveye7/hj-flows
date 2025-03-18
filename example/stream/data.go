package main

// PlanDaily 计划报表
type PlanDaily struct {
	Date     string `db:"date" json:"date,omitempty"`
	Adverid  uint64 `db:"adverid" json:"adverid,omitempty"`
	Aid      uint64 `db:"aid" json:"aid,omitempty"`
	Platform uint32 `db:"platform" json:"platform,omitempty"`
	Channel  uint32 `db:"channel" json:"channel,omitempty"`

	State        *string  `db:"state" json:"state,omitempty"`
	ActiveTarget *string  `db:"active_target" json:"active_target,omitempty"`
	ActionTarget *string  `db:"action_target" json:"action_target,omitempty"`
	Name         *string  `db:"name" json:"name,omitempty"`
	Bid          *float64 `db:"bid" json:"bid,omitempty"`
	RegisterNum  *int64   `db:"register_num" json:"register_num,omitempty"`
	Cost         *float64 `db:"cost" json:"cost,omitempty"`
	Convert      *int64   `db:"convert" json:"convert,omitempty"`
	Ad1Num       *int64   `db:"ad1_num" json:"ad1_num,omitempty"`
	Ad2Num       *int64   `db:"ad2_num" json:"ad2_num,omitempty"`
	Ad3Num       *int64   `db:"ad3_num" json:"ad3_num,omitempty"`
	Ad4Num       *int64   `db:"ad4_num" json:"ad4_num,omitempty"`
	Ad5Num       *int64   `db:"ad5_num" json:"ad5_num,omitempty"`
	Ad6Num       *int64   `db:"ad6_num" json:"ad6_num,omitempty"`
	Ad7Num       *int64   `db:"ad7_num" json:"ad7_num,omitempty"`
	Ad8Num       *int64   `db:"ad8_num" json:"ad8_num,omitempty"`
	Ad1Ecpm      *float64 `db:"ad1_ecpm" json:"ad1_ecpm,omitempty"`
	Ad2Ecpm      *float64 `db:"ad2_ecpm" json:"ad2_ecpm,omitempty"`
	Ad3Ecpm      *float64 `db:"ad3_ecpm" json:"ad3_ecpm,omitempty"`
	Ad4Ecpm      *float64 `db:"ad4_ecpm" json:"ad4_ecpm,omitempty"`
	Ad5Ecpm      *float64 `db:"ad5_ecpm" json:"ad5_ecpm,omitempty"`
	Ad6Ecpm      *float64 `db:"ad6_ecpm" json:"ad6_ecpm,omitempty"`
	Ad7Ecpm      *float64 `db:"ad7_ecpm" json:"ad7_ecpm,omitempty"`
	Ad8Ecpm      *float64 `db:"ad8_ecpm" json:"ad8_ecpm,omitempty"`
	Red          *int64   `db:"red" json:"red,omitempty"`

	Day2RetentionNum *int64 `db:"day2_retention_num" json:"day2_retention_num,omitempty"`
	Day3RetentionNum *int64 `db:"day3_retention_num" json:"day3_retention_num,omitempty"`
	Day4RetentionNum *int64 `db:"day4_retention_num" json:"day4_retention_num,omitempty"`
	Day5RetentionNum *int64 `db:"day5_retention_num" json:"day5_retention_num,omitempty"`
	Day6RetentionNum *int64 `db:"day6_retention_num" json:"day6_retention_num,omitempty"`
	Day7RetentionNum *int64 `db:"day7_retention_num" json:"day7_retention_num,omitempty"`

	EcpmTotal     *float64 `db:"ecpm_total" json:"ecpm_total,omitempty"`
	Day2EcpmTotal *float64 `db:"day2_ecpm_total" json:"day2_ecpm_total,omitempty"`
	Day3EcpmTotal *float64 `db:"day3_ecpm_total" json:"day3_ecpm_total,omitempty"`
	Day4EcpmTotal *float64 `db:"day4_ecpm_total" json:"day4_ecpm_total,omitempty"`
	Day5EcpmTotal *float64 `db:"day5_ecpm_total" json:"day5_ecpm_total,omitempty"`
	Day6EcpmTotal *float64 `db:"day6_ecpm_total" json:"day6_ecpm_total,omitempty"`
	Day7EcpmTotal *float64 `db:"day7_ecpm_total" json:"day7_ecpm_total,omitempty"`

	RedTotal     *int64 `db:"red_total" json:"red_total,omitempty"`
	Day2RedTotal *int64 `db:"day2_red_total" json:"day2_red_total,omitempty"`
	Day3RedTotal *int64 `db:"day3_red_total" json:"day3_red_total,omitempty"`
	Day4RedTotal *int64 `db:"day4_red_total" json:"day4_red_total,omitempty"`
	Day5RedTotal *int64 `db:"day5_red_total" json:"day5_red_total,omitempty"`
	Day6RedTotal *int64 `db:"day6_red_total" json:"day6_red_total,omitempty"`
	Day7RedTotal *int64 `db:"day7_red_total" json:"day7_red_total,omitempty"`
}
type Client struct {
	Date        string `validate:"required,len=10"`
	Uid         uint64 // 用户ID
	DeviceId    string `validate:"required,min=6,max=64"`
	CreateAt    string `validate:"required,len=19"`
	Channel     uint32 `validate:"gt=0,lte=1000000"`
	Version     string
	RegisterDay uint32
	Clienttime  string `validate:"required,len=19"`
	Autoid      uint32 `validate:"gt=0"`
	Eventid     uint32 `validate:"gt=0"`
	Ext1        string
	Ext2        string
	Ext3        string
	Ext4        string
	Ext5        string
	Ext6        string
	Ext7        string
	Ext8        string
	Ext9        string
	Ext10       string
}

type UserDaily struct {
	Date string `db:"date"`
	Uid  uint64 `db:"uid"`

	Days          *uint32  `db:"days"`
	OnlineTimes   *int64   `db:"online_times"`
	AdVideoNum    *uint32  `db:"ad_video_num"`
	AdVideoAmount *float64 `db:"ad_video_amount"`
	AdAmountTotal *float64 `db:"ad_amount_total"`
	Red           *int64   `db:"red"`
	GameNum       *uint32  `db:"game_num"`
}
