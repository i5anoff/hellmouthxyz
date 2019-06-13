package main

import "time"

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
		1,
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

		a.AllBindPoseTransformations[boneName] = bindPoseMatrices
	}
}

func (a *SkinnedAnimation) calculateFinalMatrices(armature *Armature) {

	for _, bone := range armature.Bones {
		a.InvertedMatrices[bone.Name] = bone.MatrixLocalInverted
		bindPoseMatrices := a.AllBindPoseTransformations[bone.Name]
		frames := NewIntToMatrix4fMap()

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
}

func (a *SkinnedAnimation) matrixWorld(bone *Bone, frame int) *Matrix4f {
	parent := a.getParent(bone)

	basis := a.AllBindPoseTransformations[bone.Name].Get(frame)

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
