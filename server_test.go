package server

import (
	"context"
	"sort"
	"testing"

	pb "pingthings.io/grpc-toy/sensors"
)

var ALL_SENSORS = []*pb.Sensor{
	&pb.Sensor{
		Uuid: "fake1",
		Collection: "fake-collection1",
		Name: "fake-name1",
		Unit: "fake-unit",
	},
}

type TestLogger struct {
        testLog *testing.T
        verbose bool
}

func (t TestLogger) Print(v ...interface{}) {
        t.testLog.Log(v...)
}

func (t TestLogger) Printf(format string, v ...interface{}) {
        t.testLog.Logf(format, v...)
}

func (t TestLogger) Println(v ...interface{}) {
        t.testLog.Log(v...)
}

func (t TestLogger) Debug(v ...interface{}) {
        if t.verbose {
                t.testLog.Log(v...)
        }
}

func (t TestLogger) Debugf(format string, v ...interface{}) {
        if t.verbose {
                t.testLog.Logf(format, v...)
        }
}

func (t TestLogger) Debugln(v ...interface{}) {
        if t.verbose {
                t.testLog.Log(v...)
        }
}

// Data Repo for test cases
type TestDataRepo struct {
	sensors []*pb.Sensor
}

func GetTestDataRepo(data []*pb.Sensor) *TestDataRepo {
	d := make([]*pb.Sensor, len(data))
	copy(d, data)
	return &TestDataRepo{
		sensors: d,
	}
}

// GetSensors is hardcoded to return all sensors for now.
func (t *TestDataRepo) GetSensors(s *pb.Sensor) ([]*pb.Sensor, error) {
	return t.sensors, nil
}

func (t *TestDataRepo) CreateSensor(s *pb.Sensor) (*pb.Sensor, error) {
	t.sensors = append(t.sensors, s)
	return s, nil
}

// UpdateSensors does not support filtering by the query Sensor object currently.
// Also currently entirely replaces annotations vs. inserting a new one.
func (t *TestDataRepo) UpdateSensors(query, update *pb.Sensor) ([]*pb.Sensor, error) {
	// much worse than corresponding sql, but better than reflection???
	if update.GetCollection() != "" {
		for _, i := range t.sensors {
			i.Collection = update.GetCollection()
		}
	}
	if update.GetName() != "" {
		for _, i := range t.sensors {
			i.Name = update.GetName()
		}
	}
	if update.GetUnit() != "" {
		for _, i := range t.sensors {
			i.Unit = update.GetUnit()
		}
	}
	if update.GetIngress() != "" {
		for _, i := range t.sensors {
			i.Ingress = update.GetIngress()
		}
	}
	if update.GetAnnotations() != nil {
		for _, i := range t.sensors {
			i.Annotations = update.GetAnnotations()
		}
	}
	return t.sensors, nil
}

func TestCreateSensor(t *testing.T) {
	logger := TestLogger{
		testLog: t,
		verbose: true,
	}
	c := context.Background()
	s := SensorServer{
		logger: logger,
		data: GetTestDataRepo(ALL_SENSORS),
	}
	n := pb.Sensor{
		Uuid: "fake2",
		Collection: "fake-collection2",
		Name: "fake-name2",
		Unit: "fake-unit",
	}
	want := append(ALL_SENSORS, &n)
	_, err := s.CreateSensor(c, &n)
	if err != nil {
		t.Errorf("got err %v", err)
	}
	q := pb.Sensor{}
	got, err := s.ReadSensors(c, &q)
	if !sensorListEqual(got.GetSensors(), want) {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetAllSensors(t *testing.T) {
	logger := TestLogger{
		testLog: t,
		verbose: true,
	}
	c := context.Background()
	s := SensorServer{
		logger: logger,
		data: GetTestDataRepo(ALL_SENSORS),
	}
	q := pb.Sensor{}
	got, err := s.ReadSensors(c, &q)
	want := ALL_SENSORS
	if err != nil {
		t.Errorf("got err %v", err)
	}
	if !sensorListEqual(got.GetSensors(), want) {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestUpdateAllSensors(t *testing.T) {
	logger := TestLogger{
		testLog: t,
		verbose: true,
	}
	c := context.Background()
	s := SensorServer{
		logger: logger,
		data: GetTestDataRepo(ALL_SENSORS),
	}
	q := pb.Sensor{}
	// add annotation to *all* sensors
	annotations := []*pb.Annotation{
		&pb.Annotation{
			Key: "foo",
			Value: "1 2 3",
		},
	}
	u := pb.Sensor{
		Annotations: annotations,
	}
	update := &pb.SensorUpdate{
		Selector: &q,
		Update: &u,
	}
	got, err := s.UpdateSensors(c, update)
	want :=  []*pb.Sensor{
		&pb.Sensor{
			Uuid: "fake1",
			Collection: "fake-collection1",
			Name: "fake-name1",
			Unit: "fake-unit",
			Annotations: annotations,
		},
	}
	if err != nil {
		t.Errorf("got err %v", err)
	}
	if !sensorListEqual(got.GetSensors(), want) {
		t.Errorf("got %s want %s", got, want)
	}
}

type SensorSlice []*pb.Sensor

func (s SensorSlice) Len() int {
	return len(s)
}

func (s SensorSlice) Less(i, j int) bool {
	return s[i].GetUuid() < s[j].GetUuid()
}

func (s SensorSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sensorListEqual(i, j []*pb.Sensor) bool {
	if len(i) != len(j) {
		return false
	}
	sort.Sort(SensorSlice(i))
	sort.Sort(SensorSlice(j))
	for k := range i {
		if SensorString(i[k]) != SensorString(j[k]) {
			return false
		}
	}
	return true
}
