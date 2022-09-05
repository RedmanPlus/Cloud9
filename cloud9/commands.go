package Cloud9

import (
	"errors"
	"fmt"
)

func run(what Command, i *Interpreter) {
	switch what.comm_id {
	case 0:
		var data []byte
		for _, one := range(what.args) {
			data = append(data, one.contents[:]...)
		}
		fmt.Println(string(data))
	case 1:
		condition := what.args
		if condition[0].data_type.type_id == 4 {
			if condition[0].contents == [8]byte{49, 0, 0, 0, 0, 0, 0, 0} {
				i.parent_scopes = append(i.parent_scopes, i.curr_scope)
				parent_scope_ptr := &i.parent_scopes[len(i.parent_scopes)-1]
				i.curr_scope = Scope{
					lvl: i.curr_scope.lvl + 1, 
					name: what.tok.data,
					parent: parent_scope_ptr,
				}
			}
		} else {
			err := errors.New("cannot use non-bulean equasion or object as a condition for if")
			handle(err)
		}

	}
}

func parse_args(tok Token, inter Interpreter) []Object {
	data := tok.data
	comm := strip_to_conf("commands.conf", data)

	args := extract_from_string(data, comm, true)
	if starts_with(args, "(") && ends_with(args, ")") {
		args = args[1:len(args)-1]
		var buffer string
		var result []Object
		for i := 0; i < len(args); i++ {
			if args[i] == byte(',') {
				obj := get(inter.unique_objects.inner, buffer)
				if obj.id == -1 {
					if contains_only(buffer, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "."}) {
						val := str_to_bytes(buffer)
						obj = Object{
							contents: val,
						}
					}
				}
				result = append(result, obj)
				buffer = ""
				continue
			}
			if args[i] == byte(' ') {
				continue
			}
			buffer += string(args[i])

			if i == len(args) - 1 {
				obj := get(inter.unique_objects.inner, buffer)
				if obj.id == -1 {
					if contains_only(buffer, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "."}) {
						val := str_to_bytes(buffer)
						obj = Object{
							contents: val,
						}
					}
				}
				result = append(result, obj)
				buffer = ""
			}
		}
		return result
	}
	return []Object{}
}