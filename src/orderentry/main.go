package main

import (
   me "orderentry/matchingengine"
)

func main() {
    actions := make(chan *me.Action)
    done := make(chan bool)

    go me.ConsoleActionHandler(actions, done)
	

    ob := me.NewOrderBook(actions)
    ob.AddOrder(me.NewOrder(1, false, 10, 100))
    ob.AddOrder(me.NewOrder(2, false, 20, 50))
    ob.AddOrder(me.NewOrder(3, false, 20, 100))
	ob.AddOrder(me.NewOrder(4, false, 30, 10))
    ob.AddOrder(me.NewOrder(5, true, 20, 50))
	ob.AddOrder(me.NewOrder(6, true, 40, 50))
	ob.AddOrder(me.NewOrder(7, true, 50, 100))
	ob.AddOrder(me.NewOrder(8, true, 70, 10))
    ob.CancelOrder(1)
    ob.Done()

    <-done
}
