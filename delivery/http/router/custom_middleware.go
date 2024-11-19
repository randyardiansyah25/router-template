package router

import (
	"net/http"
	"router-template/delivery"
	"router-template/delivery/http/router/middleware"

	"github.com/gin-gonic/gin"
)

//! Untuk custom middleware harus diletakan setelah middlware.RequestLogger dan sebelum middleware.ResponseLogger
//! Jika ada validasi dan mengharuskan me-reject request, jangan lakukan abort, agar middleware.ResponseLogger
//! masih bisa mencetak response.

//! Gunakan middleware.ResponseString() atau middleware.RersponseJSON kemudian tambahkan `return` didalam function
//! yang melakukan reject

func AuthCustomMiddleware(ctx *gin.Context) {
	myHeader := struct {
		RequestDate string `header:"X-Request-Date"`
		ClientId    string `header:"X-Client-ID"`
		Signature   string `header:"X-SIGN"`
	}{}

	if er := ctx.BindHeader(&myHeader); er != nil {
		delivery.PrintError("error while binding header: ", er)
		middleware.ResponseString(ctx, http.StatusBadRequest, "Bad Request - Invalid header")
		return
	}

	if myHeader.Signature != "+eu3YE7nfX1mMdrz6DSe3bUTkhQq/pYltvUtaWmXVOU=" {
		middleware.ResponseString(ctx, http.StatusUnauthorized, "Unauthorized - Invalid signature")
		return
	}

	ctx.Next()
}
