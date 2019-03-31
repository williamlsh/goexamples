package main

import (
	"context"
	"fmt"
	pb "goexamples/consignment-service/proto/consignment"

	"github.com/micro/go-micro"
)

const port = ":50051"

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository implements IRepository interface.
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.consignments = append(repo.consignments, consignment)
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service implements ShippingServer interface.
type service struct {
	repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, rsp *pb.Response) error {
	// Save consignment.
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	rsp.Created = true
	rsp.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, in *pb.GetRequest, rsp *pb.Response) error {
	consignments := s.repo.GetAll()
	rsp.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("consignment"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "shipping consignment",
		}),
	)

	srv.Init()

	pb.RegisterShippingHandler(srv.Server(), &service{repo: repo})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
