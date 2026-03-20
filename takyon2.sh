#!/usr/bin/env sh

set -eu
SHARED_DIR="$HOME/.local/bin/takyon"
ENV_STORE="$SHARED_DIR/images"
ENV_MOUNT="/mnt/takyon"

is_empty_str() {
    [ -z "$1" ]
}

file_exist() {
    [ -f "$1" ]
}

dir_exist() {
    [ -d "$1" ]
}

mount_exist(){
    mountpoint -q "$1"
}

store_path_for() {
    env_name="$ENV_STORE/$1"
    echo "$env_name"
}

mount_path_for() {
    mount_path="$ENV_MOUNT/$1"
    echo "$mount_path"
}

print_info() {
    echo "[*] $1"
}

error_exit() {
    echo "[!] $1"
    exit 1
}

create_disk_img() {
    env_name="$1"
    env_size="${2:-"2048"}"
    env_format="${3:-"ext4"}"
    img_path="$(store_path_for "$env_name")"
    
    is_empty_str "$env_name" && error_exit ""
    file_exist "$img_path" && error_exit ""
    
    # remove the num chars in the `env_size`
    null_size="$(echo "$env_size" | sed "s|[0-9]||g")"
    is_empty_str "$null_size" ||  error_exit ""

    truncate -s "${env_size}M" "$img_path"
    mkfs."$env_format" "$img_path"
}

mount_disk_img() {
    env_name="$1"
    env_shell="${2:-"bash"}"
    extra_cmd="${3:""}"
    img_path="$(store_path_for "$env_name")"
    mount_path="$(mount_path_for  "$env_name")"

    file_exist "$img_path" || error_exit ""
    mount_exist "$mount_path" && error_exit ""

    mkdir -p "$mount_path"
    mount -o loop "$img_path" "$mount_path"
    mkdir -p "$mount_path/{dev,etc,proc,sys}"
    mkdir -p "$mount_path/dev/pts"

    host_resolv="/etc/resolv.conf"
    env_resolv="$mount_path/etc/resolv.conf"
    
    if [ ! -f "$env_resolv" ] && [ -f "$host_resolv" ]; then
        cp "$host_resolv" "$env_resolv"
    fi
}

enter_disk_img() {
    
    env_name="$1"
    username="${2:-"root"}"
    env_shell="${3-"/bin/bash"}"
    extra_cmd="${4:""}"
    img_path="$(store_path_for "$env_name")"
    mount_path="$(mount_path_for  "$env_name")"

    file_exist "$img_path" || error_exit ""
    mount_exist "$mount_path" || mount_disk_img "$env_name"
    
    (
        trap 'sleep 1s && umount -R "$mount_path" && rm -rf "$mount_path"' EXIT
        
        final_script=""
        mount_script="
            mount -t proc proc \"$mount_path/proc\" && \
            mount -t sysfs sys \"$mount_path/sys\" && \
            mount --bind /dev \"$mount_path/dev\" && \
            mount -t devpts devpts \"$mount_path/dev/pts\"
        "
        if is_empty_str "$extra_cmd"; then
            final_script="
                $mount_script
                chroot --userspec $username $mount_path $env_shell
            "
        else
            final_script="
                $mount_script
                chroot --userspec $username $mount_path $env_shell -c \"$extra_cmd\"
            "
        fi

        unshare --mount --pid --net --fork --cgroup "$final_script"
    )
}

add_user() {
    env_name="$1"
    username="$2"
    cmd="
        useradd -mG wheel,audio,video $username
        passwd $username
    "
    enter_disk_img "$env_name" "root" "/bin/bash" "$cmd"
}