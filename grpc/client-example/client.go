/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2020 Arduino LLC (http://www.arduino.cc/)
 */

// Package main implements a simple gRPC client that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It interacts with the route guide service whose definition can be found in routeguide/route_guide.proto.
package main

import (
	"io"
	"log"

	pb "github.com/arduino/arduino-builder/grpc/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// printFeature gets the feature for the given point.
func autocomplete(client pb.BuilderClient, in *pb.BuildParams) {
	resp, err := client.Autocomplete(context.Background(), in)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	log.Println(resp)
}

// printFeatures lists all the features within the given bounding Rectangle.
func build(client pb.BuilderClient, in *pb.BuildParams) {
	stream, err := client.Build(context.Background(), in)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		line, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		log.Println(line)
	}
}

func main() {
	conn, err := grpc.Dial("localhost:12345", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewBuilderClient(conn)

	exampleParames := pb.BuildParams{
		BuiltInLibrariesFolders: "/ssd/Arduino-master/build/linux/work/libraries",
		CustomBuildProperties:   "build.warn_data_percentage=75",
		FQBN:                    "arduino:avr:mega:cpu=atmega2560",
		HardwareFolders:         "/ssd/Arduino-master/build/linux/work/hardware,/home/martino/.arduino15/packages,/home/martino/eslov-sk/hardware",
		OtherLibrariesFolders:   "/home/martino/eslov-sk/libraries",
		ArduinoAPIVersion:       "10805",
		SketchLocation:          "/home/martino/eslov-sk/libraries/WiFi101/examples/ScanNetworks/ScanNetworks.ino",
		ToolsFolders:            "/ssd/Arduino-master/build/linux/work/tools-builder,/ssd/Arduino-master/build/linux/work/hardware/tools/avr,/home/martino/.arduino15/packages",
		Verbose:                 true,
		WarningsLevel:           "all",
		BuildCachePath:          "/tmp/arduino_cache_761418/",
		CodeCompleteAt:          "/home/martino/eslov-sk/libraries/WiFi101/examples/ScanNetworks/ScanNetworks.ino:56:9",
	}

	//build(client, &exampleParames)
	autocomplete(client, &exampleParames)
}
