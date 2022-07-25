/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/loopholelabs/polyglot-go"
)

type testDataGenerator struct {
	Name         string      `json:"name"`
	Kind         byte        `json:"kind"`
	DecodedValue interface{} `json:"decodedValue"`

	EncodedValue string `json:"encodedValue"`

	generate func(*polyglot.Buffer) error `json:"-"`
}

func main() {
	out := flag.String("out", "out/polyglot-test-data.json", "File to write test data to")

	flag.Parse()

	testData := []*testDataGenerator{
		{
			Name:         "None",
			Kind:         polyglot.NilKind[0],
			DecodedValue: nil,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Nil()

				return nil
			},
		},
		{
			Name:         "true Bool",
			Kind:         polyglot.BoolKind[0],
			DecodedValue: true,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Bool(true)

				return nil
			},
		},
		{
			Name:         "false Bool",
			Kind:         polyglot.BoolKind[0],
			DecodedValue: false,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Bool(false)

				return nil
			},
		},
		{
			Name:         "U8",
			Kind:         polyglot.Uint8Kind[0],
			DecodedValue: 32,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Uint8(32)

				return nil
			},
		},
		{
			Name:         "U16",
			Kind:         polyglot.Uint16Kind[0],
			DecodedValue: 1024,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Uint16(1024)

				return nil
			},
		},
		{
			Name:         "U32",
			Kind:         polyglot.Uint32Kind[0],
			DecodedValue: 4294967290,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Uint32(4294967290)

				return nil
			},
		},
		{
			Name:         "U64",
			Kind:         polyglot.Uint64Kind[0],
			DecodedValue: uint64(18446744073709551610),

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Uint64(18446744073709551610)

				return nil
			},
		},
		{
			Name:         "I32",
			Kind:         polyglot.Int32Kind[0],
			DecodedValue: -2147483648,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Int32(-2147483648)

				return nil
			},
		},
		{
			Name:         "I64",
			Kind:         polyglot.Int64Kind[0],
			DecodedValue: -9223372036854775808,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Int64(-9223372036854775808)

				return nil
			},
		},
		{
			Name:         "F32",
			Kind:         polyglot.Float32Kind[0],
			DecodedValue: -214648.34432,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Float32(-214648.34432)

				return nil
			},
		},
		{
			Name:         "F64",
			Kind:         polyglot.Float64Kind[0],
			DecodedValue: -922337203685.2345,

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Float64(-922337203685.2345)

				return nil
			},
		},
		{
			Name:         "Array",
			Kind:         polyglot.SliceKind[0],
			DecodedValue: []string{"1", "2", "3"},

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.
					Slice(3, polyglot.StringKind).
					String("1").
					String("2").
					String("3")

				return nil
			},
		},
		{
			Name: "Map",
			Kind: polyglot.MapKind[0],
			DecodedValue: map[string]int{
				"1": 1,
				"2": 2,
				"3": 3,
			},

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.
					Map(3, polyglot.StringKind, polyglot.Uint32Kind).
					String("1").
					Uint32(1).
					String("2").
					Uint32(2).
					String("3").
					Uint32(3)

				return nil
			},
		},
		{
			Name:         "Bytes",
			Kind:         polyglot.BytesKind[0],
			DecodedValue: []byte("Test String"),

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Bytes([]byte("Test String"))

				return nil
			},
		},
		{
			Name:         "String",
			Kind:         polyglot.StringKind[0],
			DecodedValue: "Test String",

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.String("Test String")

				return nil
			},
		},
		{
			Name:         "Error",
			Kind:         polyglot.ErrorKind[0],
			DecodedValue: "Test String",

			generate: func(b *polyglot.Buffer) error {
				encoder := polyglot.Encoder(b)

				encoder.Error(errors.New("Test String"))

				return nil
			},
		},
	}

	for _, data := range testData {
		buffer := polyglot.NewBuffer()

		data.generate(buffer)

		data.EncodedValue = base64.StdEncoding.EncodeToString(*buffer)

		log.Println("Generated data for test", data.Name)
	}

	if err := os.MkdirAll(filepath.Dir(*out), os.ModePerm); err != nil {
		panic(err)
	}

	content, err := json.Marshal(testData)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(*out, content, os.ModePerm); err != nil {
		panic(err)
	}
}
