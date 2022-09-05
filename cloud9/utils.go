package Cloud9

import (
	"os"
)

type Set[O comparable] struct {
	inner []O
}

func (this *Set[O]) append(obj O) {
	for _, inn := range(this.inner) {
		if inn == obj {
			return
		}
	}
	this.inner = append(this.inner, obj)
}


func get(from []Object, by string) Object {
	for _, obj := range(from) {
		if obj.name == by {
			return obj
		}
	}

	return Object{id: -1}
}

func pad(str string) string {
	var result string
	for i := 0; i < len(str); i++ {
		if string(str[i]) != " " {
			result += string(str[i])
		}
	}

	return result
}

func extract_value(data [8]byte) []byte {
	var result []byte
	for _, one := range(data) {
		if one != byte(0) {
			result = append(result, one)
		}
	}

	if string(result) == ""{
		result = []byte{48}
	}

	return result
}

func insert_value(data []byte) [8]byte {
	var result [8]byte 
	for i, one := range(data) {
		result[i] = one
	}

	return result
}

func is_in(what string, where []string) (bool, int) {
	for i, one := range(where) {
		for j := 0; j < len(what); j++ {
			if string(one[j]) != string(what[j]) {
				break
			}
			if j == len(one) - 1 {
				return true, i
			}
		}
	}

	return false, -1
}

func contains(what string, where []string) bool {
	for _, one := range(where) {
		if one == what {
			return true
		}
	}
	return false
}

func split(what string, by string) []string {
	var result []string
	var buffer string
	for i := 0; i < len(what); i++ {
		if string(what[i]) == by {
			result = append(result, buffer)
			buffer = ""
			continue
		}
		buffer += string(what[i])
	}

	result = append(result, buffer)
	return result
}

func contains_only(where string, what []string) bool {
	for i := 0; i < len(where); i++ {
		if !contains(string(where[i]), what) {
			return false
		}
	}
	return true
}

func str_to_bytes(str string) [8]byte {
	var result [8]byte
	for i := 0; i < len(str); i++ {
		result[i] = str[i]
	}
	return result
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func parse_conf(addr string) map[string]int {
	conf, err := os.ReadFile(addr)
	handle(err)
	comms := split(string(conf), "\n")

	result := make(map[string]int)

	for i, comm := range(comms) {
		result[comm] = i 
	}

	return result
}

func strip_to_conf(addr string, what string) string {
	conf, err := os.ReadFile(addr)
	handle(err)
	comms := split(string(conf), "\n")
	var result string

	is_command, id := is_in(what, comms)

	if is_command {
		for i := 0; i < len(comms[id]); i++ {
			result += string(what[i])
		}
	}

	return result
}

func count_lines(text string) int {
	counter := 1
	for i := 0; i < len(text); i++ {
		if string(text[i]) == "\n" {
			counter++
		}
	}
	return counter
}

func get_by_id (id_line int, from []LToken) []LToken {
	var result []LToken
	for _, obj := range(from) {
		line := obj.get_line_id()
		if line == id_line {
			result = append(result, obj)
		}
	}
	return result
}

func get_by_operation_id (id int, from []Operation) Operation {
	for _, obj := range(from) {
		if obj.op_id == id {
			return obj
		}
	}

	return Operation{id: -1}
}

func get_lang_token_slice[O Object | Value | Operation | Command] (orig []O) []LToken {
	var result []LToken

	for _, obj := range(orig) {
		tok := LToken(obj)
		result = append(result, tok)
	}

	return result
}

func get_obj_slice[O Object | Value | Operation | Command] (orig []LToken) []O {
	var result []O

	for _, tok := range(orig) {
		obj := tok.(O)
		result = append(result, obj)
	}

	return result
}

func extract_from_string(from string, mask string, ex bool) string {
	var similar string
	var differ  string
	
	for i := 0; i < len(mask); i++ {
		if from[i] == mask[i] {
			similar += string(from[i])
		} else {
			differ += string(from[i])
		}

		if i == len(mask)-1 {
			differ += from[i+1:]
		}
	}

	if ex {
		return differ
	}
	return similar
}

func starts_with(what string, by string) bool {
	for i := 0; i < len(by); i++ {
		if what[i] != by[i] {
			return false
		}
	}
	return true
}

func ends_with(what string, by string) bool {
	len_what := len(what)
	len_by   := len(by)
	for i := 0; i < len_by; i++ {
		if what[len_what-len_by+i] != by[i] {
			return false
		}
	}
	return true
}