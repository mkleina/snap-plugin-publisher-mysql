// +build integration

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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

package mysql

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMySQLPublish(t *testing.T) {
	var buf bytes.Buffer
	metrics := []plugin.MetricType{
		*plugin.NewMetricType(core.NewNamespace("test", "string"), time.Now(), map[string]string{core.STD_TAG_PLUGIN_RUNNING_ON: "127.0.0.1"}, "", "example_string"),
		*plugin.NewMetricType(core.NewNamespace("test", "int"), time.Now(), map[string]string{core.STD_TAG_PLUGIN_RUNNING_ON: "127.0.0.1"}, "", 1),
		*plugin.NewMetricType(core.NewNamespace("test", "string", "slice"), time.Now(), map[string]string{core.STD_TAG_PLUGIN_RUNNING_ON: "localhost"}, "", []string{"str1", "str2"}),
		*plugin.NewMetricType(core.NewNamespace("test", "string", "slice"), time.Now(), map[string]string{core.STD_TAG_PLUGIN_RUNNING_ON: "localhost"}, "", []int{1, 2}),
	}
	config := make(map[string]ctypes.ConfigValue)
	enc := gob.NewEncoder(&buf)
	enc.Encode(metrics)
	config["username"] = ctypes.ConfigValueStr{Value: "root"}
	config["password"] = ctypes.ConfigValueStr{Value: ""}
	config["database"] = ctypes.ConfigValueStr{Value: "snap_test"}
	config["tablename"] = ctypes.ConfigValueStr{Value: "info"}
	mp := NewMySQLPublisher()
	cp, _ := mp.GetConfigPolicy()
	cfg, _ := cp.Get([]string{""}).Process(config)
	Convey("Publish metrics to MySQL instance should succeed and not throw an error", t, func() {
		err := mp.Publish(plugin.SnapGOBContentType, buf.Bytes(), *cfg)
		So(err, ShouldBeNil)
	})
}
