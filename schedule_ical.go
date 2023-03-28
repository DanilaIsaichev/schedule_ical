package schedule_ical

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Calendar struct {
	Name   string
	Events Events
}

// Структура события урока
type Event struct {
	Start    time.Time
	End      time.Time
	Summary  string
	Location string
	Alarm    int
}

/*
Функция, определяющая время начала и окончания для каждого урока

	ВАЖНО: на данный момент время начала и окончания указаны вручную в коде
*/
func (ev *Event) Set_datetime(date_string string, lesson_number int) error {

	var err error

	switch lesson_number {
	case 1:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 8:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 9:15:00")
		if err != nil {
			return err
		}

	case 2:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 9:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 10:15:00")
		if err != nil {
			return err
		}

	case 3:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 10:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 11:15:00")
		if err != nil {
			return err
		}

	case 4:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 11:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 12:15:00")
		if err != nil {
			return err
		}

	case 5:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 12:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 13:15:00")
		if err != nil {
			return err
		}

	case 6:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 13:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 14:15:00")
		if err != nil {
			return err
		}

	case 7:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 14:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 15:15:00")
		if err != nil {
			return err
		}

	case 8:
		ev.Start, err = time.Parse("2006-01-02 15:04:05", date_string+" 15:30:00")
		if err != nil {
			return err
		}

		ev.End, err = time.Parse("2006-01-02 15:04:05", date_string+" 16:15:00")
		if err != nil {
			return err
		}

	default:
		return errors.New("wrong lesson's number")

	}

	return nil
}

type Events []Event

/*
Функция, удаляющая более события с одинаковым временем начала и окончания

	Предназначена для удаления более ранней версии события при редактировании календаря
*/
func (evs *Events) remove_duplicates() {

	if len(*evs) >= 2 {

		unique_evs := Events{}

		for _, ev := range *evs {

			unique_counter := 0

			for _, u_ev := range unique_evs {

				if ev.Start != u_ev.Start && ev.End != u_ev.End {
					unique_counter++
				}

			}

			if unique_counter >= len(unique_evs) {
				unique_evs = append(unique_evs, ev)
			}
		}

		*evs = unique_evs

	}
}

// Функция, генерирующая текст события из структуры
func Generate_event(event_struct Event) string {

	return fmt.Sprint(`BEGIN:VEVENT
DTSTAMP:`, time.Now().Format("20060201T150405"), `
UID:`, strings.ToUpper(uuid.New().String()), `
DTSTART;TZID=Europe/Moscow:`, event_struct.Start.Format("20060201T150405"), `
DTEND;TZID=Europe/Moscow:`, event_struct.End.Format("20060201T150405"), `
SUMMARY:`, event_struct.Summary, `
LOCATION:`, event_struct.Location, `
BEGIN:VALARM
ACTION:DISPLAY
DESCRIPTION:`, event_struct.Summary, ` - `, event_struct.Location, `
TRIGGER:-PT`, math.Abs(float64(event_struct.Alarm)), `M
END:VALARM
END:VEVENT`)

}

