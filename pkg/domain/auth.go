package domain

type Token struct {
	Access  string
	Refresh string
}

type JwtConfig struct {
	Issuer               string
	ExpireTimeForAccess  int64
	ExpireTimeForRefresh int64
}
