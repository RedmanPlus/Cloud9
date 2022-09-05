package Cloud9

type Scope struct {
	lvl     int
	name    string
	parent  *Scope
}

func scope_manager(i *Interpreter, lines int, curr_line int) {
	var obj_toks, com_toks, val_toks, op_toks []LToken

	obj_toks = get_lang_token_slice[Object](i.objects)
	com_toks = get_lang_token_slice[Command](i.commands)
	val_toks = get_lang_token_slice[Value](i.values)
	op_toks  = get_lang_token_slice[Operation](i.operations)

	var tokens []LToken

	tokens = append(tokens, obj_toks...)
	tokens = append(tokens, com_toks...)
	tokens = append(tokens, val_toks...)
	tokens = append(tokens, op_toks...)

	is_scope := false

	for j := curr_line+1; j < lines; j++ {
		line_toks := get_by_id(j, tokens)
		if check_scopes(line_toks, i.curr_scope) {
			is_scope = true
			break
		}
	}

	if !is_scope {
		scope_ptr := i.curr_scope.parent
		i.curr_scope = *scope_ptr
	}
}

func check_scopes(ltoks []LToken, scope Scope) bool {
	for _, tok := range(ltoks) {
		scope_lvl := tok.get_scope()
		if scope_lvl.lvl == scope.lvl && scope_lvl.name == scope.name {
			return true
		}
	}
	return false
}