package api

import "squirrel-dev/internal/squ-apiserver/module/auth/api/res"

func toTokenResponse(token string) res.TokenRes {
	return res.TokenRes{
		Token: token,
	}
}
