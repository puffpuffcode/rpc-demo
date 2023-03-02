package server

import (
	"bookstore/db"
	"bookstore/model"
	"bookstore/pb"
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedBookStoreServer
	bookStore *db.BookStore
}

func StartServer() {
	// 启动端口监听
	l, err := net.Listen("tcp", ":19090")
	if err != nil {
		log.Fatalf("net.Listen failed: %v\n", err)
	}
	// 创建 RPC 服务
	s := grpc.NewServer()
	// 绑定服务
	pb.RegisterBookStoreServer(s, &server{bookStore: &db.BookStore{}})

	log.Fatalln(s.Serve(l))
}

func StartHttpProxy() {
	conn, err := grpc.DialContext(
		context.Background(),
		":19090",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("grpc.DialContext err:", conn)
	}
	defer conn.Close()

	sm := runtime.NewServeMux()
	if err = pb.RegisterBookStoreHandler(context.Background(), sm, conn); err != nil {
		log.Fatalln("pb.RegisterBookStoreHandler err:", err)
	}

	s := &http.Server{
		Addr:    ":19191",
		Handler: sm,
	}

	log.Fatalln(s.ListenAndServe())
}

// 查询所有书架
func (s *server) ListShelves(ctx context.Context, emptypb *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	shelves, err := s.bookStore.ListShelves(ctx)
	if err != nil {
		log.Printf("BookStore.ListShelves err: %v\n", err)
		return nil, status.Errorf(codes.Internal, "Get BookStore.ListShelves ERROR")
	}
	respShelves := make([]*pb.Shelf,0,len(shelves))
	for _, sf := range shelves {
		respShelves = append(respShelves, &pb.Shelf{
			Id:    sf.ID,
			Theme: sf.Theme,
			Size:  sf.Size,
		})
	}
	resp := &pb.ListShelvesResponse{
		Shelves: respShelves,
	}
	return resp, nil
}

// 创建一个书架
func (s *server) CreateShelf(ctx context.Context, req *pb.CreateShelfRequest) (*pb.Shelf, error) {
	reqShelf := req.GetShelf()
	// 参数检查
	if reqShelf.Theme == "" {
		return nil, status.Error(codes.InvalidArgument,"reqShelf.Theme is invalid.")
	}
	shelf := model.Shelf{
		ID:        reqShelf.Id,
		Theme:     reqShelf.Theme,
		Size:      reqShelf.Size,
		CreateAt:  time.Now(),
		UpdatedAt: time.Now(),
	}
	rs, err := s.bookStore.CreateShelf(ctx, &shelf)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.Shelf{
		Id:    rs.ID,
		Theme: rs.Theme,
		Size:  rs.Size,
	}, nil
}

// 获取一个书架
func (s *server) GetShelf(ctx context.Context, req *pb.GetShelfRequest) (*pb.Shelf, error) {
	shelf, err := s.bookStore.GetShelf(ctx, req.GetShelf())
	// 如果查询结果为空
	if err == gorm.ErrRecordNotFound {
		return &pb.Shelf{}, nil
	}
	// 其他错误
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Get BookStore.ListShelves ERROR: %v", err)
	}
	resp := &pb.Shelf{
		Id:    shelf.ID,
		Theme: shelf.Theme,
		Size:  shelf.Size,
	}
	return resp, nil
}

// 删除一个书架
func (s *server) DeleteShelf(ctx context.Context, req *pb.DeleteShelfRequest) (*emptypb.Empty, error) {
	if req.GetShelf() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument,"id is <= 0")
	}
	err := s.bookStore.DeleteShelf(ctx, req.GetShelf())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "BookStore.DeleteShelf err: %v", err)
	}
	return &emptypb.Empty{}, nil
}
