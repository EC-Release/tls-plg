/*
 * Copyright (c) 2016 General Electric Company. All rights reserved.
 *
 * The copyright to the computer software herein is the property of
 * General Electric Company. The software may be used and/or copied only
 * with the written permission of General Electric Company or in accordance
 * with the terms and conditions stipulated in the agreement/contract
 * under which the software has been supplied.
 *
 * author: chia.chang@ge.com
 */

package main

import (
	"os"
	"net"
	"net/url"
	util "github.build.ge.com/212359746/ecutil"
	plugin "github.build.ge.com/212359746/ecplugin"
)

var ()

const (
	YML_TLS_FLAG = "tls"
	REV = "v1"

)
func main(){
	defer func(){
		if r:=recover();r!=nil{
			util.ErrLog(r)
			os.Exit(1)
		}
	}()

	util.Init("TLS Plugin",true)

	p:=plugin.NewPlugin()
	util.DbgLog(p.Content)

	op:=p.Content

	op1:=op[YML_TLS_FLAG].([]interface{})
	op2:=op1[0].(map[interface{}]interface{})
	if op2["status"].(string)=="active"{
		
		_, err := net.ResolveTCPAddr("tcp",op2["hostname"].(string)+":"+op2["tlsport"].(string))
		if err != nil {
			panic(err)
		}

		p:=plugin.NewProxy(REV)

		u, err := url.Parse(op2["proxy"].(string))
		if err != nil {
			panic(err)
		}

		if err:=p.Init(op2["schema"].(string)+"://"+op2["hostname"].(string)+":"+op2["tlsport"].(string),u);err!=nil{
			panic(err)
		}

		p.Start(op2["port"].(string))
	}
	

}
