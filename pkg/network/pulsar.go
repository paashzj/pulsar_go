// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package network

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/paashzj/pulsar_go/pkg/api"
	pb "github.com/paashzj/pulsar_go/pkg/internal/pulsar_proto"
	"github.com/paashzj/pulsar_go/pkg/util"
	"github.com/panjf2000/gnet"
	"github.com/sirupsen/logrus"
)

func Run(networkConfig *api.NetworkConfig, impl api.PulsarServer) error {
	server := &Server{
		EventServer: nil,
		pulsarImpl:  impl,
	}
	go func() {
		err := gnet.Serve(server, fmt.Sprintf("tcp://%s:%d", networkConfig.ListenHost, networkConfig.ListenTcpPort), gnet.WithMulticore(networkConfig.MultiCore), gnet.WithCodec(util.Codec))
		logrus.Error("pulsar broker started error ", err)
	}()
	return nil
}

type Server struct {
	*gnet.EventServer
	pulsarImpl api.PulsarServer
}

func (s *Server) OnInitComplete(server gnet.Server) (action gnet.Action) {
	logrus.Info("Pulsar Server Started")
	return
}

func (s *Server) React(frame []byte, c gnet.Conn) ([]byte, gnet.Action) {
	cmd := &pb.BaseCommand{}
	err := proto.Unmarshal(frame[4:], cmd)
	if err != nil {
		logrus.Error("marshal request error ", err)
		return nil, gnet.Close
	}
	switch *cmd.Type {
	case pb.BaseCommand_CONNECT:
		connected, err := s.pulsarImpl.Connect(cmd.Connect)
		if err != nil {
			logrus.Error("execute error ", err)
			return nil, gnet.Close
		}
		marshal, err := connected.Marshal()
		if err != nil {
			logrus.Error("marshal error ", cmd.Type)
			return nil, gnet.Close
		}
		return marshal, gnet.None
	default:
		break
	}
	logrus.Error("unsupported protocol ", cmd.Type)
	return nil, gnet.Close
}

func (s *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	logrus.Info("new connection connected ", " from ", c.RemoteAddr())
	return
}

func (s *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	logrus.Info("connection closed from ", c.RemoteAddr())
	return
}
