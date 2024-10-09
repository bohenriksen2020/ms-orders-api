# From tutorial: 
https://github.com/dreamsofcode-io/golang-microservice-course-nn/blob/main/README.md

# Insert a new order
curl -X POST -H "Content-Type: application/json" -d '{"customer_id":"'"$(uuidgen)"'","line_items":[{"item_id":"'"$(uuidgen)"'","quantity":5,"price":1999}]}' http://localhost:3000/orders

# Test
curl -sS localhost:3000/orders | jq

# Update order to shipped
curl -X PUT -d '{"status":"shipped"}' -sS "localhost:3000/orders/3347639392487409127" | jq


# Update order to completed
curl -X PUT -d '{"status":"completed"}' -sS "localhost:3000/orders/3347639392487409127" | jq


# Delete order
curl -X DELETE localhost:3000/orders/3347639392487409127