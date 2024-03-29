package helpers

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
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
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/structs"
	models "github.com/oluwapaso/hd_models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"mvdan.cc/xurls/v2"
)

const APP_NAME = "Hauling Desk"
const BASE_URL = "https://local.haulingdesk.com"                                         //https://haulingdesk.com
const CRM_URL = "https://yaimalamela.com"                                                //http://local.muvcars.com //https://crm.haulingdesk.com //https://yaimalamela.com
const API_URL = "http://rqeejaczw2.execute-api.us-east-1.amazonaws.com/HaulingDeskGoAPI" //http://local.muvcars.com/api/v1/router
const DEVELOPERS_LINK = "https://main.d3bzogxvbyj700.amplifyapp.com/"
const EMAIL_TRACKER_LINK = BASE_URL + "/track-open-rates.php"
const EMAIL_CLICK_TRACKER_LINK = BASE_URL + "/track-click-rates.php"
const YYYY_MM_DD__HHMMSS string = "2006-01-02 15:04:05"
const YYYY_MM_DD string = "2006-01-02"
const HH_MM_SS string = "15:04:05"
const MM_DD_YYYY string = "01-02-2006"
const DD_MM_YYYY string = "02-01-2006"
const MM_DD_YYYY__gi_A string = "01-02-2006 3:4 PM"
const F_d_Y string = "January 02, 2006"
const AUTOMATED_EMAIL_LOGO = "https://cronetic.com/crm/img/powered-by-v2.png"
const COMP_ADDRESS = "7807 W Loop 1604 N San Antonio, TX 78254"
const SHIPPERS_PORTAL_URL = "https://shippers.hauling-desk.com" //https://shippers.hauling-desk.com //https://shippers.haulingdesk.com
const CARRIERS_PORTAL_URL = "https://carriers.hauling-desk.com" //https://carriers.hauling-desk.com //https://carriers.haulingdesk.com
const HD_SENDGRID_KEY = "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"
const HD_NOTIFICATIONS_EMAIL = "notifications@haulingdesk.com"
const SUPPORT_EMAIL = "support@haulingdesk.com"
const SUPPORT_PHONE = "+2348062744512"
const ACCOUNTS_EMAIL = "accounts@haulingdesk.com"
const GIT_TOKEN = "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"

const LOGO_BUCKETS = "hauling-desk-logos"
const AGENTS_DP_BUCKETS = "hauling-desk-agents-dp"
const DEAL_FILES_BUCKETS = "hauling-desk-deal-files"
const CONTACTS_FILES_BUCKETS = "hauling-desk-contacts-files"
const CSV_FILES_BUCKETS = "hauling-csv-files"

const TWILIO_ACCOUNT_SID = "AC9c7d25d759ec866710c2b38245c54893"
const TWILIO_AUTH_TOKEN = "d3f7c21920ca9635a2ec7ee157088997"
const TWILIO_TWIML_SID = "APd7fda7e62f8312c6db3cf71ea3c470f1"
const TWILIO_CALL_PER_MIN = 1 //Twilio: $0.0040/Min = 250 min in $1 => HD $1 = 125 Min
const TWILIO_SMS_PER_PAGE = 2 //Twilio: $0.0079/Page = 126 pages in $1 => HD $1 = 63 Pages

const IP_ADDRESS = "209.182.198.166"           //yailamela(209.182.198.166)
const CPANEL_HOME_DIRECTORY = "/home/profe160" //yailamela(/home/profe160)
const CPANEL_FULL_PATH = "../"                 //yailamela(/home2/profe160/public_html/crm) //../ (Local)
const LEAD_PIPING_SCRIPT = "/home/profe160/public_html/cron_jobs/cron_forwarder.php"
const WEBMAIL_PIPING_SCRIPT = "/home/profe160/public_html/cron_jobs/incoming_email_parser.php"

