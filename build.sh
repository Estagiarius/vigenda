#!/bin/bash

# Script to build Vigenda for multiple platforms

# Create a dist directory if it doesn't exist
mkdir -p dist

# Target platforms
PLATFORMS=(
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
)

# Go source path
SOURCE_PATH="./cmd/vigenda/"

# Iterate over platforms and build
for platform in "${PLATFORMS[@]}"
do
    # Split platform into OS and ARCH
    GOOS_VAL=$(echo $platform | cut -d'/' -f1)
    GOARCH_VAL=$(echo $platform | cut -d'/' -f2)

    # Set output name
    OUTPUT_NAME="dist/vigenda-${GOOS_VAL}-${GOARCH_VAL}"
    if [ "${GOOS_VAL}" == "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi

    echo "Building for ${GOOS_VAL}/${GOARCH_VAL}..."

    # Build command
    # The -ldflags="-s -w" flags are used to strip debug information and reduce binary size.
    # CGO_ENABLED=0 is often used for cross-compilation to avoid issues with C dependencies,
    # but go-sqlite3 requires CGO. We will manage C compilers for targets if needed.
    # For go-sqlite3, specific C compilers might be needed for cross-compilation.
    # For example, for Windows, a MinGW compiler (e.g., x86_64-w64-mingw32-gcc) is needed.
    # For macOS cross-compilation from Linux, appropriate SDKs and compilers are needed.

    # For now, let's try with CGO_ENABLED=1 and see if the default GCC can handle it or if we need specific cross-compilers.
    # The go-sqlite3 documentation mentions that for cross-compiling, you usually need to set CC environment variable.
    # Example: CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build

    # We will attempt a simple build first. If CGO issues arise, we'll address them.
    # The user has installed gcc, which might handle linux builds.
    # Cross-compilation for Windows and macOS with CGO dependencies from Linux can be complex.

    CMD_PREFIX=""
    if [ "${GOOS_VAL}" == "windows" ] && [ "${GOARCH_VAL}" == "amd64" ]; then
        CMD_PREFIX="CC=x86_64-w64-mingw32-gcc "
    elif [ "${GOOS_VAL}" == "windows" ] && [ "${GOARCH_VAL}" == "386" ]; then # Example for 32-bit windows if needed
        CMD_PREFIX="CC=i686-w64-mingw32-gcc "
    fi

    # For macOS (darwin), true cross-compilation from Linux for CGO projects typically requires
    # a macOS SDK and a specially built clang. This is hard to set up in a generic Linux CI environment.
    # Tools like 'osxcross' attempt to address this. For now, we'll skip explicit CC for darwin
    # as it's unlikely to work without a more complex setup. The previous errors for darwin confirm this.

    COMMAND="${CMD_PREFIX}GOOS=${GOOS_VAL} GOARCH=${GOARCH_VAL} CGO_ENABLED=1 go build -ldflags=\"-s -w\" -o ${OUTPUT_NAME} ${SOURCE_PATH}"
    echo "Executing: ${COMMAND}"
    eval ${COMMAND}

    if [ $? -eq 0 ]; then
        echo "Successfully built ${OUTPUT_NAME}"
    else
        echo "Error building for ${GOOS_VAL}/${GOARCH_VAL}"
        # Optionally, exit on error:
        # exit 1
    fi
    echo ""
done

echo "All builds completed. Binaries are in the dist/ directory."
