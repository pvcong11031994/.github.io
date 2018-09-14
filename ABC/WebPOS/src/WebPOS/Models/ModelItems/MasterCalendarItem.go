package ModelItems

type MasterCalendarItem struct {
	Mcyyymmdd  string `sql:"mc_yyyymmdd"`
	Mcyyyy     string `sql:"mc_yyyy"`
	Mcmm       string `sql:"mc_mm"`
	Mcdd       string `sql:"mc_dd"`
	Mcdow      string `sql:"mc_dow"`
	Mcweeknum  string `sql:"mc_weeknum"`
	Mcweekdate string `sql:"mc_weekdate"`
	McKey string
}
