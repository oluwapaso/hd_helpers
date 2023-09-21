package helpers

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/rand"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	math_rand "math/rand"
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
	"golang.org/x/text/message"
)

const YYYY_MM_DD__HHMMSS string = "2006-01-02 15:04:05"
const YYYY_MM_DD string = "2006-01-02"
const MM_DD_YYYY string = "01-02-2006"
const DD_MM_YYYY string = "02-01-2006"
const MM_DD_YYYY__gi_A string = "01-02-2006 3:4 PM"
const F_d_Y string = "January 02, 2006"
const SHIPPERS_PORTAL_URL = "https://shippers.hauling-desk.com" //https://shippers.hauling-desk.com //https://shippers.haulingdesk.com
const CARRIERS_PORTAL_URL = "https://carriers.hauling-desk.com" //https://carriers.hauling-desk.com //https://carriers.haulingdesk.com

func HandlePanic(via string) {
	if r := recover(); r != nil {
		fmt.Printf("\nRecovered from %s \npanic: %v", via, r)
		//fmt.Printf("\nRecovered from %s \npanic: %v \nstack trace: %s", via, r, string(debug.Stack()))
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

func GetDateLayout(date string) string {

	date = strings.ReplaceAll(date, "/", "-")
	splitSpace := strings.Split(date, " ")
	val := strings.Split(splitSpace[0], "-")

	if len(val) < 1 {
		return ""
	}

	var (
		format string
	)

	mm_dd_ptrn := regexp.MustCompile(`[0-9]{2}-[0-9]{2}-[0-9]{4}`)
	mm_dd_indices := mm_dd_ptrn.FindAllStringIndex(date, 1)
	mm_dd_found := len(mm_dd_indices)
	if mm_dd_found > 0 {

		val_1 := ParseInt(val[0]) //first value

		if val_1 <= 12 {
			//mm-dd-yyyy
			format = "01-02-2006"
		} else {
			//dd-mm-yyyy
			format = "02-01-2006"
		}

	}

	yy_mm_dd_ptrn := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
	yy_mm_dd_indices := yy_mm_dd_ptrn.FindAllStringIndex(date, 1)
	yy_mm_dd_found := len(yy_mm_dd_indices)
	if yy_mm_dd_found > 0 {

		val_2 := ParseInt(val[1]) //second value

		if val_2 <= 12 {
			//yyyy-mm-dd
			format = "2006-01-02"
		} else {
			//yyyy-dd-mm
			format = "2006-02-01"
		}

	}

	//Time part
	if len(splitSpace) > 1 {

		hh_mm_ss_ptrn := regexp.MustCompile(`[0-9]{2}:[0-9]{2}:[0-9]{2}`)
		hh_mm_ss_indices := hh_mm_ss_ptrn.FindAllStringIndex(date, 1)
		hh_mm_ss_found := len(hh_mm_ss_indices)
		if hh_mm_ss_found > 0 {
			format += " 15:04:05"
		}

	}

	return format

}

func ParseDateToFormat(value string, format string, layout string) string {

	date, _ := time.Parse(layout, value)
	ret_date := date.Format(format)

	return fmt.Sprint(ret_date)

}

func Date(format string) string {

	date_format := YYYY_MM_DD__HHMMSS
	if format == "YYYY-MM-DD H:i:s" {
		date_format = YYYY_MM_DD__HHMMSS
	} else if format == "YYYY-MM-DD" {
		date_format = YYYY_MM_DD
	} else if format == "F d, Y" {
		date_format = "January 02, 2006"
	} else if format == "M DD, YY" {
		date_format = "Jan 02, 2006"
	} else if format == "MM-DD-YYYY" {
		date_format = "01-02-2006"
	} else if format == "YYYY" {
		date_format = "2006"
	} else if format == "MM" {
		date_format = "01"
	} else if format == "DD" {
		date_format = "02"
	}

	date := time.Now().Format(date_format)

	return fmt.Sprint(date)

}

func GetSpecificDate(format string, offset int) string {

	date_format := YYYY_MM_DD__HHMMSS
	if format == "YYYY-MM-DD H:i:s" {
		date_format = YYYY_MM_DD__HHMMSS
	} else if format == "YYYY-MM-DD" {
		date_format = YYYY_MM_DD
	} else if format == "MM-DD-YYYY" {
		date_format = "01-02-2006"
	}

	date := time.Now().AddDate(0, 0, offset).Format(date_format)

	return fmt.Sprint(date)

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
	return fmt.Sprintf("(%s) %s-%s", part_1, part_2, part_3)

}

func LastNumOfCharacters(val string, length int) string {

	var result string
	if len(val) >= length {
		result = val[len(val)-length:]
	} else {
		result = val
	}

	return result

}

func UnformattedPhoneNumber(number string) string {

	clean := RemoveNoneNumerics(number)
	phn_number := LastNumOfCharacters(clean, 10)

	return fmt.Sprint(phn_number)

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

func StringWithCharset(length int, charset string) string {
	seededRand := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomChars(length int) string {
	const charset = "ABCDKLM028983NOEF2661982GHIJ8272PQTUV128882WXYRSZA89283BCD484EFGHI033JKLMNOPQRSTUVWXYZ0123456789"
	return StringWithCharset(length, charset)
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

func GzInflate(val []byte) string {

	reader := bytes.NewReader(val)

	gzreader, e1 := gzip.NewReader(reader)
	if e1 != nil {
		fmt.Println(e1) // Maybe panic here, depends on your error handling.
	}
	fmt.Println(gzreader)

	output, e2 := io.ReadAll(gzreader)
	if e2 != nil {
		fmt.Println(e2)
	}

	return string(output)

}

func ReadGzip(content []byte) error {
	var buf *bytes.Buffer = bytes.NewBuffer(content)
	fmt.Printf(fmt.Sprint(buf))
	gRead, err := zlib.NewReader(buf)
	if err != nil {
		return err
	}

	if t, err := io.Copy(os.Stdout, gRead); err != nil {
		fmt.Println(t)
		return err
	}

	if err := gRead.Close(); err != nil {
		return err
	}
	return nil
}

func Json_encode(data interface{}) (string, error) {
	jsons, err := json.Marshal(data)
	return string(jsons), err
}

func Json_decode(data string) (map[string]interface{}, error) {
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

func String_To_Array(data string) ([]interface{}, error) {
	var dat []interface{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

func IsJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func IsArray(val interface{}) bool {
	_, ok := val.([]interface{})
	if !ok {
		return false
	}
	return true
}

func ScanMultiRow(columns []string, rows *sql.Rows) ([]interface{}, error) {

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for i := range columns {
		valuePtrs[i] = &values[i] //set to memory address of values
	}

	err := rows.Scan(valuePtrs...)

	return values, err

}

func ColValue(scanned_val []interface{}, index int) string {

	val := scanned_val[index]
	b, ok := val.([]byte)

	var value interface{}

	if ok {
		value = string(b)
	} else {
		value = val
	}

	output := fmt.Sprint(value)
	if output == "<nil>" {
		output = ""
	}

	return output

}

func JsonColValue(scanned_val []interface{}, index int) interface{} {

	val := scanned_val[index]
	b, ok := val.([]byte)

	var value interface{}

	if ok {
		value = string(b)
	} else {
		value = val
	}

	output := fmt.Sprint(value)
	if output == "<nil>" {
		output = ""
	}

	var json_output map[string]interface{}
	if output != "" {
		json_output, _ = Json_decode(output)
	}

	return json_output

}

func ArrayColValue(scanned_val []interface{}, index int) interface{} {

	val := scanned_val[index]
	b, ok := val.([]byte)

	var value interface{}

	if ok {
		value = string(b)
	} else {
		value = val
	}

	output := fmt.Sprint(value)
	if output == "<nil>" {
		output = ""
	}

	var json_output []interface{}
	if output != "" {
		json_output, _ = String_To_Array(output)
	}

	return json_output

}

func BuildSingleSelectColumns(fields []interface{}) string {

	var fieldVal string
	var query_fields string
	for _, column := range fields {
		if column != "" && column != "*" {
			/** OrdersField is inside table_columns.go **/
			column := strings.TrimSpace(fmt.Sprint(column))
			exploded_orders_fld := strings.Split(string(models.OrdersField), ",")
			if In_array(column, exploded_orders_fld) {
				fieldVal += column + ","
			}
		}
	}

	field_val := ArrayUnique(strings.Split(fieldVal, ","))
	for _, val := range field_val {
		if val != "" && val != "*" {
			query_fields += val + ","
		}
	}

	query_fields = strings.TrimRight(query_fields, ",")
	return query_fields

}

func In_array(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func ArrayUnique(arr []string) []string {
	size := len(arr)
	result := make([]string, 0, size)
	temp := map[string]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; ok != true {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}

func ThreadError(err error, thread_id int) error {

	err_resp := models.LambdaError{
		Message:   fmt.Sprint(err),
		Thread_Id: thread_id,
	}

	jsonErr, _ := json.Marshal(err_resp)

	return errors.New(string(jsonErr))

}

func DelEmailSrcErr(err error, source_email string) error {

	err_resp := models.Delete_Src_Emails_Error{
		Message:      fmt.Sprint(err),
		Source_Email: source_email,
	}

	jsonErr, _ := json.Marshal(err_resp)

	return errors.New(string(jsonErr))

}

func Mysql_IN_Builder(values string) [2]string {

	split_vals := strings.Split(values, ",")
	join_vals := strings.Join(split_vals, `","`)
	join_vals = `"` + join_vals + `"`

	plc_holders := ""
	for _, _ = range split_vals {
		plc_holders += "?,"
	}
	plc_holders = strings.TrimRight(plc_holders, ",")

	return [2]string{
		join_vals,
		plc_holders,
	}

}

func NilToEmptyString(value string) string {
	if value == "<nil>" {
		value = ""
	}
	return value
}

func NilToEmptyJson(value string) string {
	if value == "<nil>" || value == "" {
		value = "{}"
	}
	return value
}

func ReplaceNilWith(hay, niddle string) string {
	value := hay
	if hay == "<nil>" {
		value = niddle
	}
	return value
}

func NumberFormat(value string) string {

	// English currency formatting
	em := message.NewPrinter(language.English)
	var enNumber string = em.Sprint(value)
	return enNumber

}

func IntArrayToStrJoin(ints []int, delim string) string {

	var output string
	for _, val := range ints {
		output += fmt.Sprint(val) + "" + delim
	}

	output = strings.TrimRight(output, delim)

	return output

}

func ReplaceLast(hay, niddle, val string) (x2 string) {
	i := strings.LastIndex(hay, niddle)
	excludingLast := hay[:i] + strings.Replace(hay[i:], niddle, val, 1)
	return excludingLast
}

func Send_API_Response(data models.Lambda_API_Response) string {
	jsons, _ := json.Marshal(data)
	return string(jsons)
}

func ValidateRequiredFileds(JsonData map[string]interface{}, fields []string) []string {

	var missingField []string
	for _, field := range fields {
		if field != "" {
			_, ok := JsonData[field]
			if !ok {
				missingField = append(missingField, field+" is not supplied")
			}
		}
	}

	return missingField

}

func ValidateFieldValues(fieldVal string, expectedValues []string, fieldName string) string {

	expected_vals := ""
	for _, val := range expectedValues {
		if val != "" {
			expected_vals += val + ", "
		}
	}

	expected_vals = strings.TrimRight(expected_vals, ", ")
	expected_vals = ReplaceLast(expected_vals, ",", " or")
	invalidValue := "Expected values for " + fieldName + " is " + expected_vals + ", `" + fieldVal + "` was sent"

	for _, val := range expectedValues {
		if val != "" {
			if fieldVal == val {
				/** One matched **/
				invalidValue = "Valid"
			}
		}
	}

	return invalidValue

}

func ConcatMultipleSlices[T any](slices [][]T) []T {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s) //Copy everything is source(s) into current position(i) to the end(:)
	}

	return result
}
