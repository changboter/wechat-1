// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package js

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"strconv"
)

// js api 微信支付接口 getBrandWCPayRequest 的参数.
//
//  在前端 js 中这样调用:
//
//  WeixinJSBridge.invoke('getBrandWCPayRequest', {
//      "appId": "wx2421b1c4370ec43b",                        // 公众号名称，由商户传入
//      "timeStamp": " 1395712654",                           // 时间戳，自1970 年以来的秒数
//      "nonceStr": "e61463f8efa94090b1f366cccfbbb444",       // 随机串
//      "package": "prepay_id=u802345jgfjsdfgsdg888",         // package
//      "signType": "MD5",                                    // 微信签名方式
//      "paySign": "70EA570631E4BB79628FBCA90534C63FF7FADD89" // 微信签名
//  }, function (res) {
//      if (res.err_msg == "get_brand_wcpay_request:ok") {}
//  	// 使用以上方式判断前端返回,微信团队郑重提示：res.err_msg 将在用户支付成功后返回ok，但并不保证它绝对可靠。
//  });
//
type PayRequestParameters struct {
	AppId     string `json:"appId"`            // 必须, 公众号身份的唯一标识
	NonceStr  string `json:"nonceStr"`         // 必须, 商户生成的随机字符串, 32个字符以内
	TimeStamp int64  `json:"timeStamp,string"` // 必须, unixtime, 商户生成

	Package string `json:"package"` // 必须, 订单详情组合成的字符串

	Signature  string `json:"paySign"`  // 必须, 该 PayRequestParameters 自身的签名. see PayRequestParameters.SetSignature
	SignMethod string `json:"signType"` // 必须, 签名方式, 目前仅支持 MD5
}

// 设置签名字段.
//  appKey: 商户支付密钥Key
//
//  NOTE: 要求在 para *PayRequestParameters 其他字段设置完毕后才能调用这个函数, 否则签名就不正确.
func (para *PayRequestParameters) SetSignature(appKey string) (err error) {
	var Hash hash.Hash
	var Signature []byte

	switch para.SignMethod {
	case "md5", "MD5":
		Hash = md5.New()
		Signature = make([]byte, md5.Size*2)

	default:
		err = fmt.Errorf(`unknown sign method: %q`, para.SignMethod)
		return
	}

	// 字典序
	// appid
	// appkey
	// noncestr
	// package
	// timestamp
	Hash.Write([]byte("appid="))
	Hash.Write([]byte(para.AppId))
	Hash.Write([]byte("&appkey="))
	Hash.Write([]byte(appKey))
	Hash.Write([]byte("&noncestr="))
	Hash.Write([]byte(para.NonceStr))
	Hash.Write([]byte("&package="))
	Hash.Write([]byte(para.Package))
	Hash.Write([]byte("&timestamp="))
	Hash.Write([]byte(strconv.FormatInt(para.TimeStamp, 10)))

	hex.Encode(Signature, Hash.Sum(nil))
	Signature = bytes.ToUpper(Signature)

	para.Signature = string(Signature)
	return
}
