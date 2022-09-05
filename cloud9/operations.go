package Cloud9

import (
	"errors"
	"fmt"
	"strconv"
)

func (i *Interpreter) compute_operation (id int, id_line int) ([8]byte, error){
	var op Operation
	for _, o := range(i.operations) {
		if o.id == id && o.id_line == id_line {
			op = o
		}
	}

	var left_operant  bool = false
	var right_operant bool = false

	var lo_type int
	var ro_type int

	var lo_data [8]byte
	var ro_data [8]byte

	for _, v := range(i.values) {
		if !left_operant{
			if v.id == id-1 && v.id_line == id_line {
				left_operant = true
				lo_data = v.contents
				lo_type = v.data_type.type_id
			}
		}
		if !right_operant{
			if v.id == id + 1 && v.id_line == id_line {
				right_operant = true
				ro_data = v.contents
				ro_type = v.data_type.type_id
			}
		}
	}

	for _, o := range(i.unique_objects.inner) {
		addr := i.all_objects[o]
		for _, a := range(addr.inner) {
			if !left_operant {
				if a.id == id-1 && a.id_line == id_line {
					left_operant = true
					lo_data = o.contents
					lo_type = o.data_type.type_id
				}
			}
			if !right_operant{
				if a.id == id + 1 && a.id_line == id_line {
					right_operant = true
					ro_data = o.contents
					ro_type = o.data_type.type_id
				}
			}
		}
	}

	if !(left_operant && right_operant) {
		err := errors.New("Operation cannot be performed, both operants must be either a value or an object")

		return [8]byte{}, err
	}

	if !(lo_type == ro_type) {
		err := errors.New("Operation cannot be performed, both operants must be the same type")

		return [8]byte{}, err
	}
	
	lo_data_slice := extract_value(lo_data)
	ro_data_slice := extract_value(ro_data)

	var b_slice []byte

	if lo_type == 1 {
		lo_value, err := strconv.Atoi(string(lo_data_slice))
		handle(err)
		ro_value, err := strconv.Atoi(string(ro_data_slice))
		handle(err)

		var result int

		switch op.op_id {
		case 2:
			result = lo_value + ro_value
		case 3: 
			result = lo_value - ro_value
		case 4:
			result = lo_value * ro_value
		case 5:
			result = lo_value / ro_value
		}

		b_slice = []byte(strconv.Itoa(result))
	} else if lo_type == 2 {
		lo_value, err := strconv.ParseFloat(string(lo_data_slice), 64)
		handle(err)
		ro_value, err := strconv.ParseFloat(string(ro_data_slice), 64)
		handle(err)

		var result float64

		switch op.op_id {
		case 2:
			result = lo_value + ro_value
		case 3: 
			result = lo_value - ro_value
		case 4:
			result = lo_value * ro_value
		case 5:
			result = lo_value / ro_value
		}

		b_slice = []byte(pad(fmt.Sprintf("%8.2f", result)))
	} else if lo_type == 3 {
		lo_value := string(lo_data_slice)
		ro_value := string(ro_data_slice)

		var result string

		switch op.op_id {
		case 2:
			result = lo_value + ro_value
		default:
			err := errors.New("Cannot perform operations rather than concatenation to a string")
			
			return [8]byte{}, err
		}

		b_slice = []byte(result)
	} else if lo_type == 4 {
		err := errors.New("Cannot perform mathematical operations to a boolean value")

		return [8]byte{}, err
	}

	

	b_result := insert_value(b_slice)

	return b_result, nil
}