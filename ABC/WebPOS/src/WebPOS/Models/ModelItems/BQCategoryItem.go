package ModelItems

type BQCategory struct {
	CategoryCD   string `sql:"CD"`
	CategoryName string `sql:"Name"`
	OffCheck     bool
}

type BQCategoryCDNameItem struct {
	CD   string `sql:"CD"`
	Name string `sql:"Name"`
}

type BQMediaItem struct {
	Type string `sql:"media_type"`
	CD   string `sql:"media_cd"`
	Name string `sql:"media_name"`
}

type MediaItem struct {
	MediaCd   string `sql:"media_cd"`
	MediaName string `sql:"media_name"`
	Product   string `sql:"product"`
}

//特大大中小分類
//特大分類
type MediaGroup1 struct {
	MediaGroup1Cd   string
	MediaGroup1Name string
	MediaGroup2     []*MediaGroup2
	JanGroup        string
}

//大分類
type MediaGroup2 struct {
	MediaGroup2Cd   string
	MediaGroup2Name string
	MediaGroup3     []*MediaGroup3
}

//中分類
type MediaGroup3 struct {
	MediaGroup3Cd   string
	MediaGroup3Name string
	MediaGroup4     []*MediaGroup4
}

//小分類
type MediaGroup4 struct {
	MediaGroup4Cd   string
	MediaGroup4Name string
}
