package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"math"
	"os"
	"sort"
	"time"
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

type Bones struct {
	Data map[string]map[string]map[string]*IntToMatrix4fMap // awful
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
	fmt.Println("m")

	fmt.Println(m)

	fmt.Println("r")

	fmt.Println(r)

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


func CoTangent(angle float32) float32 {
	return float32(1.0 / math.Tan(float64(angle)))
}

func DegreesToRadians(degrees float32) float32 {
	return degrees * float32(math.Pi/180)
}

func NewProjectionMatrix() *Matrix4f {
	projectionMatrix := new(Matrix4f)

	projectionMatrix.M00 = 1
	projectionMatrix.M11 = 1
	projectionMatrix.M22 = 1
	projectionMatrix.M33 = 1

	var fieldOfView = float32(80.0)
	var aspectRatio = float32(float32(Width) / float32(Height))
	near_plane := 0.1
	far_plane := 300.0

	y_scale := CoTangent(DegreesToRadians(fieldOfView / 2.0))
	x_scale := y_scale / aspectRatio
	frustum_length := far_plane - near_plane

	projectionMatrix.M00 = x_scale
	projectionMatrix.M11 = y_scale
	projectionMatrix.M22 = float32(-((far_plane + near_plane) / frustum_length))
	projectionMatrix.M32 = -1
	projectionMatrix.M23 = float32(-((2 * near_plane * far_plane) / frustum_length))
	projectionMatrix.M33 = 0

	return projectionMatrix
}


func NewSkinnedAnimation(armature *Armature, keyframes map[string]*IntToMatrix4fMap) *SkinnedAnimation {
	sa := &SkinnedAnimation{
		map[string]*IntToMatrix4fMap{},
		map[string]*Matrix4f{},
		map[string]*IntToMatrix4fMap{},
		armature,
		1,
		true,
		true,
		1.0 / 1,
		time.Now().UnixNano(),
		0.0,
		"Cube",
		1,
		3,
		2,
	}

	sa.generateFrames(keyframes)
	sa.calculateFinalMatrices(armature)

	return sa
}

func (a *SkinnedAnimation) generateFrames(exportedKeyframes map[string]*IntToMatrix4fMap) {

	for boneName, frameToMatrix := range exportedKeyframes {

		keyframes := []int{}
		bindPoses := []*Matrix4f{}

		for _, key := range frameToMatrix.Keys() {

			m := frameToMatrix.Get(key)

			keyframes = append(keyframes, key)
			bindPoses = append(bindPoses, m)
		}

		bindPoseMatrices := NewIntToMatrix4fMap()

		for i := 0; i < len(keyframes)-1; i++ {
			keyFrameA := keyframes[i]
			keyFrameB := keyframes[i+1]
			bindPoseKeyFrameA := bindPoses[i]
			bindPoseKeyFrameB := bindPoses[i+1]
			if i == 0 {
				bindPoseMatrices.Set(keyFrameA, bindPoseKeyFrameA)
			}
			bindPoseMatrices.Set(keyFrameB, bindPoseKeyFrameB)
		}

		fmt.Println("boneName")
		fmt.Println(boneName)

		a.AllBindPoseTransformations[boneName] = bindPoseMatrices
	}
}

type SkinnedAnimation struct {
	AllBindPoseTransformations map[string]*IntToMatrix4fMap
	InvertedMatrices           map[string]*Matrix4f
	BindMatrices               map[string]*IntToMatrix4fMap

	Armature        *Armature
	CurrentFrame    int64
	Playing         bool
	ForwardPlay     bool
	FrameTime       float64
	LastTime        int64
	UnprocessedTime float64
	MeshName        string
	StartFrame      int64
	EndFrame        int64
	FPS 			int64
}

func (a *SkinnedAnimation) calculateFinalMatrices(armature *Armature) {

	for _, bone := range armature.Bones {
		a.InvertedMatrices[bone.Name] = bone.MatrixLocalInverted
		fmt.Println("bone.Name")
		fmt.Println(bone.Name)

		bindPoseMatrices := a.AllBindPoseTransformations[bone.Name]
		frames := NewIntToMatrix4fMap()

		fmt.Println("bindPoseMatrices")
		fmt.Println(bindPoseMatrices)
		if bindPoseMatrices != nil {
			for _, frame := range bindPoseMatrices.Keys() {
				mw := a.matrixWorld(bone, frame)
				frames.Set(frame, mw)
			}
		} else {
			for i := a.StartFrame; i < a.EndFrame; i++ {
				frames.Set(int(i), bone.MatrixLocal)
			}
		}
		a.BindMatrices[bone.Name] = frames
	}
	fmt.Printf("Finished Setting BindMatrices")

}

func (a *SkinnedAnimation) matrixWorld(bone *Bone, frame int) *Matrix4f {
	parent := a.getParent(bone)

	basis := a.AllBindPoseTransformations[bone.Name].Get(frame)

	fmt.Println("bone.Name")
	fmt.Println(bone.Name)

	fmt.Println("basis")
	fmt.Println(basis)

	if parent == nil {
		return bone.MatrixLocal.Mul(basis)
	} else {
		return a.invertedParentChildLocal(parent, bone, frame, basis)
	}
}

func (a *SkinnedAnimation) invertedParentChildLocal(parent, child *Bone, frame int, childBasis *Matrix4f) *Matrix4f {
	parentMatrix := a.matrixWorld(parent, frame)
	invertedParentChildLocal := (parent.MatrixLocalInverted).Mul(child.MatrixLocal)
	invertedParentChildLocalBasis := invertedParentChildLocal.Mul(childBasis)
	return parentMatrix.Mul(invertedParentChildLocalBasis)
}

func (a *SkinnedAnimation) getParent(bone *Bone) *Bone {
	var parent *Bone
	if len(bone.ParentName) > 0 {
		parent = a.Armature.Bones[bone.ParentName]
	}
	return parent
}



type TextureAndBufferIds struct {
	TextureID uint32
	BufferID  uint32
}

func ListToImage2(modelMatrixElements []float32) *TextureAndBufferIds {

	var modelMatrixId uint32

	gl.GenTextures(1, &modelMatrixId)

	gl.BindTexture(gl.TEXTURE_BUFFER, modelMatrixId)

	var modelBuffer uint32

	gl.GenBuffers(1, &modelBuffer)

	gl.BindBuffer(gl.TEXTURE_BUFFER, modelBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(modelMatrixElements)*4, gl.Ptr(modelMatrixElements), gl.STATIC_DRAW)

	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, modelBuffer)
	gl.BindTexture(gl.TEXTURE_BUFFER, 0)

	gl.BindBuffer(gl.TEXTURE_BUFFER, 0)

	return &TextureAndBufferIds{
		modelMatrixId,
		modelBuffer,
	}
}

func UpdateListToImage(modelBuffer uint32, modelMatrixElements []float32) {
	//FloatBuffer modelMatrixImage = BufferUtil.listAsFloatBuffer(list);
	//
	//GL15.glBindBuffer(GL31.GL_TEXTURE_BUFFER, modelBuffer);
	//GL15.glBufferData(GL31.GL_TEXTURE_BUFFER, modelMatrixImage, GL15.GL_STATIC_DRAW);
	//
	//GL31.glTexBuffer(GL31.GL_TEXTURE_BUFFER, GL30.GL_R32F, modelBuffer);
	//
	//GL15.glBindBuffer(GL31.GL_TEXTURE_BUFFER, 0);

	gl.BindBuffer(gl.TEXTURE_BUFFER, modelBuffer)
	gl.BufferData(gl.TEXTURE_BUFFER, len(modelMatrixElements)*4, gl.Ptr(modelMatrixElements), gl.STATIC_DRAW)

	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32F, modelBuffer)

	gl.BindBuffer(gl.TEXTURE_BUFFER, 0)
}