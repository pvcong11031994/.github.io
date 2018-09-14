package ModelItems

type SystemStatusItem struct {
	CreatedAt		string `sql:"ss_created_at"`
	Chain			string `sql:"ss_chain"`
	Group			string `sql:"ss_group"`
	Detail			string `sql:"ss_detail"`
}