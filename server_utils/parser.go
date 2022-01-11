package server_utils

import (
	"errors"
	"strconv"
	"strings"
)

func get_operator(statement string, cols map[string]int) ([][]string, []string, error) { //
	statement = statement[1 : len(statement)-1]
	strarray := strings.Fields(strings.TrimSpace(statement)) // "c1", "=", "2","and", "c2", "!=", "3"
	var res [][]string                                       //[[c1, ==, 5], [c2, !=, 4]]
	var ops []string                                         //[+, *]
	var tmp []string
	for _, str := range strarray {
		_, ok := cols[str]
		if (ok || str == "*") && len(tmp) <= 0 {
			tmp = append(tmp, str)
			continue
		} else if len(tmp) <= 0 && !ok && str != "*" {
			return nil, nil, errors.New("bad query:column name doesn't exists or statement order is wrong")
		}
		if len(tmp) > 0 {
			if str == "and" {
				ops = append(ops, "*")
				tmp = tmp[:0]
			} else if str == "or" {
				ops = append(ops, "+")
				tmp = tmp[:0]
			} else if str == "==" || str == "$=" {
				tmp = append(tmp, "==")
			} else if str == "&=" || str == "!=" {
				tmp = append(tmp, "&=")
			} else {
				// fmt.Println(len(str))
				tmp = append(tmp, str[1:len(str)-1])
				slice := make([]string, len(tmp))
				copy(slice, tmp)
				res = append(res, slice)

			}
		} else {
			return nil, nil, errors.New("statement is not right, may be bad orders")
		}
	}
	return res, ops, nil

}

func Select_data(query string, col_name map[string]int, data [][]string) error {
	statement, ops, err := get_operator(query, col_name)
	if err != nil {
		return err
	}
	var stk []int
	var i, j int
	var res [][]string

	res = append(res, data[0])
	for i = 1; i < len(data); i++ {
		stk = stk[:0]
		for j = 0; j < len(statement); j++ { // "c1,==,5"
			tmp_str := make([]string, 0)
			if statement[j][0] == "*" {
				tmp_str = data[i]
			} else {
				tmp_str = append(tmp_str, data[i][col_name[statement[j][0]]])
			}
			cur_bool := 1
			for _, data_col := range tmp_str {
				//内部是and

				if statement[j][1] == "==" {
					if data_col == statement[j][2] {
						cur_bool *= 1
					} else {
						cur_bool *= 0
					}
				} else if statement[j][1] == "!=" {
					if data_col != statement[j][2] {
						cur_bool *= 1
					} else {
						cur_bool *= 0
					}
				} else if statement[j][1] == "&=" { // substring
					if strings.Contains(data_col, statement[j][2]) {
						cur_bool *= 1
					} else {
						cur_bool *= 0
					}
				}
			}
			stk = append(stk, cur_bool)
		}
		// calculate if the term > 0, then add
		// stk = [1,0,0] , ops = [+, *]
		if len(stk)-1 != len(ops) {
			return errors.New("parse failed: operator doesn't match with statement")
		}
		s := strconv.Itoa(stk[0])
		var k int
		for k = 1; k < len(stk); k++ {
			s += ops[k-1]
			s += strconv.Itoa(stk[k])
		}
		ans := calculate(s)
		if ans > 0 {
			res = append(res, data[i])
		}
	}
	// save to db
	WriteCsv(res)
	return nil

}

func calculate(s string) (ans int) {
	stack := []int{}
	preSign := '+'
	num := 0
	for i, ch := range s {
		isDigit := '0' <= ch && ch <= '9'
		if isDigit {
			num = num*10 + int(ch-'0')
		}
		if !isDigit && ch != ' ' || i == len(s)-1 {
			switch preSign {
			case '+':
				stack = append(stack, num)
			case '-':
				stack = append(stack, -num)
			case '*':
				stack[len(stack)-1] *= num
			default:
				stack[len(stack)-1] /= num
			}
			preSign = ch
			num = 0
		}
	}
	for _, v := range stack {
		ans += v
	}
	return
}
