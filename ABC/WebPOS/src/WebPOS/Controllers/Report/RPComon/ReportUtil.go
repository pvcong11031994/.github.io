package RPComon

import "strings"

// サーバ名と店舗名を分ける
func ParseServerList(serverShops []string) []string {

	listServerName := []string{}
	mapServerName := map[string]bool{}
	for _, sh := range serverShops {
		//serverName := strings.Split(sh, "|")[0]
		serverName := ""
		arrsh := strings.Split(sh, "|")
		if len(arrsh) > 0 {
			serverName = arrsh[0]
		}
		if !mapServerName[serverName] && strings.Compare(serverName, "") != 0 {
			mapServerName[serverName] = true
			listServerName = append(listServerName, serverName)
		}
	}
	return listServerName
}

func ConvertJanMakerCode(code string) (string,bool){

	if strings.TrimSpace(code) == "" {
		return "",true
	}
	if strings.HasPrefix(code, "9784") {
		return code + "%",true
	} else if !strings.HasPrefix(code, "9784") && len(code) <= 6 {
		return "9784" + code + "%",true
	} else {
		return "",false
	}
}