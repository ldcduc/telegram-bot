package node

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = 
	port     = 
	user     = 
	password = 
	dbname   = 
)

func GetPsqlInfo(offset int) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port + offset, user, password, dbname)

	return psqlInfo
}

func Check_last_consensused_height(maxNode int, last_consensused_block_height []int, failedNode []int) []int {
	log.Printf("Checking last consensused height")
	for nodeIndex := 1; nodeIndex <= maxNode; nodeIndex ++ {
		db, err := sql.Open("postgres", GetPsqlInfo(nodeIndex)) // 9000 + offset
		if err != nil {
			panic(err)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		fmt.Println("Successfully connected to node ", nodeIndex)

		var blockValue BlockValue
		blockValueSql := "SELECT key, value FROM block_values"

		err = db.QueryRow(blockValueSql).Scan(&blockValue.key, &blockValue.value)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}

		if last_consensused_block_height[nodeIndex - 1] == blockValue.value {
			fmt.Println("FALSE AT NODE ", nodeIndex)
			failedNode = append(failedNode, nodeIndex)
		} else {
			last_consensused_block_height[nodeIndex - 1] = blockValue.value
		}
	}

	return failedNode
}

type BlockValue struct {
	key   string
	value int
}