// Функция, генерирующая ical файл
func Make_calendar(cal Calendar, directory string) error {

	// Проверяем, что название календаря не пусто
	if cal.Name != "" {

		var file *os.File

		// Если файла не существует...
		if _, err := os.Stat(directory); os.IsNotExist(err) {

			// Разбиваем путь на поддиректории
			path := strings.Split(directory, "/")

			// Проверяем правильность указанного пути
			if path[len(path)-1] == "cal.ics" {

				path_str := ""

				// Проверяем по очереди существование каждой директории на пути к файлу
				for _, path_part := range path[:len(path)-1] {

					path_str += path_part + "/"

					if _, err := os.Stat(path_str); os.IsNotExist(err) {

						// Если директория не существует - создаём
						err := os.Mkdir(path_str, 0777)
						if err != nil {
							return err
						}

					}
				}

			} else {
				return errors.New("wrong path")
			}

		} else {

			// Парсим файл
			file_cal, err := parse_ical(directory)
			if err != nil {
				return err
			}

			// Добавляем необходимые события
			cal.Events = append(cal.Events, file_cal.Events...)

			// Удаляем повторяющиеся события
			cal.Events.remove_duplicates()

		}

		calendar_str := `BEGIN:VCALENDAR
PRODID:SCHOOL80-SCHEDULER
NAME:` + cal.Name + `
VERSION:2.0
CALSCALE:GREGORIAN
BEGIN:VTIMEZONE
TZID:Europe/Moscow
TZURL:http://tzurl.org/zoneinfo-outlook/Europe/Moscow
X-LIC-LOCATION:Europe/Moscow
BEGIN:STANDARD
TZNAME:MSK
TZOFFSETFROM:+0300
TZOFFSETTO:+0300
DTSTART:19700101T000000`

		// Если в структуре календаря есть события - добавлем в файл
		if len(cal.Events) != 0 {

			for _, event := range cal.Events {
				calendar_str += "\n" + Generate_event(event)
			}
		}

		calendar_str += `
END:STANDARD
END:VTIMEZONE
END:VCALENDAR`

		file, err := os.OpenFile(directory, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(calendar_str)
		if err != nil {
			return err
		} else {
			return nil
		}

	} else {
		return errors.New("calendar's name is empty")
	}
}

func parse_ical(directory string) (cal Calendar, err error) {

	// считываем данные из файла
	data, err := os.ReadFile(directory)
	if err != nil {
		return Calendar{}, err
	}

	// разбиваем данные на строки
	data_strs := strings.Split(string(data), "\n")

	parsed_cal := Calendar{}

	if len(data_strs) != 0 {

		for data_str_number, data_str := range data_strs {

			// Проверяем, что строка не пуста
			if len(data_str) != 0 {

				if len(data_str) >= 5 {
					if data_str[:5] == "NAME:" {

						// Проверяем, что имя встречается только один раз
						if parsed_cal.Name == "" {
							parsed_cal.Name = data_str[5:]
						} else {
							return Calendar{}, errors.New("found more than one NAME in file" + directory)
						}
					}
				}

				if data_str == "BEGIN:VEVENT" {

					ev := Event{}

					start_str := data_strs[data_str_number+3]

					// Парсим начало события
					if start_str[:7] == "DTSTART" {

						ev.Start, err = time.Parse("20060201T150405", start_str[len(start_str)-15:])
						if err != nil {
							return Calendar{}, err
						}

					} else {
						return Calendar{}, errors.New("wrong syntax: expected: DTSTART, got: " + start_str[:7])
					}

					end_str := data_strs[data_str_number+4]

					// Парсим начало события
					if end_str[:5] == "DTEND" {

						ev.End, err = time.Parse("20060201T150405", end_str[len(end_str)-15:])
						if err != nil {
							return Calendar{}, err
						}

					} else {
						return Calendar{}, errors.New("wrong syntax: expected: DTEND, got: " + end_str[:5])
					}

					summary_str := data_strs[data_str_number+5]

					// Парсим описание события
					if summary_str[:7] == "SUMMARY" {

						if summary_str[8:] != "" {

							ev.Summary = summary_str[8:]

						} else {
							return Calendar{}, errors.New(fmt.Sprint("no summury found for event in file ", directory, " in line ", data_str_number))
						}

					} else {
						return Calendar{}, errors.New("wrong syntax: expected: SUMMARY, got: " + summary_str[:7])
					}

					location_str := data_strs[data_str_number+6]

					// Парсим место события
					if location_str[:8] == "LOCATION" {

						if location_str[9:] != "" {

							ev.Location = location_str[9:]

						} else {
							return Calendar{}, errors.New(fmt.Sprint("no location found for event in file ", directory, " in line ", data_str_number))
						}

					} else {
						return Calendar{}, errors.New("wrong syntax: expected: LOCATION, got: " + location_str[:8])
					}

					alarm_str := data_strs[data_str_number+10]

					// Парсим триггер события
					if alarm_str[:7] == "TRIGGER" {

						if alarm_str[8:len(alarm_str)-1] != "" {

							var int_val int
							_, err := fmt.Sscanf(alarm_str, "TRIGGER:-PT%dM", &int_val)
							if err != nil {
								return Calendar{}, err
							}

							if int_val <= 60 {
								ev.Alarm = int_val
							} else {
								return Calendar{}, errors.New(fmt.Sprint("alarm trigger is greater than 60 in file ", directory, " in line ", data_str_number))
							}

						} else {
							return Calendar{}, errors.New(fmt.Sprint("no alarm found for event in file ", directory, " in line ", data_str_number))
						}

					} else {
						return Calendar{}, errors.New("wrong syntax: expected: TRIGGER, got: " + alarm_str[:7])
					}

					parsed_cal.Events = append(parsed_cal.Events, ev)

				}
			} else {
				return Calendar{}, errors.New(fmt.Sprint("an empty string has found in file ", directory, " in line ", data_str_number))
			}

		}

	} else {
		return Calendar{}, errors.New("tried to parse an empty file " + directory)
	}

	if parsed_cal.Name != "" {

		return parsed_cal, nil

	} else {
		return Calendar{}, errors.New("didn't find NAME in file " + directory)
	}
}
