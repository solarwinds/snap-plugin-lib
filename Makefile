# File managed by pluginsync

#
# Copyright 2015 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
# Copyright (c) 2020 SolarWinds Worldwide, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

default:
	$(MAKE) deps
deps:
	bash -c "./scripts/deps.sh"
test:
	bash -c "./scripts/test.sh $(TEST_TYPE)"
test-small:
	bash -c "./scripts/test.sh small"
test-medium:
	bash -c "./scripts/test.sh medium"
test-all:
	$(MAKE) test-small
	$(MAKE) test-medium
check:
	$(MAKE) test

build:
	$(MAKE) deps
	bash -c "./scripts/build.sh"

