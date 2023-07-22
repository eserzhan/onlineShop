package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/yervsil/onlineShop/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
	cartsTable = "carts"
	productsTable = "products"
	cartItemsTable = "cartItems"
	ordersTable = "orders"
	orderItemsTable = "orderItems"
)

type UsersRepository struct {
	db *sqlx.DB
}

func (r *UsersRepository) CreateUser(user domain.User) (int, error){
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash, email, phone, registered_at) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Password, user.Email, user.Phone, user.RegisteredAt)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UsersRepository) GetByCredentials(email, password string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s where email = $1 and password_hash = $2", usersTable)

	err := r.db.Get(&id, query, email, password)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UsersRepository) SetSession(userID string, session domain.Session) error {
	query := `
		UPDATE Users
		SET session = $1, last_login = $2
		WHERE id = $3
	`
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, sessionJSON, time.Now(), userID)
	return err
}


func (r *UsersRepository) GetByRefreshToken(refreshToken string) (int, error) {

	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE session ->> 'refreshToken' = $1 AND (session ->> 'expiresAt')::timestamp > NOW() LIMIT 1", usersTable)

	err := r.db.Get(&id, query, refreshToken)

	if err != nil {
		return 0, err
	}

	return id, nil
}


func (r *UsersRepository) AddToCart(quantity int, userId, productId string) error {
	var product domain.Product
	inStockQuery := fmt.Sprintf("SELECT quantity FROM %s where id = $1", productsTable)
	row := r.db.QueryRow(inStockQuery, productId)
	err := row.Scan(&product.Quantity)
	if err != nil {
		return err 
	}

	if quantity > product.Quantity {
		return errors.New("insufficient quantity available")
	}


    tx, err := r.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Проверяем, существует ли товар с указанным идентификатором
    productExistsQuery := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1", productsTable)
    var productExists bool
    err = tx.QueryRow(productExistsQuery, productId).Scan(&productExists)
    if err != nil {
        return err
    }
    if !productExists {
        return fmt.Errorf("product not found")
    }

    // Проверяем, существует ли у пользователя корзина
    cartExistsQuery := fmt.Sprintf("SELECT id FROM %s WHERE user_id = $1", cartsTable)
	var cartID int 
    err = tx.QueryRow(cartExistsQuery, userId).Scan(&cartID)
    if err != nil {
        createCartQuery := fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING id", cartsTable)
        var cartID int
        err = tx.QueryRow(createCartQuery, userId).Scan(&cartID)
        if err != nil {

            return err
        }

        insertCartItemQuery := fmt.Sprintf("INSERT INTO %s (cart_id, product_id, quantity) VALUES ($1, $2, $3)", cartItemsTable)
        _, err = tx.Exec(insertCartItemQuery, cartID, productId, quantity)
        if err != nil {

            return err
        }
    }else  {
        var productExists bool
        // Если у пользователя уже есть корзина, и добавляемый товар уже в корзине, увеличиваем количество имеющегося товара 
        itemAlreadyInCart := fmt.Sprintf("SELECT 1 FROM %s WHERE cart_id = $1 AND product_id = $2", cartItemsTable)
		_ = tx.QueryRow(itemAlreadyInCart, cartID, productId).Scan(&productExists)

        if productExists {
            updateCartItemQuery := fmt.Sprintf("UPDATE %s SET quantity = quantity + $1 WHERE cart_id = $2 AND product_id = $3", cartItemsTable)
            _, err = tx.Exec(updateCartItemQuery, quantity, cartID, productId)
            if err != nil {
    
                return err
            }
        }else{
            insertCartItemQuery := fmt.Sprintf("INSERT INTO %s (cart_id, product_id, quantity) VALUES ($1, $2, $3)", cartItemsTable)
            _, err = tx.Exec(insertCartItemQuery, cartID, productId, quantity)
            if err != nil {
    
                return err
            }
        } 
    }

    return tx.Commit()
}

func (r *UsersRepository) CreateOrder(userId, delivery, payment string) (int, error) {
    tx, err := r.db.Begin()
    if err != nil {
        return -1, err
    }
    defer tx.Rollback()

    // Создаем новый заказ
    createOrderQuery := fmt.Sprintf(`INSERT INTO %s (user_id, status, delivery_method, payment_method) VALUES ($1, $2, $3, $4) RETURNING id`, ordersTable)
    var orderID int
    err = tx.QueryRow(createOrderQuery, userId, "created", delivery, payment).Scan(&orderID)
    if err != nil {
        return -1, err
    }

    // Получаем идентификатор корзины пользователя
    cartIDQuery := fmt.Sprintf(`SELECT id FROM %s WHERE user_id = $1`, cartsTable)
    var cartID int
    err = tx.QueryRow(cartIDQuery, userId).Scan(&cartID)
    if err != nil {
        return -1, err
    }

    // Получаем товары из корзины пользователя
    cartItemsQuery := fmt.Sprintf(`SELECT product_id, quantity FROM %s WHERE cart_id = $1`, cartItemsTable)
    rows, err := r.db.Query(cartItemsQuery, cartID)
    if err != nil {
        return -1, err
    }
    defer rows.Close()

    // Вставляем записи о товарах в таблицу order_items и обновляем количество товара в таблице products
    createOrderItemQuery := fmt.Sprintf(`INSERT INTO %s (order_id, product_id, quantity) VALUES ($1, $2, $3)`, orderItemsTable)
    decreaseQuantityFromProduct := fmt.Sprintf(`UPDATE %s SET quantity = quantity - $1 WHERE id = $2`, productsTable)
    decreaseQuantityFromCart := fmt.Sprintf(`UPDATE %s ci SET quantity = CASE WHEN ci.quantity > pt.quantity THEN pt.quantity ELSE ci.quantity END FROM %s pt WHERE pt.id = ci.product_id AND cart_id != $1 AND product_id = $2`, cartItemsTable, productsTable)
    for rows.Next() {
        var productID, quantity int
        err := rows.Scan(&productID, &quantity)
        if err != nil {
            return -1, err
        }

        _, err = tx.Exec(createOrderItemQuery, orderID, productID, quantity)
        if err != nil {
            return -1, err
        }

        _, err = tx.Exec(decreaseQuantityFromProduct, quantity, productID)
        if err != nil {
            return -1, err
        }

        _, err = tx.Exec(decreaseQuantityFromCart, cartID, productID)
        if err != nil {
            return -1, err
        }
    }

    deleteCartItemsQuery := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1`, cartsTable)
    _, err = tx.Exec(deleteCartItemsQuery, userId)
    if err != nil {
        return -1, err
    }

    err = tx.Commit()
    if err != nil {
        return -1, err
    }

    return orderID, nil
}

func(r *UsersRepository) GetOrders(userId string) ([]domain.Order, error) {
    query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, ordersTable)

    var res []domain.Order 
    err := r.db.Select(&res, query, userId)
    if err != nil {
        return []domain.Order{}, err 
    }

    return res, nil 
}

func(r *UsersRepository) GetOrderById(userId, orderId string) (domain.Order, error) {
    query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1 AND id = $2`, ordersTable)

    var res domain.Order
    err := r.db.Get(&res, query, userId, orderId)
    if err != nil {
        return domain.Order{}, err
    }

    return res, nil 
}

func(r *UsersRepository) GetCart(userId string) (domain.Cart, error) {
    query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, cartsTable)

    var res domain.Cart 
    err := r.db.Get(&res, query, userId)
    if err != nil {
        return domain.Cart{}, err 
    }

    return res, nil 
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

