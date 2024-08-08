TARGET_ARCHIVE=./data/archives/orig_data.zip
TESTGROUND=./data/testground

echo "Deleting testground: ${TESTGROUND}"
rm -rf "$TESTGROUND"

echo "Recreating testground: ${TESTGROUND}"
mkdir -p "$TESTGROUND"

touch "$TESTGROUND/.gitkeep"

unzip "$TARGET_ARCHIVE" -d "$TESTGROUND"
