// Copyright 2016, Timothy Bogdala <tdb@animal-machine.com>
// See the LICENSE file for more details.

package fizzgui

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	graphics "github.com/tbogdala/fizzle/graphicsprovider"
)

const defaultTextureSampler graphics.Texture = 0

// cmdList will hold all of the information required for one draw call in the
// user interface.
type cmdList struct {
	comboBuffer  []float32 // vbo combo floats
	indexBuffer  []uint32  // vbo elements
	faceCount    uint32    // face count
	indexTracker uint32    // the offset for the next set of indexes when adding new faces
	clipRect     Rect      // clip rect [x1,y1,x2,y2] top-left to bottom-right
	// textureID    graphics.Texture // texture to bind

	isCustom     bool   // is this a custom render command?
	onCustomDraw func() // called during Manager.Draw()
}

// NewCmdList creates a new command list for rendering.
func newCmdList() *cmdList {
	const defaultBufferSize = 1024
	cmds := new(cmdList)
	cmds.comboBuffer = make([]float32, 0, defaultBufferSize)
	cmds.indexBuffer = make([]uint32, 0, defaultBufferSize)
	return cmds
}

// AddFaces takes the raw vertex attribute data in a float slice as well as the
// element indexes and adds it to the internal buffers for rendering.
func (cmds *cmdList) AddFaces(comboFloats []float32, indexInts []uint32, faceCount uint32) {
	cmds.comboBuffer = append(cmds.comboBuffer, comboFloats...)

	// manually adjust each index so that they don't collide with
	// existing element indexes
	var highestIndex uint32
	startIndex := cmds.indexTracker
	for _, idx := range indexInts {
		if idx > highestIndex {
			highestIndex = idx
		}
		cmds.indexBuffer = append(cmds.indexBuffer, startIndex+idx)
	}

	cmds.faceCount += faceCount
	cmds.indexTracker += highestIndex + 1
}

// PrefixFaces takes the raw vertex attribute data in a float slice as well as the
// element indexes and adds it to the internal buffers for rendering at the begining.
func (cmds *cmdList) PrefixFaces(comboFloats []float32, indexInts []uint32, faceCount uint32) {
	cmds.comboBuffer = append(cmds.comboBuffer, comboFloats...)

	// manually adjust each index so that they don't collide with
	// existing element indexes
	var temp []uint32
	var highestIndex uint32
	startIndex := cmds.indexTracker
	for _, idx := range indexInts {
		if idx > highestIndex {
			highestIndex = idx
		}
		temp = append(temp, startIndex+idx)
	}
	cmds.indexBuffer = append(temp, cmds.indexBuffer...)

	cmds.faceCount += faceCount
	cmds.indexTracker += highestIndex + 1
}

// DrawRectFilledDC draws a rectangle in the user interface using a solid background.
// Coordinate parameters should be passed in display coordinates.
// Returns the combo vertex data, element indexes and face count for the rect.
func (cmds *cmdList) DrawRectFilledDC(r Rect, color mgl.Vec4, tex graphics.Texture, uv mgl.Vec4) ([]float32, []uint32, uint32) {
	// uv := whitePixelUv

	// tlx, tly, brx, bry float32

	verts := [8]float32{
		r.TLX, r.BRY,
		r.BRX, r.BRY,
		r.TLX, r.TLY,
		r.BRX, r.TLY,
	}
	indexes := [6]uint32{
		0, 1, 2,
		1, 3, 2,
	}

	uvs := [8]float32{
		uv[0], uv[1],
		uv[2], uv[1],
		uv[0], uv[3],
		uv[2], uv[3],
	}

	comboBuffer := []float32{}
	indexBuffer := []uint32{}

	// add the four vertices
	for i := 0; i < 4; i++ {
		// add the vertex
		comboBuffer = append(comboBuffer, verts[i*2])
		comboBuffer = append(comboBuffer, verts[i*2+1])

		// add the uv
		comboBuffer = append(comboBuffer, uvs[i*2])
		comboBuffer = append(comboBuffer, uvs[i*2+1])

		// add the texture index to use in UV lookup
		comboBuffer = append(comboBuffer, float32(tex))

		// add the color
		comboBuffer = append(comboBuffer, color[:]...)
	}

	// define the polys with 2 faces (6 indexes)
	for i := 0; i < 6; i++ {
		indexBuffer = append(indexBuffer, indexes[i])
	}

	// return the vertex data
	return comboBuffer, indexBuffer, 2
}
