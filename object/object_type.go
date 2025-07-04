package object

type ObjectType string

const (
	INTEGER_OBJ           ObjectType = "INTEGER"
	BOOLEAN_OBJ           ObjectType = "BOOLEAN"
	NULL_OBJ              ObjectType = "NULL"
	RETURN_VALUE_OBJ      ObjectType = "RETURN"
	FUNCTION_OBJ          ObjectType = "FUNCTION"
	STRING_OBJ            ObjectType = "STRING"
	BUILTIN_OBJ           ObjectType = "BUILTIN"
	ERROR_OBJ             ObjectType = "ERROR"
	ARRAY_OBJ             ObjectType = "ARRAY"
	HASH_OBJ              ObjectType = "HASH"
	QUOTE_OBJ             ObjectType = "QUOTE"
	MACRO_OBJ             ObjectType = "MACRO"
	COMPILED_FUNCTION_OBJ ObjectType = "COMPILED_FUNCTION"
	CLOSURE_OBJ           ObjectType = "CLOSURE"
)
