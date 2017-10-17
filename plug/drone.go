package plug

type (
	Repo struct {
		Owner   string
		Name    string
		Link    string
		Avatar  string
		Branch  string
		Private bool
		Trusted bool
	}
	Build struct {
		Number   int64
		Event    string
		Status   string
		Deploy   string
		Created  int64
		Started  int64
		Finished int64
		Link     string
	}
	Commit struct {
		Sha     string
		Ref     string
		Link    string
		Branch  string
		Message string
		Author  Author
	}
	Author struct {
		Name   string
		Email  string
		Avatar string
	}
)
