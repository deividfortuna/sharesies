## Unofficial Sharesies NZ SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/deividfortuna/sharesies.svg)](https://pkg.go.dev/github.com/deividfortuna/sharesies)


### Get Started
`go get github.com/deividfortuna/sharesies`

### New Sharesies Client
```go
c := &sharesies.SharesiesCredentials{
	Username: "email@exmaple.com",
	Password: "your_password_here",
}
s, err := sharesies.New(c)
if err != nil {
	log.Fatal(err)
}
```

### Get companies/funds listed
```go
ir := &sharesies.InstrumentsRequest{
	Page:            1,
	Perpage:         100,
	Sort:            "name",
	Pricechangetime: "1y",
}
i, err := s.Instruments(ir)
if err != nil {
	log.Fatal(err)
}

fmt.Println(i)
```

## LICENSE
MIT License - Copyright (c) 2021 [Deivid Fortuna](https://github.com/deividfortuna/sharesies/blob/main/LICENSE)