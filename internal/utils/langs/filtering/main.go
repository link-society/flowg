package filtering

func Compile(input string) (Filter, error) {
	return newFilterImpl(input)
}
