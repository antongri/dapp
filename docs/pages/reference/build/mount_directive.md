---
title: Using mounts
sidebar: reference
permalink: reference/build/mount_directive.html
---

There are three types of mount you can use in dappfile:
1. Using `tmp_dir`

    The `tmp_dir` is a temporary folder which is created by dapp for every build process of a specific dimg or artifact and is deleted after the build process ends. The `tmp_dir` can be accessed only from stages of a specific dimg or artifact it was created for. The `tmp_dir` can be mounted to a specified absolute path, and it will be the one shared folder for all of stages in the time of build.

    It is convenient to use `tmp_dir` for folders or files which needn't exist in the image after the build. E.g it can be a list of applications from `apt`.

    Here are the example of using `tmp_dir`:

    ```yaml
    dimg: ~
    from: ubuntu:16.04
    mount:
    - from: tmp_dir
      to: /var/lib/apt/lists/
```

1. Using `build_dir`

    The `build_dir` is a temporary folder which is created by dapp. This special folder is shared for every dimgs and artifacts described in one dappfile (a [project]({{ site.baseurl }}/reference/glossary.html#project)). The `build_dir` can be accessed in the time of build from all stages of all dimgs and artifacts of a project and can be mounted to a specified absolute path. This folder won't be deleted after the build process ends, and dapp will use it again in the time of next build. This folder is located in the `~/.dapp/builds/<dapp name>` folder, and can be manually deleted.

    It is convenient to use `build_dir` for folders or files which you need to share between images in a one project. E.g it can be an `apt` cache of downloaded packages, located in the `/var/cache/apt` folder.

    Here are the example of using `build_dir` for caching `/var/cache/apt`:

    ```yaml
    dimg: ~
    from: ubuntu:16.04
    mount:
    - from: build_dir
      to: /var/cache/apt
```

1. Using absolute path to the file.

    By using `mount.fromPath` directive you can mount a file or a folder on your host to path in a specific image in the time of build. A file or a folder will be mounted for read/write, and will be unmounted in the end of the build process (you should copy files you want to be kept in the resulting image).

    Here are the example of mounting a file:

    ```yaml
    dimg: ~
    from: ubuntu:16.04
    mount:
    - fromPath: /buildtime/libs
      to: /app/libs
```

In any cases, mounting folder will clear the mount point in the image, and also, in the resulting image there will be a clean folder. Therefore, if you need keep mounted files or folders in the built image you should copy them.

## Syntax

```yaml
mount:
- from: build_dir
  to: <absolute_path>
- from: tmp_dir
  to: <absolute_path>
- fromPath: <absolute_path>
  to: <absolute_path>
```
