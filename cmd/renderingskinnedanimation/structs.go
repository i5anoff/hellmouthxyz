package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Mesh struct {
	Indices []uint32 `json:"indices"`
	Coordinates []Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Index int `json:"index"`
	Vertices []float32 `json:"xyz"`
	Textures []float32 `json:"uvs"`
	Skin map[string]float32 `json:"skin"`
	TotalWeight float32 `json:"totalWeight"`
}

type Armature struct {
	Name  string           `json:"name"`
	Bones map[string]*Bone `json:"bones"`
}

type Bone struct {
	Name                string  	`json:"name"`
	ParentName          string		`json:"parentName"`
	MatrixLocal         *Matrix4f 	`json:"matrix_local"`
	MatrixLocalInverted *Matrix4f 	`json:"matrix_local_inverted"`
}

type IntToMatrix4fMap struct {
	values map[int]*Matrix4f
	keys   []int
}

func NewIntToMatrix4fMap() *IntToMatrix4fMap {
	s := new(IntToMatrix4fMap)

	s.values = make(map[int]*Matrix4f)
	s.keys = []int{}

	return s
}

func (m *IntToMatrix4fMap) sort() {
	sort.Ints(m.keys)
}

func (m *IntToMatrix4fMap) Set(key int, value *Matrix4f) {
	_, present := m.values[key]

	if !present {
		m.keys = append(m.keys, key)
	}

	m.values[key] = value

	m.sort()
}

func (m *IntToMatrix4fMap) Get(key int) *Matrix4f {
	return m.values[key]
}

func (m *IntToMatrix4fMap) Keys() []int {
	return m.keys
}


func (e *IntToMatrix4fMap) UnmarshalJSON(b []byte) error {

	c := map[int]*Matrix4f{}

	if err := json.Unmarshal(b, &c); err == nil {
		x := NewIntToMatrix4fMap()

		for k, v := range c {
			x.Set(k, v)
		}

		*e = *x

		return err
	} else {
		fmt.Println("TESTTTSTSTST")
		fmt.Println(err)
		os.Exit(0)
	}

	return nil
}

type StringToFloat32Map struct {
	values map[string]float32
	keys   []string
}

func NewStringToFloat32Map() *StringToFloat32Map {
	s := new(StringToFloat32Map)

	s.values = make(map[string]float32)
	s.keys = []string{}

	return s
}

func (m *StringToFloat32Map) Set(key string, value float32) {
	_, present := m.values[key]

	if !present {
		m.keys = append(m.keys, key)
	}

	m.values[key] = value
}

func (m *StringToFloat32Map) Get(key string) float32 {
	return m.values[key]
}

func (m *StringToFloat32Map) Keys() []string {
	return m.keys
}

type TextureAndBufferIds struct {
	TextureID uint32
	BufferID  uint32
}