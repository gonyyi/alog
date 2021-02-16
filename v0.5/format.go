package alog

// Format a bit-formatFlag formatFlag options that is used for variety of configuration.
type Format uint32

func (f Format) Reset() Format {
	return Format(uint32(0))
}
func (f Format) On(item Format) Format {
	return f | item
}
func (f Format) Off(item Format) Format {
	return f &^ item
}

const (
	// Fprefix will show prefix when printing log message
	Fprefix Format = 1 << iota
	// Fdate will show both CCYY and MMDD
	Fdate
	// FdateDay will show 0-6 for JSON or (Sun-Mon)
	FdateDay
	// Ftime will show HHMMSS.000; for json, it will be HHMMSS000
	Ftime
	// FtimeUnix will show unix time
	FtimeUnix
	// FtimeUnixNano will show unix time
	FtimeUnixMs
	// FtimeUTC will show UTC time formats
	FtimeUTC
	// Flevel show Level in the log messsage.
	Flevel
	// Ftag will show tags
	Ftag
	// Fjson will print to a JSON
	Fjson

	// Fdefault will show month/day with time, and Level of logging.
	Fdefault = Fdate | Ftime | Flevel | Ftag
	fUseTime = Fdate | FdateDay | Ftime | FtimeUnix | FtimeUnixMs
)
