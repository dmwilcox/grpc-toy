package server

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "pingthings.io/grpc-toy/sensors"
)

type Log interface {
        Print(v ...interface{})
        Printf(format string, v ...interface{})
        Println(v ...interface{})
        Debug(v ...interface{})
        Debugf(format string, v ...interface{})
        Debugln(v ...interface{})
}

func SensorString(s *pb.Sensor) string {
	var annotations = ""
	if s.GetAnnotations() != nil {
		annotations = AnnotationsString(s.GetAnnotations())
	}
	return fmt.Sprintf(
		"< Sensor uuid=%s collection=%s name=%s unit=%s ingress=%s annotations=%s >",
		s.GetUuid(),
		s.GetCollection(),
		s.GetName(),
		s.GetUnit(),
		s.GetIngress(),
		annotations,
	)
}

func AnnotationString(a *pb.Annotation) string {
	return fmt.Sprintf(
		"< Annotation key=%s value=%s >",
		a.GetKey(),
		a.GetValue(),
	)
}

func AnnotationsString(a []*pb.Annotation) string {
	// TODO need to sort these first
	var as = make([]string, 1)
	for _, i := range a {
		as = append(as, AnnotationString(i))
	}
	return fmt.Sprintf(
		"<AnnoList %s >",
		strings.Join(as, ", "),
	)
}


type DataRepo interface {
	CreateSensor(s *pb.Sensor) (*pb.Sensor, error)
	GetSensors(s *pb.Sensor) ([]*pb.Sensor, error)
	UpdateSensors(q, u *pb.Sensor) ([]*pb.Sensor, error)
}

// interface for datastore/repo and pass one to sensor server to call
type SensorServer struct {
	logger Log
	data DataRepo
}

// TODO use real gRPC code + status for errors (wrap)

func (s *SensorServer) CreateSensor(c context.Context, n *pb.Sensor) (*pb.Sensor, error) {
	n, err := s.data.CreateSensor(n)
	s.logger.Printf("Created %+v in data repo", n)
	return n, err
}

func (s *SensorServer) ReadSensors(c context.Context, q *pb.Sensor) (*pb.Sensors, error) {
	l, err := s.data.GetSensors(q)
	s.logger.Printf("Got %+v from data repo", l)
	return &pb.Sensors{Count: int32(len(l)), Sensors: l}, err
}

func (s *SensorServer) UpdateSensors(c context.Context, u *pb.SensorUpdate) (*pb.Sensors, error) {
	selector := u.GetSelector()
	update := u.GetUpdate()
	if selector == nil || update == nil {
		return nil, status.Errorf(codes.InvalidArgument, "selector and update must be populated")
	}
	updated, err := s.data.UpdateSensors(selector, update)
	return &pb.Sensors{Count: int32(len(updated)), Sensors: updated}, err
}

//func (*UnimplementedSensorServiceServer) DeleteSensor(context.Context, *Sensor) (*Sensors, error) {
//        return nil, status.Errorf(codes.Unimplemented, "method DeleteSensor not implemented")
//}
