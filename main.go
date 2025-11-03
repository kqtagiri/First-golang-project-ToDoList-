package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type task struct {
	name     string
	desc     string
	date     string
	priority int
	status   bool
}

func (t *task) String() string {
	return fmt.Sprint(t.name, " | ", t.desc, " | ", t.date, " | ", t.priority, " | ", t.status)
}

func EnterTheDate() (string, error) {
	var date string
	fmt.Println("Введите дату(формат ввода даты - DD.MM.YYYY):")
	fmt.Scan(&date)
	err := errors.New("Неправильный формат ввода даты")
	if len(date) != 10 {
		return "", err
	}
	if dd, err := strconv.Atoi(date[:2]); err != nil {
		return "", err
	} else if mm, err := strconv.Atoi(date[3:5]); err != nil {
		return "", err
	} else if yyyy, err := strconv.Atoi(date[6:]); err != nil {
		return "", err
	} else if yyyy < 2000 || yyyy > 2025 || dd <= 0 || dd > 31 || mm <= 0 || mm > 12 {
		err := errors.New("Неправильный формат ввода даты")
		return "", err
	} else {
		return date, nil
	}
}

func (t *task) PrintTasks() {
	var temp, temp2 string
	if t.priority == 1 {
		temp = "Низкий"
	} else if t.priority == 2 {
		temp = "Средний"
	} else {
		temp = "Высокий"
	}
	if t.status == false {
		temp2 = "[-]"
	} else {
		temp2 = "[✓]"
	}
	spaces1 := " "
	spaces2 := " "
	spaces3 := " "
	spaces4 := "     "
	if len(t.name) < 10 {
		spaces1 = strings.Repeat(" ", 10-len(t.name))
	}
	if len(t.desc) < 9 {
		spaces2 = strings.Repeat(" ", 9-len(t.desc))
	}
	fmt.Println(t.name, spaces1, t.desc, spaces2, t.date, spaces3, temp, spaces4, temp2)
}

func PrintTable() {
	spaces := strings.Repeat(" ", 10)
	fmt.Println(spaces, "Название     Описание     Дата     Приоритет     Статус   ")
}

