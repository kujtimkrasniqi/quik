package route

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Wallet struct {
	Wallet_id int `json:"wallet_id"`
	Balance   int `json:"balance"`
}

func getBalance() {
	router := gin.Default()
	router.GET("/api/v1/wallets/:wallet_id/balance", func(c *gin.Context) {
		wallet_id := c.Param("wallet_id")
		dsn := "root@tcp(127.0.0.1:3306)/test"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}

		var wallet Wallet

		db.Raw("SELECT wallet_id, balance FROM wallet_balance where wallet_id = ?", wallet_id).Scan(&wallet)
		c.JSON(200, gin.H{
			"wallet_id": &wallet.Wallet_id,
			"balance":   &wallet.Balance,
		})
	})
}

func creditBalance() {
	router := gin.Default()
	router.POST("/api/v1/wallets/:wallet_id/credit", func(c *gin.Context) {
		wallet_id := c.Param("wallet_id")
		credit := c.PostForm("credit")
		credit_int, err := strconv.Atoi(credit)

		dsn := "root@tcp(127.0.0.1:3306)/test"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if credit_int < 0 {
			c.JSON(400, "Negative value it's not allowed")
			sqlDB, err := db.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.Close()
		}

		if err != nil {
			panic("failed to connect database")
		}
		var wallet Wallet
		db.Raw("SELECT wallet_id, balance FROM wallet_balance where wallet_id = ?", wallet_id).Scan(&wallet)
		c.JSON(200, &wallet.Balance)
		if credit_int > *&wallet.Balance {
			c.JSON(400, "You don't have enough money, please deposit to continue with buying")
			sqlDB, err := db.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.Close()
		}
		updated_balance := *&wallet.Balance - credit_int

		db.Exec("INSERT INTO wallet_credit (wallet_id, credit) VALUES (?,?)", wallet_id, credit_int)
		// db.Exec("INSERT INTO wallet_balance (wallet_id, balance, debit, credit) VALUES (?,?,?,?)", wallet_id, updated_balance, 0, credit_int)
		db.Exec("UPDATE wallet_balance SET balance =? WHERE wallet_id = ?", updated_balance, wallet_id)
	})
}

func debitBalance() {
	router := gin.Default()
	router.POST("/api/v1/wallets/:wallet_id/debit", func(c *gin.Context) {
		wallet_id := c.Param("wallet_id")
		debit := c.PostForm("debit")
		debit_int, err := strconv.Atoi(debit)

		dsn := "root@tcp(127.0.0.1:3306)/test"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if debit_int < 0 {
			c.JSON(400, "Negative value it's not allowed")
			sqlDB, err := db.DB()
			if err != nil {
				panic("failed to connect database")
			}
			sqlDB.Close()
		}

		if err != nil {
			panic("failed to connect database")
		}
		var wallet Wallet
		db.Raw("SELECT wallet_id, balance FROM wallet_balance where wallet_id = ?", wallet_id).Scan(&wallet)

		c.JSON(200, &wallet.Wallet_id)

		if *&wallet.Wallet_id == 0 {
			c.JSON(200, "here")
			db.Exec("INSERT INTO wallet_balance (wallet_id, balance) VALUES (?,?)", wallet_id, debit_int)
		}

		updated_balance := *&wallet.Balance + debit_int
		db.Exec("INSERT INTO wallet_debit (wallet_id, debit) VALUES (?,?)", wallet_id, debit_int)
		// db.Exec("INSERT INTO wallet_balance (wallet_id, balance, debit, credit) VALUES (?,?,?,?)", wallet_id, updated_balance, debit_int, 0)
		db.Exec("UPDATE wallet_balance SET balance =? WHERE wallet_id = ?", updated_balance, wallet_id)

	})
}

func handleRequests() {
	router := gin.Default()
	getBalance()
	creditBalance()
	debitBalance()
	router.Run(":8081")
}

func main() {
	handleRequests()
}
