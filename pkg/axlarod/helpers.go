package axlarod

func InitalizeLarod() (*Larod, error) {
	l := NewLarod()
	if err := l.Initalize(); err != nil {
		return nil, err
	}
	return l, nil
}
