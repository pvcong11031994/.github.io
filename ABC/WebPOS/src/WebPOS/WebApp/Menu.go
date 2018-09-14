package WebApp

import (
	"WebPOS/Models/ModelItems"
	"bytes"
	"fmt"
	"github.com/goframework/gf"
	"io"
	"strings"
)

type Menu struct {
	Path     string //URL path
	Name     string //View menu name
	SubMenu  MenuList
	Level    int
	ParentID string
}

type MenuList []*Menu

func (menu *Menu) write(w io.Writer) {

	if menu.SubMenu == nil || len(menu.SubMenu) == 0 {
		// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT START
		//if strings.EqualFold(menu.Path, "/download/makersalestock/shop_list_download") {
		//	fmt.Fprintf(w, `<a style="background-image: url(/static/img/icon/download_white.png);
    		//								background-repeat: no-repeat;
    		//								background-position: 4px 15px;
    		//								padding-left: 26px; "href="%s">
    		//					%s</a>`, menu.Path, menu.Name)
		//}
		if strings.Contains(menu.Path, "/download/makersalestock/shop_list_download") {
			fmt.Fprintf(w, `<a style="background-image: url(/static/img/icon/download_white.png);
    										background-repeat: no-repeat;
    										background-position: 4px 15px;
    										padding-left: 26px; "href="%s">
    							%s</a>`, menu.Path, menu.Name)
			// ASO-5929 [BA]mBAWEB-v09f 店舗一覧ダウンロード-複数ファイル対応 - EDIT START
		} else {
			// ASO-5721 メニューURLが「http」で始まる場合の対応 - EDIT START
			//fmt.Fprintf(w, `<a href="%s">%s</a>`, menu.Path, menu.Name)
			if len(menu.Path) > len(CONSTANT_URL_HTTP) {
				if strings.Compare(menu.Path[0:len(CONSTANT_URL_HTTP)], CONSTANT_URL_HTTP) == 0 {
					fmt.Fprintf(w, `<a href="%s" target="_blank">%s</a>`, menu.Path, menu.Name)
				} else {
					fmt.Fprintf(w, `<a href="%s">%s</a>`, menu.Path, menu.Name)
				}
			} else {
				fmt.Fprintf(w, `<a href="%s">%s</a>`, menu.Path, menu.Name)
			}
			// ASO-5721 メニューURLが「http」で始まる場合の対応 - EDIT END
		}

	} else {
		fmt.Fprintf(w, `<a>%s</a>`, menu.Name)
		menu.SubMenu.write(w)
	}
}

func (menuList MenuList) write(w io.Writer) {

	if menuList == nil {
		return
	}
	fmt.Fprint(w, "<ul>")
	for i := range menuList {
		if menuList[i].SubMenu == nil || len(menuList[i].SubMenu) == 0 {
			fmt.Fprint(w, "<li>")
		} else {
			fmt.Fprint(w, `<li class="expand">`)
		}
		menuList[i].write(w)

		fmt.Fprint(w, "</li>")
	}
	fmt.Fprint(w, "</ul>")
}

func (menuList MenuList) ToString() string {
	buf := bytes.NewBuffer(make([]byte, 8192))
	buf.Reset()
	menuList.write(buf)

	return buf.String()
}

func GetEnableMenu(enablePathMap map[string]bool, listMenuLevel []ModelItems.MenuLevelItem, ctx *gf.Context) MenuList {

	//var iMenuLevel MenuList

	iMenuLevel := GetMenuLevel(listMenuLevel)
	//GetMenuLevel(listMenuLevel,&iMenuLevel,"",1)
	return iMenuLevel
}

func GetMenuLevel(listMenuLevel []ModelItems.MenuLevelItem) MenuList {

	mapMenu := map[string]*Menu{
		"": &Menu{
			SubMenu: MenuList{},
			Level:   0,
		},
	}

	for _, v := range listMenuLevel {
		newMenu := Menu{}
		newMenu.Path = "/" + v.MenuUrl
		// ASO-5721 メニューURLが「http」で始まる場合の対応 - ADD START
		if len(v.MenuUrl) > len(CONSTANT_URL_HTTP) {
			if strings.Compare(v.MenuUrl[0:len(CONSTANT_URL_HTTP)], CONSTANT_URL_HTTP) == 0 {
				newMenu.Path = v.MenuUrl
			}
		}
		// ASO-5721 メニューURLが「http」で始まる場合の対応 - ADD END
		newMenu.Name = v.MenuName
		newMenu.ParentID = v.ParentMenuId
		newMenu.Level = v.MenuLevel
		newMenu.SubMenu = MenuList{}
		mapMenu[v.MenuId] = &newMenu
		parentId := v.ParentMenuId
		parentMenu := mapMenu[parentId]
		if parentMenu != nil && v.MenuLevel == parentMenu.Level+1 {
			parentMenu.SubMenu = append(parentMenu.SubMenu, &newMenu)
		}
	}

	return mapMenu[""].SubMenu
}
