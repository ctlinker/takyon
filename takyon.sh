#!/usr/bin/env bash
set -eu

ENV_STORAGE_FOLDER="$HOME/.local/share/tackyon/"
ENV_MOUNT_DIR="/mnt/"

# --- UTILS ---

setup_dirs() {
    mkdir -p "$ENV_STORAGE_FOLDER"
    mkdir -p "$ENV_MOUNT_DIR"
}

is_empty_str() {
    [ -z "$1" ]
}

file_exist() {
    [ -f "$1" ]
}

dir_exist() {
    [ -d "$1" ]
}

print_info() {
    echo "[*] $1"
}

exit_error() {
    echo "[!] $1" > /dev/tty
    exit 1
}

internal_err() {
    echo "[!!!] Internal Error: $1" > /dev/tty
    exit 1
}

img_path_for() {
    env_name="$1"
    echo "$ENV_STORAGE_FOLDER/$env_name.img"
}

mount_path_for() {
    env_name="$1"
    echo "$ENV_MOUNT_DIR/$env_name"
}

# --- FUNCTIONS ---

create_disk_image() {
    [ $# -lt 1 ] && internal_err "Usage: create_disk_image [NAME] [SIZE] [FORMAT=ext4]"
    env_name="$1"
    env_size="${2:-2048}"
    disk_format="${3:-ext4}"
    disk_img_path="$(img_path_for "$env_name")"

    if is_empty_str "$env_name"; then
        exit_error "Env name is empty"
    fi

    if file_exist "$disk_img_path"; then
        exit_error "Env of this name already exists"
    fi

    if ! [[ "$env_size" =~ ^[0-9]+$ ]]; then
        exit_error "Invalid disk image size. Expecting an integer."
    fi

    case "$disk_format" in
        ext4|ext3|ext2)
            print_info "Creating disk image \"$env_name\" of size ${env_size}MB"
            ;;
        *)
            exit_error "Invalid Disk Image Format"
            ;;
    esac

    truncate -s "${env_size}M" "$disk_img_path" > /dev/null
    mkfs."$disk_format" "$disk_img_path" > /dev/null
}

resize_disk_image() {
    [ $# -lt 1 ] && internal_err "Usage: resize_disk_image [NAME] [SIZE] [FORMAT=ext4]"
    env_name="$1"
    env_size="${2:-2048}"
    disk_format="${3:-ext4}"
    disk_img_path="$(img_path_for "$env_name")"

    if is_empty_str "$env_name"; then
        exit_error "Env name is empty"
    fi

    if ! file_exist "$disk_img_path"; then
        exit_error "Env of this name does not exist"
    fi

    if ! [[ "$env_size" =~ ^[0-9]+$ ]]; then
        exit_error "Invalid disk image size. Expecting an integer."
    fi

    case "$disk_format" in
        ext4|ext3|ext2)
            print_info "Resizing disk image \"$env_name\" to size ${env_size}MB"
            ;;
        *)
            exit_error "Invalid Disk Image Format"
            ;;
    esac

    truncate -s "${env_size}M" "$disk_img_path" > /dev/null

        # Grow the filesystem to fill the new space
    if [[ "$disk_format" =~ ext[2-4] ]]; then
        print_info "Checking filesystem..."
        e2fsck -f "$disk_img_path" > /dev/null 2>&1 || internal_err "Filesystem check failed"
        print_info "Resizing filesystem..."
        resize2fs "$disk_img_path" > /dev/null || internal_err "Filesystem resize failed"
    fi
}

