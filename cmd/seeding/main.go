package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository/postgres"
)

const (
	dirPath = "./db/data"
)

var (
	folders = []string{
		"users", "partners", "provinces", "cities", "districts", "sub_districts",
		"pharmacies", "categories", "manufacturers", "product_classifications",
		"product_forms", "products", "logistics", "pharmacy_logistics", "pharmacy_products", "pharmacist_details", "others",
	}
)

func main() {
	config := config.InitConfig()
	logger.SetLogrusLogger(config.App)

	db, err := postgres.Init(config)
	if err != nil {
		logger.Log.Fatal("failed to connect db")
	}
	defer db.Close()

	for _, folder := range folders {
		path := fmt.Sprintf("./db/data/%s", folder)
		entries, err := os.ReadDir(path)
		if err != nil {
			logger.Log.Fatal("failed to read folder")
		}
		logger.Log.Infof("Seeding %s table", folder)
		seeding(db, len(entries), folder)
	}
	logger.Log.Info("finished seeding")
}

func seeding(db *sql.DB, steps int, data string) {
	for i := 1; i < steps+1; i++ {
		logger.Log.Infof("Steps %d...", i)
		path := filepath.Join(dirPath, data, fmt.Sprintf("%d_%s.sql", i, data))
		c, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
			logger.Log.Fatal("failed to read SQL file")
		}

		sql := string(c)
		// rewrite(i, sql, data)
		_, err = db.Exec(sql)
		if err != nil {
			log.Println(err)
			logger.Log.Fatal("failed to exec SQL file")
		}
	}
}

// func rewrite(steps int, sql, dataFolder string) string {
// 	arr := strings.Split(sql, "\n")
// 	head := arr[0]
// 	leftHead := head[:len(arr[0])-9]
// 	rightHead := head[len(arr[0])-9:]
// 	arr[0] = leftHead + ",created_at,updated_at" + rightHead

// 	for i := 1; i < len(arr)-1; i++ {
// 		left := arr[i][:len(arr[i])-2]
// 		right := arr[i][len(arr[i])-2:]
// 		date := timeStamp()
// 		arr[i] = left + ",'" + date + "'" + ",'" + date + "'" + right
// 	}

// 	newSql := strings.Join(arr, "\n")
// 	f, err := os.Create(fmt.Sprintf("%s/%s/%d_%s.sql", dirPath, dataFolder, steps, dataFolder))
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		return ""
// 	}
// 	defer f.Close()

// 	_, err = f.WriteString(newSql)
// 	if err != nil {
// 		fmt.Println("Error writing to file:", err)
// 		return ""
// 	}

// 	fmt.Println("File written successfully")

// 	return newSql
// }

// func timeStamp() string {
// 	min := time.Date(2010, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
// 	max := time.Date(2024, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
// 	delta := max - min

// 	sec := rand.Int63n(delta) + min
// 	return time.Unix(sec, 0).Format(time.RFC3339)
// }
