package main

import (
	"context"
	"fmt"

	protoGo "github.com/kandevarg/deliveryapp.stockservice/proto/protoGo"
	micro "github.com/micro/go-micro"
)

//Repository ...
type Repository interface {
	GetStockInfo(*protoGo.GetStockInfoRequest) (int32, error)
}

type memoryRepository struct {
	stocks []*protoGo.Product
}

// Our grpc service handler
type service struct {
	repo Repository
}

//GetStockInfo ...
func (repo *memoryRepository) GetStockInfo(req *protoGo.GetStockInfoRequest) (int32, error) {

	for _, product := range repo.stocks {
		if product.Id == req.ProductId {
			return product.StockQty, nil
		}
	}

	return 0, nil
}

func (s *service) GetStockInfo(ctx context.Context, req *protoGo.GetStockInfoRequest, res *protoGo.GetStockInfoResponse) error {

	// // Find the next available vessel
	stockQty, err := s.repo.GetStockInfo(req)
	if err != nil {
		return err
	}

	res = &protoGo.GetStockInfoResponse{
		ProductId: req.ProductId,
		StockQty:  stockQty,
	}

	return nil
}

func main() {

	stocks := []*protoGo.Product{
		&protoGo.Product{
			Id:       "1",
			StockQty: 10,
		},
		&protoGo.Product{
			Id:       "2",
			StockQty: 30,
		},
		&protoGo.Product{
			Id:       "3",
			StockQty: 200,
		},
	}

	fmt.Println("stocks", stocks)

	repo := &memoryRepository{stocks}

	microService := micro.NewService(
		micro.Name("deliveryapp.stockservice"),
	)

	microService.Init()

	// Register our implementation with
	protoGo.RegisterStockServiceHandler(microService.Server(), &service{repo})

	// Run the servers
	if err := microService.Run(); err != nil {
		fmt.Println(err)
	}
}
