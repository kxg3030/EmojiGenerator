package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

type SslMiddleware struct {
	
}

func NewSslMiddleware() *SslMiddleware{
	return &SslMiddleware{

	}
}

func (this *SslMiddleware)Render()gin.HandlerFunc  {
	return func(context *gin.Context) {
		if context.Request.TLS != nil{
			context.Set("protocol","https://")
		}else{
			context.Set("protocol","http://")
		}
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost    : "0.0.0.0:9527",
		})
		err := secureMiddleware.Process(context.Writer, context.Request)
		if err != nil {
			return
		}
		context.Next()
	}
}