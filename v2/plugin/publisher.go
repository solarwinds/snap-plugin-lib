/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

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

package plugin

type Publisher interface {
	Publish(ctx PublishContext) error
}

type LoadablePublisher interface {
	Publisher
	Load(ctx Context) error
}

type UnloadablePublisher interface {
	Publisher
	Unload(ctx Context) error
}

type DefinablePublisher interface {
	Publisher
	PluginDefinition(def PublisherDefinition) error
}

type CustomizableInfoPublisher interface {
	Publisher
	CustomInfo(ctx Context) interface{}
}

type PublishContext interface {
	Context

	ListAllMetrics() []Metric
	Count() int
}

// PublisherDefinition provides API for specifying plugin (publisher) metadata (supported metrics, descriptions etc)
type PublisherDefinition interface {
	Definition

	// Define example config (which will be presented when example task is printed)
	DefineExampleConfig(cfg string) error
}
