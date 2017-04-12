package umb_lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
	"regexp"
)

const (
	default_filename = "./config.txt"
)

func ToPrettyNameForm(s string) string {
	return strings.Title(strings.ToLower(s))
}

// Конвертируем привычную дату в time.Time
func D2T(d string) (time.Time, error) {
	return time.Parse("2.1.2006", d)
}

// Конвертируем time.Time в привычную дату
func T2D(d time.Time) string {
	return d.Format("02.01.2006")
}

// читаем конфиг. По умолчанию default_filename
func ReadConf(filename string) map[string]string {
	Config := make(map[string]string)

	if filename == "" {
		filename = default_filename
	}

	data, err := ioutil.ReadFile(filename)
	Check_err(err, 1)

	var rows = strings.Split(string(data), "\n")
	for i := 0; i < len(rows); i++ {
		if strings.HasPrefix(rows[i], "#") != true && strings.Index(rows[i], "=") != -1 {
			addConf(Config, rows[i], "")
		}
	}
	return Config
}

func addConf(C map[string]string, row string, suffix string) {
	parts := strings.SplitN(row, "=", 2)
	if parts[0] != "" && parts[1] != "" {
		if strings.HasPrefix(parts[0], "DBNAME") {
			// особый случай DBNAME/DBNAME2
			subparts := strings.Split(parts[1], ";")
			C[strings.ToUpper(parts[0])+suffix] = subparts[0]
			s := ""
			if parts[0] == "DBNAME2" {
				s = "2"
			}
			for i := 1; i < len(subparts); i++ {
				addConf(C, subparts[i], s)
			}
		} else {
			C[strings.ToUpper(parts[0])+suffix] = parts[1]
		}
	}
}

// Получаем конфиг для отправки почты
func GetEMailConf(conf map[string]string, key string) (ret map[string]string, err error) {

	err = json.Unmarshal([]byte(conf[strings.ToUpper(key)]), &ret)

	return ret, err
}

// парсинг данный для SMTP из конфига
func ParseConfEmail(data string) map[string]string {
	smtp_server_conf := make(map[string]string)
	err := json.Unmarshal([]byte(data), &smtp_server_conf)
	Check_err(err, 1)

	return smtp_server_conf
}

// Функция для обработки ошибок. Если t=1, то panic, иначе просто выводим сообщение об ошибке
func Check_err(e error, t int64) {
	if e != nil {
		if t == 1 {
			panic(e)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s!", e)
		}
	}
}

func Check_errs(errs []error, t int64) {
	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "Error:")
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "\t%s\n", e)
		}
		if t == 1 {
			panic("Error")
		}
	}
}

func Round(val float64, places int) float64 {
	var pow float64 = 1
	for i := 0; i < places; i++ {
		pow *= 10
	}
	return float64(int((val*pow)+0.5)) / pow
}

// Получени значения из структуры по имени
// пример https://gist.github.com/nikonor/2eebddb68f8925c2e8b0ed405e5bb8da
func GetValueByName(v interface{}, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f
}

// Установка значения поля структуры по имени
// пример https://gist.github.com/nikonor/2eebddb68f8925c2e8b0ed405e5bb8da
func SetValueByName(v interface{}, field string, newval interface{}) {
	r := reflect.ValueOf(v).Elem().FieldByName(field)
	r.Set(reflect.ValueOf(newval))
}

// Проверка email на валидность. Тупая, конечно.
func ValidateEmail(email string) bool {
    email = strings.ToLower(email)
    Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return Re.MatchString(email)
}