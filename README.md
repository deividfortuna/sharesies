## Unofficial Go Sharesies NZ SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/deividfortuna/sharesies.svg)](https://pkg.go.dev/github.com/deividfortuna/sharesies)

![Sharesies NZ](https://images.squarespace-cdn.com/content/58bc788c59cc68b9696b9ee0/1543372882154-5E6PGXVJGOIQU30NTJKJ/sharesies.png?content-type=image%2Fpng)

The project is under heavy development so interfaces and structure might/will change.
Library initially written to be used by the project [Sharesies Bot](https://github.com/deividfortuna/sharesies-bot) to apply Dollar-cost averaging on NZ Market Exchange.

### Installation
`go get github.com/deividfortuna/sharesies`

### Get Started
```go
ctx := context.Background()
s, _ := sharesies.New(nil)
```

### Authenticate
```go
p, err := s.Authenticate(ctx, &sharesies.SharesiesCredentials{
	Username: "email@exmaple.com",
	Password: "your_password_here",
})
if err != nil {
	log.Fatal(err)
}
```

### Listed Companies/Funds
```go
i, err := s.Instruments(ctx, &sharesies.InstrumentsRequest{
	Page:            1,
	Perpage:         100,
	Sort:            "name",
	Pricechangetime: "1y",
})
if err != nil {
	log.Fatal(err)
}

fmt.Println(i)
```

### Buy Transaction
```go
fundId := "0545fbc5-b579-4944-9057-55d01849a493"
costBuy, err := s.CostBuy(ctx, fundId, 100.00)
if err != nil {
	log.Fatal(err)
}

b, err := s.Buy(ctx, costBuy)
if err != nil {
	log.Fatal(err)
}

fmt.Println(b)
```

### Sell Transaction
```go
fundId := "0545fbc5-b579-4944-9057-55d01849a493"
shares := 1.5 //number of shares to sell

costSell, err := s.CostSell(ctx, fundId, shares)
if err != nil {
	log.Fatal(err)
}

sr, err := s.Sell(ctx, costSell)
if err != nil {
	log.Fatal(err)
}

fmt.Println(b)
```

## LICENSE
MIT License - Copyright (c) 2021 [Deivid Fortuna](https://github.com/deividfortuna/sharesies/blob/main/LICENSE)