mount_disk_img() {
    [ $# -lt 1 ] && internal_err "Usage: mount_disk_img [ENV_NAME]"
    env_name="$1"
    disk_img_path="$(img_path_for "$env_name")"
    mount_path="$(mount_path_for "$env_name")"

    if ! file_exist "$disk_img_path"; then
        exit_error "Mount failed: Env of this name does not exist"
    fi

    if mountpoint -q "$mount_path"; then
        print_info "Target is already mounted"
        return 0
    fi

    mkdir -p "$mount_path" && \
    mount -o loop "$disk_img_path" "$mount_path" 
    
    if mountpoint -q "$mount_path"; then
        mkdir -p "$mount_path"/{dev,proc,sys,etc}
        mkdir -p "$mount_path/dev/pts"

        host_resolv="/etc/resolv.conf"
        mount_resolv="$mount_path/$host_resolv"

        if [ ! -f "$mount_resolv" ] && [ -f "$host_resolv" ]; then
            cp "$host_resolv" "$mount_resolv"
        fi

    else
        exit_error "Mount verification failed"
    fi

    print_info "Mounted $env_name at $mount_path"
}

umount_disk_image() {
    [ $# -lt 1 ] && internal_err "Usage: umount_disk_image [ENV_NAME]"

    env_name="$1"
    mount_path="$(mount_path_for "$env_name")"

    if ! mountpoint -q "$mount_path"; then
        exit_error "Target is not a mountpoint. Aborting to save your host disk!"
    fi

    umount -R "$mount_path"
    rm -rf "$mount_path"
}

enter_disk_image() {
    env_name="$1"
    extra_cmd="${2:-""}"
    mount_path="$(mount_path_for "$env_name")"

    if ! mountpoint -q "$mount_path"; then
        mount_disk_img "$env_name"
    fi

    # Use a subshell to ensure cleanup
    (
        trap 'sleep 1 && umount -R "$mount_path" && rm -rf $mount_path' EXIT
        
        mount_script="
            mount -t proc proc \"$mount_path/proc\" && \
            mount -t sysfs sys \"$mount_path/sys\" && \
            mount --bind /dev \"$mount_path/dev\" && \
            mount -t devpts devpts \"$mount_path/dev/pts\"
        "

        final_script=""

        if [ -z "$extra_cmd" ]; then
            final_script="
                $mount_script
                chroot \"$mount_path\" /bin/bash
            "
        else
            final_script="
                $mount_script
                chroot \"$mount_path\" /bin/bash -c \"$extra_cmd\"
            "
        fi
        
        unshare --mount --cgroup --pid --fork bash -c "$final_script"
    )
}

remove_disk_image(){
    disk_img_path="$(img_path_for "$1")"
    mount_path="$(mount_path_for "$1")"
    
    if ! file_exist "$disk_img_path"; then
        exit_error "Env of this name does not exist"
    fi

    if mountpoint -q "$mount_path"; then
        exit_error "This env is currently mounted, umount it first"
    fi 

    rm "$disk_img_path"
}

list_env() {
    [ -d "$ENV_STORAGE_FOLDER" ] || return

    for item in "$ENV_STORAGE_FOLDER"/*.img; do
        [ -e "$item" ] || continue

        env_name="$(basename "$item" .img)"
        mnt="$(mount_path_for "$env_name")"

        # Get size in MB (rounded)
        env_size_mb="$(du -m "$item" | cut -f1)"

        # Check mount status
        if mountpoint -q "$mnt" 2>/dev/null; then
            status="mounted"
        else
            status="idle"
        fi

        printf "  - %s (%s) - %sMB\n" "$env_name" "$status" "$env_size_mb"
    done
}

# --- COMMANDS ---

cmd_create() {
    env_name="$1"
    env_size="${2:-2048}"
    create_disk_image "$env_name" "$env_size"
}

cmd_mount() {
    env_name="$1"
    mount_disk_img "$env_name"
}

cmd_resize() {
    env_name="$1"
    env_size="${2:-2048}"
    resize_disk_image "$env_name" "$env_size"
}

cmd_enter() {
    env_name="$1"
    enter_disk_image "$env_name"
}

cmd_flash() {
    env_name="$1"
    mount_path="$(mount_path_for "$env_name")"
    shift

    if ! mountpoint -q "$mount_path"; then
        exit_error "Target is not a mountpoint. Aborting to save your host disk!"
    fi

    if [ $# -eq 0 ]; then
        pacstrap -K "$mount_path" base bash coreutils nano git sudo
    else
        pacstrap -K "$mount_path" "$@"
    fi
}

cmd_install_desktop() {
    env_name="$1"
    env_de="$2"
    packages=""

    case "$env_de" in
        xfce4)
            packages="xorg xorg-xinit xfce4 xfce4-goodies lightdm lightdm-gtk-greeter"
        ;;
        kde)
            packages="plasma kde-applications sddm xorg xorg-xinit"
            ;;
        *)
        exit_error "Unsupport desktop envirement"
        ;;
    esac

    pacman -S xorg-server-xephyr --needed
    # shellcheck disable=SC2086
    cmd_flash "$env_name" $packages
}

cmd_start() {
    env_name="$1"
    env_de="$2"
    display_num=":4"
    launcher=""

    case "$env_de" in
        xfce4)
            launcher="dbus-launch --exit-with-session startxfce4"
            ;;
        gnome)
            launcher="
            export XDG_RUNTIME_DIR=/tmp/runtime-root
            export XDG_SESSION_TYPE=x11
            mkdir -p \$XDG_RUNTIME_DIR
            dbus-run-session gnome-session
            "
            ;;
        *)
            exit_error "Unsupport desktop envirement"
            ;;
    esac

    print_info "Waiting for Xephyr to be ready..."
    while [ ! -e "/tmp/.X11-unix/X${display_num#:}" ]; do
        sleep 0.1
    done


    
    enter_disk_image "$env_name" "
    export DISPLAY=$display_num
    $launcher
    "

    print_info "Session closed"
}

cmd_nested_server() {
    if [ "$EUID" -eq 0 ]; then
        exit_error "DO NOT RUN THIS AS ROOT"
    fi
    display_num=":4"
    print_info "Launching Xephyr as $(whoami)..."
    Xephyr -br -ac -noreset -screen 1280x720 "$display_num" > /dev/null &
}

cmd_remount() {
    env_name="$1"
    umount_disk_image "$env_name"
    mount_disk_img "$env_name"
}

cmd_umount() {
    env_name="$1"
    umount_disk_image "$env_name"
}

cmd_remove() {
    env_name="$1"
    remove_disk_image "$env_name"
}

cmd_list() {
    list_env
}

cmd_help() {
    cat << EOF
Usage: $0 [command] [args]

Commands:
    create [NAME] [SIZE]    Create a disk image
    mount [NAME]            Mount a disk image
    umount [NAME]           Umount a disk image
    enter [NAME]            Enter a mounted disk image
    resize [NAME] [SIZE]    Resize a disk image
    remove [NAME]           Remove a disk image
    list                    List the current disk images
    flash [NAME] [pkgs...]  Install packages to a mounted disk image
EOF
}

# --- SETUP ---

dispatch() {
    cmd="cmd_$(echo "$1" | tr '-' '_')"
    shift
    if command -v "$cmd" >/dev/null 2>&1; then
        $cmd "$@"
    else
        echo "Unknown command: $cmd"
        cmd_help
        exit 1
    fi
}

main() {
    case "${1:-}" in
        -h|--help|help|"")
            cmd_help
            exit 0
            ;;
    esac
    setup_dirs
    dispatch "$@"
}

main "$@"
