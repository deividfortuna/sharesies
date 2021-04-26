## Unofficial Go Sharesies NZ SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/deividfortuna/sharesies.svg)](https://pkg.go.dev/github.com/deividfortuna/sharesies)

![Sharesies NZ](https://images.squarespace-cdn.com/content/58bc788c59cc68b9696b9ee0/1543372882154-5E6PGXVJGOIQU30NTJKJ/sharesies.png?content-type=image%2Fpng)

### Installation
`go get github.com/deividfortuna/sharesies`

### Get Started
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

### Companies/funds listed
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

### Buy stock
```go
fundId := "0545fbc5-b579-4944-9057-55d01849a493"
costBuy, err := s.CostBuy(fundId, 100.00)
if err != nil {
	log.Fatal(err)
}
p, err := s.Buy(costBuy)
if err != nil {
	log.Fatal(err)
}

fmt.Println(p)
```

## LICENSE
MIT License - Copyright (c) 2021 [Deivid Fortuna](https://github.com/deividfortuna/sharesies/blob/main/LICENSE)