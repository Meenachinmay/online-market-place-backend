package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	grpctransportsoda_finance "soda-interview/app/services/soda-interview-grpc/internal/transport/grpc/soda-finance"
	grpctransportorder "soda-interview/app/services/soda-interview-grpc/internal/transport/grpc/order"
	grpctransportproduct "soda-interview/app/services/soda-interview-grpc/internal/transport/grpc/product"
	grpctransportreferral_blog "soda-interview/app/services/soda-interview-grpc/internal/transport/grpc/referral-blog"

	"soda-interview/business/core/finance"
	"soda-interview/business/core/order"
	"soda-interview/business/core/product"
	"soda-interview/business/core/referral-blog"

	orderstore "soda-interview/business/data/stores/order"
	productstore "soda-interview/business/data/stores/product"
	blogstore "soda-interview/business/data/stores/referral-blog"
	financestore "soda-interview/business/data/stores/soda-finance"

	"soda-interview/foundation/bootstrap"
	"soda-interview/foundation/logger"
	orderv1 "soda-interview/foundation/proto/order/v1"
	productv1 "soda-interview/foundation/proto/product/v1"
	referralblogv1 "soda-interview/foundation/proto/referral-blog/v1"
	financev1 "soda-interview/foundation/proto/soda-finance/v1"
)

func main() {
	bootstrap.Run(func(log *logger.Logger, db *pgxpool.Pool, grpcServer *grpc.Server) {
		// Stores
		productSt := productstore.NewStore(log, db)
		blogSt := blogstore.NewStore(log, db)
		orderSt := orderstore.NewStore(log, db)
		financeSt := financestore.NewStore(log, db)

		// Core Services
		// Note: Product and Blog services accept interfaces. Order and Finance accept concrete stores for TX handling.
		productService := product.NewService(log, productSt)
		blogService := referralblog.NewService(log, blogSt)
		financeService := finance.NewService(log, db, financeSt)
		orderService := order.NewService(log, db, orderSt, productSt, blogSt, financeSt)

		// Transport Handlers
		productHandler := &grpctransportproduct.Handler{Service: productService}
		blogHandler := &grpctransportreferral_blog.Handler{Service: blogService}
		financeHandler := &grpctransportsoda_finance.Handler{Service: financeService}
		orderHandler := &grpctransportorder.Handler{Service: orderService}

		// Registration
		productv1.RegisterProductServiceServer(grpcServer, productHandler)
		referralblogv1.RegisterBlogServiceServer(grpcServer, blogHandler)
		financev1.RegisterFinanceServiceServer(grpcServer, financeHandler)
		orderv1.RegisterOrderServiceServer(grpcServer, orderHandler)

		log.Info("All services registered")
	})
}
