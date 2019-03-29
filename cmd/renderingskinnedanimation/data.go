package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Bones struct {
	Data map[string]map[string]map[string]*IntToMatrix4fMap // awful
}

type Meshes struct {
	Data map[string]Mesh
}

type Mesh struct {
	Indices []float32 `json:"indices"`
	Vertices []Vertex `json:"coordinates"`
}

type Vertex struct {
	Skin map[string]float32 `json:"skin"`
	XYZ []float32 `json:"xyz"`
	UV []float32 `json:"uvs"`
	Index float32 `json:"index"`
	TotalWeight float32 `json:"totalWeight"`
}

type Armature struct {
	Name  string           `json:"name"`
	Bones map[string]*Bone `json:"bones"`
}

type Bone struct {
	Name                string  	`json:"name"`
	ParentName          string		`json:"parentName"`
	MatrixLocal         *Matrix4f 	`json:"matrixLocal"`
	MatrixLocalInverted *Matrix4f 	`json:"matrixLocalInverted"`
}

type Matrix4f struct {
	M00 float32
	M01 float32
	M02 float32
	M03 float32

	M10 float32
	M11 float32
	M12 float32
	M13 float32

	M20 float32
	M21 float32
	M22 float32
	M23 float32

	M30 float32
	M31 float32
	M32 float32
	M33 float32
}

func (e *Matrix4f) UnmarshalJSON(b []byte) error {

	c := []float32{}

	if err := json.Unmarshal(b, &c); err == nil {
		x := &Matrix4f{	c[0],c[1],c[2],c[3],
						c[4],c[5],c[6],c[7],
						c[8],c[9],c[10],c[11],
						c[12],c[13],c[14],c[15],}
		*e = *x
	} else {
		return err
	}

	return nil
}

func (m *Matrix4f) Mul(r *Matrix4f) *Matrix4f {
	res := new(Matrix4f)

	res.M00 = (m.M00 * r.M00) + (m.M01 * r.M10) + (m.M02 * r.M20) + (m.M03 * r.M30)
	res.M01 = (m.M00 * r.M01) + (m.M01 * r.M11) + (m.M02 * r.M21) + (m.M03 * r.M31)
	res.M02 = (m.M00 * r.M02) + (m.M01 * r.M12) + (m.M02 * r.M22) + (m.M03 * r.M32)
	res.M03 = (m.M00 * r.M03) + (m.M01 * r.M13) + (m.M02 * r.M23) + (m.M03 * r.M33)

	res.M10 = (m.M10 * r.M00) + (m.M11 * r.M10) + (m.M12 * r.M20) + (m.M13 * r.M30)
	res.M11 = (m.M10 * r.M01) + (m.M11 * r.M11) + (m.M12 * r.M21) + (m.M13 * r.M31)
	res.M12 = (m.M10 * r.M02) + (m.M11 * r.M12) + (m.M12 * r.M22) + (m.M13 * r.M32)
	res.M13 = (m.M10 * r.M03) + (m.M11 * r.M13) + (m.M12 * r.M23) + (m.M13 * r.M33)

	res.M20 = (m.M20 * r.M00) + (m.M21 * r.M10) + (m.M22 * r.M20) + (m.M23 * r.M30)
	res.M21 = (m.M20 * r.M01) + (m.M21 * r.M11) + (m.M22 * r.M21) + (m.M23 * r.M31)
	res.M22 = (m.M20 * r.M02) + (m.M21 * r.M12) + (m.M22 * r.M22) + (m.M23 * r.M32)
	res.M23 = (m.M20 * r.M03) + (m.M21 * r.M13) + (m.M22 * r.M23) + (m.M23 * r.M33)

	res.M30 = (m.M30 * r.M00) + (m.M31 * r.M10) + (m.M32 * r.M20) + (m.M33 * r.M30)
	res.M31 = (m.M30 * r.M01) + (m.M31 * r.M11) + (m.M32 * r.M21) + (m.M33 * r.M31)
	res.M32 = (m.M30 * r.M02) + (m.M31 * r.M12) + (m.M32 * r.M22) + (m.M33 * r.M32)
	res.M33 = (m.M30 * r.M03) + (m.M31 * r.M13) + (m.M32 * r.M23) + (m.M33 * r.M33)

	return res
}