const CPANEL_USERNAME = "profe160"     //yailamela(profe160)
const CPANEL_PASSWORD = "7i(&V)+BO)G+" //yailamela(7i(&V)+BO)G+)
const EMAIL_DOMAIN = "yaimalamela.com"
const XML_API_PORT = "2083"

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

	var delimeter string
	if strings.Contains(date, "/") {
		delimeter = "/"
	} else if strings.Contains(date, "-") {
		delimeter = "-"
	} else {
		return ""
	}

	//date = strings.ReplaceAll(date, "/", "-")
	splitSpace := strings.Split(date, " ")
	val := strings.Split(splitSpace[0], delimeter)

	if len(val) < 1 {
		return ""
	}

	var (
		format string
	)

	mm_dd_ptrn := regexp.MustCompile(`[0-9]{2}` + delimeter + `[0-9]{2}` + delimeter + `[0-9]{4}`)
	mm_dd_indices := mm_dd_ptrn.FindAllStringIndex(date, 1)
	mm_dd_found := len(mm_dd_indices)
	if mm_dd_found > 0 {

		val_1 := ParseInt(val[0]) //first value

		if val_1 <= 12 {
			//mm-dd-yyyy
			format = "01" + delimeter + "02" + delimeter + "2006"
		} else {
			//dd-mm-yyyy
			format = "02" + delimeter + "01" + delimeter + "2006"
		}

	}

	yy_mm_dd_ptrn := regexp.MustCompile(`[0-9]{4}` + delimeter + `[0-9]{2}` + delimeter + `[0-9]{2}`)
	yy_mm_dd_indices := yy_mm_dd_ptrn.FindAllStringIndex(date, 1)
	yy_mm_dd_found := len(yy_mm_dd_indices)
	if yy_mm_dd_found > 0 {

		val_2 := ParseInt(val[1]) //second value

		if val_2 <= 12 {
			//yyyy-mm-dd
			format = "2006" + delimeter + "01" + delimeter + "02"
		} else {
			//yyyy-dd-mm
			format = "2006" + delimeter + "02" + delimeter + "01"
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

func Parse_Date_To_A_Format(value, format string) string {

	layout, _ := dateparse.ParseFormat(value)
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
	} else if format == "H:i:s" {
		date_format = HH_MM_SS
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

func DateInLoc(format string, loc *time.Location) string {

	date_format := YYYY_MM_DD__HHMMSS
	if format == "YYYY-MM-DD H:i:s" {
		date_format = YYYY_MM_DD__HHMMSS
	} else if format == "YYYY-MM-DD" {
		date_format = YYYY_MM_DD
	} else if format == "H:i:s" {
		date_format = HH_MM_SS
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

	date := time.Now().In(loc).Format(date_format)

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

func GetSpecificDateInLoc(format string, offset int, loc *time.Location) string {

	date_format := YYYY_MM_DD__HHMMSS
	if format == "YYYY-MM-DD H:i:s" {
		date_format = YYYY_MM_DD__HHMMSS
	} else if format == "YYYY-MM-DD" {
		date_format = YYYY_MM_DD
	} else if format == "MM-DD-YYYY" {
		date_format = "01-02-2006"
	}

	date := time.Now().In(loc).AddDate(0, 0, offset).Format(date_format)

	return fmt.Sprint(date)

}

func GetDateOffsetInLoc(format string, offset int, loc *time.Location, date time.Time) string {

	date_format := YYYY_MM_DD__HHMMSS
	if format == "YYYY-MM-DD H:i:s" {
		date_format = YYYY_MM_DD__HHMMSS
	} else if format == "YYYY-MM-DD" {
		date_format = YYYY_MM_DD
	} else if format == "MM-DD-YYYY" {
		date_format = "01-02-2006"
	}

	new_date := date.In(loc).AddDate(0, 0, offset).Format(date_format)

	return new_date

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

func RemoveNoneAlphabets(value string) string {
	alphabets := regexp.MustCompile(`[^A-Za-z]+`)
	return alphabets.ReplaceAllString(value, "")
}

func PregReplace(value, pattern, replace_with string) string {
	none_numeric := regexp.MustCompile(pattern)
	return none_numeric.ReplaceAllString(value, replace_with)
}

func PregMatch(value, pattern string) string {

	// Compile the regular expression
	re := regexp.MustCompile(pattern)
	// Find the match in the input string
	match := re.FindStringSubmatch(value)
	// Check if there is a match
	if len(match) > 0 {
		return match[0]
	} else {
		return ""
	}

}

func GetTextBetween(input, start, end string) (string, error) {
	// Create the regular expression pattern
	pattern := fmt.Sprintf("%s(.*?)%s", regexp.QuoteMeta(start), regexp.QuoteMeta(end))

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find the match in the input string
	match := re.FindStringSubmatch(input)

	if len(match) >= 2 {
		// Return the text between the start and end strings
		return match[1], nil
	}

	// Return an error if no match is found
	return "", fmt.Errorf("no match found between '%s' and '%s'", start, end)
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

func ParsePhoneForSMS(number string) string {

	number = UnformattedPhoneNumber(number)
	if number != "" {
		number = "1" + number
	}

	return number

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
	models.Agents | models.LeadSource | models.MappedImportStatus | models.MappedLeadSource | models.ImportDealHeader | models.Vehicles |
		models.Issues
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

func StringWithCharset(length int, charset string, seededRand *math_rand.Rand) string {
	math_rand.Seed(time.Now().UnixNano())
	seededRand = math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomChars(length int) string {
	seededRand := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	const charset = "ABCDKLM028983NOEF2661982GHIJ8272PQTUV128882WXYRSZA89283BCD484EFGHI033JKLMNOPQRSTUVWXYZ0123456789"
	return StringWithCharset(length, charset, seededRand)
}

func RandomInts(length int) string {
	seededRand := math_rand.New(math_rand.NewSource(time.Now().UnixNano()))
	const charset = "92-391729812-8387378309833-97209-37923-38-182273877-9238732-9281292391" +
		"20423636838-9309893323-7376389347-467849474"
	return StringWithCharset(length, charset, seededRand)
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
	dat := map[string]interface{}{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

func Json_encode_decode(data interface{}) (map[string]interface{}, error) {
	jsons, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	dat := map[string]interface{}{}
	err = json.Unmarshal(jsons, &dat)
	return dat, err
}

func Json_decode_array(data string) []interface{} {
	jsonArray := []interface{}{}
	json.Unmarshal([]byte(data), &jsonArray)
	return jsonArray
}

func String_To_Array(data string) ([]interface{}, error) {
	dat := []interface{}{}
	err := json.Unmarshal([]byte(data), &dat)
	return dat, err
}

func Array_To_String(array []interface{}, delimeter string) string {

	result := ""
	for _, arr := range array {
		result += fmt.Sprint(arr) + "" + delimeter
	}

	result = strings.TrimRight(result, delimeter)
	return result
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

func JsonArrayColValue(scanned_val []interface{}, index int) interface{} {

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
		json_output = Json_decode_array(output)
	}

	return json_output

}

func ValToEmptyJson(val string) string {
	if val == "" {
		val = "{}"
	}
	return val
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

func BuildSingleSelectColumns(fields []interface{}, model_fields string) string {

	var fieldVal string
	var query_fields string
	for _, column := range fields {
		if column != "" && column != "*" {
			/** OrdersField is inside table_columns.go **/
			column := strings.TrimSpace(fmt.Sprint(column))
			exploded_orders_fld := strings.Split(string(model_fields), ",")
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

func MultiTableSelCols(joiner, columns string) string {

	var output string
	split := strings.Split(columns, ",")
	for _, val := range split {
		val = strings.TrimSpace(val)
		if val != "" {
			output += joiner + "." + val + ","
		}
	}
	output = strings.TrimRight(output, ",")

	return output

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

func ArrayUniqueIntfc(arr []interface{}) []interface{} {
	size := len(arr)
	result := make([]interface{}, 0, size)
	temp := map[interface{}]struct{}{}
	for i := 0; i < size; i++ {
		if _, ok := temp[arr[i]]; ok != true {
			temp[arr[i]] = struct{}{}
			result = append(result, arr[i])
		}
	}
	return result
}

func RemoveArrayByKeyValue(json_val, key, value string) string {

	var output []interface{}
	var jsonData []interface{}
	json.Unmarshal([]byte(json_val), &jsonData)

	for _, line := range jsonData {
		line_map, ok := line.(map[string]interface{})
		if ok {
			if line_map[key] != value {
				output = append(output, line_map)
			}
		}
	}

	outputJson, _ := Json_encode(output)
	return outputJson

}

func RemoveArrayByIndex(array []interface{}, pos int) []interface{} {
	return append(array[:pos], array[pos+1:]...)
}

func ArrayEnd(array []interface{}) interface{} {
	return array[len(array)-1]
}

func CountJson(json string) int {
	jsonArray := Json_decode_array(json)
	count := len(jsonArray)
	return count
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
	if value == "<nil>" || value == "" || value == "null" {
		value = ""
	}
	return value
}

func NilToEmptyJson(value string) string {
	if value == "<nil>" || value == "" || value == "null" {
		value = "{}"
	}
	return value
}

func NilToEmptyArray(value string) string {
	if value == "<nil>" || value == "" || value == "null" {
		value = "[]"
	}
	return value
}

func NilToFloat(value interface{}) float64 {

	var val string
	var output float64

	if value == "<nil>" || value == "" || value == nil {
		output = 0
	} else {
		val = fmt.Sprint(value)
		output, _ = strconv.ParseFloat(val, 64)
	}

	return output

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
	to_use := NilToFloat(value)
	em := message.NewPrinter(language.English)
	var enNumber string = em.Sprint(to_use)
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

func Implode(array []interface{}, delim string) string {

	var output string
	for _, val := range array {
		output += fmt.Sprint(val) + "" + delim
	}

	output = strings.TrimRight(output, delim)

	return output

}

func ImplodeStrings(array []string, delim string) string {

	var output string
	for _, val := range array {
		output += fmt.Sprint(val) + "" + delim
	}

	output = strings.TrimRight(output, delim)

	return output

}

func StringToInterfaceSlice(items, delimeter string) []interface{} {

	var output []interface{}
	exploded_val := strings.Split(items, delimeter)

	for _, val := range exploded_val {
		val = strings.TrimSpace(val)
		if val != "" {
			output = append(output, val)
		}
	}

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
		field = strings.TrimSpace(field)
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
	lent_of_vals := len(expectedValues)
	for _, val := range expectedValues {
		if val != "" {
			expected_vals += val + ", "
		}
	}

	expected_vals = strings.TrimRight(expected_vals, ", ")
	if lent_of_vals > 1 {
		expected_vals = ReplaceLast(expected_vals, ",", " or")
	}
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

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidateValueLength(value string, min int, max int, feild string) []string {

	if min > max {
		max = min + 1
		/** Making sure max is always largest **/
	}

	var respMsg []string
	if len(value) < min || len(value) > max {

		if len(value) < min {
			respMsg = append(respMsg, feild+" can't be less than "+fmt.Sprint(min)+" characters.")
		}

		if len(value) > max {
			respMsg = append(respMsg, feild+" can't be more than "+fmt.Sprint(max)+" characters.")
		}
	}

	return respMsg
}

func ConcatMultipleSlices[T any](slices [][]T) []T {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]T, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s) //Copy everything in source(s) into current position(i) to the end(:)
	}

	return result
}

func ConcatMultipleMaps(slice_of_maps []interface{}) map[string]interface{} {

	newMap := map[string]interface{}{}
	for _, s := range slice_of_maps {
		map_line, ok := s.(map[string]interface{})
		if ok {
			for map_key, map_val := range map_line {
				newMap[map_key] = map_val
			}
		}
	}

	return newMap

}

func ReplaceCommonSysCodes(temp string) string {

	if temp != "" {
		temp = strings.ReplaceAll(temp, "{{app_name}}", APP_NAME)
		temp = strings.ReplaceAll(temp, "{{base_url}}", BASE_URL)
		temp = strings.ReplaceAll(temp, "{{crm_url}}", CRM_URL)
		temp = strings.ReplaceAll(temp, "{{AUTOMATED_EMAIL_LOGO}}", APP_NAME)
		temp = strings.ReplaceAll(temp, "{{comp_address}}", COMP_ADDRESS)
	}
	return temp

}

func GetAllVehiclesTotalNum(vehicles []interface{}) int {

	total := 0
	for _, veh := range vehicles {
		veh_line, ok := veh.(map[string]interface{})
		if ok {
			quantity := ParseInt(veh_line["quantity"])
			if quantity <= 0 {
				quantity = 1
			}

			total += quantity
		}
	}

	return total

}

func GetCarType(car string) string {

	car = strings.ToLower(car)
	car = PregReplace(car, "[^a-z]", " ")

	car_pat := regexp.MustCompile(`car|coupe|convertible|sedan|sedan small|sedan midsize|sedan large`)
	car_incs := car_pat.FindAllStringIndex(car, 1)
	car_found := len(car_incs)
	if car_found > 0 {
		return "Car"
	}

	pickup_pat := regexp.MustCompile(`dually|pickup|pick up|pickup small|pickup crew cab|pickup fullsize|pickup full size|Pickup Ext Cab`)
	pickup_incs := pickup_pat.FindAllStringIndex(car, 1)
	pickup_found := len(pickup_incs)
	if pickup_found > 0 {
		return "Pickup"
	}

	suv_pat := regexp.MustCompile(`jeep|suv|suv small|suv mid size|suv midsize|suv large`)
	suv_incs := suv_pat.FindAllStringIndex(car, 1)
	suv_found := len(suv_incs)
	if suv_found > 0 {
		return "SUV"
	}

	van_pat := regexp.MustCompile(`van|van mini|van full size|van fullsize|van extd length|van pop top`)
	van_incs := van_pat.FindAllStringIndex(car, 1)
	van_found := len(van_incs)
	if van_found > 0 {
		return "Van"
	}

	rv_pat := regexp.MustCompile(`rv`)
	rv_incs := rv_pat.FindAllStringIndex(car, 1)
	rv_found := len(rv_incs)
	if rv_found > 0 {
		return "RV"
	}

	trailer_pat := regexp.MustCompile(`travel|trailer|travel trailer|rv trailer`)
	trailer_incs := trailer_pat.FindAllStringIndex(car, 1)
	trailer_found := len(trailer_incs)
	if trailer_found > 0 {
		return "Travel Trailer"
	}

	mc_pat := regexp.MustCompile(`motorcycle`)
	mc_incs := mc_pat.FindAllStringIndex(car, 1)
	mc_found := len(mc_incs)
	if mc_found > 0 {
		return "Motorcycle"
	}

	boat_pat := regexp.MustCompile(`boat`)
	boat_incs := boat_pat.FindAllStringIndex(car, 1)
	boat_found := len(boat_incs)
	if boat_found > 0 {
		return "Boat"
	}

	atv_pat := regexp.MustCompile(`atv|atv utv`)
	atv_incs := atv_pat.FindAllStringIndex(car, 1)
	atv_found := len(atv_incs)
	if atv_found > 0 {
		return "ATV"
	}

	he_pat := regexp.MustCompile(`heavy equipment`)
	he_incs := he_pat.FindAllStringIndex(car, 1)
	he_found := len(he_incs)
	if he_found > 0 {
		return "Heavy Equipment"
	}

	yatch_pat := regexp.MustCompile(`large yacht|yacht`)
	yatch_incs := yatch_pat.FindAllStringIndex(car, 1)
	yatch_found := len(yatch_incs)
	if yatch_found > 0 {
		return "Large Yacht"
	}

	other_pat := regexp.MustCompile(`other`)
	other_incs := other_pat.FindAllStringIndex(car, 1)
	other_found := len(other_incs)
	if other_found > 0 {
		return "Other"
	}

	return car

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func PutPresignURL(cfg aws.Config, bucket string, file_loc string, ftype string) (string, error) {
	s3client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3client)
	presignedUrl, err := presignClient.PresignPutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(file_loc),
			ContentType: &ftype,
			ACL:         "public-read",
		},
		s3.WithPresignExpires(time.Minute*15),
	)

	if err != nil {
		return "Error", err
	}

	return presignedUrl.URL, nil
}

func BuildVehicles(vehicles []interface{}) []interface{} {
	parsed_vehicles := []interface{}{}

	for i, vehicle := range vehicles {
		veh, ok := vehicle.(map[string]interface{})
		if ok {

			qty := ParseInt(veh["quantity"])
			if qty < 1 {
				qty = 1
			}

			full_vehicle := fmt.Sprint(veh["year"]) + ` ` + Ucwords(fmt.Sprint(veh["make"])) + ` ` + Ucwords(fmt.Sprint(veh["model"]))
			vehicleLine := map[string]interface{}{
				"vehicle_id":   i,
				"year":         ParseInt(veh["year"]),
				"make":         Ucwords(fmt.Sprint(veh["make"])),
				"model":        Ucwords(fmt.Sprint(veh["model"])),
				"quantity":     qty,
				"type":         GetCarType(fmt.Sprint(veh["type"])),
				"running":      fmt.Sprint(veh["running"]),
				"ship_via":     fmt.Sprint(veh["ship_via"]),
				"full_vehicle": full_vehicle,
			}

			parsed_vehicles = append(parsed_vehicles, vehicleLine)
		}
	}

	return parsed_vehicles

}

func GetAllLinks(content string) []string {

	resp := []string{}
	links := xurls.Strict().FindAllString(content, -1)
	resp = ArrayUnique(links)
	return resp

}

func GetAllAnchor(content string) []string {

	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var f func(*html.Node)
	links := []string{}
	f = func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if a.Val != "" {
						links = append(links, a.Val)
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	links = ArrayUnique(links)
	return links

}

func FindFloa64Max(arr []float64) float64 {
	if len(arr) == 0 {
		// Handle the case when the array is empty
		return 0
	}

	max := arr[0]
	for _, value := range arr {
		if value > max {
			max = value
		}
	}

	return max
}

func FindFLoat64Min(arr []float64) float64 {
	if len(arr) == 0 {
		// Handle the case when the array is empty
		return 0
	}

	min := arr[0]
	for _, value := range arr {
		if value < min {
			min = value
		}
	}

	return min
}

func InvokeLocalSamFunction(func_name string, json_event string) string {

	//command := "echo " + json_event + " | sam local invoke " + func_name + " -e -"
	//command = `echo {"body": "{\"message\": \"hello mr elpaso\"}"} | sam local invoke -e -`
	//command = "./sam local invoke"
	//commands := []string{"echo", `{"body": "{\"message\": \"hello mr elpaso\"}"}`, "|", "sam", "local", "invoke", "-e", "-"}

	// Set up the command
	cmd := exec.Command("sam", "local", "invoke", func_name, "--no-event", "-e", "post-event.json")

	// Set working directory to the SAM application directory
	cmd.Dir = `C:\xampp\apps\sam-example\test-app`

	// Set the environment variables if needed
	cmd.Env = os.Environ()

	// Buffer to store the standard output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return ""
	}

	// Extract response from the output
	response := out.String()
	return response

	// Run the command and capture output
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return "Error invoking function:" + fmt.Sprint(err)
	// }

	// return string(output)

}
