package helpers

import (
	"crypto/rand"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/fatih/structs"
	models "github.com/oluwapaso/hd_models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const YYYY_MM_DD__HHMMSS string = "2006-01-02 15:04:05"
const YYYY_MM_DD string = "2006-01-02"

func HandlePanic(via string) {
	if r := recover(); r != nil {
		fmt.Printf("\nRecovered from %s panic: %v", via, r)
	}
}

func ParseInt(val interface{}) int {
	intVal, _ := strconv.Atoi(fmt.Sprint(val))
	return intVal
}

func NumWithDot(value interface{}) string {
	pattern := regexp.MustCompile(`[^0-9.]`)
	return pattern.ReplaceAllString(fmt.Sprint(value), "")
}

func MappStructToScannedFields(rows *sql.Rows, columns []string) (map[string]interface{}, error) {

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i, _ := range values {
		pointers[i] = &values[i]
	}

	resultMap := make(map[string]interface{})

	err := rows.Scan(pointers...)
	if err != nil {
		return resultMap, err
	}

	for i, val := range values {
		resultMap[columns[i]] = val
	}

	return resultMap, nil

}

func ReadCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func ReadCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	csvReader := csv.NewReader(resp.Body)
	//csvReader.Comma = ';'
	csvReader.LazyQuotes = true

	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil

}

func ParseToDatetime(value string) string {

	layout, _ := dateparse.ParseFormat(value)
	date, _ := time.Parse(layout, value)
	ret_date := date.Format(YYYY_MM_DD__HHMMSS)

	return fmt.Sprint(ret_date)

}

func ParseToDate(value string) string {

	layout, _ := dateparse.ParseFormat(value)
	date, _ := time.Parse(layout, value)
	ret_date := date.Format(YYYY_MM_DD)

	return fmt.Sprint(ret_date)

}

func StringToTime(value string) int64 {

	layout, _ := dateparse.ParseFormat(value)
	date, _ := time.Parse(layout, value)
	timestamp := date.Unix()
	return timestamp

}

func Ucwords(input string) string {
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(input))
}

func RemoveNoneNumerics(value string) string {
	none_numeric := regexp.MustCompile(`[^0-9]+`)
	return none_numeric.ReplaceAllString(value, "")
}

func UsaPhoneNumber(value string) string {

	number := RemoveNoneNumerics(value)
	sn := strings.Split(number, "")
	var (
		part_1 string
		part_2 string
		part_3 string
	)
	for i, val := range sn {
		if i <= 2 {
			part_1 += val
		}

		if i > 2 && i <= 5 {
			part_2 += val
		}

		if i > 5 && i <= 9 {
			part_3 += val
		}
	}
	return fmt.Sprintf("%s-%s-%s", part_1, part_2, part_3)

}

func ParseVehicles(vehicles string) []models.Vehicles {
	/** Returns array **/
	var parsed_vehicles []models.Vehicles
	if vehicles != "" {

		eploded_vehicles := strings.Split(vehicles, ",")
		vehicle_id := 1

		for _, veh := range eploded_vehicles {

			veh = strings.TrimSpace(veh)
			var vehicleLine models.Vehicles

			if veh != "" {
				eachVeh := strings.Split(veh, " ")
				vehicle_year := strings.TrimSpace(eachVeh[0])
				vehicle_make := strings.TrimSpace(eachVeh[1])
				yr_mk := vehicle_year + " " + vehicle_make
				vehicle_model := strings.TrimSpace(strings.ReplaceAll(veh, yr_mk, ""))

				vehicleLine.Vehicle_id = ParseInt(vehicle_id)
				vehicleLine.Year = ParseInt(vehicle_year)
				vehicleLine.Make = Ucwords(vehicle_make)
				vehicleLine.Model = Ucwords(vehicle_model)
				vehicleLine.Full_vehicle = fmt.Sprintf("%d %s %s", ParseInt(vehicle_year), Ucwords(vehicle_make), Ucwords(vehicle_model))
				vehicle_id++

				parsed_vehicles = append(parsed_vehicles, vehicleLine)

			}

		}

		return parsed_vehicles

	} else {
		return parsed_vehicles
	}

}