func (m *Matrix4f) Get1D() []float32 {
	return []float32{m.M00,m.M01,m.M02,m.M03,   m.M10,m.M11,m.M12,m.M13,   m.M20,m.M21,m.M22,m.M23,  m.M30,m.M31,m.M32,m.M33 }
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

const SkinnedVertices = `{
  "Cube": {
    "indices": [
      0,
      1,
      2,
      3,
      4,
      7,
      6,
      5,
      8,
      12,
      13,
      9,
      14,
      15,
      16,
      10,
      17,
      18,
      19,
      11,
      20,
      21,
      22,
      23,
      24,
      20,
      23,
      25,
      26,
      17,
      11,
      27,
      28,
      14,
      10,
      29,
      30,
      8,
      9,
      31
    ],
    "coordinates": [
      {
        "totalWeight": 1,
        "index": 0,
        "xyz": [
          1,
          0.9999999403953552,
          0
        ],
        "uvs": [
          0.33333349227905273,
          0.5000000596046448
        ],
        "skin": {
          "Bone": 0.9098,
          "Bone.001": 0.0902
        }
      },
      {
        "totalWeight": 1,
        "index": 1,
        "xyz": [
          1,
          -1,
          0
        ],
        "uvs": [
          0.6666666865348816,
          0.5000000596046448
        ],
        "skin": {
          "Bone": 0.91731,
          "Bone.001": 0.08269
        }
      },
      {
        "totalWeight": 1,
        "index": 2,
        "xyz": [
          -1.0000001192092896,
          -0.9999998211860657,
          0
        ],
        "uvs": [
          0.6666666269302368,
          0.7500000596046448
        ],
        "skin": {
          "Bone": 0.9098,
          "Bone.001": 0.0902
        }
      },
      {
        "totalWeight": 1,
        "index": 3,
        "xyz": [
          -0.9999996423721313,
          1.0000003576278687,
          0
        ],
        "uvs": [
          0.33333340287208557,
          0.75
        ],
        "skin": {
          "Bone": 0.91731,
          "Bone.001": 0.08269
        }
      },
      {
        "totalWeight": 1,
        "index": 4,
        "xyz": [
          1.0000004768371582,
          0.999999463558197,
          4
        ],
        "uvs": [
          0.9999999403953552,
          0.5000001788139343
        ],
        "skin": {
          "Bone": 0.0902,
          "Bone.001": 0.9098
        }
      },
      {
        "totalWeight": 1,
        "index": 5,
        "xyz": [
          0.9999993443489075,
          -1.0000005960464478,
          4
        ],
        "uvs": [
          0.6666667461395264,
          0.5000002980232239
        ],
        "skin": {
          "Bone": 0.08269,
          "Bone.001": 0.91731
        }
      },
      {
        "totalWeight": 1,
        "index": 6,
        "xyz": [
          -1.0000003576278687,
          -0.9999996423721313,
          4
        ],
        "uvs": [
          0.6666668653488159,
          0.7500002384185791
        ],
        "skin": {
          "Bone": 0.0902,
          "Bone.001": 0.9098
        }
      },
      {
        "totalWeight": 1,
        "index": 7,
        "xyz": [
          -0.9999999403953552,
          1,
          4
        ],
        "uvs": [
          1,
          0.7500001788139343
        ],
        "skin": {
          "Bone": 0.08269,
          "Bone.001": 0.91731
        }
      },
      {
        "totalWeight": 1,
        "index": 8,
        "xyz": [
          1.000000238418579,
          0.9999997019767761,
          2
        ],
        "uvs": [
          0.33333325386047363,
          0.7500000596046448
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 9,
        "xyz": [
          0.9999996423721313,
          -1.000000238418579,
          2
        ],
        "uvs": [
          3.9736413270929916e-08,
          0.75
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 10,
        "xyz": [
          -1.000000238418579,
          -0.9999997615814209,
          2
        ],
        "uvs": [
          0.33333343267440796,
          0.2500000596046448
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 11,
        "xyz": [
          -0.9999997615814209,
          1.000000238418579,
          2
        ],
        "uvs": [
          0.6666668057441711,
          0.2500000298023224
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 12,
        "xyz": [
          1.0000004768371582,
          0.999999463558197,
          4
        ],
        "uvs": [
          0.33333322405815125,
          1
        ],
        "skin": {
          "Bone": 0.0902,
          "Bone.001": 0.9098
        }
      },
      {
        "totalWeight": 1,
        "index": 13,
        "xyz": [
          0.9999993443489075,
          -1.0000005960464478,
          4
        ],
        "uvs": [
          0,
          1
        ],
        "skin": {
          "Bone": 0.08269,
          "Bone.001": 0.91731
        }
      },
      {
        "totalWeight": 1,
        "index": 14,
        "xyz": [
          0.9999996423721313,
          -1.000000238418579,
          2
        ],
        "uvs": [
          0.6666666865348816,
          0.25
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 15,
        "xyz": [
          0.9999993443489075,
          -1.0000005960464478,
          4
        ],
        "uvs": [
          0.6666666865348816,
          0.5000000596046448
        ],
        "skin": {
          "Bone": 0.08269,
          "Bone.001": 0.91731
        }
      },
      {
        "totalWeight": 1,
        "index": 16,
        "xyz": [
          -1.0000003576278687,
          -0.9999996423721313,
          4
        ],
        "uvs": [
          0.33333349227905273,
          0.5000001788139343
        ],
        "skin": {
          "Bone": 0.0902,
          "Bone.001": 0.9098
        }
      },
      {
        "totalWeight": 1,
        "index": 17,
        "xyz": [
          -1.000000238418579,
          -0.9999997615814209,
          2
        ],
        "uvs": [
          1,
          0.25
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 18,
        "xyz": [
          -1.0000003576278687,
          -0.9999996423721313,
          4
        ],
        "uvs": [
          1,
          0.5
        ],
        "skin": {
          "Bone": 0.0902,
          "Bone.001": 0.9098
        }
      },
      {
        "totalWeight": 1,
        "index": 19,
        "xyz": [
          -0.9999999403953552,
          1,
          4
        ],
        "uvs": [
          0.6666668653488159,
          0.5
        ],
        "skin": {
          "Bone": 0.08269,
          "Bone.001": 0.91731
        }
      },
      {
        "totalWeight": 1,
        "index": 20,
        "xyz": [
          1.000000238418579,
          0.9999997019767761,
          2
        ],
        "uvs": [
          0.3333333730697632,
          0.25
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 21,
        "xyz": [
          1,
          0.9999999403953552,
          0
        ],
        "uvs": [
          0.33333340287208557,
          0.4999999701976776
        ],
        "skin": {
          "Bone": 0.9098,
          "Bone.001": 0.0902
        }
      },
      {
        "totalWeight": 1,
        "index": 22,
        "xyz": [
          -0.9999996423721313,
          1.0000003576278687,
          0
        ],
        "uvs": [
          2.384184796255795e-07,
          0.5000000596046448
        ],
        "skin": {
          "Bone": 0.91731,
          "Bone.001": 0.08269
        }
      },
      {
        "totalWeight": 1,
        "index": 23,
        "xyz": [
          -0.9999997615814209,
          1.000000238418579,
          2
        ],
        "uvs": [
          1.5894565308371966e-07,
          0.2500000596046448
        ],
        "skin": {
          "Bone": 0.5,
          "Bone.001": 0.5
        }
      },
      {
        "totalWeight": 1,
        "index": 24,
        "xyz": [
          1.0000004768371582,
          0.999999463558197,
          4
        ],
        "uvs": [
          0.3333333134651184,
          0
        ],
        "skin": {
          "Bone": 0.0902,
          "Bone.001": 0.9098
        }
      },
      {
        "totalWeight": 1,
        "index": 25,
        "xyz": [
          -0.9999999403953552,
          1,
          4
        ],
        "uvs": [
          0,
          5.960463766996327e-08
        ],
        "skin": {
          "Bone": 0.08269,
          "Bone.001": 0.91731
        }
      },
      {
        "totalWeight": 1,
        "index": 26,
        "xyz": [
          -1.0000001192092896,
          -0.9999998211860657,
          0
        ],
        "uvs": [
          1,
          0
        ],
        "skin": {
          "Bone": 0.9098,
          "Bone.001": 0.0902
        }
      },
      {
        "totalWeight": 1,
        "index": 27,
        "xyz": [
          -0.9999996423721313,
          1.0000003576278687,
          0
        ],
        "uvs": [
          0.6666667461395264,
          1.887447709236767e-08
        ],
        "skin": {
          "Bone": 0.91731,
          "Bone.001": 0.08269
        }
      },
      {
        "totalWeight": 1,
        "index": 28,
        "xyz": [
          1,
          -1,
          0
        ],
        "uvs": [
          0.6666666865348816,
          0
        ],
        "skin": {
          "Bone": 0.91731,
          "Bone.001": 0.08269
        }
      },
      {
        "totalWeight": 1,
        "index": 29,
        "xyz": [
          -1.0000001192092896,
          -0.9999998211860657,
          0
        ],
        "uvs": [
          0.33333340287208557,
          6.193295121192932e-08
        ],
        "skin": {
          "Bone": 0.9098,
          "Bone.001": 0.0902
        }
      },
      {
        "totalWeight": 1,
        "index": 30,
        "xyz": [
          1,
          0.9999999403953552,
          0
        ],
        "uvs": [
          0.3333333134651184,
          0.5000000596046448
        ],
        "skin": {
          "Bone": 0.9098,
          "Bone.001": 0.0902
        }
      },
      {
        "totalWeight": 1,
        "index": 31,
        "xyz": [
          1,
          -1,
          0
        ],
        "uvs": [
          9.934103673003847e-08,
          0.5000000596046448
        ],
        "skin": {
          "Bone": 0.91731,
          "Bone.001": 0.08269
        }
      }
    ]
  }
}
`

const BoneMatrices = `{
  "Cube": {
    "ArmatureAction": {
      "Bone": {
        "1": [
          1,
          0,
          0,
          0,
          0,
          1,
          0,
          0,
          0,
          0,
          1,
          0,
          0,
          0,
          0,
          1
        ],
        "2": [
          1,
          0,
          0,
          0,
          0,
          0.93969,
          -0.34202,
          0,
          0,
          0.34202,
          0.93969,
          0,
          0,
          0,
          0,
          1
        ],
        "3": [
          1,
          0,
          0,
          0,
          0,
          0.76604,
          -0.64279,
          0,
          0,
          0.64279,
          0.76604,
          0,
          0,
          0,
          0,
          1
        ]
      },
      "Bone.001": {
        "1": [
          1,
          0,
          0,
          0,
          0,
          1,
          0,
          0,
          0,
          0,
          1,
          0,
          0,
          0,
          0,
          1
        ],
        "2": [
          1,
          0,
          0,
          0,
          0,
          0.93969,
          -0.34202,
          0,
          0,
          0.34202,
          0.93969,
          0,
          0,
          0,
          0,
          1
        ],
        "3": [
          1,
          0,
          0,
          0,
          0,
          0.76604,
          -0.64279,
          0,
          0,
          0.64279,
          0.76604,
          0,
          0,
          0,
          0,
          1
        ]
      }
    }
  }
}
`

const ArmatureData = `{
  "Armature": {
    "name": "Armature",
    "bones": {
      "Bone": {
        "name": "Bone",
        "matrixLocal": [
          1,
          0,
          0,
          0,
          0,
          0,
          -1,
          0,
          0,
          1,
          0,
          0,
          0,
          0,
          0,
          1
        ],
        "matrixLocalInverted": [
          1,
          -0,
          0,
          -0,
          -0,
          0,
          1,
          0,
          0,
          -1,
          0,
          -0,
          -0,
          0,
          -0,
          1
        ]
      },
      "Bone.001": {
        "name": "Bone.001",
        "matrixLocal": [
          1,
          0,
          0,
          0,
          0,
          0,
          -1,
          2,
          0,
          1,
          0,
          0,
          0,
          0,
          0,
          1
        ],
        "parentName": "Bone",
        "matrixLocalInverted": [
          1,
          -0,
          0,
          -0,
          -0,
          0,
          1,
          0,
          0,
          -1,
          0,
          2,
          -0,
          0,
          -0,
          1
        ]
      }
    },
    "matrixWorld": [
      1,
      0,
      0,
      0,
      0,
      1,
      0,
      0,
      0,
      0,
      1,
      0,
      0,
      0,
      0,
      1
    ]
  }
}
`
