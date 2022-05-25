//go:build bench
// +build bench

/*
Benchmark test mainly measures:
- the time needed to create validation tree (1-33 metrics)
- the time needed to validate metrics based only on definition tree (100-10000)
- the time needed to validate metrics based on definition and filtering tree
*/

/*
 Copyright (c) 2022 SolarWinds Worldwide, LLC

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

package runner

import (
	"strings"
	"testing"
	"time"

	"github.com/solarwinds/snap-plugin-lib/v2/internal/plugins/collector/proxy"
	"github.com/solarwinds/snap-plugin-lib/v2/plugin"
)

const collectTimeout = 10 * time.Second

var metricDefinition = []string{
	"/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Pending",
	"/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Running",
	"/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Succeeded",
	"/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Failed",
	"/kubernetes/pod/[node]/[namespace]/[pod]/status/phase/Unknown",
	"/kubernetes/pod/[node]/[namespace]/[pod]/status/condition/ready",
	"/kubernetes/pod/[node]/[namespace]/[pod]/status/condition/scheduled",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/restarts",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/ready",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/waiting",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/running",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/status/terminated",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/requested/cpu/cores",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/requested/memory/bytes",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/limits/cpu/cores",
	"/kubernetes/container/[namespace]/[node]/[pod]/[container]/limits/memory/bytes",
	"/kubernetes/node/[node]/spec/unschedulable",
	"/kubernetes/node/[node]/status/outofdisk",
	"/kubernetes/node/[node]/status/allocatable/cpu/cores",
	"/kubernetes/node/[node]/status/allocatable/memory/bytes",
	"/kubernetes/node/[node]/status/allocatable/pods",
	"/kubernetes/node/[node]/status/capacity/cpu/cores",
	"/kubernetes/node/[node]/status/capacity/memory/bytes",
	"/kubernetes/node/[node]/status/capacity/pods",
	"/kubernetes/deployment/[namespace]/[deployment]/metadata/generation",
	"/kubernetes/deployment/[namespace]/[deployment]/status/observedgeneration",
	"/kubernetes/deployment/[namespace]/[deployment]/status/targetedreplicas",
	"/kubernetes/deployment/[namespace]/[deployment]/status/availablereplicas",
	"/kubernetes/deployment/[namespace]/[deployment]/status/unavailablereplicas",
	"/kubernetes/deployment/[namespace]/[deployment]/status/updatedreplicas",
	"/kubernetes/deployment/[namespace]/[deployment]/status/deploynotfinished",
	"/kubernetes/deployment/[namespace]/[deployment]/spec/desiredreplicas",
	"/kubernetes/deployment/[namespace]/[deployment]/spec/paused",
}

var nodesDef = []string{ // len = 30
	"Node-Dale-5fe0ba", "Node-Brinda-4dcd07", "Node-Leora-486548", "Node-Ernestine-07bb0f", "Node-Velvet-cc08e8",
	"Node-Kamilah-b005c0", "Node-Hyon-9188da", "Node-Velva-3240d3", "Node-Julian-8f0f48", "Node-Fausto-e5b66c",
	"Node-Lynn-3a8117", "Node-January-93bf5e", "Node-Jeramy-80c412", "Node-Marci-f804e3", "Node-Donald-b6f2b9",
	"Node-Lindsey-736c74", "Node-Gwyneth-cdff3d", "Node-Lilly-70a14b", "Node-Sun-cf799d", "Node-Thurman-9de7aa",
	"Node-Pandora-d2c4dc", "Node-Luke-005eb8", "Node-Sonia-b87d0e", "Node-Kiara-389341", "Node-Kacie-0e560f",
	"Node-Britni-8781e0", "Node-Yon-8add7e", "Node-Irene-6840c8", "Node-Heidy-80a200", "Node-Donald-f12bac",
}

var nsDef = []string{ // len = 45
	"ns-moph-35", "ns-sightly-77", "ns-grumous-38", "ns-pragmatize-25", "ns-schistocoelia-54",
	"ns-anthroponomy-7", "ns-astrut-35", "ns-clonus-3", "ns-sulphureosuffused-85", "ns-scrappet-95",
	"ns-rebasis-73", "ns-parasternal-26", "ns-frangent-33", "ns-unquelled-53", "ns-blowhole-39",
	"ns-Atka-29", "ns-upsilon-56", "ns-epidermous-53", "ns-Tonna-16", "ns-Akhissar-71",
	"ns-yokeage-60", "ns-semianatomical-73", "ns-atmologist-84", "ns-nonconterminous-66", "ns-undulatingly-74",
	"ns-amoebian-87", "ns-panhuman-18", "ns-symphonically-12", "ns-toileted-24", "ns-congruently-59",
	"ns-saccharulmic-59", "ns-panosteitis-49", "ns-shellman-34", "ns-eximiously-70", "ns-platyglossate-95",
	"ns-unbalance-29", "ns-pleurodiscous-3", "ns-Ghaznevid-58", "ns-bedress-79", "ns-sectarial-18",
	"ns-Chamidae-9", "ns-relay-21", "ns-spherical-7", "ns-perchlorate-56", "ns-hoyle-53",
}

var podDef = []string{
	"pod-France-6e74a0de", "pod-Kyrgyzstan-85ce76de", "pod-Indonesia-193f0132", "pod-SriLanka-6916239b", "pod-Portugal-b26229f2",
	"pod-Mozambique-42fe8f00", "pod-Bahamas-89cc5d5c", "pod-Lesotho-2d9f2585", "pod-Lesotho-054cd6cb", "pod-Grenada-695e912b",
}

var contDef = []string{ // len = 40
	"cont-Hungary-2020-07-11", "cont-Mongolia-2025-02-20", "cont-Colombia-2023-11-13", "cont-Syria-2016-01-20", "cont-SriLanka-2020-09-03",
	"cont-Guinea-2009-10-24", "cont-PapuaNewGuinea-2029-08-21", "cont-Guinea-Bissau-2019-03-27", "cont-Oman-2020-07-27", "cont-Pakistan-2019-04-20",
	"cont-Jordan-2013-04-02", "cont-Malta-2019-12-11", "cont-Thailand-2019-09-15", "cont-Tanzania-2012-07-20", "cont-Germany-2028-09-12",
	"cont-Australia-2029-05-03", "cont-Niger-2010-01-18", "cont-Botswana-2018-12-07", "cont-Monaco-2024-04-02", "cont-Burundi-2011-04-22",
	"cont-SouthSudan-2027-07-09", "cont-Australia-2023-03-21", "cont-Cambodia-2015-11-17", "cont-Sweden-2015-01-04", "cont-Tajikistan-2027-12-15",
	"cont-AntiguaBarbuda-2021-08-19", "cont-Zimbabwe-2026-10-08", "cont-Guatemala-2013-04-25", "cont-Laos-2009-10-28", "cont-Turkey-2018-04-14",
	"cont-Portugal-2029-04-11", "cont-Timor-Leste-2014-06-06", "cont-StVincentandTheGrenadines-2018-11-16", "cont-Netherlands-2014-03-11", "cont-Lebanon-2029-04-19",
	"cont-Qatar-2011-06-21", "cont-MarshallIslands-2018-06-16", "cont-Canada-2026-07-11", "cont-Andorra-2011-04-19", "cont-Zambia-2009-05-10",
}

var metricsToValidateAll []string

func init() {
	// create any possible metric combination
	for _, mt := range metricDefinition {
		metricsToValidateAll = append(metricsToValidateAll, mt)
	}

	metricsToValidateAll = generateCombination(metricsToValidateAll, "[node]", nodesDef)
	metricsToValidateAll = generateCombination(metricsToValidateAll, "[namespace]", nsDef)
	metricsToValidateAll = generateCombination(metricsToValidateAll, "[pod]", podDef)
	metricsToValidateAll = generateCombination(metricsToValidateAll, "[container]", contDef)
}

func generateCombination(mts []string, s string, choices []string) []string {
	resMt := []string{}
	for _, mt := range mts {
		if strings.Index(mt, s) != -1 {
			for _, ch := range choices {
				resMt = append(resMt, strings.Replace(mt, s, ch, 1))
			}
		}
	}
	return resMt
}

///////////////////////////////////////////////////////////////////////////////

type benchCollector struct {
	metricDefined  int // numbers of metrics to define
	addMetricRatio int // every n metric will be taken to result
}

func (bc *benchCollector) DefineMetrics(colDef plugin.CollectorDefinition) error {
	for i := 0; i < bc.metricDefined; i++ {
		colDef.DefineMetric(metricDefinition[i], "", true, "")
	}

	return nil
}

func (bc *benchCollector) Collect(ctx plugin.CollectContext) error {
	for i := 0; i < len(metricsToValidateAll); i += bc.addMetricRatio {
		err := ctx.AddMetric(metricsToValidateAll[i], i)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func genParseDefinitionN(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchColl := &benchCollector{
			metricDefined:  n,
			addMetricRatio: 1,
		}
		proxy.NewContextManager(benchColl, "benchmark collector", "0.0.1") // build metrics definition tree
	}
}

func BenchmarkParseDefinition_1Metric(b *testing.B)           { genParseDefinitionN(1, b) }
func BenchmarkParseDefinition_1Group_7Metrics(b *testing.B)   { genParseDefinitionN(7, b) }
func BenchmarkParseDefinition_2Groups_16Metrics(b *testing.B) { genParseDefinitionN(16, b) }
func BenchmarkParseDefinition_3Groups_24Metrics(b *testing.B) { genParseDefinitionN(24, b) }
func BenchmarkParseDefinition_4Groups_33Metrics(b *testing.B) { genParseDefinitionN(33, b) }

///////////////////////////////////////////////////////////////////////////////

func genMetricAddition(addRatio int, b *testing.B) {
	const taskId = 1

	for i := 0; i < b.N; i++ {
		benchColl := &benchCollector{
			metricDefined:  33,
			addMetricRatio: addRatio,
		}
		ctxMan := proxy.NewContextManager(benchColl, "benchmark collector", "0.0.2")
		err := ctxMan.LoadTask(taskId, []byte("{}"), []string{})
		if err != nil {
			panic(err)
		}

		b.ResetTimer()
		chunkCh = ctxMan.RequestCollect(taskId)

		select {
		case <-chunkCh:
		// ok
		case <-time.After(collectTimeout):
			panic("timeout occurred")
		}

		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFilterMetrics_All_25p(b *testing.B)  { genMetricAddition(4, b) }
func BenchmarkFilterMetrics_All_50p(b *testing.B)  { genMetricAddition(2, b) }
func BenchmarkFilterMetrics_All_100p(b *testing.B) { genMetricAddition(1, b) }

///////////////////////////////////////////////////////////////////////////////
