package gowl

type Flag uint32

func (f *Flag) Set(flag Flag)     { *f |= flag }
func (f *Flag) Clear(flag Flag)   { *f &= ^flag }
func (f *Flag) Toggle(flag Flag)  { *f ^= flag }
func (f Flag) Has(flag Flag) bool { return f&flag != 0 }
