package main

type DataMapBuilder struct {
	data map[string]interface{}
}

func NewDataMapBuilder() *DataMapBuilder {
	return &DataMapBuilder{
		data: make(map[string]interface{}),
	}
}

func (d *DataMapBuilder) WithField(name string, value interface{}) *DataMapBuilder {
	d.data[name] = value
	return d
}

func (d *DataMapBuilder) Build() map[string]interface{} {
	return d.data
}
