#!/usr/bin/env sh

# --- Configuration ---

set -eu
ENV_IMAGE_STORE="./temp/envs"

# --- Helpers ---

setup(){
    mkdir -p "$ENV_IMAGE_STORE"
}

print_info(){
    echo "[*] $1" 
}

error_exit() {
    echo "[!] $1"
    exit 1
}

create_sparse_disk_img() {
    name="$1"
    format="$2"
    size="${3:-2048}"

    image="$ENV_IMAGE_STORE/$name.img"

    print_info "Creating sparse disk image $image (${size}MB)"

    truncate -s "${size}M" "$image"

    print_info "Formatting $image as $format"

    mkfs."$format" "$image" >/dev/null
}

bind_virtual_filesystem() {
    target="$1"
    for item in dev proc sys; do
        mkdir -p "$target/$item"
    done
    mount --bind /dev "$target/dev"
    mount -t proc proc "$target/proc"
    mount -t sysfs sys "$target/sys"
}

get_env_path() {
    env="$1"
    image="$ENV_IMAGE_STORE/$env.img"
    [ -f "$image" ] && echo "$image"
}

get_env_mount_path() {
    env="$1"
    image="$ENV_IMAGE_STORE/$env.img"
    mount_dir="/mnt/$env"
    [ -f "$image" ] && [ -d "$mount_dir" ] && echo "$mount_dir"
}

# --- Commands ---

cmd_create() {
    env=$1
    size="${2:-2048}" # $2 or 2048MB
    create_sparse_disk_img "$env" ext4 "$size"
}

cmd_list() {
    print_info "Available envs:"
    for item in "$ENV_IMAGE_STORE"/*; do
        [ -e "$item" ] || continue
        echo "  - $(basename "$item" .img)"
    done
}

cmd_mount() {
    env="$1"
    mount_dir="/mnt/$env"
    image="$ENV_IMAGE_STORE/$env.img"
    [ -f "$image" ] || error_exit "Env $env does not exist"

    mkdir -p "$mount_dir"
    mount -o loop "$image" "$mount_dir"
    bind_virtual_filesystem "$mount_dir"
    print_info "Mounted $env at $mount_dir"
}

cmd_enter(){

    env="$1"
    image="$ENV_IMAGE_STORE/$env.img"
    mount_dir="/mnt/$env"

    [ -f "$image" ] || error_exit "Env does not exist"
    
    sudo unshare --mount --fork --pid -- bash -c "
        set -e

        mkdir -p $mount_dir
        mount -o loop $image $mount_dir

        mkdir -p $mount_dir/{dev,proc,sys,dev/pts}

        mount --bind /dev $mount_dir/dev
        mount -t proc proc $mount_dir/proc
        mount -t sysfs sys $mount_dir/sys
        mount -t devpts devpts $mount_dir/dev/pts
        mount --bind /etc/resolv.conf $mount_dir/etc/resolv.conf # net resolver
        chroot $mount_dir /bin/bash
    "

}

cmd_umount() {
    env="$1"
    mount_dir="/mnt/$env"
    image="$ENV_IMAGE_STORE/$env.img"
    [ -f "$image" ] || error_exit "Env $env does not exist"
    [ -d "$mount_dir" ] || error_exit "Env $env is not mounted"

    umount -R "$mount_dir"    
    rm -rf "$mount_dir"
    print_info "Umounted $env at $mount_dir"
}

cmd_init() {
    [ $# -lt 1 ] && error_exit "Usage: takyon init <env> [packages]"

    env_name="$1"
    env="$(get_env_path "$env_name")"
    mount_dir="$(get_env_mount_path "$env_name")"

    [ -z "$env" ] && error_exit "Env does not exist"
    [ -z "$mount_dir" ] && error_exit "Env is not mounted"

    shift
    if [ $# -eq 0 ]; then
        pacstrap -K "$mount_dir" base bash coreutils nano git sudo
    else
        pacstrap -K "$mount_dir" "$@"
    fi
}

cmd_help(){
    cat << EOF
takyon: An env isolation/management tool

Usage: tackyon [COMMAND]

Commands:
    list                    list the given available env
    create [NAME] [SIZE]    create a new env empty
    mount [NAME]            mount an env
    umount [NAME]           umount an env
EOF
}

# --- Infrastructure ---

dispatch() {
    raw="$1"
    cmd="$(echo "cmd_$raw" | tr "-" "_")"
    shift # remove the 1st arg of `dispath`

    if command -v "$cmd" >/dev/null 2>&1; then
        $cmd "$@"
    else
        echo "Unknown command: $raw"
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

    setup
    dispatch "$@"
}

main "$@"