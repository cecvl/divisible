# Divisible

 A small reflex / number game — decide whether a number is divisible by 3.

 Features
- Play quick rounds against a timer.
- Simple keyboard controls and bonus questions.
- Scores persisted locally using SQLite; release binaries published via GitHub Releases.

 Installation

- Install from releases (recommended):

   ```bash
   curl -sSL https://raw.githubusercontent.com/cecvl/divisible/main/install.sh | bash -s --
   ```

- Install a specific release tag:

   ```bash
   curl -sSL https://raw.githubusercontent.com/cecvl/divisible/main/install.sh | bash -s -- v0.1.0
   ```

 Build from source

 ```bash
 git clone https://github.com/cecvl/divisible.git
 cd divisible
 go build -o divisible ./cmd/game
 ./divisible
 ```

 Run (dev)

- Use `make run` (project Makefile) or run the built binary.

 Controls

- `Y`: answer "divisible by 3".
- `N`: answer "not divisible by 3".
- `1` / `2`: choose bonus answers when prompted.
- `P`: pause/resume.
- `R`: restart after game over.

 SQLite score DB

- The game stores scores locally in `scores.db` (in the working directory). It records each play and maintains basic stats (plays, total_score, last_score, best_score).

 CI / Releases

- GitHub Actions builds release binaries for supported OS/arch pairs and uploads them to GitHub Releases on tag pushes.
- The included `install.sh` expects release assets named like `divisible_<tag>_<os>_<arch>.tar.gz`.

 Contributing

- Open issues or PRs on the repository. For build/release changes see `.github/workflows/release.yml`.

 License

- MIT (see LICENSE file)
