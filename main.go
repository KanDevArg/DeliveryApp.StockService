package main

import (
	"context"
	"fmt"

	stockServiceProtoGo "github.com/kandevarg/deliveryapp.stockservice/proto/protoGo"
	micro "github.com/micro/go-micro"
)

//Repository ...
type Repository interface {
	GetStockInfo(*stockServiceProtoGo.GetStockInfoRequest) (int32, error)
}

type memoryRepository struct {
	stocks []*stockServiceProtoGo.Product
}

// Our grpc service handler
type service struct {
	repo Repository
}

//GetStockInfo ...
func (repo *memoryRepository) GetStockInfo(req *stockServiceProtoGo.GetStockInfoRequest) (int32, error) {

	for _, product := range repo.stocks {
		if product.Id == req.ProductId {
			return product.StockQty, nil
		}
	}

	return 0, nil
}

func (s *service) GetStockInfo(ctx context.Context, req *stockServiceProtoGo.GetStockInfoRequest, res *stockServiceProtoGo.GetStockInfoResponse) error {

	// // Find the next available vessel
	stockQty, err := s.repo.GetStockInfo(req)
	if err != nil {
		return err
	}

	res = &stockServiceProtoGo.GetStockInfoResponse{
		ProductId: req.ProductId,
		StockQty:  stockQty,
	}

	return nil
}

func main() {

	stocks := []*stockServiceProtoGo.Product{
		&stockServiceProtoGo.Product{
			Id:       "1",
			StockQty: 10,
		},
		&stockServiceProtoGo.Product{
			Id:       "2",
			StockQty: 30,
		},
		&stockServiceProtoGo.Product{
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
	stockServiceProtoGo.RegisterStockServiceHandler(microService.Server(), &service{repo})

	// Run the servers
	if err := microService.Run(); err != nil {
		fmt.Println(err)
	}
}
