package store

type Mem struct {
	data map[string]map[string]string
}

func init() {
	registerStore("mem", func() Store {
		return &Mem{
			data: map[string]map[string]string{},
		}
	})
}

// Complete do Initialize
func (m *Mem) Complete() error {
	return nil
}

// CreateSpace create a new namespace for specific data set
func (m *Mem) CreateSpace(name string) (created bool, err error) {
	_, e := m.data[name]
	if e {
		return false, nil
	}

	m.data[name] = map[string]string{}
	return true, nil
}

// Set update a value of key
func (m *Mem) Set(space string, key, value string) error {
	d, e := m.data[space]
	if e {
		return SpaceNotFound
	}

	d[key] = value
	return nil
}

// Get return target value of key
func (m *Mem) Get(space string, key string) (value string, exist bool, err error) {
	d, e := m.data[space]
	if e {
		return "", false, SpaceNotFound
	}

	v, exist := d[key]
	return v, true, nil
}

// Delete delete target key
func (m *Mem) Delete(space string, key string) error {
	d, e := m.data[space]
	if e {
		return SpaceNotFound
	}
	delete(d, key)
	return nil
}

// DeleteSpace Delete whole namespace
func (m *Mem) DeleteSpace(name string) error {
	_, e := m.data[name]
	if e {
		return SpaceNotFound
	}
	delete(m.data, name)
	return nil
}
