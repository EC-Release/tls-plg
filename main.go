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
	util "github.com/wzlib/wzutil"
	plugin "github.com/wzlib/wzplugin"
	model "github.com/wzlib/wzschema"
	"gopkg.in/yaml.v2"
	"encoding/base64"

)

var (
	REV string = "v1.2beta"
	log *util.AppLog
)

const (
	EC_LOGO = `
           ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄
          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
          ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀▀▀
          ▐░▌          ▐░▌   
          ▐░█▄▄▄▄▄▄▄▄▄ ▐░▌
          ▐░░░░░░░░░░░▌▐░▌
          ▐░█▀▀▀▀▀▀▀▀▀ ▐░▌
          ▐░▌          ▐░▌
          ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄▄▄ 
          ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌
           ▀▀▀▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀▀▀▀  @Enterprise-Connect 
`
	COPY_RIGHT = "Enterprise-Connect,  @General Electric"
	ISSUE_TRACKER = "https://github.com/EC-Release/sdk/issues"

	AUTH_HEADER = "Authorization"

	EC_SUB_HEADER  = "Predix-Zone-Id"

	CF_INS_IDX_EV  = "CF_INSTANCE_INDEX"
	CF_INS_HEADER  = "X-CF-APP-INSTANCE"
	EC_INS_IDX_EV  = "EC_INSTANCE_INDEX"
	EC_INS_HEADER  = "X-EC-APP-INSTANCE"

	CA_URL = "https://github.com/EC-Release/certifactory"
)


func init(){
	bc:=&model.BrandingConfig{
		CONFIG_MAIN: "/.ec",
		BRAND_CONFIG: "EC",
		PASSPHRASE_EXT: "PPS",
		ART_NAME: "agent",
		LOGO: EC_LOGO,
		COPY_RIGHT: COPY_RIGHT,
		HEADER_PLUGIN: "ec-plugin",
		HEADER_CONFIG: "ec-config",
		STREAM_PATH: "/agent",
		HEADER_AUTH: AUTH_HEADER,
		HEADER_SUB_ID: EC_SUB_HEADER,
		HEADER_CF_INST: CF_INS_HEADER,
		HEADER_INST: EC_INS_HEADER,
		ENV_CF_INST_IDX: CF_INS_IDX_EV,
		ENV_INST_IDX: EC_INS_IDX_EV,
		URL_CA: CA_URL,
		URL_ISSUE_TRACKER: ISSUE_TRACKER,
	}
	
	util.Branding(bc)
	log = util.NewAppLog("tls")
}

func GetTLSSetting()(map[string]interface{}, error){
	
	plg:=flag.String("plg","","Enable support for EC TLS Plugin.")
	ver:=flag.Bool("ver", false, "Show current tls revision.")

	flag.Parse()
	
	if *ver {
		log.InfoLog("Rev:"+REV)
		os.Exit(0)
		return nil,nil
	}

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

func main(){

	defer func(){
		if r:=recover();r!=nil{
			util.PanicRecovery(r)
		} else {
			log.InfoLog("plugin undeployed.")
		}
	}()

	t,err:=GetTLSSetting()
	
	if err!=nil{
		panic(err)
	}
	
	log.DbgLog(t)
	
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