func SubStr(str string, start, end int) string {
	return strings.TrimSpace(str[start:end])
}

func RemoveMultipleSpace(value string) string {
	space := regexp.MustCompile(`!\s+!`)
	result := space.ReplaceAllString(value, " ")
	result = strings.TrimSpace(result)
	return result
}

func ParseStateCityZip(value string) string {

	var (
		zipCode string = ""
		state   string = ""
		city    string = ""
	)

	if value != "" {

		zip_ptrn := regexp.MustCompile(`[0-9]{6}|[0-9]{5}`)
		zip_indices := zip_ptrn.FindAllStringIndex(value, 1)
		zipFound := len(zip_indices)
		if zipFound > 0 {
			zip_start := zip_indices[0][0]
			zip_end := zip_indices[0][1]
			zipCode = SubStr(value, zip_start, zip_end)
		}

		state_ptrn := regexp.MustCompile(`\s[A-Za-z]{2}\s|\s[A-Za-z]{2},|,[A-Za-z]{2},|,[A-Za-z]{2}\s,|,[A-Za-z]{2}\s|,\s[A-Za-z]{2}$|,[A-Za-z]{2}$`)
		state_indices := state_ptrn.FindAllStringIndex(value, 1)
		stateFound := len(state_indices)
		if stateFound > 0 {
			state_start := state_indices[0][0]
			state_end := state_indices[0][1]
			state = SubStr(value, state_start, state_end)
			state = strings.ReplaceAll(state, ",", "")
			state = RemoveMultipleSpace(state)
		}

		city = strings.ReplaceAll(value, zipCode, "")
		city = strings.ReplaceAll(city, state, "")
		city = strings.ReplaceAll(city, ",", "")
		city = RemoveMultipleSpace(city)

	}

	return fmt.Sprintf("%s, %s, %d", Ucwords(city), strings.ToUpper(state), ParseInt(zipCode))

}

// map[string]map[string]interface{}
// func ArrayColumn(input []models.Agents, columnKey string) []interface{} {

// 	var columns []interface{}
// 	columns = make([]interface{}, 0, len(input))
// 	for _, val := range input {
// 		mappedAgnts := structs.Map(val)
// 		fmt.Printf("\nValue: %v - Key: %s", mappedAgnts[columnKey], columnKey)
// 		if v, ok := mappedAgnts[columnKey]; ok {
// 			columns = append(columns, v)
// 		}
// 	}

// 	return columns

// }

type ArrColInterface interface {
	models.Agents | models.LeadSource | models.MappedImportStatus | models.MappedLeadSource | models.ImportDealHeader
}

func ArrayColumn[T ArrColInterface](input []T, columnKey string) []interface{} {

	var columns []interface{}
	columns = make([]interface{}, 0, len(input))
	for _, val := range input {
		to_map := structs.Map(val)
		if v, ok := to_map[columnKey]; ok {
			columns = append(columns, v)
		}
	}
	return columns

}

func ArraySearch(needle interface{}, hystack interface{}) (index int) {
	index = -1

	switch reflect.TypeOf(hystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(hystack)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				index = i
				return
			}
		}
	}
	return
}

func GetArrayKeyIndex[T ArrColInterface](value interface{}, array *[]T, column string) int {
	index := ArraySearch(value, ArrayColumn(*array, column))
	return index
}

func ValidateEmailAddress(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func ParseCSVField(row map[string]string, key string, dataHeaders []models.ImportDealHeader) string {

	var isSkipped string = "true"
	var isMatched string = "false"
	var fieldVal string = ""

	header_index := GetArrayKeyIndex(key, &dataHeaders, "Id")
	if header_index > -1 {
		isSkipped = dataHeaders[header_index].Skip
		isMatched = dataHeaders[header_index].Matched
	}

	if isSkipped == "false" && isMatched == "true" {
		fieldVal = row[key]
	}

	return strings.TrimSpace(fieldVal)

}
