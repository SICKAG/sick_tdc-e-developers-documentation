/* Package created 22.09.2023. for the purposes of SICK Mobilisis d.o.o. */
/* Created to connect to a database and add / delete values from it */

/* Package was modified 04.10.2023. to fit MQTT needs */

package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

/* MQTT Data struct */
type MQTTData struct {
	Id                 int
	Altitude           float32
	Course             *string
	Fix                int
	GpsFixAvailable    bool
	Hdop               float32
	Latitude           float32
	Longitude          float32
	NumberOfSatellites int
	SpeedKnots         float32
	SpeedMph           float32
	Time               string
	Rssi               int
	DataLinkType       string
	Rfid               string
}

/* Function connects to a database using the specified dbms and connection string */
/* In case of not being able to connect, it prints the error received by the service. Otherwise, returns the connection that is used for other operations. */
/* Use with defer db.Close() */
func Connect(dbms string, connectionString string) *sql.DB {
	db, err := sql.Open(dbms, connectionString)
	if err != nil {
		fmt.Println("Couldn't connect to SQL database: ", err)
	}
	return db
}

/* Processes the given string command by connecting to the sql.DB object */
func ModifyDB(db *sql.DB, action string) {
	_, err := db.Exec(action)
	if err != nil {
		fmt.Println("Database error: ", err)
	}
	fmt.Println("Action successfully completed.")
}

/* Selects values from database specified in query */
func SelectValues(db *sql.DB, query string) []MQTTData {
	var results []MQTTData

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error finding rows: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var data MQTTData
		err := rows.Scan(
			&data.Id, &data.Altitude, &data.Course, &data.Fix, &data.GpsFixAvailable,
			&data.Hdop, &data.Latitude, &data.Longitude, &data.NumberOfSatellites,
			&data.SpeedKnots, &data.SpeedMph, &data.Time, &data.Rssi, &data.DataLinkType, &data.Rfid,
		)
		if err != nil {
			fmt.Println("Error fetching data: ", err)
		}
		results = append(results, data)
	}
	return results
}
