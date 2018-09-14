package ApiInit

import (
	"WebPOS/ControllersApi/Utils"
	"bytes"
	"encoding/json"
	"github.com/goframework/gf"
	"net"
	"strings"
)

const (
	MSG_ERR_IP_NOT_ALLOW = "このIPアドレスでアクセスできません。"
)

const (
	CFG_KEY_API_ENABLE_FILTER_IP = "API.EnableFilterIP"
	CFG_KEY_API_IP_ALLOW_FROM    = "API.IpAllowFrom"
	CFG_KEY_API_IP_ALLOW_TO      = "API.IpAllowTo"
	CFG_KEY_API_IP_ALLOW_SUBNET  = "API.IpAllowSubnet"
	CFG_KEY_API_IP_ALLOW_LIST    = "API.IpAllows"
)

var (
	_IpListDefault = []string{
		"::1",
	}
	_EnableFilterIP = false
	_IpAllowFrom    = net.IP{}
	_IpAllowTo      = net.IP{}
	_IpAllowSubnet  = &net.IPNet{}
	_IpAllowList    = []string{}
)

func FilterIp() {
	//
	loadIpAllowCfg()
	//
	gf.Filter(ApiUtils.ROUTE_API_FILTER, filterIp)
}

func loadIpAllowCfg() {
	cfg := ApiUtils.GetConfig()
	_EnableFilterIP = cfg.BoolOrFalse(CFG_KEY_API_ENABLE_FILTER_IP)
	_IpAllowFrom = net.ParseIP(cfg.StrOrEmpty(CFG_KEY_API_IP_ALLOW_FROM))
	_IpAllowTo = net.ParseIP(cfg.StrOrEmpty(CFG_KEY_API_IP_ALLOW_TO))
	_, _IpAllowSubnet, _ = net.ParseCIDR(cfg.StrOrEmpty(CFG_KEY_API_IP_ALLOW_SUBNET))
	_IpAllowList = cfg.List(CFG_KEY_API_IP_ALLOW_LIST)
}

func filterIp(ctx *gf.Context) {
	//
	ctx.ViewBases = nil
	//
	if !_EnableFilterIP ||
		((_IpAllowFrom == nil || _IpAllowTo == nil) && _IpAllowSubnet == nil && strings.TrimSpace(strings.Join(_IpAllowList, "")) == "") {
		return
	}
	//
	requestIP := net.ParseIP(ctx.GetRequestIP())
	// Check by range
	if _IpAllowFrom != nil && _IpAllowTo != nil {
		if 0 <= bytes.Compare(requestIP, _IpAllowFrom) && bytes.Compare(requestIP, _IpAllowTo) <= 0 {
			return
		}
	}
	// Check by subnet
	if _IpAllowSubnet != nil {
		if _IpAllowSubnet.Contains(requestIP) {
			return
		}
	}
	// Check by list defined
	if len(_IpAllowList) > 0 {
		for _, ip := range _IpAllowList {
			if bytes.Compare(requestIP, net.ParseIP(ip)) == 0 {
				return
			}
		}
	}
	// Check by default
	if strings.TrimSpace(strings.Join(_IpAllowList, "")) != "" {
		for _, ip := range _IpListDefault {
			if bytes.Compare(requestIP, net.ParseIP(ip)) == 0 {
				return
			}
		}
	}
	// Force
	ctx.ForcedStatus = true
	responseByte, _ := json.Marshal(map[string]string{
		"resultCode":    ApiUtils.PROCESS_ERROR_APP,
		"resultMessage": MSG_ERR_IP_NOT_ALLOW,
	})
	//
	ctx.JsonPResponse = ApiUtils.ResponseWithCallbackEmbed("callback", responseByte)
}
