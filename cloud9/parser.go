package Cloud9

import (
	"fmt"
	"os"
)

type Token struct {
	id   int
	data string
}

func parse_line(line string) []Token {
	var buffer string
	var buffer_id int = 0
	var result []Token
	for i := 0; i < len(line); i++ {
		if string(line[i]) == "\t" {
			continue
		}
		if line[i] == ' ' {
			if buffer == "" {
				continue
			}
			tok := Token{
				id: buffer_id,
				data: buffer,
			}
			result = append(result, tok)
			buffer_id++
			buffer = ""
			i++
		}

		buffer += string(line[i])
	}

	tok := Token{
		id: buffer_id,
		data: buffer,
	}

	result = append(result, tok)

	return result
}

type Type struct {
	type_id int
}

type Object struct {
	id        int
	id_line   int
	name   	  string
	data_type Type
	contents  [8]byte
	scope     Scope
}

type Operation struct {
	id      int
	id_line int
	op_id   int
	scope   Scope
}

type Value struct {
	id        int
	id_line   int
	data_type Type
	contents  [8]byte
	scope     Scope
}

type Command struct {
	id      int
	id_line int
	comm_id int
	tok     Token
	args    []Object
	scope   Scope
}

func (this Object) get_line_id() int {
	return this.id_line
}

func (this Operation) get_line_id() int {
	return this.id_line
}

func (this Value) get_line_id() int {
	return this.id_line
}

func (this Command) get_line_id() int {
	return this.id_line
}

func (this Object) get_scope() Scope {
	return this.scope
}

func (this Operation) get_scope() Scope {
	return this.scope
}

func (this Value) get_scope() Scope {
	return this.scope
}

func (this Command) get_scope() Scope {
	return this.scope
}

type LToken interface {
	get_line_id()   int
	get_scope()     Scope
}

func extract_objects(tokens []Token, scope *int, scope_name *string) ([]Object, []Value, []Operation, []Command) {
	var obj_result []Object
	var val_result []Value
	var op_result  []Operation
	var com_result []Command

	comm_map := parse_conf("commands.conf")

	for j, tok := range(tokens) {
		
		tok.data = pad(tok.data)

		if tok.data == "{" {
			*scope = *scope + 1
			*scope_name = tokens[j-1].data
			continue
		}

		if tok.data == "}" {
			*scope = *scope - 1
			*scope_name = ""
			continue
		}

		if contains(tok.data, []string{"=", "+", "-", "/", "*", "==", ">", "<", ">=", "<=", "!="}) {
			var op_id int
			switch tok.data {
			case "=":
				op_id = 1
			case "+":
				op_id = 2
			case "-":
				op_id = 3
			case "*":
				op_id = 4
			case "/":
				op_id = 5
			case "==":
				op_id = 6
			case ">":
				op_id = 7
			case "<":
				op_id = 8
			case ">=":
				op_id = 9
			case "<=":
				op_id = 10
			case "!=":
				op_id = 11
			}
			op := Operation{
				id:  tok.id,
				op_id: op_id,
				scope: Scope{lvl: *scope, name: *scope_name},
			}
			op_result = append(op_result, op)
			continue
		}
		if contains_only(tok.data, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ".", "'", "true", "false"}) {
			var d_type int
			if contains(tok.data, []string{"."}) {
				d_type = 2
			} else if contains(tok.data, []string{"'"}) {
				d_type = 3
			} else if contains(tok.data, []string{"true"}) || contains(tok.data, []string{"false"}) {
				d_type = 4
				if tok.data == "true" {
					tok.data = "1"
				} else if tok.data == "false" {
					tok.data = "0"
				}
			} else {
				d_type = 1
			}
			
			val := Value{
				id:        tok.id,
				contents:  str_to_bytes(tok.data),
				data_type: Type{type_id: d_type},
				scope: Scope{lvl: *scope, name: *scope_name},
			}
			val_result = append(val_result, val)
			continue
		}

		comm := strip_to_conf("commands.conf", tok.data)

		if comm != "" {
			comm_id := comm_map[comm]

			com := Command{
				id:      tok.id,
				comm_id: comm_id,
				tok:     tok,
				scope: Scope{lvl: *scope, name: *scope_name},
			}

			com_result = append(com_result, com)
			continue
		}

		if contains(tok.data, []string{"int", "flt", "str", "bln"}) {
			continue
		}

		var t Type

		if j != len(tokens)-1 {
			if contains(tokens[j+1].data, []string{"int", "flt", "str", "bln"}) {
				var type_id int
				switch tokens[j+1].data {
				case "int":
					type_id = 1
				case "flt":
					type_id = 2
				case "str":
					type_id = 3
				case "bln":
					type_id = 4
				}
				t = Type{type_id}
				
			}
		}

		obj := Object{
			id:        tok.id,
			name:      tok.data,
			data_type: t,
			scope:     Scope{lvl: *scope, name: *scope_name},
		}
		obj_result = append(obj_result, obj)
	}

	return obj_result, val_result, op_result, com_result
}

