package manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	shorturlpb "github.com/username/shorturl/internal/rpc/proto"
	"google.golang.org/grpc"
)

// ClientManager 集中管理所有 gRPC 连接和客户端存根
type ClientManager struct {
	// 1. Connection Map: 存储并复用连接 (Key: Address, Value: *grpc.ClientConn)
	connMap map[string]*grpc.ClientConn

	// 2. Client Map: 存储客户端存根 (Key: ServiceName, Value: *ServiceEntry)
	clientMap map[string]interface{}

	mu sync.RWMutex // 互斥锁，用于保护 map 的并发访问
}

func NewClientManer() *ClientManager {
	return &ClientManager{
		connMap:   make(map[string]*grpc.ClientConn),
		clientMap: make(map[string]interface{}),
	}
}

func (m *ClientManager) InitServices(ctx context.Context, configs []string, opts ...grpc.DialOption) error {

	log.Println("初始化", configs)
	// 互斥锁应该是 Lock() 而不是 RLock()，因为 InitServices 会修改 m.connMap 和 m.clientMap
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, cfg := range configs {
		// 1. 检查连接是否存在并复用
		conn, ok := m.connMap[cfg]
		var err error
		if !ok {
			// 2. 建立新连接
			conn, err = grpc.NewClient(cfg, opts...)
			if err != nil {
				return fmt.Errorf("failed to dial %s: %w", cfg, err)
			}
			m.connMap[cfg] = conn
		}
		// 3. 检查是否有对应的客户端创建闭包
		operation, ok := clientFactoryMap[cfg]

		if !ok {
			// 如果配置的 key 在 clientMap 中不存在，应该返回错误或跳过
			return fmt.Errorf("unknown service configuration key: %s", cfg)
		}

		// 4. 执行闭包创建客户端存根，并存储
		// 这里不需要取地址符号 &，因为 operation(conn) 返回的就是存根对象
		temp := operation(conn)
		reuslt, ok1 := temp.(shorturlpb.ShortenerServiceClient)
		log.Println(reuslt, ok1)
		m.clientMap[cfg] = temp
	}
	log.Println(m)
	go func() {
		<-ctx.Done()
		log.Println("manager start close")
		m.Close()
	}()
	return nil
}

func (m *ClientManager) GetClient(serviceName string) (interface{}, error) {
	log.Println(m)
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.clientMap[serviceName]
	log.Println(m.clientMap, serviceName)
	if !ok {
		return nil, fmt.Errorf("client for service %s not registered", serviceName)
	}
	return client, nil
}

// Close safely closes all unique grpc.ClientConn instances.
func (m *ClientManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var allErrMsg strings.Builder
	// 遍历 Conn Map，只关闭唯一的连接实例
	for addr, conn := range m.connMap {
		if err := conn.Close(); err != nil {
			allErrMsg.WriteString(fmt.Sprintf("failed to close conn for %s: %s| \n", addr, err))
		}
	}

	// 清空 maps
	m.connMap = make(map[string]*grpc.ClientConn)
	m.clientMap = make(map[string]interface{})
	// 清空 maps 放在这里，确保即使有错误，maps 也会被清空

	if allErrMsg.Len() > 0 {
		return errors.New(allErrMsg.String())
	}

	return nil // 如果没有错误，返回 nil
}
