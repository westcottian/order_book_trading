package matchingengine

import (
			"fmt"
			"time"
		)

const MAX_PRICE = 10000000

type PricePoint struct {
    orderHead *Order
    orderTail *Order
}

func (pp *PricePoint) Insert(o *Order) {
    if pp.orderHead == nil {
        pp.orderHead = o
        pp.orderTail = o
    } else {
        pp.orderTail.next = o
        pp.orderTail = o
    }
}

type OrderStatus int
const (
    OS_NEW OrderStatus = iota
    OS_OPEN
    OS_PARTIAL
    OS_FILLED
    OS_CANCELLED
)

type Order struct {
    id uint64
    isBuy bool
    price uint32
    volume uint32
    status OrderStatus
    next *Order
	timestamp time.Time
}

func (o *Order) String() string {
    return fmt.Sprintf("\nOrder{id:%v,isBuy:%v,price:%v,volume:%v}",
        o.id, o.isBuy, o.price, o.volume)
}

func NewOrder(id uint64, isBuy bool, price uint32, volume uint32) *Order {
    return &Order{id: id, isBuy: isBuy, price: price, volume: volume,
        status: OS_NEW}
}

type OrderBook struct {
    // These are more estimates than reportable values
    ask uint32
    bid uint32
    orderIndex map[uint64]*Order
    prices [MAX_PRICE]*PricePoint
    actions chan<- *Action
}



func NewOrderBook(actions chan<- *Action) *OrderBook {
    ob := new(OrderBook)
    ob.bid = 0
    ob.ask = MAX_PRICE
    for i := range ob.prices {
        ob.prices[i] = new(PricePoint)
    }
    ob.actions = actions
    ob.orderIndex = make(map[uint64]*Order)
    return ob
}

func (ob *OrderBook) AddOrder(o *Order) {

	tm := time.Now()
	//go bx.TimeLimitHandler(o)
    // Try to fill immediately
    if o.isBuy {
        ob.actions <- NewBuyAction(o)
        ob.FillBuy(o)
    } else {
        ob.actions <- NewSellAction(o)
        ob.FillSell(o)
    }

    // Into the book
    if o.volume > 0 {
        ob.Match(o)
    }
}

func (ob *OrderBook) Match(o *Order) {
    pp := ob.prices[o.price]
    pp.Insert(o)
    o.status = OS_OPEN
    if o.isBuy && o.price >= ob.bid {
        ob.bid = o.price
    } else if !o.isBuy && o.price < ob.ask {
        ob.ask = o.price
    }
    ob.orderIndex[o.id] = o
}

func (ob *OrderBook) FillBuy(o *Order) {
    for ob.ask < o.price && o.volume > 0 {
        pp := ob.prices[ob.ask]
        ppOrderHead := pp.orderHead
        for ppOrderHead != nil {
            ob.Exchange(o, ppOrderHead)
            ppOrderHead = ppOrderHead.next
            pp.orderHead = ppOrderHead
        }
        ob.ask++
    }
}

func (ob *OrderBook) FillSell(o *Order) {
    for ob.bid >= o.price && o.volume > 0 {
        pp := ob.prices[ob.bid]
        ppOrderHead := pp.orderHead
        for ppOrderHead != nil {
            ob.Exchange(o, ppOrderHead)
            ppOrderHead = ppOrderHead.next
            pp.orderHead = ppOrderHead
        }
        ob.bid--
    }
}

func (ob *OrderBook) Exchange(o, ppOrderHead *Order) {
    if ppOrderHead.volume >= o.volume {
        ob.actions <- NewFilledAction(o, ppOrderHead)
        ppOrderHead.volume -= o.volume
        o.volume = 0
        o.status = OS_FILLED
        return
    } else {
        // Partial fill
        if ppOrderHead.volume > 0 {
            ob.actions <- NewPartialFilledAction(o, ppOrderHead)
            o.volume -= ppOrderHead.volume
            o.status = OS_PARTIAL
            ppOrderHead.volume = 0
        }
    }
}

func (ob *OrderBook) CancelOrder(id uint64) {
    ob.actions <- NewCancelAction(id)
    if o, ok := ob.orderIndex[id]; ok {
        // If this is the last order at a particular price point
        // we need to update the bid/ask...
        o.volume = 0
        o.status = OS_CANCELLED
    }
    ob.actions <- NewCancelledAction(id)
}

func (ob *OrderBook) Done() {
    ob.actions <- NewDoneAction()
}