func parse_code(code string) ([]Object, []Value, []Operation, []Command) {
	var code_lines []string

	var all_objects    []Object
	var all_values     []Value
	var all_operations []Operation
	var all_commands   []Command
	var scope          int
	var scope_name     string

	code_lines = split(code, "\n")
	for i := 0; i < len(code_lines); i++ {
		var objects    []Object
		var values     []Value
		var operations []Operation
		var commands   []Command
		tokens := parse_line(code_lines[i])
		objects, values, operations, commands = extract_objects(tokens, &scope, &scope_name)
		for j := 0; j < len(objects); j ++ {
			objects[j].id_line = i
		}
		for j := 0; j < len(values); j ++ {
			values[j].id_line = i
		}
		for j := 0; j < len(operations); j ++ {
			operations[j].id_line = i
		}
		for j := 0; j < len(commands); j ++ {
			commands[j].id_line = i
		}
		all_objects    = append(all_objects, objects...)
		all_values     = append(all_values, values...)
		all_operations = append(all_operations, operations...)
		all_commands   = append(all_commands, commands...)
	}

	return all_objects, all_values, all_operations, all_commands
}

func debug(i Interpreter) {
	fmt.Printf("DEBUG Cloud9: Objects        - %v\n", i.objects)
	fmt.Printf("DEBUG Cloud9: Values         - %v\n", i.values)
	fmt.Printf("DEBUG Cloud9: Operations     - %v\n", i.operations)
	fmt.Printf("DEBUG Cloud9: Commands       - %v\n", i.commands)
	fmt.Printf("DEBUG Cloud9: Object map     - %v\n", i.all_objects)
	fmt.Printf("DEBUG Cloud9: Unique objects - %v\n", i.unique_objects)
	fmt.Printf("DEBUG Cloud9: Current scope  - %v\n", i.curr_scope)
	fmt.Printf("DEBUG Cloud9: Parent scopes  - %v\n", i.parent_scopes)
	fmt.Printf("\n")
} 

func Mainloop(addr string, d bool) {
	file, err := os.ReadFile(addr)
	handle(err)
	objects, values, operations, comms  := parse_code(string(file))
	i := Interpreter{
		objects: objects,
		values: values,
		operations: operations,
		commands: comms,
	}
	i.init()
	lines := count_lines(string(file))

	for j := 0; j < lines; j++ {

		toks := get_lang_token_slice[Operation](i.operations)
		ops  := get_by_id(j, toks)
		objs := get_obj_slice[Operation](ops)
		eq := get_by_operation_id(1, objs)
		if eq.id != -1 {
			if eq.scope.lvl != i.curr_scope.lvl && eq.scope.name != i.curr_scope.name {
				continue
			}
			i.gather_value(0, j)
			i.gather_objects()

			if d {
				fmt.Printf("DEBUG Cloud9: line %v\n", j)
				debug(i)
			}
			continue
		} else {
			toks  := get_lang_token_slice[Command](i.commands)
			comms := get_by_id(j, toks)
			objs  := get_obj_slice[Command](comms)
			if objs[0].scope.lvl != i.curr_scope.lvl && objs[0].scope.name != i.curr_scope.name {
				continue
			}
			tok   := objs[0].tok
			args  := parse_args(tok, i)
			objs[0].args = args
			run(objs[0], &i)

			if d {
				fmt.Printf("DEBUG Cloud9: line %v\n", j)
				debug(i)
			}
		}

		scope_manager(&i, lines, j)
	}

	if d {
		fmt.Println("-------------------------------------")
		debug(i)
	}	
}
