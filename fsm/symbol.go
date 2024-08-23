package fsm

type (
	LogicSymbol      uint
	ArithmeticSymbol uint
)

var (
	EQUAL                LogicSymbol = 0 // ==
	EQ                   LogicSymbol = 0 // ==
	NOT_EQUAL            LogicSymbol = 1 // !=
	NE                   LogicSymbol = 1 // !=
	GRATER_THAN          LogicSymbol = 2 // >
	GT                   LogicSymbol = 2 // >
	GRATER_THAN_OR_EQUAL LogicSymbol = 3 // >=
	GE                   LogicSymbol = 3 // >=
	LESS_THAN            LogicSymbol = 4 // <
	LT                   LogicSymbol = 4 // <
	LESS_THAN_OR_EQUAL   LogicSymbol = 5 // <=
	LE                   LogicSymbol = 5 // <=

	ASSIGN     ArithmeticSymbol = 0 // =
	ADD_ASSIGN ArithmeticSymbol = 1 // +=
	SUB_ASSIGN ArithmeticSymbol = 2 // -=
	MUL_ASSIGN ArithmeticSymbol = 3 // *=
	DIV_ASSIGN ArithmeticSymbol = 4 // /=

)

func (symbol LogicSymbol) String() string {
	switch symbol {
	case EQ, EQUAL:
		return "=="
	case NE, NOT_EQUAL:
		return "!="
	case GT, GRATER_THAN:
		return ">"
	case GE, GRATER_THAN_OR_EQUAL:
		return ">="
	case LT, LESS_THAN:
		return "<"
	case LE, LESS_THAN_OR_EQUAL:
		return "<="
	default:
		return ""
	}
}

func (symbol ArithmeticSymbol) String() string {
	switch symbol {
	case ASSIGN:
		return "="
	case ADD_ASSIGN:
		return "+="
	case SUB_ASSIGN:
		return "-="
	case MUL_ASSIGN:
		return "*="
	case DIV_ASSIGN:
		return "/="
	default:
		return ""
	}
}
