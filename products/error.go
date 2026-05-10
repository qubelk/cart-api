package products

import "errors"

var ItemNotExists error = errors.New("Error when deleting item: Item not exists")
