# Takyon

> *A filesystem-driven container runtime for multi-environment Linux workflows.*

Takyon is a lightweight environment manager built on top of **disk images + namespaces**.
It lets you run multiple isolated Linux environments — dev, pentest, gaming, desktop — without polluting your base system.

No daemons. No heavy abstraction.
Just **images, mounts, and control.**

---

## ✨ Concept

Takyon treats environments as **first-class filesystem objects**:

* A container = **a disk image**
* Running = **mounted + namespaced**
* Isolation = **Linux primitives (mount, PID, etc.)**

This means:

> You don’t manage containers.
> You manage *states of filesystems*.

---

## 🧱 Architecture

```tree
/var/lib/takyon/
    └── images/
        ├── kali.img
        ├── dev.img
        └── arch.img

/mnt/takyon/
└── kali/
└── dev/
```

* Images are stored in `/var/lib/takyon`
* Mounted environments live in `/mnt/takyon`
* A shared directory can be injected across environments

---

## ⚙️ Installation

```bash
git clone https://github.com/yourname/takyon
cd takyon
go build -o takyon
sudo mv takyon /usr/local/bin/
```

---

## 🚀 Usage

### Create a container

```bash
takyon create my-env -s 2048 -f ext4
```

* `-s` → size in MB
* `-f` → filesystem format

---

### Flash a container (install a distro)

```bash
takyon flash my-env -d debian
takyon flash kali-env -d kali
```

Supported:

* `debian`
* `ubuntu`
* `kali`
* `arch`

---

### Mount a container

```bash
takyon mount my-env
```

---

### List containers

```bash
takyon list
```

Example output:

```plaintext
[→] dev           (mounted)   -> /mnt/takyon/dev
[→] kali          (idle)      -> /mnt/takyon/kali
[→] broken-env    (corrupted) -> /mnt/takyon/broken-env
```

---

## Lifecycle

| State     | Meaning                      |
| --------- | ---------------------------- |
| idle      | image exists, not mounted    |
| mounted   | environment is active        |
| corrupted | invalid or broken filesystem |

---

## Internals

Takyon relies on:

* `truncate` → create disk images
* `mkfs.*` → format filesystem
* `mount (loop)` → attach image
* `debootstrap` / `pacstrap` → install distro

No virtualization. No emulation.
Just native Linux.

---

## Sudo & Paths

Takyon avoids `$HOME` to prevent permission issues when using `sudo`.

Defaults:

```plaintext
/var/lib/takyon        → storage
/mnt/takyon            → mountpoints
```

Override with:

```bash
export TAKYON_CONTAINER_DIR_PATH=...
export TAKYON_CONTAINER_MOUNT_PATH=...
export TAKYON_SHARED_DIR=...
```

---

## Philosophy

Takyon is built on a simple idea:

> *Isolation should be explicit, not hidden behind layers.*

Instead of abstracting Linux, Takyon embraces it:

* You can inspect everything
* You can manually mount images
* You can debug with standard tools

It’s closer to **Arch + chroot + namespaces** than to Docker.

---

## Roadmap

* [ ] Namespace isolation (PID, NET, UTS)
* [ ] Overlay support (layered environments)
* [ ] Shared directory binding
* [ ] Snapshot / rollback system
* [ ] `takyon doctor` (health checks)
* [ ] CLI UX polish

---

## 🧠 Why Takyon?

Because sometimes you don’t want:

* containers that hide everything
* or VMs that cost too much

You just want:

> *multiple clean Linux worlds, coexisting, without conflict.*

---

## ⚠️ Disclaimer

Takyon is experimental.
You are working directly with filesystems and mounts.

Break things. Learn things. Fix them.

---

## 🧃 Closing note

Takyon sits in a strange place:

* lighter than Docker
* more structured than raw chroot
* more transparent than most tools
