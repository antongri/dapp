---
title: Supported directives
sidebar: reference
permalink: reference/build/supported_directives.html
---

<div class="tab">
  <button class="tablinks active" onclick="openTab(event, 'Dimg')">Dimg</button>
  <button class="tablinks" onclick="openTab(event, 'Artifact')">Artifact</button>
</div>

<div id="Dimg" class="tabcontent active" markdown="1">

```yaml
dimg: <dimg_name... || ~>
from: <image>
fromCacheVersion: <version>
fromDimg: <dimg_name>
fromDimgArtifact: <artifact_name>
git:
# local git
- as: <custom_name>
  add: <absolute_path>
  to: <absolute_path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative_path or glob>
  excludePaths:
  - <relative_path or glob>
  stageDependencies:
    install:
    - <relative_path or glob>
    beforeSetup:
    - <relative_path or glob>
    setup:
    - <relative_path or glob>
# remote git
- url: <git_repo_url>
  branch: <branch_name>
  commit: <commit>
  tag: <tag>
  as: <custom_name>
  add: <absolute_path>
  to: <absolute_path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative_path or glob>
  excludePaths:
  - <relative_path or glob>
  stageDependencies:
    install:
    - <relative_path or glob>
    beforeSetup:
    - <relative_path or glob>
    setup:
    - <relative_path or glob>
import:
- artifact: <artifact_name>
  before: <install || setup>
  after: <install || setup>
  add: <absolute_path>
  to: <absolute_path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative_path or glob>
  excludePaths:
  - <relative_path or glob>
shell:
  beforeInstall:
  - <cmd>
  install:
  - <cmd>
  beforeSetup:
  - <cmd>
  setup:
  - <cmd>
  cacheVersion: <version>
  beforeInstallCacheVersion: <version>
  installCacheVersion: <version>
  beforeSetupCacheVersion: <version>
  setupCacheVersion: <version>
ansible:
  beforeInstall:
  - <task>
  install:
  - <task>
  beforeSetup:
  - <task>
  setup:
  - <task>
  cacheVersion: <version>
  beforeInstallCacheVersion: <version>
  installCacheVersion: <version>
  beforeSetupCacheVersion: <version>
  setupCacheVersion: <version>
mount:
- from: build_dir
  to: <absolute_path>
- from: tmp_dir
  to: <absolute_path>
- fromPath: <absolute_path>
  to: <absolute_path>
docker:
  VOLUME:
  - <volume>
  EXPOSE:
  - <expose>
  ENV:
    <env_name>: <env_value>
  LABEL:
    <label_name>: <label_value>
  ENTRYPOINT:
  - <entrypoint>
  CMD:
  - <cmd>
  ONBUILD:
  - <onbuild>
  WORKDIR: <workdir>
  USER: <user>
asLayers: <false || true>
```

</div>

<div id="Artifact" class="tabcontent" markdown="1">

```yaml
artifact: <artifact_name>
from: <image>
fromCacheVersion: <version>
fromDimg: <dimg_name>
fromDimgArtifact: <artifact_name>
git:
# local git
- as: <custom_name>
  add: <absolute_path>
  to: <absolute_path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative_path or glob>
  excludePaths:
  - <relative_path or glob>
  stageDependencies:
    install:
    - <relative_path or glob>
    beforeSetup:
    - <relative_path or glob>
    setup:
    - <relative_path or glob>
    buildArtifact:
    - <relative_path or glob>
# remote git
- url: <git_repo_url>
  branch: <branch_name>
  commit: <commit>
  tag: <tag>
  as: <custom_name>
  add: <absolute_path>
  to: <absolute_path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative_path or glob>
  excludePaths:
  - <relative_path or glob>
  stageDependencies:
    install:
    - <relative_path or glob>
    beforeSetup:
    - <relative_path or glob>
    setup:
    - <relative_path or glob>
    buildArtifact:
    - <relative_path or glob>
shell:
  beforeInstall:
  - <cmd>
  install:
  - <cmd>
  beforeSetup:
  - <cmd>
  setup:
  - <cmd>
  buildArtifact:
  - <cmd>
  cacheVersion: <version>
  beforeInstallCacheVersion: <version>
  installCacheVersion: <version>
  beforeSetupCacheVersion: <version>
  setupCacheVersion: <version>
  buildArtifactCacheVersion: <version>
ansible:
  beforeInstall:
  - <task>
  install:
  - <task>
  beforeSetup:
  - <task>
  setup:
  - <task>
  buildArtifact:
  - <task>
  cacheVersion: <version>
  beforeInstallCacheVersion: <version>
  installCacheVersion: <version>
  beforeSetupCacheVersion: <version>
  setupCacheVersion: <version>
  buildArtifactCacheVersion: <version>
mount:
- from: build_dir
  to: <absolute_path>
- from: tmp_dir
  to: <absolute_path>
- fromPath: <absolute_path>
  to: <absolute_path>
asLayers: <false || true>
```

</div>
