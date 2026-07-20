# SteamFlush

SteamFlush is a lightweight, high-performance utility designed to identify and remove orphan shadercache folders from your Steam installation. Primarily targeted at **Steam Deck** users and Linux gaming enthusiasts, it helps reclaim significant disk space by cleaning up leftover cache files from uninstalled games.

## Features

- **Automated Scanning:** Quickly identifies orphan shadercache directories.
- **Mass Deletion:** Supports batch processing to free up disk space in seconds.
- **Steam Integration:** Fetches game names via Steam API to provide context for the files you are deleting.
- **Safe Operations:** Implements strict path validation to prevent path traversal and accidental system file deletion.
- **Modern UI:** Built with Svelte 5 for a responsive, clean, and intuitive interface.

## Technical Stack

- **Backend:** Go (Wails framework)
- **Frontend:** Svelte 5 (TypeScript)
- **Styling:** CSS
- **API:** Steam Web API

## Installation

### Prerequisites
- Go 1.21+
- [Wails](https://wails.io/) installed (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- Node.js & npm/pnpm

### Build
1. Clone the repository:
   ```bash
   git clone <your-repository-url>
   cd steamflush
   ```

2. Install frontend dependencies:
    ```bash
    cd frontend
    npm install
    ```

3. Build the application:
    ```bash
    # From the project root
    wails build
    ```

## Steam Deck Optimization

SteamFlush is designed to work seamlessly on the Steam Deck, both in Desktop Mode and as a "Non-Steam Game".

#### Desktop Mode

The interface is optimized for touch and trackpad interaction with 40px+ clickable targets. Use the trackpad as a mouse to navigate and click.

#### Steam Gaming Mode (Non-Steam Game)

To use SteamFlush directly from your Game Library, we recommend adding it as a "Non-Steam Game" and configure **Steam Input**:

1. **Add to Library:** In Desktop Mode, right-click the compiled `steamflush` binary and "Add to Steam".
2. **Configure Controller Layout:**
    * Open Steam, go to the application properties, and edit the **Controller Layout**.
    * **Recommended Mapping:**
        * **D-Pad:** Map to `Tab` (Forward focus) and `Shift+Tab` (backward focus).
        * **A Button:** Map to `Enter` (Click/Confirm).
        * **B Button:** Map to `Escape` (Cancel/Close).
    * This allows you to navigate the application entirely using the physical controls of the Steam Deck.

## Security & Reliability

* **Sandboxing:** Path operations are strictly restricted to your specific Steam shadercache directory.
* **Atomic Operations:** Uses batch processing with error handling to ensure data integrity during mass deletions.
* **Rate Limiting:** Includes robust handling for Steam API rate limits and local caching to minimize network overhead.

## License

This project is licensed under the **GNU General Public License v3.0 (GPL-3.0)**. See the [LICENCE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please follow the standard workflow:

1. Fork the repo.
2. Create a feature branch.
3. Commit your changes.
4. Open a Pull Request.