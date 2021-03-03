package tls

import "github.com/muzin/go_rt/try"

// 服务 Listen 异常
var X509KeyPairException = try.DeclareException("X509KeyPairException")
