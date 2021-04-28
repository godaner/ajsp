package main

import (
	"fmt"
	"regexp"
	"testing"
)

func TestParse(t *testing.T) {
	s := `
<!DOCTYPE html>
<html>
<head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=no"/>
        <link rel="stylesheet" type="text/css" href="/mobile/css/base.css"/>
        <script language="JavaScript" type="text/javascript" src="/mobile/js/rem_fix.js"></script>
        <title>安居区机构编制二维码在线履职监督评价管理系统</title>
        <style type="text/css">
                body{margin:0;padding:0;background:url(/mobile/images/error_bg.jpg) no-repeat;height: 100vh;background-size: cover;overflow: hidden;}
        </style>
</head>
<body>
        <div class="success_container" style="text-align: center;height: 100%;">
                <img class="png" style="margin-top:2rem;display: inline-block;" src="/mobile/images/exception.png">
                <p style="font-size: 0.4rem;color:#ffffff;">1分钟内不能重复评价！</p>
                <a href="/app/zxing/809f8a537caf460981234975658ce420" style="color:#ffffff;font-size:0.4rem;position: absolute;bottom: 0.1rem;left: calc(50% - 1rem);width: 2rem;">返回首页</a>
        </div>
</body>
</html>`
	reg, _ := regexp.Compile(`<p style="font-size: 0.[0-9]*rem;color:#ffffff;">([\w\W]+)</p>\n`)
	res := reg.FindSubmatch([]byte(s))
	for _, r := range res {
		fmt.Println(string(r))
	}
}