func main() {
	fileName := "file.txt"
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	count := 0
	var name, desc string
	var status bool
	var choose, index, pr int
	list := []task{}
	for check := true; check != false; {
		fmt.Println("1.Добавить запись\n2.Удалить запись\n3.Редактировать запись\n4.Посмотреть все записи\n5.Выход")
		fmt.Scan(&choose)
		if choose == 5 {
			break
		}
		switch choose {
		case 1:
			fmt.Println("Введите новую запись:")
			fmt.Println("Введите название:")
			fmt.Scan(&name)
			fmt.Println("Введите описание:")
			fmt.Scan(&desc)
			date, err := EnterTheDate()
			for err != nil {
				date, err = EnterTheDate()
			}
			fmt.Println("Введите приоритет(1 - Низкий, 2 - Средний, 3 - Высокий):")
			fmt.Scan(&pr)
			for pr < 1 || pr > 3 {
				fmt.Println("Ошибка. Введите цифру от 1 до 3")
				fmt.Scan(&pr)
			}
			fmt.Println("Введите статус(0 - Не выполнен, 1 - Выполнен):")
			fmt.Scan(&status)
			t := task{name, desc, date, pr, status}
			list = append(list, t)
			_, err = file.WriteString(t.String() + "\n")
			if err != nil {
				panic(err)
			}
			count++
			break
		case 2:
			if count == 0 {
				fmt.Println("Нет данных для удаления!")
				break
			}
			var ch2 int
			fmt.Println("1 - Удалить все выполненные записи\n2 - Удалить конкретную запись")
			fmt.Scan(&ch2)
			if ch2 == 1 {
				for i := 1; i <= count; {
					if list[i-1].status == true {
						list = append(list[:i-1], list[i:]...)
						count--
						continue
					}
					i++
				}
			} else if ch2 == 2 {
				fmt.Println("Введите номер для удаления, от 1 до ", count, ":")
				fmt.Scan(&index)
				for index < 1 || index > count {
					fmt.Println("Ошибка. Введите цифру от 1 до ", count)
					fmt.Scan(&index)
				}
				list = append(list[:index-1], list[index:]...)
				count--
			} else {
				fmt.Println("Нет такого выбора")
			}
			file.Close()
			file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				panic(err)
			}
			for i := 0; i < count; i++ {
				_, err = file.WriteString(list[i].String() + "\n")
				if err != nil {
					panic(err)
				}
			}
			break
		case 3:
			if count == 0 {
				fmt.Println("Нет данных для редактирования!")
				break
			}
			fmt.Println("Введите номер для редактирования, от 1 до ", count, ":")
			fmt.Scan(&index)
			fmt.Println("Введите другую запись:")
			fmt.Println("Введите название:")
			fmt.Scan(&name)
			fmt.Println("Введите описание:")
			fmt.Scan(&desc)
			fmt.Println("Введите дату добавления:")
			date, err := EnterTheDate()
			for err != nil {
				date, err = EnterTheDate()
			}
			fmt.Println("Введите приоритет(1 - Низкий, 2 - Средний, 3 - Высокий):")
			fmt.Scan(&pr)
			for pr < 1 || pr > 3 {
				fmt.Println("Ошибка. Введите цифру от 1 до 3")
				fmt.Scan(&pr)
			}
			fmt.Println("Введите статус(0 - Не выполнен, 1 - Выполнен):")
			fmt.Scan(&status)
			t := task{name, desc, date, pr, status}
			list[index-1] = t
			file.Close()
			file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				panic(err)
			}
			for i := 0; i < count; i++ {
				_, err = file.WriteString(list[i].String() + "\n")
				if err != nil {
					panic(err)
				}
			}
			break
		case 4:
			if count == 0 {
				fmt.Println("Нет данных для просмотра!")
				break
			}
			PrintTable()
			for i := 0; i < count; i++ {
				fmt.Print("Запись №", i+1, ": ")
				list[i].PrintTasks()
			}
			fmt.Println("Введите:\n1 - для сортировки\n2 - для поиска\n3 - выход")
			var ch2 int
			fmt.Scan(&ch2)
			switch ch2 {
			case 1:
				fmt.Println("Введите:\n1 - для сортировки по приоритету\n2 - для сортировки по дате добавления")
				var ch3 int
				fmt.Scan(&ch3)
				if ch3 == 1 {
					sort.Slice(list, func(i, j int) bool {
						return list[i].priority > list[j].priority
					})
				}
				if ch3 == 2 {
					sort.Slice(list, func(i, j int) bool {
						if list[i].date[6:] == list[j].date[6:] {
							if list[i].date[3:5] == list[j].date[3:5] {
								return list[i].date[:2] > list[j].date[:2]
							} else {
								return list[i].date[3:5] > list[j].date[3:5]
							}
						} else {
							return list[i].date[6:] > list[j].date[6:]
						}
					})
				}
				fmt.Println("Список после сортировки:")
				for i := range count {
					fmt.Print("Запись №", i+1, ": ")
					list[i].PrintTasks()
				}
				break
			case 2:
				fmt.Println("Введите название для поиска:")
				fmt.Scan(&name)
				equal := false
				for i := range count {
					if strings.EqualFold(list[i].name, name) {
						fmt.Println("Найдено совпадение!")
						list[i].PrintTasks()
						equal = true
					}
				}
				if !equal {
					fmt.Println("Совпадений не найдено!")
				}
				break
			case 3:
				break
			default:
				fmt.Println("Такой функционал отсутсвует")
			}
		default:
			fmt.Println("Такой функционал отсутствует")
		}
	}
}