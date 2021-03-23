package kinematics

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/edaniels/golog"
)

type AutoGenerated struct {
	Model struct {
		Manufacturer string `json:"manufacturer"`
		Name         string `json:"name"`
		Framecount   int    `json:"framecount"`
		Bodies       []struct {
			ID      int   `json:"id"`
			Ignores []int `json:"ignores"`
		} `json:"bodies"`
		Fixeds []struct {
			ID    int `json:"id"`
			Frame struct {
				A string `json:"a"`
				B string `json:"b"`
			} `json:"frame"`
			Rotation struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"rotation"`
			Translation struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"translation"`
		} `json:"fixeds"`
		Joints []struct {
			ID    int    `json:"id"`
			Type  string `json:"type"`
			Frame struct {
				A string `json:"a"`
				B string `json:"b"`
			} `json:"frame"`
			Axis struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"axis"`
			Max float64 `json:"max"`
			Min float64 `json:"min"`
		} `json:"joints"`
		// Home position of joints. Optional, if not provided will be set to all 0.
		Home []float64 `json:"home"`
	} `json:"model"`
}

func ParseJSONFile(filename string, logger golog.Logger) (*Model, error) {
	model := NewModel()
	id2frame := make(map[string]*Frame)
	m := AutoGenerated{}

	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error("failed to read json file")
	}

	err = json.Unmarshal(jsonData, &m)
	if err != nil {
		logger.Error("failed to unmarshall json file")
	}

	model.manufacturer = m.Model.Manufacturer
	model.name = m.Model.Name

	// Create world frame
	wFrame := NewFrame()
	wFrame.IsWorld = true
	setOrient([6]float64{0, 0, 0, 0, 0, 0}, &wFrame.i)
	model.Add(wFrame)
	id2frame["world"] = wFrame

	for _, body := range m.Model.Bodies {
		bFrame := NewFrame()
		bFrame.IsBody = true
		model.Add(bFrame)
		nodeID := "body" + strconv.Itoa(body.ID)
		id2frame[nodeID] = bFrame
		bFrame.Name = nodeID
	}
	for i := 0; i < m.Model.Framecount; i++ {
		frame := NewFrame()
		model.Add(frame)
		nodeID := "frame" + strconv.Itoa(i)
		id2frame[nodeID] = frame
		frame.Name = nodeID
	}

	// Iterate over bodies a second time, setting which ones should ignore one another now that they're all in id2frame
	for _, body := range m.Model.Bodies {
		nodeID := "body" + strconv.Itoa(body.ID)

		b1 := id2frame[nodeID]
		for _, ignore := range body.Ignores {
			ignoreID := "body" + strconv.Itoa(ignore)
			b2 := id2frame[ignoreID]
			b1.selfcollision[b2] = true
			b2.selfcollision[b1] = true
		}
	}

	for _, fixed := range m.Model.Fixeds {
		frameA := id2frame[fixed.Frame.A]
		frameB := id2frame[fixed.Frame.B]

		fixedT := NewTransform()
		fixedT.SetName("fixed" + strconv.Itoa(fixed.ID))

		fixedT.SetEdgeDescriptor(model.AddEdge(frameA, frameB))
		model.Edges[fixedT.GetEdgeDescriptor()] = fixedT
		setOrient([6]float64{fixed.Rotation.X, fixed.Rotation.Y, fixed.Rotation.Z, fixed.Translation.X, fixed.Translation.Y, fixed.Translation.Z}, fixedT)

		fixedT.x.Translation = fixedT.t.Translation()
		fixedT.x.Rotation = fixedT.t.Linear()
	}

	// Now we add all of the transforms. Will eventually support: "cylindrical|fixed|helical|prismatic|revolute|spherical"
	for _, joint := range m.Model.Joints {

		// TODO(pl): Make this a switch once we support more than one joint type
		if joint.Type == "revolute" {
			// TODO(pl): Add speed, wraparound, etc
			frameA := id2frame[joint.Frame.A]
			frameB := id2frame[joint.Frame.B]

			rev := NewJoint(1, 1)
			rev.SetEdgeDescriptor(model.AddEdge(frameA, frameB))
			model.Edges[rev.GetEdgeDescriptor()] = rev

			rev.max = append(rev.max, joint.Max*math.Pi/180)
			rev.min = append(rev.min, joint.Min*math.Pi/180)

			// TODO(pl): Add default on z
			// TODO(pl): Enforce between 0 and 1
			rev.SpatialMat.Set(0, 0, joint.Axis.X)
			rev.SpatialMat.Set(1, 0, joint.Axis.Y)
			rev.SpatialMat.Set(2, 0, joint.Axis.Z)

			rev.SetName("joint" + strconv.Itoa(joint.ID))
		} else {
			logger.Error("Unsupported joint type detected:", joint.Type)
		}
	}

	model.Update()
	if m.Model.Home != nil {
		model.Home = m.Model.Home
	} else {
		for i := 0; i < len(model.Joints); i++ {
			model.Home = append(model.Home, 0)
		}
	}

	return model, err
}

func setOrient(orient [6]float64, trans *Transform) {
	// Important: always rotate in ZYX order
	// Why? Matrix math is not commutative, so do it once one way, and that's the way it needs to be done everywhere
	trans.t.RotZ(orient[2])
	trans.t.RotY(orient[1])
	trans.t.RotX(orient[0])

	trans.t.SetX(orient[3])
	trans.t.SetY(orient[4])
	trans.t.SetZ(orient[5])
}
