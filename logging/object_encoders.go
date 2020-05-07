/*
Copyright 2020 The Knative Authors

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

package logging

import (
	"strings"

	"go.uber.org/zap/zapcore"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
)

// This file contains the specific object encoders for use in Knative
// to optimize logging experience and performance.

// StringSet returns a marshaler for the set of strings.
func StringSet(s sets.String) zapcore.ObjectMarshalerFunc {
	return func(enc zapcore.ObjectEncoder) error {
		enc.AddString("keys", strings.Join(s.UnsortedList(), ","))
		return nil
	}
}

// NamespacedName returns a marshaler for NamespacedName.
func NamespacedName(n types.NamespacedName) zapcore.ObjectMarshalerFunc {
	return func(enc zapcore.ObjectEncoder) error {
		enc.AddString("key", n.Namespace+"/"+n.Name)
		return nil
	}
}
