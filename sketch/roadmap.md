# Takyon V2 Roadmap

## Conatiner Storage & Representation

### 1. Store layout

Eventually, will respect the OCI specification, for the store should look like :

```tree
store/
  images/
    <image-id>/
      rootfs.img   # readonly
      meta.yaml

  containers/
    <container-id>/
      meta.yaml
      writable.img
      parent: <image-id>

  registry.yaml
```

> For human readability `yaml` is prefered hover `json` for metadata

### 2. Id & Registry File

An `<id>` is a unique identifier used by `takyon` to reference a container in the
`registry.yaml`

```yaml
containers:
  black-arch: <id>
  ubuntu-minimal: <id>

images:
    arch-linux: <id>

entries:
  <id>:
    created-at: ...
    type: container

  <id>:
    created-at: ...
    type: image
```

### 3. Meta.yaml File

tis file (`meta.yaml`) is used to store data about a container img,
and how to run it.

```yaml
name: black-arch
config:
    entrypoint: /bin/bash
    default-user: root
    env:
        DISPLAY: null
    runtime:
        namespace:
            pid: true
            net: true
            mnt: true
```

### 4. Export Format

Container sould be exportable as `tar` file

## Mvp specification

### 1. Objective

On the moment takyon is meant as a multi os (linux) runner, isolation, testing,
and it being easy to use are the main the goals

### 2. Spec

**Basic**:

```sh
# a pull create a readonly fs
takyon pull arch-linux http://rootfs-download-link
takyon pull ubuntu ./ubuntu-am64-rootfs.tar
```

```sh
takyon create balck-arch 16G
takyon create white-arch 16G --from arch-linux
```

```sh
takyon import ./ubuntu-min.tar
takyon export white-arch ./w-arch.tar
```

```sh
takyon exec black-arch ./myscript.sh
takyon exec black-arch "echo hi"
```

```sh
takyon enter black-arch --env MYVAR=VAl,V2=VAL2
takyon enter black-arch --instance
```

```sh
# foking freeze the source as readonly, create a new empty image that has it for parent
takyon fork black-arch white-arch
```

```sh
takyon config black-arch --default-user lix
takyon config black-arch --entrypoint /bin/fish
takyon config black-arch --isolate net,mount,pid
```

```sh
takyon display-server # never use as root
takyon run black-arch --app wireshark
takyon run minimal-ubuntu --desktop kde
takyon run minimal-ubuntu --desktop kde --display :1
```

