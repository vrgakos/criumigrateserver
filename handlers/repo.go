package handlers

type RequestPreDump struct {
	Pid       int       `json:"pid"`
	Path      string    `json:"path"`
	TrackMem  bool		`json:"trackMem"`
}

type RequestDump struct {
	Pid           int       `json:"pid"`
	Path          string    `json:"path"`
	PrevPath      string    `json:"prevPath"`
	TrackMem      bool		`json:"trackMem"`
	Lazy          bool		`json:"lazy"`
	LazyPort      int		`json:"lazyPort"`
	ShellJob	  bool		`json:"shellJob"`
}

type RequestRestore struct {
	Path          string    `json:"path"`
	ShellJob	  bool		`json:"shellJob"`
	Lazy          bool		`json:"lazy"`
}

type RequestLazyPages struct {
	Path          string    `json:"path"`
	Address       string	`json:"address"`
	Port          int		`json:"port"`
}

/*type ResponsePreDump struct {
	Request	  RequestPreDump   `json:"request"`
}*/


type RequestBaseScp struct {
	From  string    `json:"from"`
	To    string    `json:"to"`
}

type RequestBaseClean struct {
	Path  string    `json:"path"`
}

type jsonOk struct {
	Text string `json:"text"`
}

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}


// Give us some seed data
func init() {

}

