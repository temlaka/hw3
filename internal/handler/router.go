package handler

import (
	"net/http"
)

func Routes(authH *AuthHandler, cryptoH *CryptoHandler, secret string) http.Handler {

	cryptoMux := http.NewServeMux()

	cryptoMux.HandleFunc("GET /crypto", cryptoH.GetCryptos)
	cryptoMux.HandleFunc("POST /crypto", cryptoH.AddCrypto)
	cryptoMux.HandleFunc("GET /crypto/{symbol}", cryptoH.GetCryptoInfo)
	cryptoMux.HandleFunc("PUT /crypto/{symbol}/refresh", cryptoH.UpdateCryptoInfo)
	cryptoMux.HandleFunc("GET /crypto/{symbol}/history", cryptoH.GetCryptoHistory)
	cryptoMux.HandleFunc("GET /crypto/{symbol}/stats", cryptoH.GetStatsCrypto)
	cryptoMux.HandleFunc("DELETE /crypto/{symbol}", cryptoH.DeleteCrypto)

	protected := AuthCheck(secret, cryptoMux)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", authH.Register)
	mux.HandleFunc("POST /auth/login", authH.Login)
	mux.Handle("/crypto", protected)
	mux.Handle("/crypto/", protected)

	return mux

}
