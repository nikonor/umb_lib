# umb_lib

package umb_lib
    import "github.com/nikonor/umb_lib"



FUNCTIONS

func Check_err(e error, t int64)
    Функция для обработки ошибок. Если t=1, то panic, иначе просто выводим
    сообщение об ошибке

func Check_errs(errs []error, t int64)

func D2T(d string) (time.Time, error)
    Конвертируем привычную дату в time.Time

func GetEMailConf(conf map[string]string, key string) (ret map[string]string, err error)
    Получаем конфиг для отправки почты

func GetValueByName(v interface{}, field string) interface{}
    Получени значения из структуры по имени пример
    https://gist.github.com/nikonor/2eebddb68f8925c2e8b0ed405e5bb8da

func ParseConfEmail(data string) map[string]string
    парсинг данный для SMTP из конфига

func ReadConf(filename string) map[string]string
    читаем конфиг. По умолчанию default_filename

func Round(val float64, places int) float64

func SetValueByName(v interface{}, field string, newval interface{})
    Установка значения поля структуры по имени пример
    https://gist.github.com/nikonor/2eebddb68f8925c2e8b0ed405e5bb8da

func T2D(d time.Time) string
    Конвертируем time.Time в привычную дату

func ToPrettyNameForm(s string) string
