package umb_lib

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	default_filename = "./config.txt"
)

type AttachFile struct {
	Name string
	Body []byte
}

type AttachDoc struct {
	Doc_id int64
	Name   string
	Type   int64
}

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

// Проверка сущестования файла "/tmp/<name>.pid"
// Если он есть, то возвращаем  true
// Если его нет, то создаем, записываем PID и возвращаем false
func CheckPidFile(name string) bool {
	var (
		PidFile *os.File
		err     error
	)
	pidfilename := fmt.Sprintf("/tmp/%s.pid", name)
	PidFile, err = os.Open(pidfilename)
	if err != nil {
		PidFile, err = os.Create(pidfilename)
		defer PidFile.Close()
		if err != nil {
			panic(err)
		}
		if PidFile.WriteString(fmt.Sprintf("%d\n", os.Getpid())); err != nil {
			panic(err)
		}
		return false
	}
	return true
}

// Отправка письма через SMTP
// параметр id может быть любым, но лучше случайным или id какоей-либо таблицы
func SendEMail(id int64, smtp_conf map[string]string, from, to, subj, body string, selfcopy bool, attachefiles []AttachFile, attachedocs []AttachDoc) error {
	var (
		pdf      string
		err      error
		array_to []string
	)

	// конвертируем To в массив
	for _, t := range strings.Split(to, ",") {
		array_to = append(array_to, strings.TrimSpace(t))
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", array_to...)
	if selfcopy {
		m.SetAddressHeader("Bcc", from, "")
	}
	m.SetHeader("Subject", subj)
	m.SetBody("text/html", body)

	for _, f := range attachefiles {
		if f.Body != nil {
			tmpfilename := fmt.Sprintf("%s", f.Name)
			if err = ioutil.WriteFile(tmpfilename, []byte(f.Body), 0644); err != nil {
				return err
			}
		}
		m.Attach(f.Name)
	}

	for _, d := range attachedocs {
		pdf, err = GetPDF(d.Doc_id, d.Type, id, d.Name)
		if err != nil {
			return err
		}
		m.Attach(pdf)
	}

	d := gomail.NewDialer(smtp_conf["smtp"], 587, smtp_conf["login"], smtp_conf["password"])
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	// Чистим все лишние файлы.
	os.Remove(pdf)
	for _, f := range attachefiles {
		if f.Body != nil {
			os.Remove(f.Name)
		}
	}

	return nil
}

func GetPDF(doc_id, type_id, mail_id int64, filename string) (string, error) {
	tmpfilename := fmt.Sprintf("%s/mail%d/%s", os.TempDir(), mail_id, filename)
	// pdf, err := exec.Command(fmt.Sprintf("%s/print.pl", default_print_pl_path), fmt.Sprintf("%s_id=%d", TypeRel[type_id], doc_id), fmt.Sprintf("do_what=%s", TypeRel[type_id]), "pdf=1", "podp=1").Output()
	// if err != nil {
	// 	return "", err
	// }
	// if err = ioutil.WriteFile(tmpfilename, pdf, 0644); err != nil {
	// 	return "", err
	// }
	return tmpfilename, nil
}
