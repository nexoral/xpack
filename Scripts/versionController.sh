#!/bin/bash
# Fetch remote version
remote_url="https://raw.githubusercontent.com/nexoral/xpack/main/VERSION"
remote_version=$(curl -s "$remote_url")
if [[ -z "$remote_version" ]]; then
    echo "Error: Unable to fetch remote version."
    exit 1
fi
echo "Current GitHub version: $remote_version"

# Read local version
local_version=$(cat "$(dirname "$0")/../VERSION" 2>/dev/null || echo "0.0.0")
echo "Local version: $local_version"

# Compare versions: returns 0 if first > second
ver_gt() {
    local IFS=.
    local raw1 raw2 i ver1 ver2
    # Strip suffix from versions (ignore -beta or -stable)
    raw1="${1%%-*}"
    raw2="${2%%-*}"
    # parse numeric parts
    read -ra ver1 <<<"$raw1"
    read -ra ver2 <<<"$raw2"
    # pad shorter array with zeros
    for ((i = ${#ver1[@]}; i < ${#ver2[@]}; i++)); do ver1[i]=0; done
    for ((i = ${#ver2[@]}; i < ${#ver1[@]}; i++)); do ver2[i]=0; done
    # compare each segment
    for ((i = 0; i < ${#ver1[@]}; i++)); do
        if ((10#${ver1[i]} > 10#${ver2[i]})); then return 0; fi
        if ((10#${ver1[i]} < 10#${ver2[i]})); then return 1; fi
    done
    return 1
}

# Exit if local is already ahead of remote
if ver_gt "$local_version" "$remote_version"; then
    echo "Local version ($local_version) is ahead of remote ($remote_version). Skipping update."
    exit 0
fi

# ---- now select version type ----
options=("Stable" "Beta")
selected=0

tput civis # hide cursor
show_menu() {
    for i in "${!options[@]}"; do
        if [[ $i -eq $selected ]]; then
            echo -e "\033[7m> ${options[$i]}\033[0m"
        else
            echo "  ${options[$i]}"
        fi
    done
}
refresh() {
    tput cup 4 0
    tput ed
    show_menu
}

show_menu
while true; do
    read -rsn3 key
    case "$key" in
    $'\x1b[A') # Up arrow
        ((selected = (selected - 1 + ${#options[@]}) % ${#options[@]}))
        refresh
        ;;
    $'\x1b[B') # Down arrow
        ((selected = (selected + 1) % ${#options[@]}))
        refresh
        ;;
    "") # Enter key
        version_type="${options[$selected]}"
        break
        ;;
    esac
done
tput cnorm # show cursor
echo -e "\nSelected: $version_type"

# Prompt for new version
read -p "Enter new version: " new_version

# Validate version format
if ! [[ "$new_version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Invalid version format. Use X.Y.Z (e.g., 1.0.0)."
    exit 1
fi

# Update local version file
suffix=$(echo "$version_type" | tr '[:upper:]' '[:lower:]')
echo "${new_version}-${suffix}" >"$(dirname "$0")/../VERSION"
echo "Local version updated to ${new_version}-${suffix}"


# Update static version in main.go (var VERSION) and Banner.go (const Version)
main_go="$(dirname "$0")/../src/Core/main.go"
# Replace lines like: var VERSION = "..."
sed -i "s|^\([[:space:]]*\)var VERSION = .*|\1var VERSION = \"${new_version}-${suffix}\"|" "$main_go"

banner_go="$(dirname "$0")/../src/base/banner.go"
sed -i "s|^\([[:space:]]*\)const Version = .*|\1const Version = \"${new_version}-${suffix}\"|" "$banner_go"

echo "Updated main.go and banner.go with new version."

# Update version in README.md installation URLs
readme_file="$(dirname "$0")/../README.md"
sed -i "s|^\(\s*wget .*/releases/download/v\)[^/]*\(/xpack_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2${new_version}-${suffix}\3|" "$readme_file"
sed -i "s|^\(\s*sudo dpkg -i xpack_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2|" "$readme_file"

# Update version in installation.md file
installation_file="$(dirname "$0")/../INSTALLATION.md"
sed -i "s|^\(\s*wget .*/releases/download/v\)[^/]*\(/xpack_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2${new_version}-${suffix}\3|" "$installation_file"
sed -i "s|^\(\s*sudo dpkg -i xpack_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2|" "$installation_file"

echo "Updated INSTALLATION.md with new version in installation URLs."

# Update version in LEARN.md file
learn_file="$(dirname "$0")/../LEARN.md"
sed -i "s|^\(\s*wget .*/releases/download/v\)[^/]*\(/xpack_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2${new_version}-${suffix}\3|" "$learn_file"
sed -i "s|^\(\s*sudo dpkg -i xpack_\)[^_]*\(_amd64\.deb\)|\1${new_version}-${suffix}\2|" "$learn_file"

echo "Updated LEARN.md with new version in installation URLs."

# Update version in SCRIPTS/installer.sh
installer_file="$(dirname "$0")/../Scripts/installer.sh"
sed -i "s|^\(VERSION=\"\)[^\"']*\(\".*\)$|\1${new_version}-${suffix}\2|" "$installer_file"
echo "Updated installer.sh with new version."