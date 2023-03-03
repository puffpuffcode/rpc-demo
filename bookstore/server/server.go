package server

import (
	"bookstore/db"
	"bookstore/model"
	"bookstore/pb"
	"context"
	"errors"
	"net/http"
	"strings"

	"strconv"

	"log"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

var (
	ErrNotFound error = errors.New("source not found")

	DefaultPageSize   int64  = 5
	DefaultPageNextID string = "0"
)

type server struct {
	pb.UnimplementedBookStoreServer
	bookStore *db.BookStore
}

func StartServer() {
	// 创建 RPC 服务
	s := grpc.NewServer()
	// 绑定服务
	pb.RegisterBookStoreServer(s, &server{bookStore: &db.BookStore{}})

	// 同一个端口处理 http 和 grpc 请求
	// gateway 服务，转发到 19090 端口
	gwmux := runtime.NewServeMux()
	dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterBookStoreHandlerFromEndpoint(context.Background(), gwmux, ":19090", dops); err != nil {
		log.Fatalln("RegisterBookStoreHandlerFromEndpoint err:", err.Error())
	}
	http.ListenAndServe(
		":19090",
		grpcHandlerFunc(s, gwmux),
	)
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// func StartHttpProxy() {
// 	conn, err := grpc.DialContext(
// 		context.Background(),
// 		":19090",
// 		grpc.WithBlock(),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		log.Fatalln("grpc.DialContext err:", conn)
// 	}
// 	defer conn.Close()

// 	sm := runtime.NewServeMux()
// 	if err = pb.RegisterBookStoreHandler(context.Background(), sm, conn); err != nil {
// 		log.Fatalln("pb.RegisterBookStoreHandler err:", err)
// 	}

// 	s := &http.Server{
// 		Addr:    ":19191",
// 		Handler: sm,
// 	}

// 	log.Fatalln(s.ListenAndServe())
// }

// 查询所有书架
func (s *server) ListShelves(ctx context.Context, emptypb *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	shelves, err := s.bookStore.ListShelves(ctx)
	if err != nil {
		log.Printf("BookStore.ListShelves err: %v\n", err)
		return nil, status.Errorf(codes.Internal, "Get BookStore.ListShelves ERROR")
	}
	respShelves := make([]*pb.Shelf, 0, len(shelves))
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
		return nil, status.Error(codes.InvalidArgument, "reqShelf.Theme is invalid.")
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
		return nil, status.Errorf(codes.InvalidArgument, "id is <= 0")
	}
	err := s.bookStore.DeleteShelf(ctx, req.GetShelf())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "BookStore.DeleteShelf err: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// 创建一本书
func (s *server) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	// 检查参数
	if req.GetShelf() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "req.Shelf <= 0")
	}
	// 将书加入书架
	bookResp, err := s.bookStore.CreateBook(ctx, req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &pb.CreateBookResponse{}, status.Errorf(codes.Internal, "not a shelf")
	} else if errors.Is(err, gorm.ErrRegistered) {
		return &pb.CreateBookResponse{}, status.Error(codes.AlreadyExists, err.Error())
	}
	return &pb.CreateBookResponse{
		Book: &pb.Book{
			Id:      bookResp.ID,
			Author:  bookResp.Author,
			Title:   bookResp.Title,
			ShelfId: bookResp.ShelfID,
		},
	}, nil
}

// 查询书
func (s *server) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	// 检查参数
	shelfID := req.GetShelf()
	if shelfID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	var resBooks []*model.Book
	var page model.Page
	var err error
	// 如果 page_token == ""
	if req.GetPageToken() == "" {
		// 默认第一页
		resBooks, err = s.bookStore.ListBooks(ctx, shelfID, DefaultPageNextID, DefaultPageSize+1)
		page.PageSize = DefaultPageSize
	} else {
		// page_token 无效
		page = model.Token(req.GetPageToken()).Decode()
		if page.IsInVaild() {
			return nil, status.Errorf(codes.InvalidArgument, "invaild token")
		}
		// 有效，获取分页
		resBooks, err = s.bookStore.ListBooks(ctx, shelfID, page.NextID, page.PageSize+1)
	}
	// 统一处理 error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}
	// 拼接响应结果
	respBooks := make([]*pb.Book, 0, len(resBooks))
	for _, v := range resBooks {
		respBooks = append(respBooks, &pb.Book{
			Id:      v.ID,
			Author:  v.Author,
			Title:   v.Title,
			ShelfId: v.ShelfID,
		})
	}
	// 还有数据，更新 token
	if len(resBooks) > int(page.PageSize) {
		page.NextTimeAtUTC = time.Now().Unix()
		page.NextID = strconv.Itoa(int(respBooks[len(respBooks)-1].Id))
		return &pb.ListBooksResponse{
			Books:     respBooks,
			PageToken: string(page.Encode()),
		}, nil
	}
	return &pb.ListBooksResponse{
		Books:     respBooks,
		PageToken: "",
	}, nil
}
