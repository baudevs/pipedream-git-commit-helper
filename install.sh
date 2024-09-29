#!/bin/bash

# Name of the tool
TOOL_NAME="pdcommit"
TOOL_MAN_PAGE="src/man/pdcommit.1"

# Define user-level directories
USER_BIN_DIR="$HOME/bin"
USER_MAN_DIR="$HOME/man/man1"

# Create the user-level directories if they don't exist
mkdir -p $USER_BIN_DIR
mkdir -p $USER_MAN_DIR

# Build the Go binary
echo "Building the $TOOL_NAME binary..."
if go build -o $TOOL_NAME; then
  echo "$TOOL_NAME binary built successfully."
else
  echo "Failed to build $TOOL_NAME. Please ensure Go is installed."
  exit 1
fi

# Move the binary to ~/bin
echo "Installing $TOOL_NAME to $USER_BIN_DIR..."
if mv $TOOL_NAME $USER_BIN_DIR/; then
  echo "$TOOL_NAME installed successfully in $USER_BIN_DIR."
else
  echo "Failed to install $TOOL_NAME. Please check your permissions."
  exit 1
fi

# Check if the man page exists and install it
if [ -f $TOOL_MAN_PAGE ]; then
  echo "Installing $TOOL_NAME man page to $USER_MAN_DIR..."
  if cp $TOOL_MAN_PAGE $USER_MAN_DIR/; then
    echo "Man page installed successfully."
  else
    echo "Failed to install the man page. Please check your permissions."
    exit 1
  fi
else
  echo "Man page ($TOOL_MAN_PAGE) not found. Skipping man page installation."
fi

# Check if ~/bin is in the user's PATH, and add it if it's not
if [[ ":$PATH:" != *":$USER_BIN_DIR:"* ]]; then
  echo "Adding $USER_BIN_DIR to your PATH..."
  echo "export PATH=\"$USER_BIN_DIR:\$PATH\"" >> ~/.bashrc
  echo "export PATH=\"$USER_BIN_DIR:\$PATH\"" >> ~/.zshrc
  export PATH="$USER_BIN_DIR:$PATH"
  echo "Your PATH has been updated."
fi

# Check if ~/man is in the user's MANPATH, and add it if it's not
if [[ ":$MANPATH:" != *":$HOME/man:"* ]]; then
  echo "Adding $HOME/man to your MANPATH..."
  echo "export MANPATH=\"$HOME/man:\$MANPATH\"" >> ~/.bashrc
  echo "export MANPATH=\"$HOME/man:\$MANPATH\"" >> ~/.zshrc
  export MANPATH="$HOME/man:$MANPATH"
  echo "Your MANPATH has been updated."
fi

echo "$TOOL_NAME installation complete!"
echo "Please restart your terminal or run 'source ~/.bashrc' or 'source ~/.zshrc' to apply the changes."