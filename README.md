# order_book_trading
A matching or trading engine is a piece of software that keeps a record of all open orders in a market and generates new trades if the two orders can be fulfilled by each other.

For this project, my basic trading engine supports:

    adding, canceling, and executing buy/sell orders
    Execute all operations in as close to O(1) as possible.
    
Each Order is either a buy or sell, has a limit price and amount, and a status that lets us know whether the order is open, partially filled, filled, etc. 
Lastly, each order is linked to the next order at the same price point so we can ensure orders are examined in the order they are entered.

The OrderBook does most of the heavy lifting by keeping track of the current maximum bid and minimum ask, an index mapping order IDs to pointers so we can easily cancel outstanding orders, an array of all possible price points, as well as a channel to report actions (fills, cancellations, etc.) to some handler as they occur.

Adding orders

Inserting a new order into the book is just a matter of appending it to the list of orders at its price point, updating the order book's bid or ask if necessary, and adding an entry in the order index. These are all O(1).

Canceling orders

Canceling is done by looking up the order by ID from the index and setting the amount to 0, also O(1). 

Performance

Testing performance of the system was done by pre-generating large numbers of random orders, varying the number of orders, mean and standard deviation of price, and maximum amount. The number of actions per second was then timed under different configurations.

