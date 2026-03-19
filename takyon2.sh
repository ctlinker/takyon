#!/usr/bin/env sh
set -eu
ENV_STORAGE_FOLDER="./temp/envs"
ENV_MOUNT_DIR="/mnt"

# --- UTILS ---

setup_dirs() {
    mkdir -p "$ENV_STORAGE_FOLDER"
}

is_empty_str() {
    [ -z "$1" ] && return 0
    return 1
}

file_exist() {
    [ -f "$1" ] && return 0
    return 1
}

dir_exist() {
    [ -f "$1" ] && return 0
    return 1
}

print_info() {
    echo "[*] $1"
}

exit_error() {
    echo "[!] $1" > /dev/tty
    exit 1
}

internal_err(){
    echo "[!!!] Internal Error : $1" > /dev/tty
    exit 1
}

ns_path_for() { echo "/tmp/ns_$1"; }

img_path_for() {
    env_name="$1"
    disk_img_path="$ENV_STORAGE_FOLDER/$env_name.img"
    echo "$disk_img_path"
}

mount_path_for() {
    env_name="$1"
    mount_path="$ENV_MOUNT_DIR/$env_name"
    echo "$mount_path"
}

unshared_from() {
    env_name="$1"
    mnt_ns="$(ns_path_for "$env_name")_mnt"
    pid_ns="$(ns_path_for "$env_name")_pid"

    # Create empty files to act as anchors
    touch "$mnt_ns" "$pid_ns"

    # 1. We create the Mount namespace using the built-in helper
    # 2. We must manually bind the PID namespace because unshare won't do it for us
    # We run a background 'sleep' to keep the namespaces alive
    sudo unshare --mount="$mnt_ns" bash -c "
        mount --bind /proc/self/ns/pid $pid_ns
        exec sleep infinity
    " &
    
    # Give the background process a moment to initialize the bind mount
    sleep 0.1
}

use_namespace() {
    env_name="$1"
    cmd="$2"
    mnt_ns="$(ns_path_for "$env_name")_mnt"
    pid_ns="$(ns_path_for "$env_name")_pid"

    # nsenter DOES take the file paths for each specific namespace
    sudo nsenter --mount="$mnt_ns" "$cmd"
}

# --- FUNCTION ---



create_disk_image() {

    [ $# = 0 ] && internal_err "Usage: create_disk_image [NAME] [SIZE] [FORMAT=ext4]"

    env_name="$1"
    env_size="${2:-2048}"
    disk_format="${3:-"ext4"}"
    disk_img_path="$(img_path_for "$env_name")"

    if is_empty_str "$env_name"; then 
        exit_error "Env name is empty"
    fi

    if file_exist "$disk_img_path"; then 
        exit_error "Env of this name already exist"
    fi

    if ! is_empty_str "$(echo "$env_size" | sed "s|[0-9]||g")"; then
        exit_error "Invalid disk image size, Expecting an integer"
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

    [ $# = 0 ] && internal_err "Usage: resize_disk_image [NAME] [SIZE] [FORMAT=ext4]"

    env_name="$1"
    env_size="${2:-2048}"
    disk_format="${3:-"ext4"}"
    disk_img_path="$(img_path_for "$env_name")"

    if is_empty_str "$env_name"; then 
        exit_error "Env name is empty"
    fi

    if ! file_exist "$disk_img_path"; then 
        exit_error "Env of this name does not exist"
    fi

    if ! is_empty_str "$(echo "$env_size" | sed "s|[0-9]||g")"; then
        exit_error "Invalid disk image size, Expecting an integer"
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
}

mount_disk_img() {
    
    [ $# = 0 ] && internal_err "Usage: mount_disk_img [ENV_NAME]"

    env_name="$1"
    disk_img_path="$(img_path_for "$env_name")" 
    mount_path="$(mount_path_for "$env_name")"

    if ! file_exist "$disk_img_path"; then 
        exit_error "Mount Failed, Env of this name does not exist"
    fi

    if mountpoint -q "$mount_path"; then
        print_info "Target is already not a mountpoint"
        return 0
    fi

    unshared_from "$env_name"

    use_namespace "$env_name" bash -c "
            mkdir -p \"$mount_path\"
            mkdir -p \"$mount_path/{dev,proc,sys,etc,dev/pts}\"
            mount -o loop \"$disk_img_path\" \"$mount_path\"
            mount --bind /dev \"$mount_path\"/dev
            mount -t proc proc \"$mount_path\"/proc
            mount -t sysfs sys \"$mount_path\"/sys
            mount -t devpts devpts \"$mount_path\"/dev/pts
    "

    if [ ! -f "$mount_path/etc/resolv.conf" ] && [ -f /etc/resolv.conf ]; then
        cp /etc/resolv.conf "$mount_path/etc/resolv.conf"
    fi
    
    print_info "Mounted $env_name at $mount_path"

}

enter_disk_image() {
    env_name="$1"
    mount_path="$(mount_path_for "$env_name")"

    mountpoint -q "$mount_path" || mount_disk_img "$env_name"

    unshare --fork --pid bash -c "chroot \"$mount_path\" /bin/bash"
    umount -R "$mount_path"
}

# --- COMMANDS ---

cmd_create() {
    env_name="$1"
    env_size="${2:-""}"
    create_disk_image "$env_name" "$env_size"
}

cmd_mount() {
    env_name="$1"
    env_size="${2:-""}"
    mount_disk_img "$env_name" "$env_size"
}

cmd_resize() {
    env_name="$1"
    env_size="${2:-""}"
    resize_disk_image "$env_name" "$env_size"
}

cmd_enter() {
    env_name="$1"
    env_size="${2:-""}"
    enter_disk_image "$env_name"
}

cmd_flash() {
    mount_path="$(mount_path_for "$1")"
    shift

    if ! mountpoint -q "$mount_path"; then
        exit_error "Target is not a mountpoint. Aborting to save your host disk!"
    fi

    sleep 1000

    if [ $# -eq 0 ]; then
        pacstrap -K "$mount_path" base bash coreutils nano git sudo
    else
        pacstrap -K "$mount_path" "$@"
    fi
}

cmd_help() {
    echo ""
}

# --- SETUP ---

dispatch(){
    cmd="cmd_$(echo "$1" | tr "-" "_")"
    shift # remove the first arument of this

    if command -v "$cmd" >/dev/null 2>&1; then
        $cmd "$@"
    else 
        echo "Unknown command: $cmd"
        cmd_help
        exit 1
    fi
}

main () {
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