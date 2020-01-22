/*
 * Copyright (c) 2016 General Electric Company. All rights reserved.
 *
 * The copyright to the computer software herein is the property of
 * General Electric Company. The software may be used and/or copied only
 * with the written permission of General Electric Company or in accordance
 * with the terms and conditions stipulated in the agreement/contract
 * under which the software has been supplied.
 *
 * author: apolo.yasuda@ge.com
 */

package main

import (
	"os"
	"flag"
	"errors"
	//"net"
	"net/url"
	util "github.build.ge.com/212359746/wzutil"
	plugin "github.build.ge.com/212359746/wzplugin"
	"gopkg.in/yaml.v2"
	"encoding/base64"

)

var (
	REV string = "beta"
)

const (

	EC_LOGO = `
           ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄
          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░
          ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀
          ▐░▌          ▐░▌   
          ▐░█▄▄▄▄▄▄▄▄▄ ▐░▌
          ▐░░░░░░░░░░░▌▐░▌
          ▐░█▀▀▀▀▀▀▀▀▀ ▐░▌
          ▐░▌          ▐░▌
          ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄ 
          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
           ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  @Digital Connect 
`
	COPY_RIGHT = "Digital Connect,  @GE Corporate"
	ISSUE_TRACKER = "https://github.com/Enterprise-connect/sdk/issues"


	//agent authorization header
	AUTH_HEADER = "Authorization"
	
	//ec service header in predix
	EC_SUB_HEADER  = "Predix-Zone-Id"

	//app index available in a cf environment
	CF_INS_IDX_EV  = "CF_INSTANCE_INDEX"
	//forwarding header targeting a cf environment
	CF_INS_HEADER  = "X-CF-APP-INSTANCE"

	//app index available in a watcher environment
	EC_INS_IDX_EV  = "EC_INSTANCE_INDEX"
	//forwarding header targeting a watcher environement
	EC_INS_HEADER  = "X-EC-APP-INSTANCE"

	//xcalr url
	CA_URL = "https://xcalr.apps.ge.com/v2beta"

	//watcher url
	WATCHER_URL = "https://raw.githubusercontent.com/Enterprise-connect/sdk/v1.1beta.watcher/watcher.yml"

)

func GetTLSSetting()(map[string]interface{}, error){
	
	plg:=flag.String("plg","","Enable support for EC TLS Plugin.")
	ver:=flag.Bool("ver", false, "Show current tls revision.")

	flag.Parse()
	
	if *ver {
		util.InfoLog("Rev:"+REV)
		os.Exit(0)
		return nil,nil
	}

	//util.DbgLog(*plg)
	f, err := base64.StdEncoding.DecodeString(*plg)
	if err!=nil{
		return nil, err
	}
	t:=make(map[string]interface{})
	
	err=yaml.Unmarshal(f, &t)
	if err!=nil{
		return nil,err
	}

	if len(t)<1 {
		return nil, errors.New("invalid file format in plugins.yml")
	}
	
	return t,nil

}

func init(){
}

func main(){

	defer func(){
		if r:=recover();r!=nil{
			util.PanicRecovery(r)
		} else {
			util.InfoLog("plugin undeployed.")
		}
	}()

	bc:=&util.BrandingConfig{
		CONFIG_MAIN: "/.ec",
		BRAND_CONFIG: "EC", 
		LOGO: EC_LOGO,
		COPY_RIGHT: COPY_RIGHT,
		HEADER_PLUGIN: "ec-plugin",
		HEADER_CONFIG: "ec-config",
		HEADER_AUTH: AUTH_HEADER,
		HEADER_SUB_ID: EC_SUB_HEADER,
		HEADER_CF_INST: CF_INS_HEADER,
		HEADER_INST: EC_INS_HEADER,
		ENV_CF_INST_IDX: CF_INS_IDX_EV,
		ENV_INST_IDX: EC_INS_IDX_EV,
		URL_CA: CA_URL,
		URL_WATCHER_CONF: WATCHER_URL,
		//URL_WATCHER_REPO: WATCHER_CONTR_URL,
		URL_ISSUE_TRACKER: ISSUE_TRACKER,
	}
	
	util.Branding(bc)

	util.Init("tls",true)

	t,err:=GetTLSSetting()
	
	if err!=nil{
		panic(err)
	}
	
	util.DbgLog(t)
	
	//tcp resolve is irrelevant in this plugin, use http proxy instead
	//_, err= net.ResolveTCPAddr("tcp",t["hostname"].(string)+":"+t["tlsport"].(string))
	//if err != nil {
	//	panic(err)
	//}

	p:=plugin.NewProxy(REV)
	
	u, err := url.Parse(t["proxy"].(string))
	if err != nil {
		panic(err)
	}

	if err:=p.Init(t["schema"].(string)+"://"+t["hostname"].(string)+":"+t["tlsport"].(string),u);err!=nil{
		panic(err)
	}

	p.Start(t["port"].(string))
		
}
