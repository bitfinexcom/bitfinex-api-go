package subs

type Subs struct {
	Subs map[string]map[string]string
}

const subsLimit = 30

// New returns pointer to instacne of Subscriptions
func New() *Subs {
	return &Subs{
		Subs: make(map[string]map[string]string),
	}
}

// LimitReached returns true if number of subs > subsLimit
func (s *Subs) LimitReached() bool {
	return len(s.Subs) == subsLimit
}

// Added checks if given subscription is already added. Used to
// avoid duplicate subscriptions per client
func (s *Subs) Added(sub map[string]string) (res bool) {
	_, res = s.Subs[s.getID(sub)]
	return
}

// Add adds new subscription to the list
func (s *Subs) Add(sub map[string]string) bool {
	s.Subs[s.getID(sub)] = sub
	return true
}

// GetAll returns all subscriptions
func (s *Subs) GetAll() map[string]map[string]string {
	return s.Subs
}

func (s *Subs) getID(sub map[string]string) (key string) {
	for _, v := range sub {
		key = key + "#" + v
	}
	return
}
