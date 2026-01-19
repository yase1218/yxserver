package charge

import (
	"context"
	"gameserver/internal/game/model"
)

var orderManager *OrderManager

//var orderManagerOnce sync.Once

type OrderManager struct {
	ByOrderId map[string]*model.OrderInfo
	ByUid     map[uint64][]*model.OrderInfo
}

func InitOrderManager() {
	orderManager = new(OrderManager)
	orderManager.ByOrderId = make(map[string]*model.OrderInfo)
	orderManager.ByUid = make(map[uint64][]*model.OrderInfo)
	LoadAllOrder(context.Background())
}

func GetOrderManager() *OrderManager {
	//orderManagerOnce.Do(func() {
	//	if orderManager == nil {
	//		orderManager = &OrderManager{}
	//		orderManager.ByOrderId = make(map[string]*model.OrderInfo)
	//		orderManager.ByUid = make(map[uint64][]*model.OrderInfo)
	//	}
	//})
	return orderManager
}

func (manager *OrderManager) GetOrderById(orderId string) *model.OrderInfo {
	return manager.ByOrderId[orderId]
}

func (manager *OrderManager) GetOrderByUid(uid uint64) []*model.OrderInfo {
	return manager.ByUid[uid]
}

func (manager *OrderManager) AddOrder(order *model.OrderInfo) {
	manager.ByOrderId[order.OrderId] = order
	manager.ByUid[order.UserId] = append(manager.ByUid[order.UserId], order)
}

// 添加玩家订单
func (manager *OrderManager) AddUserOrders(orders []*model.OrderInfo) {
	var uid uint64
	for i := 0; i < len(orders); i++ {
		uid = orders[i].UserId
		manager.ByOrderId[orders[i].OrderId] = orders[i]
	}
	if uid > 0 {
		manager.ByUid[uid] = append(manager.ByUid[uid], orders...)
	}
}

func (manager *OrderManager) DelUserOrder(uid uint64) {
	for _, info := range manager.ByUid[uid] {
		delete(manager.ByOrderId, info.OrderId)
	}
	delete(manager.ByUid, uid)
}
